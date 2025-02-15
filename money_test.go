package money

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

var one = decimal.NewFromInt(1)

func TestNew(t *testing.T) {
	m := New(one, EUR)

	if !m.amount.Equal(decimal.NewFromInt(1)) {
		t.Errorf("Expected %d got %d", 1, m.amount)
	}

	if m.currency.Code != EUR {
		t.Errorf("Expected currency %s got %s", EUR, m.currency.Code)
	}

	m = New(decimal.NewFromInt(-100), EUR)

	if !m.amount.Equal(decimal.NewFromInt(-100)) {
		t.Errorf("Expected %v got %v", "-100", m.amount.String())
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
		{decimal.NewFromInt(1), AED, "1.00 .\u062f.\u0625"},
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

func Test_Add(t *testing.T) {
	one := New(decimal.NewFromFloat(1.23), "CNY")
	two := New(decimal.NewFromFloat(1.23), "CNY")
	three := New(decimal.NewFromFloat(1.23), "CNY")

	result, _ := one.Add(two, three)

	expect := decimal.NewFromFloat(3.69)
	if !result.amount.Equal(expect) {
		t.Errorf("Expected %s got %s", expect.String(), result.amount.String())
	}
}

func TestDecimal(t *testing.T) {
	fmt.Println(decimal.NewFromFloat(1234.234).String())
}

func Test_Convert(t *testing.T) {
	one := New(decimal.NewFromFloat(1), "SGD")

	exchangeRate := 5.3
	result := one.Convert("CNY", &exchangeRate)

	fmt.Println(result.Currency().get().Code)
	fmt.Println(result.amount.String())
}

func Test_MarshalJson(t *testing.T) {
	dec := decimal.NewFromFloat(12.33).MarshalJSON
	json.Marshal(dec)

	money := NewFromFloat(12.33, "CNY")
	result, _ := json.Marshal(money)
	fmt.Println(string(result))
}

func Test_UnMarshalJson(t *testing.T) {
	dec := "{\"amount\":\"12.33\",\"currency\":{\"code\":\"CNY\",\"numeric_code\":\"156\",\"fraction\":2,\"grapheme\":\"元\",\"template\":\"1 $\",\"decimal\":\".\",\"thousand\":\",\"}}"

	var money Money
	err:= json.Unmarshal([]byte(dec), &money)
	if err != nil {
		panic("Unmarshal error")

	}
	fmt.Println(money.amount)
	fmt.Println(money.currency)
}
