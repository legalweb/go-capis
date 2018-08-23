package capis

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type (
	// BuildConfigurationResponse contains options used to build embeds.
	BuildConfigurationResponse struct {
		Data []struct {
			Type    ProductType `json:"type"`
			Columns []struct {
				Value string `json:"value"`
				Label string `json:"label"`
			} `json:"columns"`
			Filters []struct {
				Value   string `json:"value"`
				Label   string `json:"label"`
				Mutli   bool   `json:"multi"`
				Choices []struct {
					Value string `json:"value"`
					Label string `json:"label"`
				} `json:"choices"`
			} `json:"filters"`
			AvailableGroups []string `json:"available_groups"`
		} `json:"data"`
	}
)

// GetBuildConfiguration will query capis for the available options to
// build embeds with.
func (c *Client) GetBuildConfiguration(ctx context.Context) (*BuildConfigurationResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.GetBuildConfigurations")
	defer span.End()

	req, err := c.newRequest("GET", "/v1/info/build-configurations", nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	}

	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err := statusCodeToError(res.StatusCode); err != nil {
		return nil, err
	}

	var obj *BuildConfigurationResponse
	return obj, unmarshalResponse(res, obj)
}
