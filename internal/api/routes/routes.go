package routes

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/handlers"
	"kite/internal/api/handlers/fund_pay"
	"kite/internal/database"
	"kite/internal/repositories/accounts"
)

func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", handlers.HealthCheck)

	v1 := e.Group("/api/v1")

	// 支付提现相关路由
	fundPayRoutes := v1.Group("/fund-pay")

	userRepo := accounts.NewUserRepository(database.GetDB())
	withdrawHandler := fund_pay.NewWithdrawHandler(userRepo)
	fundPayRoutes.POST("/withdraw/apply", withdrawHandler.Apply)
}
