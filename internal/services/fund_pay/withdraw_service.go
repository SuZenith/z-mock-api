package fund_pay

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	KiteError "kite/internal/errors"
	"kite/internal/services/account"
)

type WithdrawService struct {
	userService account.UserService
}

func NewWithdrawService(userService account.UserService) *WithdrawService {
	return &WithdrawService{
		userService,
	}
}

func (s *WithdrawService) Apply(ctx echo.Context, payload payloads.ApplyPayload) error {
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
