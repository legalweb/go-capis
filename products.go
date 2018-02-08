package capis

type (
	ProductsService struct {
		c *Client
	}

	ProductFilters struct {
	}
)

func (c *Client) Products() *ProductsService {
	return &ProductsService{c}
}
