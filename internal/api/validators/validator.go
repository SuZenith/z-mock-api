package validators

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	KiteError "kite/internal/errors"
	"reflect"
	"strings"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// 注册自定义标签处理函数，使用结构体字段的 json 标签作为错误消息的字段名
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return cv.translateError(err)
	}
	return nil
}

func (cv *CustomValidator) translateError(err error) error {
	if err == nil {
		return nil
	}
	// 检查是否为验证错误
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return KiteError.NewWithMessage(KiteError.ValidationError, "Invalid request data", err)
	}
	// 构建详细的错误消息
	errorMessages := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		errorMessages = append(errorMessages, cv.buildErrorMessage(e))
	}
	// 返回单个错误
	if len(errorMessages) == 1 {
		return KiteError.NewWithMessage(KiteError.ValidationError, errorMessages[0], err)
	}
	// 返回多个错误
	message := fmt.Sprintf("validation failed: %s", strings.Join(errorMessages, "; "))
	return KiteError.NewWithMessage(KiteError.ValidationError, message, err)
}

func (cv *CustomValidator) buildErrorMessage(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, param)
	case "numeric":
		return fmt.Sprintf("%s must be numeric", field)
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric", field)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, param)
	default:
		return fmt.Sprintf("%s failed validation: %s", field, tag)
	}
}

// BindAndValidate 绑定请求数据到结构体并验证
func BindAndValidate(c echo.Context, i interface{}) error {
	// 绑定请求数据到结构体
	if err := c.Bind(i); err != nil {
		return KiteError.NewWithMessage(KiteError.BadRequestError, "Invalid request format", err)
	}
	// 验证请求数据
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}
