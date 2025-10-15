package domain

import (
	"errors"
)

// Sentinel errors
var (
	ErrMultipleMatches = errors.New("multiple matches found")

	ErrBadRequest        = errors.New("bad request")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("conflict")
	ErrPrecondition      = errors.New("precondition failed")
	ErrUnsupported       = errors.New("unsupported media type")
	ErrValidation        = errors.New("validation failed")
	ErrRateLimited       = errors.New("rate limited")
	ErrGatewayTimeout    = errors.New("gateway timeout")
	ErrUnavailable       = errors.New("service unavailable")
	ErrUpstream          = errors.New("upstream error")
	ErrInternal          = errors.New("internal error")
	ErrMalformedResponse = errors.New("malformed response")

	ErrUnknownStateFilter = errors.New("is unknown (valid: all, active, canceled, completed)")
)
