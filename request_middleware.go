package capis

import (
	"net/http"
)

type (
	// RequestMiddlewareFunc is a function that will use or manipulated the request
	// and return a request.
	RequestMiddlewareFunc func(*http.Request) *http.Request

	// RequestMiddleware contains request middleware funcs.
	RequestMiddleware []RequestMiddlewareFunc
)

// Apply the request middleware functions to the request and return the end result.
func (m RequestMiddleware) Apply(r *http.Request) *http.Request {
	for _, f := range m {
		r = f(r)
	}
	return r
}
