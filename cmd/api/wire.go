//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	handle_fund_pay "kite/internal/api/handlers/fund_pay"
	"kite/internal/configs"
	"kite/internal/database"
	repo_accounts "kite/internal/repositories/accounts"
	repo_fund_pay "kite/internal/repositories/fund_pay"
	repo_orders "kite/internal/repositories/orders"

	service_account "kite/internal/services/account"
	service_fund_pay "kite/internal/services/fund_pay"
)

var RepositorySet = wire.NewSet(
	repo_accounts.NewUserRepository,
	repo_fund_pay.NewWithdrawOrderRepository,
	repo_orders.NewCreateOrdersRepository,
)

var ServiceSet = wire.NewSet(
	service_account.NewUserService,
	service_fund_pay.NewWithdrawService,
)

var HandlerSet = wire.NewSet(
	handle_fund_pay.NewWithdrawHandler,
)

func InitializeApp(cfg *configs.MySQLConfig, echo *echo.Echo) (*Server, error) {
	wire.Build(
		database.NewMySQLConnection,
		HandlerSet,
		ServiceSet,
		RepositorySet,
		NewServer,
	)
	return nil, nil
}
