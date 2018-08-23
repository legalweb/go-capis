package capis

import (
	"bytes"
	"context"
	"encoding/json"

	querystring "github.com/google/go-querystring/query"
	"go.opencensus.io/trace"
)

type (
	// ListGroupsResponse ...
	ListGroupsResponse struct {
		Data []*Group `json:"data"`
	}

	// FindGroupResponse ...
	FindGroupResponse struct {
		Data *Group `json:"data"`
	}

	// GroupFilters ...
	GroupFilters struct {
		Type string `json:"type"`
	}

	// Group ...
	Group struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	}
)

// ListGroups ...
func (c *Client) ListGroups(ctx context.Context, filters *GroupFilters) (*ListGroupsResponse, error) {
	obj := &ListGroupsResponse{}
	qs, _ := querystring.Values(filters)

	req, err := c.newRequest("GET", "/v1/groups?"+qs.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}

// FindGroup ...
func (c *Client) FindGroup(ctx context.Context, name string) (*FindGroupResponse, error) {
	obj := &FindGroupResponse{}

	req, err := c.newRequest("GET", "/v1/groups/"+name, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}

// NewGroupRequest ...
type NewGroupRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// NewGroup ...
func (c *Client) NewGroup(ctx context.Context, opts *NewGroupRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.NewGroup")
	defer span.End()

	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", "/v1/groups", bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		return err
	}

	return nil
}

// SetGroupProductsRequest ...
type SetGroupProductsRequest struct {
	GroupID  string   `json:"-"`
	Products []string `json:"product_ids"`
}

// SetGroupProducts ...
func (c *Client) SetGroupProducts(ctx context.Context, opts *SetGroupProductsRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.SetGroupProducts")
	defer span.End()

	b, err := json.Marshal(opts)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", "/v1/groups/"+opts.GroupID+"/products", bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		return err
	}

	return nil
}

// IsType ...
func (g *Group) IsType(t ProductType) bool {
	return ProductType(g.Type) == t
}
