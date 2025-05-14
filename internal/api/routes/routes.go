package routes

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/handlers"
	"kite/internal/api/handlers/fund_pay"
)

func RegisterRoutes(e *echo.Echo, withdrawHandler *fund_pay.WithdrawHandler) {
	e.GET("/health", handlers.HealthCheck)

	v1 := e.Group("/api/v1")

	// 支付提现相关路由
	fundPayRoutes := v1.Group("/fund-pay")
	fundPayRoutes.POST("/withdraw/apply", withdrawHandler.Apply)
}
