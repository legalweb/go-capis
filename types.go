package capis

type (
	// Rate is a generic representation from capis.
	Rate struct {
		Value       float64 `json:"value"`
		Description string  `json:"description"`
	}

	// RatePeriod is a generic representation from capis.
	RatePeriod struct {
		Rate

		Period Months `json:"period"`
	}

	// Months is a generic representation from capis.
	Months struct {
		Value       int64  `json:"value"`
		Description string `json:"description"`
	}

	// Years is a generic representation from capis.
	Years struct {
		Value       int64  `json:"value"`
		Description string `json:"description"`
	}

	// Money is a generic representation from capis.
	Money struct {
		Currency    string `json:"currency"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
	}

	// Fee is a generic representation from capis.
	Fee struct {
		Fixed       *Money  `json:"fixed,omitempty"`
		Variable    float64 `json:"variable,omitempty"`
		Description string  `json:"description"`
	}
)

// NewRate returns a new struct representing Rate.
func NewRate(v float64, d string) Rate {
	return Rate{v, d}
}

// NewRatePeriod returns a new struct representing RatePeriod.
func NewRatePeriod(v float64, d string, m Months) RatePeriod {
	return RatePeriod{
		Rate:   NewRate(v, d),
		Period: m,
	}
}

// NewMonths returns a new struct representing Months.
func NewMonths(m int64, d string) Months {
	return Months{
		Value:       m,
		Description: d,
	}
}

// NewYears returns a new struct representing Years.
func NewYears(y int64, d string) Years {
	return Years{
		Value:       y,
		Description: d,
	}
}

// NewMoney returns a new struct representing Money.
func NewMoney(c string, a int64, d string) Money {
	return Money{
		Currency:    c,
		Amount:      a,
		Description: d,
	}
}

// NewFixedFee returns a new struct representing FixedFee.
func NewFixedFee(money *Money, description string) Fee {
	return Fee{
		Fixed:       money,
		Description: description,
	}
}

// NewVariableFee returns a new struct representing VariableFee.
func NewVariableFee(pcent float64, description string) Fee {
	return Fee{
		Variable:    pcent,
		Description: description,
	}
}
