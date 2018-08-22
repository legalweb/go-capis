package capis

// Healthy will determine if we can talk to comparisonapis.com and if
// we can check it's health.
func (c *Client) Healthy() bool {
	req, err := c.newRequest("GET", "/healthz", nil)
	res, err := c.Do(req)
	if err != nil {
		return false
	}

	return res.StatusCode == 200
}
