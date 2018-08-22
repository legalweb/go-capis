package capis

import (
	"io"
	"net/http"
	"strings"
)

const (
	// DefaultBaseURL is the known production url for capis.
	DefaultBaseURL = "https://comparisonapis.com"
)

type (
	// Client will talk to comparisonapis.com
	Client struct {
		*http.Client

		base  string
		token string
	}
)

// New will return a client with the default http client.
func New(base, token string) *Client {
	return NewWithHTTPClient(base, token, http.DefaultClient)
}

// NewWithHTTPClient will return a new comparisonapis.com client using the
// http client provided.
func NewWithHTTPClient(base, token string, hc *http.Client) *Client {
	return &Client{
		Client: hc,
		base:   strings.TrimRight(base, "/"),
		token:  token,
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}

	return req, nil
}
