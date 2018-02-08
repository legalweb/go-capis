package capis

import (
	"io"
	"net/http"
	"strings"
)

type (
	Client struct {
		*http.Client

		base  string
		token string
	}
)

func New(base, token string) *Client {
	return &Client{
		Client: http.DefaultClient,
		base:   strings.TrimRight(base, "/"),
		token:  token,
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer "+c.token)
	return req, nil
}
