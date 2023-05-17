package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"lwebco.de/go-capis"
)

var (
	username = flag.String("username", "", "comparisonapis.com username")
	password = flag.String("password", "", "comparisonapis.com password")
	token    = flag.String("token", "", "comparisonapis.com token")

	groupName = flag.String("group-name", "", "product group name")
)

type (
	mortgageProductsRemote interface {
		ListMortgages(ctx context.Context, filters *capis.MortgageProductFilters) (*capis.ListMortgagesResponse, error)
	}

	groupsRemote interface {
		FindGroup(ctx context.Context, name string) (*capis.FindGroupResponse, error)
	}

	mortgageProductsRepository struct {
		sync.RWMutex

		gr        groupsRemote
		mpr       mortgageProductsRemote
		groupName string
		products  []*capis.Mortgage
	}
)

func (r *mortgageProductsRepository) SyncEvery(ctx context.Context, t *time.Ticker) {
	for range t.C {
		if err := r.Sync(ctx); err != nil {
			log.Println("unable to sync products", err)
		}
	}
}

func (r *mortgageProductsRepository) Sync(ctx context.Context) error {
	r.RWMutex.Lock()
	defer r.RWMutex.Unlock()

	grp, err := r.gr.FindGroup(ctx, r.groupName)
	if err != nil {
		return fmt.Errorf("unable to get group %w", err)
	}

	resp, err := r.mpr.ListMortgages(ctx, &capis.MortgageProductFilters{
		ID: grp.Data.Products,
	})
	if err != nil {
		return fmt.Errorf("unable to get products list %w", err)
	}

	r.products = resp.Data

	return nil
}

func (r *mortgageProductsRepository) All() []*capis.Mortgage {
	r.RWMutex.RLock()
	defer r.RWMutex.RUnlock()

	return r.products
}

func getCapisClient() *capis.Client {
	var auth capis.AuthProvider

	if *token != "" {
		auth = capis.StaticToken(*token)
	} else {
		auth = &capis.PasswordAuthentication{
			Username: *username,
			Password: *password,
		}
	}

	client, err := capis.New(capis.WithAuthProvider(auth))
	if err != nil {
		log.Fatalf("error initialising client %v\n", err)
	}

	return client
}

func repositoryFromClient(c *capis.Client, groupName string) *mortgageProductsRepository {
	return &mortgageProductsRepository{
		gr:        c,
		mpr:       c.Products(),
		groupName: groupName,
		products:  make([]*capis.Mortgage, 0),
	}
}

type sourcingRun struct {
	// request loan amount
	loanAmount int

	// maximum cost of loan in whole pounds
	maxCost int
}

func getFeeCost(loanAmount int, fee capis.Fee) int {
	if fee.Fixed != nil {
		return int(fee.Fixed.Amount)
	}

	costF := (float64(loanAmount) / 100) * fee.Variable
	return int(math.Ceil(costF))
}

func (f sourcingRun) match(mortgage *capis.Mortgage) bool {
	if f.maxCost > 0 && getFeeCost(f.loanAmount, mortgage.Fee) > f.maxCost {
		return false
	}

	return true
}

func findMatching(repo *mortgageProductsRepository, matchFunc func(mortgage *capis.Mortgage) bool) (out []capis.Mortgage) {
	items := repo.All()
	out = make([]capis.Mortgage, 0, len(items))

	for _, item := range items {
		if matchFunc(item) {
			out = append(out, *item)
		}
	}

	return
}

func main() {
	flag.Parse()

	ctx := context.Background()
	client := getCapisClient()

	repo := repositoryFromClient(client, *groupName)

	if err := repo.Sync(ctx); err != nil {
		log.Fatalf("initial sync failied %v", err)
	}
	tick := time.NewTicker(time.Minute)
	defer tick.Stop()
	go repo.SyncEvery(ctx, tick)

	params := sourcingRun{
		loanAmount: 200000,
		maxCost:    10000,
	}

	products := findMatching(repo, params.match)

	if len(products) == 0 {
		fmt.Println("sorry, no products matched the criteria")
		return
	}

	fmt.Printf("%d products where found:\n")
	fmt.Println("===========")

	for _, p := range products {
		fmt.Printf(
			"\t - [%s] %s (%.2f for %d months) + £%d in fees\n",
			p.ID,
			p.Name,
			p.OfferInterestRate.Rate.Value,
			p.OfferInterestRate.Period.Value,
			getFeeCost(params.loanAmount, p.Fee),
		)
	}

	fmt.Println()
}
