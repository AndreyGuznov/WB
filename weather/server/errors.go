package server

const (
	InternalServerError = "SERVER_ERROR"
	BadRequestError     = "BAD_REQUEST"
)

var (
	ErrInternalServerError = NewError(InternalServerError, "Unable to process the request due to server error")
)

// Error object struct
type ServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewError(code, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}
}
