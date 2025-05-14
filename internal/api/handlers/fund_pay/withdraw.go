package fund_pay

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	"kite/internal/api/validators"
	"kite/internal/services/fund_pay"
	"kite/pkg/response"
)

type WithdrawHandler struct {
	service *fund_pay.WithdrawService
}

func NewWithdrawHandler(service *fund_pay.WithdrawService) *WithdrawHandler {
	return &WithdrawHandler{service}
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
