package orders

import (
	"context"
	"gorm.io/gorm"
	"kite/internal/database"
)

type CreateOrdersRepository interface {
	GetCashOrdersCountByUserId(ctx context.Context, userId uint) (uint, error)
}

type createOrdersRepository struct {
	db *gorm.DB
}

func NewCreateOrdersRepository(connection *database.MySQLConnection) CreateOrdersRepository {
	return &createOrdersRepository{db: connection.GetDB()}
}

func (r createOrdersRepository) GetCashOrdersCountByUserId(ctx context.Context, userId uint) (uint, error) {
	var count uint

	result := r.db.WithContext(ctx).Where("userId = ?", userId).
		Where("couponFlag = ?", 0).
		Select("count(*)").
		Scan(&count)

	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
