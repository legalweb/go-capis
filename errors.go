package capis

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	// ErrUnknown is returned when we don't know what the error is.
	ErrUnknown struct {
		statusCode int
	}
)

var (
	// ErrUnauthorized is returned when the token is invalid.
	ErrUnauthorized = errors.New("unauthorized access, is the token valid")
	// ErrNotFound is returned when the resource is not found.
	ErrNotFound = errors.New("not found")
	// ErrUnreachable is returned when we cannot connect to the capis server.
	ErrUnreachable = errors.New("unreachable")
)

func statusCodeToError(sc int) error {
	if sc < 400 {
		return nil
	}

	switch sc {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusNotFound:
		return ErrNotFound
	default:
		return &ErrUnknown{sc}
	}
}

func (e *ErrUnknown) Error() string {
	return fmt.Sprintf("unknown status code: %d", e.statusCode)
}
