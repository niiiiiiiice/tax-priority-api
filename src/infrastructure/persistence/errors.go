package persistence

type RepositoryError struct {
	Code    string
	Message string
	Cause   error
}

func (e *RepositoryError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *RepositoryError) Unwrap() error {
	return e.Cause
}

// Стандартные коды ошибок
const (
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeAlreadyExists = "ALREADY_EXISTS"
	ErrCodeInvalidInput  = "INVALID_INPUT"
	ErrCodeInternal      = "INTERNAL_ERROR"
	ErrCodeConstraint    = "CONSTRAINT_VIOLATION"
	ErrCodeTimeout       = "TIMEOUT"
	ErrCodeTransaction   = "TRANSACTION_ERROR"
)

// Фабричные функции для ошибок
func NewNotFoundError(message string, cause error) *RepositoryError {
	return &RepositoryError{Code: ErrCodeNotFound, Message: message, Cause: cause}
}

func NewAlreadyExistsError(message string, cause error) *RepositoryError {
	return &RepositoryError{Code: ErrCodeAlreadyExists, Message: message, Cause: cause}
}

func NewInvalidInputError(message string, cause error) *RepositoryError {
	return &RepositoryError{Code: ErrCodeInvalidInput, Message: message, Cause: cause}
}

func NewInternalError(message string, cause error) *RepositoryError {
	return &RepositoryError{Code: ErrCodeInternal, Message: message, Cause: cause}
}
