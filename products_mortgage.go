package capis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	querystring "github.com/google/go-querystring/query"
	"go.opencensus.io/trace"
)

type (
	NewMortgageRequest struct {
		ID                       string                 `json:"id"`
		Issuer                   string                 `json:"issuer"`
		Name                     string                 `json:"name"`
		Description              string                 `json:"description"`
		URLApply                 string                 `json:"url_apply"`
		URLLogo                  string                 `json:"url_logo"`
		HighlightedPoints        []string               `json:"highlighted_points"`
		TechnicalPoints          []string               `json:"technical_points"`
		Type                     string                 `json:"type"`
		OfferInterestRate        RatePeriod             `json:"offer_interest_rate"`
		OfferInterestRateType    string                 `json:"offer_interest_rate_type"`
		StandardInterestRate     Rate                   `json:"standard_interest_rate"`
		StandardInterestRateType string                 `json:"standard_interest_rate_type"`
		LoanToValue              Rate                   `json:"loan_to_value"`
		Fee                      Money                  `json:"fee"`
		MinimumLoan              Money                  `json:"minimum_loan"`
		MaximumLoan              Money                  `json:"maximum_loan"`
		MinimumTerm              Months                 `json:"minimum_term"`
		MaximumTerm              Months                 `json:"maximum_term"`
		EarlyRedemptionCharge    Money                  `json:"early_redemption_charge"`
		IsConsumer               bool                   `json:"is_consumer"`
		IsCommercial             bool                   `json:"is_commercial"`
		BrokerOnly               bool                   `json:"broker_only"`
		Active                   bool                   `json:"active"`
		Meta                     map[string]interface{} `json:"metadata"`
	}

	Mortgage struct {
		ID                       string                 `json:"id"`
		Issuer                   string                 `json:"issuer"`
		Name                     string                 `json:"name"`
		Description              string                 `json:"description"`
		URLApply                 string                 `json:"url_apply"`
		URLLogo                  string                 `json:"url_logo"`
		HighlightedPoints        []string               `json:"highlighted_points"`
		TechnicalPoints          []string               `json:"technical_points"`
		Type                     string                 `json:"type"`
		OfferInterestRate        RatePeriod             `json:"offer_interest_rate"`
		OfferInterestRateType    string                 `json:"offer_interest_rate_type"`
		StandardInterestRate     Rate                   `json:"standard_interest_rate"`
		StandardInterestRateType string                 `json:"standard_interest_rate_type"`
		LoanToValue              Rate                   `json:"loan_to_value"`
		Fee                      Money                  `json:"fee"`
		MinimumLoan              Money                  `json:"minimum_loan"`
		MaximumLoan              Money                  `json:"maximum_loan"`
		MinimumTerm              Months                 `json:"minimum_term"`
		MaximumTerm              Months                 `json:"maximum_term"`
		EarlyRedemptionCharge    Money                  `json:"early_redemption_charge"`
		IsConsumer               bool                   `json:"is_consumer"`
		IsCommercial             bool                   `json:"is_commercial"`
		BrokerOnly               bool                   `json:"broker_only"`
		Active                   bool                   `json:"active"`
		Meta                     map[string]interface{} `json:"metadata"`
		Created                  time.Time              `json:"created"`
	}

	ListMortgagesResponse struct {
		Data []*Mortgage `json:"data"`
	}
)

func (s *ProductsService) FindMortgage(ctx context.Context, id string) (*Mortgage, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.FindMortgage")
	defer span.End()

	req, err := s.c.newRequest("GET", fmt.Sprintf("/v1/mortgages/%s", id), nil)

	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		s.c.logError(err)
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	prd := &Mortgage{}
	return prd, unmarshalResponse(res, prd)
}

func (s *ProductsService) UpdateMortgage(ctx context.Context, mortgage *Mortgage) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.UpdateMortgage")
	defer span.End()

	if len(mortgage.ID) == 0 {
		return errors.New("can only update an existing mortgage")
	}

	rb, _ := json.Marshal(mortgage)
	req, err := s.c.newRequest("PUT", fmt.Sprintf("/v1/mortgages/%s", mortgage.ID), bytes.NewReader(rb))

	if err != nil {
		return err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		return ErrUnreachable
	}
	defer res.Body.Close()

	return statusCodeToError(res.StatusCode)
}

func (s *ProductsService) NewMortgage(ctx context.Context, opts *NewMortgageRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.NewMortgage")
	defer span.End()

	rb, _ := json.Marshal(opts)
	req, err := s.c.newRequest("POST", "/v1/mortgages", bytes.NewReader(rb))

	if err != nil {
		return err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		return ErrUnreachable
	}
	defer res.Body.Close()

	return statusCodeToError(res.StatusCode)
}

func (s *ProductsService) ListMortgages(ctx context.Context, filters *ProductFilters) (*ListMortgagesResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.ListMortgages")
	defer span.End()

	obj := &ListMortgagesResponse{}
	qs, _ := querystring.Values(filters)

	req, err := s.c.newRequest("GET", "/v1/mortgages?"+qs.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		s.c.logError(err)
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err = statusCodeToError(res.StatusCode); err != nil {
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}
