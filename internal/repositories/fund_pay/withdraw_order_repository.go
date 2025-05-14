package fund_pay

import (
	"context"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"kite/internal/database"
)

type WithdrawOrderRepository interface {
	GetSumAmountByUserId(ctx context.Context, userId uint) (decimal.Decimal, error)
}

type withdrawOrderRepository struct {
	db *gorm.DB
}

func NewWithdrawOrderRepository(connection *database.MySQLConnection) WithdrawOrderRepository {
	return &withdrawOrderRepository{db: connection.GetDB()}
}

func (r *withdrawOrderRepository) GetSumAmountByUserId(ctx context.Context, userId uint) (decimal.Decimal, error) {
	var sum decimal.Decimal

	result := r.db.WithContext(ctx).Where("userId = ?", userId).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&sum)
	if result.Error != nil {
		return decimal.Zero, result.Error
	}
	return sum, nil
}
