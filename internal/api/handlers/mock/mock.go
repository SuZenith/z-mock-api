package mock

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"kite/internal/api/payloads"
	"kite/internal/api/validators"
	KiteError "kite/internal/errors"
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
	uid := ctx.Param("uid")
	path := fmt.Sprintf("/%s", ctx.Param("*"))
	contentType, responseBody, err := h.srv.Request(ctx, uid, path, "GET")
	if err != nil {
		return KiteError.New(KiteError.InternalServerError, err)
	}
	if *contentType == "application/json" {
		return ctx.JSONBlob(200, []byte(*responseBody))
	}

	return response.Success(ctx, responseBody)
}

func (h *ApiHandler) Post(ctx echo.Context) error {
	uid := ctx.Param("uid")
	path := fmt.Sprintf("/%s", ctx.Param("*"))
	contentType, responseBody, err := h.srv.Request(ctx, uid, path, "POST")
	if err != nil {
		return KiteError.New(KiteError.InternalServerError, err)
	}
	if *contentType == "application/json" {
		return ctx.JSONBlob(200, []byte(*responseBody))
	}

	return response.Success(ctx, responseBody)
}

func (h *ApiHandler) Put(ctx echo.Context) error {
	uid := ctx.Param("uid")
	path := fmt.Sprintf("/%s", ctx.Param("*"))
	contentType, responseBody, err := h.srv.Request(ctx, uid, path, "PUT")
	if err != nil {
		return KiteError.New(KiteError.InternalServerError, err)
	}
	if *contentType == "application/json" {
		return ctx.JSONBlob(200, []byte(*responseBody))
	}

	return response.Success(ctx, responseBody)
}

func (h *ApiHandler) Delete(ctx echo.Context) error {
	uid := ctx.Param("uid")
	path := fmt.Sprintf("/%s", ctx.Param("*"))
	contentType, responseBody, err := h.srv.Request(ctx, uid, path, "DELETE")
	if err != nil {
		return KiteError.New(KiteError.InternalServerError, err)
	}
	if *contentType == "application/json" {
		return ctx.JSONBlob(200, []byte(*responseBody))
	}

	return response.Success(ctx, responseBody)
}
