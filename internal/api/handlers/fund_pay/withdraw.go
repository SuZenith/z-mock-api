package fund_pay

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	"kite/internal/api/validators"
	"kite/internal/repositories/accounts"
	"kite/internal/services/fund_pay"
	"kite/pkg/response"
)

type WithdrawHandler struct {
	service *fund_pay.WithdrawService
}

func NewWithdrawHandler(userRepo accounts.UserRepository) *WithdrawHandler {
	return &WithdrawHandler{
		service: fund_pay.NewWithdrawService(userRepo),
	}
}

// Apply 申请提现
func (h *WithdrawHandler) Apply(c echo.Context) error {
	var payload payloads.ApplyPayload
	err := validators.BindAndValidate(c, &payload)
	if err != nil {
		return err
	}
	err = h.service.Apply(c, payload)
	if err != nil {
		return err
	}

	return response.SuccessWithoutData(c)
}
