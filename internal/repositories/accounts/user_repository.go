package accounts

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kite/internal/models"
)

type UserRepository interface {
	FindByUserId(ctx context.Context, userId uint) (*models.Users, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUserId(ctx context.Context, userId uint) (*models.Users, error) {
	var user models.Users

	result := r.db.WithContext(ctx).Where("userId = ?", userId).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}
