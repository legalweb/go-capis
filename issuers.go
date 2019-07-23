package capis

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	querystring "github.com/google/go-querystring/query"
	"go.opencensus.io/trace"
)

type (
	// ListIssuersResponse ...
	ListIssuersResponse struct {
		Data []*IssuerSummary `json:"data"`
	}

	IssuerFilters struct {
		Label string `json:"label"`
	}

	// IssuerSummary ...
	IssuerSummary struct {
		ID    string `json:"issuer_id"`
		Label string `json:"label"`
		Logo  string `json:"logo"`
	}

	// Issuer ...
	Issuer struct {
		ID          string `json:"issuer_id"`
		Label       string `json:"label"`
		Logo        string `json:"logo"`
		Description string `json:"description"`
	}

	// NewIssuerRequest ...
	type NewIssuerRequest struct {
		ID    string `json:"issuer_id"`
		Label string `json:"label"`
	}
)

// ListIssuers ...
func (c *Client) ListIssuers(ctx context.Context, filters *IssuerFilters, start, limit int) (*ListIssuersResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.ListIssuers")
	defer span.End()

	qs, _ := querystring.Values(filters)
	qs.Set("start", strconv.Itoa(start))
	qs.Set("limit", strconv.Itoa(limit))

	obj := &ListIssuersResponse{}

	req, err := c.newRequest("GET", "/v1/issuers?"+qs.Encode(), nil)
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
