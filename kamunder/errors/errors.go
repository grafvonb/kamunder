package errors

import "errors"

var (
	ErrNotFound     = errors.New("process instance not found")
	ErrInvalidState = errors.New("invalid process instance state")
)
