//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"kite/internal/api/handlers/mock"
	"kite/internal/configs"
	"kite/internal/database"
	"kite/internal/repositories"
	"kite/internal/services"
)

var RepositorySet = wire.NewSet(
	repositories.NewApiRepository,
)

var ServiceSet = wire.NewSet(
	services.NewApiService,
)

var HandlerSet = wire.NewSet(
	mock.NewApiHandler,
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
