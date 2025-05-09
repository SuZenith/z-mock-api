package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	kiteError "kite/internal/errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	var (
		code    = http.StatusInternalServerError
		message = "Internal Server Error"
		errCode = int(kiteError.InternalServerError)
	)

	if appErr, ok := kiteError.IsAppError(err); ok {
		code = appErr.HTTPStatus
		message = appErr.Message
		errCode = int(appErr.Code)
	}
	var e *echo.HTTPError
	if errors.As(err, &e) {
		// Echo 框架的 HTTP 错误
		code = e.Code
		message = fmt.Sprintf("%v", e.Message)
		errCode = code
	} else {
		// 未知的错误
		log.Printf("Unexpected error: %v", err)
	}
	// 如果响应已经提交，直接返回
	if c.Response().Committed {
		return
	}
	// 发送统一的 JSON 响应
	if err := c.JSON(code, ErrorResponse{
		Code:      errCode,
		Message:   message,
		RequestId: c.Response().Header().Get(echo.HeaderXRequestID),
	}); err != nil {
		log.Printf("failed to send error response: %v", err)
	}
}
