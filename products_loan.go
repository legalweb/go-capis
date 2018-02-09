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
		ID                string                 `json:"id"`
		Issuer            string                 `json:"issuer"`
		Name              string                 `json:"name"`
		Description       string                 `json:"description"`
		URLApply          string                 `json:"url_apply"`
		URLLogo           string                 `json:"url_logo"`
		HighlightedPoints []string               `json:"highlighted_points"`
		TechnicalPoints   []string               `json:"technical_points"`
		InterestRate      Rate                   `json:"interest_rate"`
		MonthlyFee        Money                  `json:"monthly_fee"`
		SetupFee          Money                  `json:"setup_fee"`
		MinimumLoan       Money                  `json:"minimum_loan"`
		MaximumLoan       Money                  `json:"maximum_loan"`
		MinimumTerm       Months                 `json:"minimum_term"`
		MaximumTerm       Months                 `json:"maximum_term"`
		GuarantorAllowed  bool                   `json:"guarantor_allowed"`
		GuarantorCriteria []string               `json:"guarantor_criteria"`
		IsConsumer        bool                   `json:"is_consumer"`
		IsCommercial      bool                   `json:"is_commercial"`
		BrokerOnly        bool                   `json:"broker_only"`
		Active            bool                   `json:"active"`
		Meta              map[string]interface{} `json:"metadata"`
	}

	Loan struct {
		ID                string                 `json:"id"`
		Issuer            string                 `json:"issuer"`
		Name              string                 `json:"name"`
		Description       string                 `json:"description"`
		URLApply          string                 `json:"url_apply"`
		URLLogo           string                 `json:"url_logo"`
		HighlightedPoints []string               `json:"highlighted_points"`
		TechnicalPoints   []string               `json:"technical_points"`
		InterestRate      Rate                   `json:"interest_rate"`
		MonthlyFee        Money                  `json:"monthly_fee"`
		SetupFee          Money                  `json:"setup_fee"`
		MinimumLoan       Money                  `json:"minimum_loan"`
		MaximumLoan       Money                  `json:"maximum_loan"`
		MinimumTerm       Months                 `json:"minimum_term"`
		MaximumTerm       Months                 `json:"maximum_term"`
		GuarantorAllowed  bool                   `json:"guarantor_allowed"`
		GuarantorCriteria []string               `json:"guarantor_criteria"`
		IsConsumer        bool                   `json:"is_consumer"`
		IsCommercial      bool                   `json:"is_commercial"`
		BrokerOnly        bool                   `json:"broker_only"`
		Active            bool                   `json:"active"`
		Meta              map[string]interface{} `json:"metadata"`
		Created           time.Time              `json:"created"`
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

func (s *ProductsService) UpdateLoan(ctx context.Context, loan *Loan) error {
	if len(loan.ID) == 0 {
		return errors.New("can only update an existing loan")
	}

	rb, _ := json.Marshal(loan)
	req, err := s.c.newRequest("PUT", fmt.Sprintf("/v1/loans/%s", loan.ID), bytes.NewReader(rb))

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
