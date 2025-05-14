package fund_pay

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"kite/internal/database"
	"kite/internal/models"
)

type RechargeOrderRepository interface {
	// QuerySuccessRechargeOrderCountByUserId 查询用户成功充值的订单数量
	QuerySuccessRechargeOrderCountByUserId(ctx context.Context, userId uint) (int64, error)
	// QueryFirstOrderByUserId 查询用户最早一笔订单
	QueryFirstOrderByUserId(ctx context.Context, userId uint) (*models.RechargeOrders, error)
}

type rechargeOrderRepository struct {
	db *gorm.DB
}

func NewRechargeOrderRepository(connection *database.MySQLConnection) RechargeOrderRepository {
	return &rechargeOrderRepository{db: connection.GetDB()}
}

func (r *rechargeOrderRepository) QuerySuccessRechargeOrderCountByUserId(ctx context.Context, userId uint) (int64, error) {
	var count int64

	result := r.db.WithContext(ctx).Where("userId = ?", userId).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (r *rechargeOrderRepository) QueryFirstOrderByUserId(ctx context.Context, userId uint) (*models.RechargeOrders, error) {
	var rechargeOrders models.RechargeOrders

	result := r.db.WithContext(ctx).Where("userId = ?", userId).First(&rechargeOrders)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &rechargeOrders, nil
}
