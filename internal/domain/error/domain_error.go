package domainerror

import "fmt"

type ErrorCode string

const (
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrInvalidInput  ErrorCode = "INVALID_INPUT"
)

type DomainError struct {
	Code    ErrorCode
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewNotFoundError(resource string) *DomainError {
	return &DomainError{
		Code:    ErrNotFound,
		Message: fmt.Sprintf("%s が見つかりません", resource),
	}
}

func NewUnauthorizedError() *DomainError {
	return &DomainError{
		Code:    ErrUnauthorized,
		Message: "アクセス権限がありません",
	}
}

func NewInvalidInputError(message string) *DomainError {
	return &DomainError{
		Code:    ErrInvalidInput,
		Message: message,
	}
}