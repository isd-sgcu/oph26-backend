package usecase

import (
	"oph26-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

type UserUsecase interface {
	PutAttendeesUseCase(c *fiber.Ctx) error
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUsecaseImpl) PutAttendeesUseCase(c *fiber.Ctx) error {
	return nil
}
