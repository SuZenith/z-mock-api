package fund_pay

import (
	"context"
	"github.com/shopspring/decimal"
	"kite/internal/models"
	"kite/internal/repositories/fund_pay"
)

type RechargeService interface {
	QuerySuccessRechargeOrderCount(ctx context.Context, userId uint) (int64, error)
	QueryFirstRechargeOrder(ctx context.Context, userId uint) (*models.RechargeOrders, error)
	QueryRechargeTotalAmount(ctx context.Context, userId uint) (decimal.Decimal, error)
}

type rechargeService struct {
	repo fund_pay.RechargeOrderRepository
}

func NewRechargeService(repo fund_pay.RechargeOrderRepository) RechargeService {
	return &rechargeService{repo}
}

func (s *rechargeService) QuerySuccessRechargeOrderCount(ctx context.Context, userId uint) (int64, error) {
	return s.repo.QuerySuccessRechargeOrderCountByUserId(ctx, userId)
}

func (s *rechargeService) QueryFirstRechargeOrder(ctx context.Context, userId uint) (*models.RechargeOrders, error) {
	return s.repo.QueryFirstOrderByUserId(ctx, userId)
}

func (s *rechargeService) QueryRechargeTotalAmount(ctx context.Context, userId uint) (decimal.Decimal, error) {
	return s.repo.QueryRechargeTotalAmountByUserId(ctx, userId)
}
