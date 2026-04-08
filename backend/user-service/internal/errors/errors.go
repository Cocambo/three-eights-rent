package errors

import goerrors "errors"

var (
	ErrNotFound           = goerrors.New("resource not found")
	ErrDuplicateEmail     = goerrors.New("duplicate email")
	ErrConflict           = goerrors.New("resource conflict")
	ErrValidation         = goerrors.New("validation failed")
	ErrInvalidCredentials = goerrors.New("invalid credentials")
	ErrInvalidToken       = goerrors.New("invalid token")
	ErrUnauthorized       = goerrors.New("unauthorized")
)
