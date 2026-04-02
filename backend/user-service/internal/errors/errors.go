package errors

import goerrors "errors"

var (
	ErrNotFound = goerrors.New("resource not found")
)
