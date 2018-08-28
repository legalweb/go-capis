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
	NewBankAccountRequest struct {
		ID                    string                 `json:"id"`
		Issuer                string                 `json:"issuer"`
		Name                  string                 `json:"name"`
		Description           string                 `json:"description"`
		URLApply              string                 `json:"url_apply"`
		URLLogo               string                 `json:"url_logo"`
		HighlightedPoints     []string               `json:"highlighted_points"`
		TechnicalPoints       []string               `json:"technical_points"`
		OfferInterestRate     RatePeriod             `json:"offer_interest_rate"`
		StandardInterestRate  Rate                   `json:"standard_interest_rate"`
		InterestPaid          string                 `json:"interest_paid"`
		OfferOverdraftRate    RatePeriod             `json:"offer_overdraft_rate"`
		StandardOverdraftRate Rate                   `json:"standard_overdraft_rate"`
		StandardChargeRate    Rate                   `json:"standard_charge_rate"`
		OfferChargeRate       Rate                   `json:"offer_charge_rate"`
		MinimumDeposit        Money                  `json:"deposit_minimum"`
		MaximumDeposit        Money                  `json:"deposit_maximum"`
		AnnualFee             Money                  `json:"annual_fee"`
		MonthlyFee            Money                  `json:"monthly_fee"`
		ApprovalCriteria      string                 `json:"approval_criteria"`
		IsISA                 bool                   `json:"is_isa"`
		IsCapitalProtected    bool                   `json:"is_capital_protected"`
		HasTransactionFees    bool                   `json:"has_transaction_fees"`
		HasOnlineBanking      bool                   `json:"has_online_banking"`
		BrokerOnly            bool                   `json:"broker_only"`
		Active                bool                   `json:"active"`
		Meta                  map[string]interface{} `json:"metadata"`
	}

	BankAccount struct {
		ID                    string                 `json:"id"`
		Issuer                string                 `json:"issuer"`
		Name                  string                 `json:"name"`
		Description           string                 `json:"description"`
		URLApply              string                 `json:"url_apply"`
		URLLogo               string                 `json:"url_logo"`
		HighlightedPoints     []string               `json:"highlighted_points"`
		TechnicalPoints       []string               `json:"technical_points"`
		OfferInterestRate     RatePeriod             `json:"offer_interest_rate"`
		StandardInterestRate  Rate                   `json:"standard_interest_rate"`
		InterestPaid          string                 `json:"interest_paid"`
		OfferOverdraftRate    RatePeriod             `json:"offer_overdraft_rate"`
		StandardOverdraftRate Rate                   `json:"standard_overdraft_rate"`
		StandardChargeRate    Rate                   `json:"standard_charge_rate"`
		OfferChargeRate       Rate                   `json:"offer_charge_rate"`
		MinimumDeposit        Money                  `json:"deposit_minimum"`
		MaximumDeposit        Money                  `json:"deposit_maximum"`
		AnnualFee             Money                  `json:"annual_fee"`
		MonthlyFee            Money                  `json:"monthly_fee"`
		ApprovalCriteria      string                 `json:"approval_criteria"`
		IsISA                 bool                   `json:"is_isa"`
		IsCapitalProtected    bool                   `json:"is_capital_protected"`
		HasTransactionFees    bool                   `json:"has_transaction_fees"`
		HasOnlineBanking      bool                   `json:"has_online_banking"`
		BrokerOnly            bool                   `json:"broker_only"`
		Active                bool                   `json:"active"`
		Meta                  map[string]interface{} `json:"metadata"`
		Created               time.Time              `json:"created"`
	}

	ListBankAccountsResponse struct {
		Data []*BankAccount `json:"data"`
	}
)

func (s *ProductsService) NewBankAccount(ctx context.Context, opts *NewBankAccountRequest) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.NewBankAccount")
	defer span.End()

	rb, _ := json.Marshal(opts)
	req, err := s.c.newRequest("POST", "/v1/bankaccounts", bytes.NewReader(rb))

	if err != nil {
		return err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		s.c.logError(err)
		return ErrUnreachable
	}
	defer res.Body.Close()

	return statusCodeToError(res.StatusCode)
}

func (s *ProductsService) FindBankAccount(ctx context.Context, id string) (*BankAccount, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.FinkBankAccount")
	defer span.End()

	req, err := s.c.newRequest("GET", fmt.Sprintf("/v1/bankaccounts/%s", id), nil)

	if err != nil {
		return nil, err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		s.c.logError(err)
		return nil, ErrUnreachable
	}
	defer res.Body.Close()

	ba := &BankAccount{}
	return ba, unmarshalResponse(res, ba)
}

func (s *ProductsService) UpdateBankAccount(ctx context.Context, bankAccount *BankAccount) error {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.UpdateBankAccount")
	defer span.End()

	if len(bankAccount.ID) == 0 {
		return errors.New("can only update an existing bank account")
	}

	rb, _ := json.Marshal(bankAccount)
	req, err := s.c.newRequest("PUT", fmt.Sprintf("/v1/bankaccounts/%s", bankAccount.ID), bytes.NewReader(rb))

	if err != nil {
		return err
	}

	res, err := s.c.Do(req.WithContext(ctx))
	if err != nil {
		s.c.logError(err)
		return ErrUnreachable
	}
	defer res.Body.Close()

	return statusCodeToError(res.StatusCode)
}

func (s *ProductsService) ListBankAccounts(ctx context.Context, filters *ProductFilters) (*ListBankAccountsResponse, error) {
	ctx, span := trace.StartSpan(ctx, "lwebco.de/go-capis/ProductsService.ListBankAccounts")
	defer span.End()

	obj := &ListBankAccountsResponse{}
	qs, _ := querystring.Values(filters)

	req, err := s.c.newRequest("GET", "/v1/bankaccounts?"+qs.Encode(), nil)
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
