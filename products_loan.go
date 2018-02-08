package capis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	querystring "github.com/google/go-querystring/query"
)

type (
	NewLoanRequest struct {
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
		EarlyRedemptionCharge    Money                  `json:"early_redemption_charge"`
		IsConsumer               bool                   `json:"is_consumer"`
		IsCommercial             bool                   `json:"is_commercial"`
		BrokerOnly               bool                   `json:"broker_only"`
		Active                   bool                   `json:"active"`
		Meta                     map[string]interface{} `json:"metadata"`
	}

	Loan struct {
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
		EarlyRedemptionCharge    Money                  `json:"early_redemption_charge"`
		IsConsumer               bool                   `json:"is_consumer"`
		IsCommercial             bool                   `json:"is_commercial"`
		BrokerOnly               bool                   `json:"broker_only"`
		Active                   bool                   `json:"active"`
		Meta                     map[string]interface{} `json:"metadata"`
		Created                  time.Time              `json:"created"`
	}

	ListLoansResponse struct {
		Data []*Loan `json:"data"`
	}
)

func (s *ProductsService) FindLoan(ctx context.Context, id string) (*Loan, error) {
	req, err := s.c.newRequest("GET", fmt.Sprintf("/v1/loans/%s", id), nil)

	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	prd := &Loan{}
	return prd, unmarshalResponse(res, prd)
}

func (s *ProductsService) UpdateLoan(ctx context.Context, mortgage *Loan) error {
	if len(mortgage.ID) == 0 {
		return errors.New("can only update an existing mortgage")
	}

	rb, _ := json.Marshal(mortgage)
	req, err := s.c.newRequest("PUT", fmt.Sprintf("/v1/loans/%s", mortgage.ID), bytes.NewReader(rb))

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

func (s *ProductsService) NewLoan(ctx context.Context, opts *NewLoanRequest) error {
	rb, _ := json.Marshal(opts)
	req, err := s.c.newRequest("POST", "/v1/loans", bytes.NewReader(rb))

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

func (s *ProductsService) ListLoans(ctx context.Context, filters *ProductFilters) (*ListLoansResponse, error) {
	obj := &ListLoansResponse{}
	qs, _ := querystring.Values(filters)

	req, err := s.c.newRequest("GET", "/v1/loans?"+qs.Encode(), nil)
	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	if err = statusCodeToError(res.StatusCode); err != nil {
		return nil, err
	}

	return obj, unmarshalResponse(res, obj)
}
