package capis

import (
	"context"

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
		ID       []string `json:"embed_ids"`
		Metadata []string `json:"metadata"`
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
