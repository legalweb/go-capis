package capis

import (
	"bytes"
	"context"
	"encoding/json"

	querystring "github.com/google/go-querystring/query"
	"go.opencensus.io/trace"
)

type (
	// EmbedCreateRequest ...
	EmbedCreateRequest struct {
		ID        string         `json:"id"`
		GroupID   string         `json:"group_id"`
		Filters   []string       `json:"filters"`
		Columns   []string       `json:"columns"`
		Theme     EmbedTheme     `json:"theme"`
		Overrides EmbedOverrides `json:"overrides"`
	}

	// EmbedUpdateRequest ...
	EmbedUpdateRequest struct {
		id      string
		Filters []string   `json:"filters"`
		Columns []string   `json:"columns"`
		Theme   EmbedTheme `json:"theme"`
	}

	// ListEmbedsResponse ...
	ListEmbedsResponse struct {
		Data []*Embed `json:"data"`
	}

	// Embed ...
	Embed struct {
		ID         string               `json:"id"`
		Introducer string               `json:"introducer"`
		Theme      EmbedTheme           `json:"theme"`
		Overrides  EmbedOverrides       `json:"overrides"`
		Filters    []string             `json:"filters"`
		Columns    []string             `json:"columns"`
		Source     EmbedProductSelector `json:"source"`
	}

	// CreateEmbedRequest ...
	CreateEmbedRequest struct {
		ID        string         `json:"id"`
		Theme     EmbedTheme     `json:"theme"`
		Overrides EmbedOverrides `json:"overrides"`
		Filters   []string       `json:"filters"`
		Columns   []string       `json:"columns"`
		Group     string         `json:"group_id"`
	}

	// EmbedTheme ...
	EmbedTheme struct {
		MainColor                         string `json:"mainColor"`
		MainFontFamily                    string `json:"mainFontFamily"`
		MainNormalFontWeight              string `json:"mainNormalFontWeight"`
		MainFontSize                      string `json:"mainFontSize"`
		MainBoldFontWeight                string `json:"mainBoldFontWeight"`
		ProductMaskBackground             string `json:"productMaskBackground"`
		ProductEmptyBackground            string `json:"productEmptyBackground"`
		ProductOutlineBackground          string `json:"productOutlineBackground"`
		ProductOutlineColor               string `json:"productOutlineColor"`
		ProductColBackground              string `json:"productColBackground"`
		ProductHighlightOutlineBackground string `json:"productHighlightOutlineBackground"`
		ProductHighlightOutlineColor      string `json:"productHighlightOutlineColor"`
		ApplyButtonBackground             string `json:"applyButtonBackground"`
		ApplyButtonColor                  string `json:"applyButtonColor"`
		InfoButtonBackground              string `json:"infoButtonBackground"`
		InfoButtonColor                   string `json:"infoButtonColor"`
		InfoCheckColor                    string `json:"infoCheckColor"`
		FilterHeaderBorder                string `json:"filterHeaderBorder"`
		FilterHeaderColor                 string `json:"filterHeaderColor"`
		FilterChosenBackground            string `json:"filterChosenBackground"`
		FilterChosenColor                 string `json:"filterChosenColor"`
	}

	// EmbedOverrides ...
	EmbedOverrides struct {
		ButtonText string                 `json:"button_text"`
		ApplyURL   string                 `json:"apply_url"`
		Meta       map[string]interface{} `json:"metadata"`
	}

	// EmbedProductSelector ...
	EmbedProductSelector struct {
		GroupID   string   `json:"group_id,omitempty"`
		ProductID []string `json:"product_ids,omitempty"`
	}

	// DetailedEmbed ...
	DetailedEmbed struct {
		ID         string               `json:"id"`
		Introducer string               `json:"introducer"`
		Theme      EmbedTheme           `json:"theme"`
		Overrides  EmbedOverrides       `json:"overrides"`
		Filters    []string             `json:"filters"`
		Columns    []string             `json:"columns"`
		Source     EmbedProductSelector `json:"source"`

		Details struct {
			ProductCount int64  `json:"product_count"`
			ProductType  string `json:"product_type"`
			Snippet      string `json:"snippet"`
		} `json:"details"`
	}

	// EmbedFilters ...
	EmbedFilters struct {
		ID       []string `json:"embed_ids" url:"embed_ids"`
		Metadata []string `json:"metadata" url:"metadata"`
	}
)

// ListEmbeds ...
func (c *Client) ListEmbeds(ctx context.Context, filters *EmbedFilters) (*ListEmbedsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.ListEmbeds")
	defer span.End()

	obj := &ListEmbedsResponse{}
	qs, _ := querystring.Values(filters)

	req, err := c.newRequest("GET", "/v1/embeds?"+qs.Encode(), nil)
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

// FindEmbed ...
func (c *Client) FindEmbed(ctx context.Context, id string) (*Embed, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.FindEmbed")
	defer span.End()

	obj := &Embed{}

	req, err := c.newRequest("GET", "/v1/embeds/"+id, nil)
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

// FindEmbedDetailed ...
func (c *Client) FindEmbedDetailed(ctx context.Context, id string) (*DetailedEmbed, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.FindEmbedDetailed")
	defer span.End()

	obj := &DetailedEmbed{}

	req, err := c.newRequest("GET", "/v1/embeds/"+id+"/detailed", nil)
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

// NewCreateEmbedRequestForGroup will return a create embed request with the
// product source of a group.
func NewCreateEmbedRequestForGroup(id string, theme EmbedTheme, overrides EmbedOverrides, filters, columns []string, group string) *CreateEmbedRequest {
	return &CreateEmbedRequest{
		ID:        id,
		Theme:     theme,
		Overrides: overrides,
		Filters:   filters,
		Columns:   columns,
		Group:     group,
	}
}

// CreateEmbed ...
func (c *Client) CreateEmbed(ctx context.Context, embed *CreateEmbedRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.CreateEmbed")
	defer span.End()

	b, err := json.Marshal(embed)
	if err != nil {
		c.logError(err)
		return err
	}

	req, err := c.newRequest("POST", "/v1/embeds", bytes.NewReader(b))
	if err != nil {
		c.logError(err)
		return err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return ErrUnreachable
	}
	res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		c.logError(err)
		return err
	}

	return nil
}

func (e *Embed) Update() *EmbedUpdateRequest {
	return &EmbedUpdateRequest{
		id:      e.ID,
		Filters: e.Filters,
		Columns: e.Columns,
		Theme:   e.Theme,
	}
}

// SetFilters will update the update embed request to make changes to the Filters.
func (e *EmbedUpdateRequest) SetFilters(in []string) *EmbedUpdateRequest {
	e.Filters = in
	return e
}

// SetColumns will update the update embed request to make changes to the Columns.
func (e *EmbedUpdateRequest) SetColumns(in []string) *EmbedUpdateRequest {
	e.Columns = in
	return e
}

// SetTheme will update the update embed request to make changes to the Theme.
func (e *EmbedUpdateRequest) SetTheme(in EmbedTheme) *EmbedUpdateRequest {
	e.Theme = in
	return e
}

// UpdateEmbed will send the request to update the embed.
func (c *Client) UpdateEmbed(ctx context.Context, euq *EmbedUpdateRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.UpdateEmbed")
	defer span.End()

	b, err := json.Marshal(euq)
	if err != nil {
		c.logError(err)
		return err
	}

	req, err := c.newRequest("PUT", "/v1/embeds/"+euq.id, bytes.NewReader(b))
	if err != nil {
		c.logError(err)
		return err
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return ErrUnreachable
	}
	res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		c.logError(err)
		return err
	}

	return nil
}
