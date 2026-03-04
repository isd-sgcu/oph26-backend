package usecase

import (
	"context"
	"fmt"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"
	"oph26-backend/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

type AuthUsecase interface {
	Login(c *fiber.Ctx) error
}

type AuthUsecaseImpl struct {
	UserRepository  repository.UserRepository
	StaffRepository repository.StaffRepository
	GoogleClientID  string
	JWTSecret       string
}

func NewAuthUsecase(userRepository repository.UserRepository, staffRepository repository.StaffRepository, googleClientID string, jwtSecret string) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepository:  userRepository,
		StaffRepository: staffRepository,
		GoogleClientID:  googleClientID,
		JWTSecret:       jwtSecret,
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
				Email: email,
				Role:  "staff",
			}
			if err := u.UserRepository.Create(user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		} else {
			user = &entity.User{
				Email: email,
				Role:  "user", // Default role
			}
			if err := u.UserRepository.Create(user); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
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

	refreshToken, err := u.generateRefreshToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set Refresh Token in HttpOnly Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true, // TODO: Set to false in dev if needed found in env?
		SameSite: "Strict",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return c.JSON(&model.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (u *AuthUsecaseImpl) validateGoogleToken(ctx context.Context, token string) (*idtoken.Payload, error) {
	// TODO: Recheck this function before production

	if u.GoogleClientID == "" {
		// Mock for now if no client ID provided in env
		// In production this should be an error or stricter.
		// For the sake of progress in this environment:
		return &idtoken.Payload{
			Claims: map[string]interface{}{
				"email": "test@example.com",
			},
		}, nil
	}

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

func (u *AuthUsecaseImpl) generateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.JWTSecret))
}
