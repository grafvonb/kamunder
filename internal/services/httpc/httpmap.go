package httpc

import (
	"fmt"
	"net/http"

	d "github.com/grafvonb/camunder/internal/domain"
)

func HttpStatusErr(hr *http.Response, body []byte) error {
	if hr == nil || hr.Request == nil || hr.Request.URL == nil {
		return fmt.Errorf("%w: invalid http response; body=%s", d.ErrUpstream, string(body))
	}
	if err := MapHTTPToDomain(hr.StatusCode); err != nil {
		return fmt.Errorf("%w: %s %s; status: %d; body: %s",
			err,
			hr.Request.Method,
			hr.Request.URL.String(),
			hr.StatusCode,
			string(body),
		)
	}
	return nil
}

func MapHTTPToDomain(status int) error {
	switch status {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return d.ErrBadRequest
	case http.StatusUnauthorized:
		return d.ErrUnauthorized
	case http.StatusForbidden:
		return d.ErrForbidden
	case http.StatusNotFound:
		return d.ErrNotFound
	case http.StatusConflict:
		return d.ErrConflict
	case http.StatusPreconditionFailed:
		return d.ErrPrecondition
	case http.StatusUnsupportedMediaType:
		return d.ErrUnsupported
	case http.StatusUnprocessableEntity:
		return d.ErrValidation
	case http.StatusTooManyRequests:
		return d.ErrRateLimited
	case http.StatusGatewayTimeout:
		return d.ErrGatewayTimeout
	case http.StatusServiceUnavailable:
		return d.ErrUnavailable
	case http.StatusBadGateway:
		return d.ErrUpstream
	case http.StatusInternalServerError:
		return d.ErrInternal
	default:
		return d.ErrUpstream
	}
}
