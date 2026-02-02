package repository

import (
	"oph26-backend/internal/config"
	"oph26-backend/internal/entity"
	"oph26-backend/internal/model"

	"github.com/google/uuid"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAllUsers() ([]entity.User, error) {
	var dbUsers []model.User
	err := config.DB.Find(&dbUsers).Error
	if err != nil {
		return nil, err
	}

	// Convert models to entities
	entities := make([]entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		entities[i] = entity.User{
			ID:         dbUser.ID,
			Email:      dbUser.Email,
			Role:       dbUser.Role,
			AttendeeId: dbUser.AttendeeId,
			StaffId:    dbUser.StaffId,
			CreatedAt:  dbUser.CreatedAt,
			UpdatedAt:  dbUser.UpdatedAt,
		}
	}

	return entities, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var dbUser model.User
	err := config.DB.First(&dbUser, id).Error
	if err != nil {
		return nil, err
	}

	// Convert model to entity
	entity := &entity.User{
		ID:         dbUser.ID,
		Email:      dbUser.Email,
		Role:       dbUser.Role,
		AttendeeId: dbUser.AttendeeId,
		StaffId:    dbUser.StaffId,
		CreatedAt:  dbUser.CreatedAt,
		UpdatedAt:  dbUser.UpdatedAt,
	}

	return entity, nil
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	// Convert entity to model
	dbUser := &model.User{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		AttendeeId: user.AttendeeId,
		StaffId:    user.StaffId,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	return config.DB.Create(dbUser).Error
}

func (r *UserRepository) UpdateUser(user *entity.User) error {
	// Convert entity to model
	dbUser := &model.User{
		ID:         user.ID,
		Email:      user.Email,
		Role:       user.Role,
		AttendeeId: user.AttendeeId,
		StaffId:    user.StaffId,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}

	return config.DB.Save(dbUser).Error
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	return config.DB.Delete(&model.User{}, id).Error
}
