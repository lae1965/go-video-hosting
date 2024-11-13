package errors

type ErrType int

type AppError struct {
	Type    ErrType
	Message string
}

const (
	Unauthorization ErrType = iota
	InvalidData
	AlreadyExist
	NotFound
	NotUnique
	NotEqual
	EmptyField
	UnknownError
)

func New(errType ErrType, message string) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
	}
}
