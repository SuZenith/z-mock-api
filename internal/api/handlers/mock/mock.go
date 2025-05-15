package mock

import (
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	"kite/internal/api/validators"
	"kite/internal/services"
	"kite/pkg/response"
)

type ApiHandler struct {
	srv services.ApiService
}

func NewApiHandler(srv services.ApiService) *ApiHandler {
	return &ApiHandler{srv}
}

func (h *ApiHandler) Create(ctx echo.Context) error {
	var payload payloads.MockApiPayload
	err := validators.BindAndValidate(ctx, &payload)
	if err != nil {
		return err
	}
	err = h.srv.Create(ctx, payload)
	if err != nil {
		return err
	}

	return response.SuccessWithoutData(ctx)
}

func (h *ApiHandler) Get(ctx echo.Context) error {

}
