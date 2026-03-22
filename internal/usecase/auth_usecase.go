package usecase

import (
	"context"
	"fmt"
	"log"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	"oph26-backend/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/api/idtoken"
)

type AuthUsecase interface {
	Login(c *fiber.Ctx) error
	GetCurrentUser(c *fiber.Ctx) error
	RefreshToken(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
}

type AuthUsecaseImpl struct {
	UserRepository         repository.UserRepository
	StaffRepository        repository.StaffRepository
	AttendeeRepository     repository.AttendeeRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	GoogleClientID         string
	JWTSecret              string
	AppEnv                 string
}

type AuthUsecaseConfig struct {
	UserRepository         repository.UserRepository
	StaffRepository        repository.StaffRepository
	AttendeeRepository     repository.AttendeeRepository
	RefreshTokenRepository repository.RefreshTokenRepository
	GoogleClientID         string
	JWTSecret              string
	AppEnv                 string
}

func NewAuthUsecase(config AuthUsecaseConfig) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepository:         config.UserRepository,
		StaffRepository:        config.StaffRepository,
		AttendeeRepository:     config.AttendeeRepository,
		RefreshTokenRepository: config.RefreshTokenRepository,
		GoogleClientID:         config.GoogleClientID,
		JWTSecret:              config.JWTSecret,
		AppEnv:                 config.AppEnv,
	}
}

func (u *AuthUsecaseImpl) Login(c *fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	ctx := c.Context()

	// 1. Validate Google ID Token
	payload, err := u.validateGoogleToken(ctx, request.IDToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("invalid google token: %v", err),
		})
	}

	email := payload.Claims["email"].(string)

	// 2. Find or Create User
	user, err := u.UserRepository.FindByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		// Find user in staff database by email, if not found create new user with default role
		staffUser, err := u.StaffRepository.FindByEmail(email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if staffUser != nil {
			user = &entity.User{
				Email:   email,
				Role:    "staff",
				StaffId: &staffUser.ID,
			}
			if err := u.UserRepository.Create(user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			// Update staff with user_id
			staffUser.UserID = &user.ID
			if err := u.StaffRepository.Update(staffUser); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		} else {
			user = &entity.User{
				Email: email,
				Role:  "attendee", // Default role
			}
			if err := u.UserRepository.Create(user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}
	} else {
		// If user exists, still check if they are staff to set the role correctly (in case they were created before as attendee)
		if user.Role != "staff" {
			staffUser, err := u.StaffRepository.FindByEmail(email)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if staffUser != nil {
				user.Role = "staff"
				user.StaffId = &staffUser.ID
				if err := u.UserRepository.Update(user); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				// Update staff with user_id if not already set
				if staffUser.UserID == nil {
					staffUser.UserID = &user.ID
					if err := u.StaffRepository.Update(staffUser); err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"error": err.Error(),
						})
					}
				}
			}
		}
		// Check if user is attendee and set AttendeeId if not set
		if user.Role == "attendee" && user.AttendeeId == nil {
			attendee, err := u.AttendeeRepository.FindByUserID(user.ID)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if attendee != nil {
				user.AttendeeId = &attendee.ID
				if err := u.UserRepository.Update(user); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
			}
		}
	}

	// 3. Generate Tokens
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	refreshToken, expiresAt, err := u.generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := u.RefreshTokenRepository.Store(&entity.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to persist refresh token",
		})
	}

	// Set Refresh Token in HttpOnly Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   u.AppEnv == "production", // Only set Secure flag in production
		SameSite: "Strict",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return c.JSON(&model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (u *AuthUsecaseImpl) validateGoogleToken(ctx context.Context, token string) (*idtoken.Payload, error) {
	if u.AppEnv == "development" {
		// In development, allow user to bypass Google token validation for easier testing.
		// Token is jwt but we won't validate it, just parse the claims for email.
		jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			// We won't validate the signature in development, so just return a dummy key.
			return []byte("dummy"), nil
		})

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("(dev) invalid token claims")
		}

		log.Printf("(dev) bypassing Google token validation, claims: %v", claims)

		return &idtoken.Payload{
			Claims: map[string]any{
				"email": claims["email"],
			},
			Audience: "development-audience",
		}, nil
	}

	// In production, validate the token properly against Google's OAuth2 service.
	return idtoken.Validate(ctx, token, u.GoogleClientID)
}

func (u *AuthUsecaseImpl) generateAccessToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"iss":         "isd-oph26-backend",
		"sub":         user.ID,
		"user_id":     user.ID,
		"email":       user.Email,
		"role":        user.Role,
		"attendee_id": user.AttendeeId,
		"staff_id":    user.StaffId,
		"exp":         time.Now().Add(time.Minute * 15).Unix(), // 15 minutes
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.JWTSecret))
}

func (u *AuthUsecaseImpl) generateRefreshToken(user *entity.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Hour * 24 * 7)
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiresAt.Unix(), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.JWTSecret))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expiresAt, nil
}

func (u *AuthUsecaseImpl) GetCurrentUser(c *fiber.Ctx) error {
	// Extract user_id from context (set by auth middleware)
	userIDValue := c.Locals("user_id")
	if userIDValue == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing user_id in token",
		})
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid user_id in token",
		})
	}

	// Fetch user from database
	user, err := u.UserRepository.FindByID(userID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(&model.UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	})
}

func (u *AuthUsecaseImpl) RefreshToken(c *fiber.Ctx) error {
	// Get refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing refresh token",
		})
	}

	// Validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(u.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid or expired refresh token",
		})
	}

	isActive, err := u.RefreshTokenRepository.IsActive(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to validate refresh token",
		})
	}

	if !isActive {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Refresh token revoked or unknown",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Invalid token claims",
		})
	}

	// Extract user_id from refresh token claims
	userID := claims["user_id"]
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: Missing user_id in refresh token",
		})
	}

	// Fetch user from database
	user, err := u.UserRepository.FindByID(userID.(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Generate new access token
	accessToken, err := u.generateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(&model.TokenResponse{
		AccessToken: accessToken,
	})
}

func (u *AuthUsecaseImpl) SignOut(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken != "" {
		if err := u.RefreshTokenRepository.Revoke(refreshToken); err != nil {
			log.Printf("failed to revoke refresh token: %v", err)
		}
	}

	// Clear the refresh_token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   u.AppEnv == "production",
		SameSite: "Strict",
		MaxAge:   -1, // Immediately delete the cookie
	})

	return c.SendStatus(fiber.StatusNoContent)
}
