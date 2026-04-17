package errors

import stderrors "errors"

var (
	ErrNotFound     = stderrors.New("not found")
	ErrValidation   = stderrors.New("validation failed")
	ErrUnauthorized = stderrors.New("unauthorized")
	ErrConflict     = stderrors.New("conflict")
)

type AppError struct {
	message string
	err     error
}

func New(err error, message string) error {
	if err == nil {
		return nil
	}

	return &AppError{
		message: message,
		err:     err,
	}
}

func (e *AppError) Error() string {
	return e.message
}

func (e *AppError) Unwrap() error {
	return e.err
}
