package capis

type (
	// ProductsService ...
	ProductsService struct {
		c *Client
	}

	// ProductFilters ...
	ProductFilters struct {
	}
)

// Products will return a new products service using
// the client for transport comunication.
func (c *Client) Products() *ProductsService {
	return &ProductsService{c}
}
