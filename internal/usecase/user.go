package usecase

import (
	"oph26-backend/internal/model/user"
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type UserUsecase struct {
	userRepo *repository.UserRepository
}

func NewUserUsecase(userRepo *repository.UserRepository) *UserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func (u *UserUsecase) GetAllUsers(c *fiber.Ctx) error {
	// Use Case calls Repository with Entity data
	entities, err := u.userRepo.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	// Use Case creates Model for response
	response := user.GetAllUsersResponse{
		Users: make([]user.UserResponse, len(entities)),
	}

	for i, entity := range entities {
		response.Users[i] = user.UserResponse{
			ID:         entity.ID,
			Email:      entity.Email,
			Role:       entity.Role,
			AttendeeId: entity.AttendeeId,
			StaffId:    entity.StaffId,
			CreatedAt:  entity.CreatedAt,
			UpdatedAt:  entity.UpdatedAt,
		}
	}

	return c.JSON(response)
}
