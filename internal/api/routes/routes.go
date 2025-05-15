package routes

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/handlers"
	"kite/internal/api/handlers/mock"
)

func RegisterRoutes(e *echo.Echo, mockHandler *mock.ApiHandler) {
	e.GET("/health", handlers.HealthCheck)

	v1 := e.Group("/api/v1")

	mockRoutes := v1.Group("/mock")
	mockRoutes.POST("/create", mockHandler.Create)
	mockRoutes.GET("/:uid/*", mockHandler.Get)

}
