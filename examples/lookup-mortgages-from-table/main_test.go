package main

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"lwebco.de/go-capis"
)

type (
	stubGroup struct {
		groupName string
		products  []*capis.Mortgage
	}
)

func (g *stubGroup) FindGroup(_ context.Context, name string) (*capis.FindGroupResponse, error) {
	if name != g.groupName {
		return nil, capis.ErrNotFound
	}

	ids := make([]string, len(g.products))
	for i, p := range g.products {
		ids[i] = p.ID
	}

	return &capis.FindGroupResponse{
		Data: &capis.DetailedGroup{
			ID:       name,
			Type:     "mortgage",
			Products: ids,
		},
	}, nil
}

func (g *stubGroup) ListMortgages(_ context.Context, filters *capis.MortgageProductFilters) (*capis.ListMortgagesResponse, error) {
	out := make([]*capis.Mortgage, 0, len(filters.ID))

	for _, expected := range filters.ID {
		for _, p := range g.products {
			if expected == p.ID {
				cp := *p // not a deep copy but hey
				out = append(out, &cp)
			}
		}
	}

	return &capis.ListMortgagesResponse{Data: out}, nil
}

func TestSync(t *testing.T) {
	expected := []*capis.Mortgage{
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e2",
		},
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e3",
		},
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e4",
		},
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e5",
		},
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e6",
		},
		{
			ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e7",
		},
	}

	group := &stubGroup{
		groupName: "testing",
		products:  expected,
	}

	sut := &mortgageProductsRepository{
		RWMutex:   sync.RWMutex{},
		gr:        group,
		mpr:       group,
		groupName: "testing",
		products:  make([]*capis.Mortgage, 0),
	}

	assert.NoError(t, sut.Sync(context.Background()))
	assert.Equal(t, sut.All(), expected)
}

func TestMatch(t *testing.T) {
	sut := &mortgageProductsRepository{
		groupName: "testing",
		products: []*capis.Mortgage{
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e2",
				Fee: capis.Fee{
					Variable:    5.1,
					Description: "",
				},
			},
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e3",
				Fee: capis.Fee{
					Variable:    3.1,
					Description: "",
				},
			},
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e4",
				Fee: capis.Fee{
					Variable:    8.1,
					Description: "",
				},
			},
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e5",
				Fee: capis.Fee{
					Variable:    5.1,
					Description: "",
				},
			},
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e6",
				Fee: capis.Fee{
					Variable:    7.1,
					Description: "",
				},
			},
			{
				ID: "e2db12d8-af0d-4bd8-8df2-1cdba3a5b7e7",
				Fee: capis.Fee{
					Variable:    5.6,
					Description: "",
				},
			},
		},
	}

	assert.Len(t, findMatching(sut, func(m *capis.Mortgage) bool {
		return m.Fee.Variable <= 5.1
	}), 3)
}

func TestSourcingRun(t *testing.T) {
	tcs := []struct {
		Params   sourcingRun
		Mortgage *capis.Mortgage
		Expected bool
	}{
		{
			Params: sourcingRun{
				loanAmount: 200000,
				maxCost:    10000, // 5%
			},
			Mortgage: &capis.Mortgage{
				Fee: capis.Fee{
					Variable:    6.0,
					Description: "",
				},
			},
			Expected: false,
		},
		{
			Params: sourcingRun{
				loanAmount: 200000,
				maxCost:    10000, // 5%
			},
			Mortgage: &capis.Mortgage{
				Fee: capis.Fee{
					Variable:    4.9,
					Description: "",
				},
			},
			Expected: true,
		},
		{
			Params: sourcingRun{
				loanAmount: 200000,
				maxCost:    10000, // 5%
			},
			Mortgage: &capis.Mortgage{
				Fee: capis.Fee{
					Fixed: &capis.Money{
						Amount: 99999,
					},
					Description: "",
				},
			},
			Expected: false,
		},
		{
			Params: sourcingRun{
				loanAmount: 200000,
				maxCost:    10000, // 5%
			},
			Mortgage: &capis.Mortgage{
				Fee: capis.Fee{
					Fixed: &capis.Money{
						Amount: 9999,
					},
					Description: "",
				},
			},
			Expected: true,
		},
	}

	for tci, tc := range tcs {
		assert.Equal(t, tc.Expected, tc.Params.match(tc.Mortgage), "test case %d did not match", tci)
	}
}
