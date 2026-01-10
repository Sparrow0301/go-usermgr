package errorx

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// AppError carries HTTP metadata for consistent error responses.
type AppError struct {
	Status  int
	Code    string
	Message string
	Details interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

// WithDetails attaches structured details (e.g., validation errors).
// We clone the value to keep AppError reusable across requests.
func (e *AppError) WithDetails(details interface{}) *AppError {
	clone := *e
	clone.Details = details
	return &clone
}

// Convenience constructors.
func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func Wrap(err error, status int, code, message string) *AppError {
	if err == nil {
		return New(status, code, message)
	}
	return &AppError{Status: status, Code: code, Message: message, Details: err.Error()}
}

// Domain specific errors.
var (
	ErrValidation         = New(http.StatusBadRequest, "VALIDATION_FAILED", "请求参数不合法")
	ErrUserExists         = New(http.StatusConflict, "USER_EXISTS", "用户名或邮箱已存在")
	ErrInvalidCredentials = New(http.StatusUnauthorized, "INVALID_CREDENTIALS", "用户名或密码错误")
	ErrUserDisabled       = New(http.StatusForbidden, "USER_DISABLED", "用户已被禁用")
	ErrForbidden          = New(http.StatusForbidden, "FORBIDDEN", "无访问权限")
	ErrUserNotFound       = New(http.StatusNotFound, "USER_NOT_FOUND", "用户不存在")
	ErrInternal           = New(http.StatusInternalServerError, "INTERNAL_ERROR", "服务器内部错误")
)

// Is checks whether err matches target *AppError (by Code).
func Is(err error, target *AppError) bool {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return false
	}
	return appErr.Code == target.Code
}

// ValidationDetails turns validator errors into a map for responses.
type ValidationErrorItem struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func FromValidationError(err error) *AppError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		items := make([]ValidationErrorItem, 0, len(ve))
		for _, fe := range ve {
			items = append(items, ValidationErrorItem{
				Field: fe.Field(),
				Tag:   fe.Tag(),
				Param: fe.Param(),
			})
		}
		return ErrValidation.WithDetails(items)
	}
	return ErrValidation.WithDetails(err.Error())
}
