package capis

import (
	"context"

	"go.opencensus.io/trace"
)

// Healthy will determine if we can talk to comparisonapis.com and if
// we can check it's health.
func (c *Client) Healthy(ctx context.Context) bool {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/Client.Healthy")
	defer span.End()

	req, err := c.newRequest("GET", "/healthz", nil)
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		c.logError(err)
		return false
	}

	return res.StatusCode == 200
}
