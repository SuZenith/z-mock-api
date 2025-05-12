package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	CodeSuccess    = 0
	MessageSuccess = "success"
)

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    data,
	})
}

func SuccessWithoutData(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: MessageSuccess,
		Data:    nil,
	})
}
