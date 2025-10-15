package services

import "errors"

var (
	ErrNoConfig     = errors.New("no config provided")
	ErrNoHTTPClient = errors.New("no http client provided")
	ErrNoLogger     = errors.New("no logger provided")

	ErrUnknownAPIVersion = errors.New("unknown API version")
	ErrCycleDetected     = errors.New("cycle detected in process instance ancestry")
)
