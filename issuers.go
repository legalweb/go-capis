package capis

import (
	"bytes"
	"context"
	"encoding/json"

	"go.opencensus.io/trace"
)

type (
	// ListIssuersResponse ...
	ListIssuersResponse struct {
		Data []*IssuerSummary `json:"data"`
	}

	// IssuerSummary ...
	IssuerSummary struct {
		ID    string `json:"issuer_id"`
		Label string `json:"label"`
	}

	// Issuer ...
	Issuer struct {
		ID    string `json:"issuer_id"`
		Label string `json:"label"`
	}
)

// ListIssuers ...
func (c *Client) ListIssuers(ctx context.Context) (*ListIssuersResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.ListIssuers")
	defer span.End()

	obj := &ListIssuersResponse{}

	req, err := c.newRequest("GET", "/v1/issuers", nil)
	if err != nil {
		c.logError(err)
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		c.logError(err)
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}

// FindIssuer ...
func (c *Client) FindIssuer(ctx context.Context, id string) (*Issuer, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.FindIssuer")
	defer span.End()

	obj := &Issuer{}

	req, err := c.newRequest("GET", "/v1/issuers/"+id, nil)
	if err != nil {
		c.logError(err)
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		c.logError(err)
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}

// NewIssuerRequest ...
type NewIssuerRequest struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// NewIssuer ...
func (c *Client) NewIssuer(ctx context.Context, opts *NewIssuerRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.NewIssuer")
	defer span.End()

	b, err := json.Marshal(opts)
	if err != nil {
		c.logError(err)
		return err
	}

	req, err := c.newRequest("POST", "/v1/issuers", bytes.NewReader(b))
	if err != nil {
		c.logError(err)
		return err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		c.logError(err)
		return err
	}

	return nil
}
