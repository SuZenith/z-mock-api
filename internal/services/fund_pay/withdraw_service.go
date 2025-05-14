package fund_pay

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"kite/internal/api/payloads"
	KiteError "kite/internal/errors"
	"kite/internal/repositories/fund_pay"
	"kite/internal/services/account"
)

type WithdrawService interface {
	Apply(ctx echo.Context, payload payloads.ApplyPayload) error
	QueryWithdrawTotalAmount(ctx context.Context, userId uint) (decimal.Decimal, error)
}

type withdrawService struct {
	userService account.UserService
	repo        fund_pay.WithdrawOrderRepository
}

func NewWithdrawService(userServer account.UserService, repo fund_pay.WithdrawOrderRepository) WithdrawService {
	return &withdrawService{userServer, repo}
}

func (s *withdrawService) Apply(ctx echo.Context, payload payloads.ApplyPayload) error {
	user, err := s.userService.FindByUserId(ctx.Request().Context(), 545252)
	if err != nil {
		return KiteError.New(KiteError.InternalServerError, err)
	}
	if user == nil {
		return KiteError.New(KiteError.UserNotFoundError, fmt.Errorf("user not found"))
	}
	fmt.Printf("%s\n", user.Name)
	return nil
}

func (s *withdrawService) QueryWithdrawTotalAmount(ctx context.Context, userId uint) (decimal.Decimal, error) {
	return s.repo.QueryWithdrawTotalAmountByUserId(ctx, userId)
}
