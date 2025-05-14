package accounts

import (
	"context"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"kite/internal/database"
	"kite/internal/models"
)

type UserRepository interface {
	FindByUserId(ctx context.Context, userId uint) (*models.Users, error)
	RefreshUserRealNameInfo(ctx context.Context, userId string, idCard string, idName string, bankCard string) error
	SetFundPassword(ctx context.Context, userId string, password string) error
	UpdateLabel(ctx context.Context, userId uint, label []byte) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(connection *database.MySQLConnection) UserRepository {
	return &userRepository{db: connection.GetDB()}
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

func (r *userRepository) RefreshUserRealNameInfo(ctx context.Context, userId string, idCard string, idName string, bankCard string) error {
	user := &models.Users{
		IdCard:   idCard,
		IdName:   idName,
		BankCard: bankCard,
	}

	result := r.db.WithContext(ctx).
		Model(&models.Users{}).
		Where("userId = ?", userId).
		Updates(user)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) SetFundPassword(ctx context.Context, userId string, password string) error {
	user := &models.Users{
		FundPassword: password,
	}
	result := r.db.WithContext(ctx).
		Model(&models.Users{}).
		Where("userId = ?", userId).
		Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) UpdateLabel(ctx context.Context, userId uint, label []byte) error {
	user := &models.Users{
		Label: json.RawMessage(label),
	}
	result := r.db.WithContext(ctx).
		Model(&models.Users{}).
		Where("userId = ?", userId).
		Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
