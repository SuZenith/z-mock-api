package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode int

// 预定义错误码
const (
	// InternalServerError 系统级错误
	InternalServerError ErrorCode = -1000 - iota
	DatabaseError
	ConfigError
	ValidationError
	DataError
	MarshalError

	// BadRequestError 客户端错误
	BadRequestError ErrorCode = -2000 - iota
	UnauthorizedError
	ForbiddenError
	NotFoundError

	// UserNotFoundError 业务级的错误
	UserNotFoundError ErrorCode = -3000 - iota
)

// 错误码到 HTTP 状态码的映射
var errorCodeToHTTPStatus = map[ErrorCode]int{
	InternalServerError: http.StatusInternalServerError,
	DatabaseError:       http.StatusInternalServerError,
	ConfigError:         http.StatusInternalServerError,
	DataError:           http.StatusInternalServerError,
	MarshalError:        http.StatusInternalServerError,
	ValidationError:     http.StatusBadRequest,
	BadRequestError:     http.StatusBadRequest,
	UnauthorizedError:   http.StatusUnauthorized,
	ForbiddenError:      http.StatusForbidden,
	NotFoundError:       http.StatusNotFound,
	UserNotFoundError:   http.StatusNotFound,
}

// 错误码到消息的映射
var errorCodeToMessage = map[ErrorCode]string{
	InternalServerError: "Internal Server Error",
	DatabaseError:       "Internal Server Error",
	ConfigError:         "Internal Server Error",
	DataError:           "Internal Server Error",
	MarshalError:        "Internal Server Error",
	ValidationError:     "Validation Error",
	BadRequestError:     "Bad Request",
	UnauthorizedError:   "Unauthorized",
	ForbiddenError:      "Forbidden",
	NotFoundError:       "Not Found",
	UserNotFoundError:   "User Not Found",
}

type AppError struct {
	Code       ErrorCode
	Message    string
	Detail     string
	Err        error
	HTTPStatus int
}

func (e *AppError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, err error) *AppError {
	message := errorCodeToMessage[code]
	httpStatus := errorCodeToHTTPStatus[code]

	return &AppError{
		Code:       code,
		Message:    message,
		Err:        err,
		HTTPStatus: httpStatus,
	}
}

func NewWithMessage(code ErrorCode, message string, err error) *AppError {
	appErr := New(code, err)
	appErr.Message = message

	return appErr
}

func (e *AppError) WithDetail(detail string) *AppError {
	e.Detail = detail
	return e
}

func IsAppError(err error) (*AppError, bool) {
	if err == nil {
		return nil, false
	}
	var appErr *AppError
	ok := errors.As(err, &appErr)
	return appErr, ok
}
