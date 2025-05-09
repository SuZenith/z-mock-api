package fund_pay

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type WithdrawHandler struct{}

func NewWithdrawHandler() *WithdrawHandler {
	return &WithdrawHandler{}
}

// Apply 申请提现
func (h *WithdrawHandler) Apply(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
