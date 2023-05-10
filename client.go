package capis

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the known production url for capis.
	DefaultBaseURL = "https://comparisonapis.com"
)

type (
	AuthProvider interface {
		AuthorizeRequest(context.Context, *http.Request) error
	}

	// Client will talk to comparisonapis.com
	Client struct {
		httpC             *http.Client
		base              string
		authProvider      AuthProvider
		logError          func(error)
		requestMiddleware RequestMiddleware
	}

	// Option customises the client.
	Option func(*Client) error
)

// New will return a client with the options provided.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		httpC:             http.DefaultClient,
		base:              DefaultBaseURL,
		logError:          func(error) {},
		requestMiddleware: RequestMiddleware{},
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// WithHTTPClient returns an option to pass to New()
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) error {
		c.httpC = hc
		return nil
	}
}

// WithAuthProvider returns an option to pass to New()
func WithAuthProvider(ap AuthProvider) Option {
	return func(c *Client) error {
		c.authProvider = ap
		return nil
	}
}

// WithBase returns an option to pass to New()
func WithBase(base string) Option {
	return func(c *Client) error {
		c.base = strings.TrimRight(base, "/")
		return nil
	}
}

// WithErrorLog returns an option to pass to New()
func WithErrorLog(out func(error)) Option {
	return func(c *Client) error {
		c.logError = out
		return nil
	}
}

// WithRequestMiddleware returns an option to pass to New()
func WithRequestMiddleware(m RequestMiddlewareFunc) Option {
	return func(c *Client) error {
		c.requestMiddleware = append(c.requestMiddleware, m)
		return nil
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	// @todo pass context to this method :)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return nil, err
	}

	if err := c.authProvider.AuthorizeRequest(ctx, req); err != nil {
		return nil, fmt.Errorf("authorization error %w", err)
	}

	return req, nil
}

// Do forwards the request to be handled by the HTTP client provided.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.httpC.Do(c.requestMiddleware.Apply(req))
}
