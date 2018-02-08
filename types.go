package capis

type (
	Rate struct {
		Value       float64 `json:"value"`
		Description string  `json:"description"`
	}

	RatePeriod struct {
		Rate

		Period Months `json:"period"`
	}

	Months struct {
		Value       int64  `json:"value"`
		Description string `json:"description"`
	}

	Years struct {
		Value       int64  `json:"value"`
		Description string `json:"description"`
	}

	Money struct {
		Currency    string `json:"currency"`
		Amount      int64  `json:"amount"`
		Description string `json:"description"`
	}
)

func NewRate(v float64, d string) Rate {
	return Rate{v, d}
}

func NewRatePeriod(v float64, d string, m Months) RatePeriod {
	return RatePeriod{
		Rate:   NewRate(v, d),
		Period: m,
	}
}

func NewMonths(m int64, d string) Months {
	return Months{
		Value:       m,
		Description: d,
	}
}

func NewYears(y int64, d string) Years {
	return Years{
		Value:       y,
		Description: d,
	}
}

func NewMoney(c string, a int64, d string) Money {
	return Money{
		Currency:    c,
		Amount:      a,
		Description: d,
	}
}
