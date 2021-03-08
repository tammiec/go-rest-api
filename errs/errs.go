package errs

import "errors"

var (
	ErrBadRequestParams = errors.New("bad request params")
	ErrInvalidId        = errors.New("invalid id")
)
