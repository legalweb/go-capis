package capis

import (
	"errors"
	"fmt"
	"net/http"
)

type (
	ErrUnknown struct {
		statusCode int
	}
)

var (
	ErrUnauthorized = errors.New("unauthorized access, is the token valid")
	ErrNotFound     = errors.New("not found")
	ErrUnreachable  = errors.New("unreachable")
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
