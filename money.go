package money

import (
	"errors"

	"github.com/shopspring/decimal"
)

// Injection points for backward compatibility.
// If you need to keep your JSON marshal/unmarshal way, overwrite them like below.
//
//	money.UnmarshalJSON = func (m *Money, b []byte) error { ... }
//	money.MarshalJSON = func (m Money) ([]byte, error) { ... }
var (
	// ErrCurrencyMismatch happens when two compared Money don't have the same currency.
	ErrCurrencyMismatch = errors.New("currencies don't match")

	// ErrInvalidJSONUnmarshal happens when the default money.UnmarshalJSON fails to unmarshal Money because of invalid data.
	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")
)

// Amount is a data structure that stores the amount being used for calculations.
type Amount = decimal.Decimal

// Money represents monetary value information, stores
// currency and amount value.
type Money struct {
	amount   Amount
	currency *Currency
}

// New creates and returns new instance of Money.
func New(amount decimal.Decimal, code string) *Money {
	return &Money{
		amount:   amount,
		currency: newCurrency(code).get(),
	}
}

// New creates and returns new instance of Money.
func NewFromFloat(amount float64, code string) *Money {
	return &Money{
		amount:   decimal.NewFromFloat(amount),
		currency: newCurrency(code).get(),
	}
}

func (m *Money) GetAmount() *Amount {
	return &m.amount
}

// Currency returns the currency used by Money.
func (m *Money) Currency() *Currency {
	return m.currency
}

// Display lets represent Money struct as string in given Currency value.
func (m *Money) Display() string {
	c := m.currency.get()
	return c.Formatter().Format(m.amount)
}

func (m *Money) Add(a ...*Money) (*Money, error) {
	currency := a[0].currency
	result := m.amount
	for _, v := range a {
		if !v.currency.equals(currency) {
			return nil, errors.New("something went wrong")
		}
		result = result.Add(v.amount)
	}
	return New(result, currency.Code), nil
}

func (m *Money) Convert(currencyCode string, exchangeRate *float64) *Money {
	rate := func() decimal.Decimal {
		if exchangeRate == nil {
			return getExchangeRate(m.Currency().Code, currencyCode)
		} else {
			return decimal.NewFromFloat(*exchangeRate)
		}
	}()

	return New(m.amount.Mul(rate), currencyCode)
}

func getExchangeRate(currencyCode, targetCurrencyCode string) decimal.Decimal {
	if currencyCode == "SGD" && targetCurrencyCode == "CNY" {
		return decimal.NewFromFloat(5.1)
	}

	if currencyCode == "CNY" && targetCurrencyCode == "CNY" {
		return decimal.NewFromFloat(1)
	}

	if currencyCode == "USD" && targetCurrencyCode == "CNY" {
		return decimal.NewFromFloat(7.1)
	}

	return decimal.NewFromFloat(1)
}
