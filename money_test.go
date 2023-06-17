package money

import (
	"testing"

	"github.com/shopspring/decimal"
)

var one = decimal.NewFromInt(1)

func TestNew(t *testing.T) {
	m := New(one, EUR)

	if m.amount != decimal.NewFromInt(1) {
		t.Errorf("Expected %d got %d", 1, m.amount)
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

  
	m = New(decimal.NewFromInt(-100), EUR)

	if m.amount != decimal.NewFromInt(-100) {
		t.Errorf("Expected %d got %d", -100, m.amount)
	}
}

func TestCurrency(t *testing.T) {
	code := "MOCK"
	decimals := 5
	AddCurrency(code, "M$", "1 $", ".", ",", decimals)
	m := New(one, code)
	c := m.Currency().Code
	if c != code {
		t.Errorf("Expected %s got %s", code, c)
	}
	f := m.Currency().Fraction
	if f != decimals {
		t.Errorf("Expected %d got %d", decimals, f)
	}
}

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		amount   decimal.Decimal
		code     string
		expected string
	}{
		{decimal.NewFromInt(-100), AED, "100 .\u062f.\u0625"},
		{decimal.NewFromFloat(0.01), USD, "$0.01"},
	}

	for _, tc := range tcs {
		m := New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}
