package money

import (
	"math"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// Formatter stores Money formatting information.
type Formatter struct {
	Fraction int
	Decimal  string
	Thousand string
	Grapheme string
	Template string
}

// NewFormatter creates new Formatter instance.
func NewFormatter(fraction int, decimal, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// Format returns string of formatted integer using given currency template.
func (f *Formatter) Format(amount decimal.Decimal) string {
	// Work with absolute amount value

	intAmount := amount.Mul(decimal.NewFromInt(int64(math.Pow10(f.Fraction)))).Abs().IntPart()
	sa := strconv.FormatInt(intAmount, 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}

	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount.
	if amount.LessThan(decimal.Zero) {
		sa = "-" + sa
	}

	return sa
}
