package errs

import "errors"

var (
	ErrBadRequest         = errors.New("bad request")
	ErrServiceUnavailable = errors.New("service unavailable")
	ErrNoRows             = errors.New("no rows found")
)
