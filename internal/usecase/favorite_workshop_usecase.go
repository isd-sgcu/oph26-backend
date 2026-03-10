package usecase

import (
	"oph26-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FavoriteWorkshopUseCaseImpl struct {
	userRepo     repository.UserRepository
	attendeeRepo repository.AttendeeRepository
	validate     *validator.Validate
}

type FavoriteWorkshopUseCase interface {
	GetMyFavWorkshops(c *fiber.Ctx) error
	PutMyFavWorkshops(c *fiber.Ctx) error
}

func NewFavoriteWorkshopUsecase(userRepo repository.UserRepository, attendeeRepo repository.AttendeeRepository) FavoriteWorkshopUseCase {
	return &FavoriteWorkshopUseCaseImpl{
		userRepo:     userRepo,
		attendeeRepo: attendeeRepo,
		validate:     validator.New(),
	}
}

func (u *FavoriteWorkshopUseCaseImpl) GetMyFavWorkshops(c *fiber.Ctx) error {
	return nil
}

func (u *FavoriteWorkshopUseCaseImpl) PutMyFavWorkshops(c *fiber.Ctx) error {
	return nil
}
