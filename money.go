package money

import (
	"errors"
	"github.com/shopspring/decimal"
	"strings"
	"encoding/json"
)

// Money represents monetary value information, stores
// currency and amount value
type Money struct {
	amount   decimal.Decimal
	currency *Currency
}

func (m *Money) UnmarshalJSON(data []byte) error {
	s := &struct {
		Amount decimal.Decimal `json:"amount"`
		Currency string `json:"currency"`
	}{}
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	val :=  NewFromDecimal(s.Amount, s.Currency)
	*m = *val

	return nil
}

func (m Money) MarshalJSON() ([]byte, error) {
	var currency string
	if m.currency != nil {
		currency = m.Currency().Code
	}
	s := &struct {
		Amount decimal.Decimal `json:"amount"`
		Currency string `json:"currency"`
	}{
		Amount: m.amount,
		Currency: currency,
	}

	return json.Marshal(s)
}



// New creates and returns new instance of Money
// amount should be in cents for currency
// Example: New(100, "EUR") = 1 EUR
func New(amount int64, code string) *Money {
	c := newCurrency(code).get()
	return &Money{
		amount:   decimal.New(amount, -int32(c.Fraction)),
		currency: c,
	}
}

// NewFromDecimal creates Money instance from decimal.Decimal
// and rounds it by currency Fraction
func NewFromDecimal(amount decimal.Decimal, code string) *Money {
	c := newCurrency(code).get()
	return &Money{
		amount:   amount.Round(int32(c.Fraction)),
		currency: c,
	}
}

// Currency returns the currency used by Money
func (m *Money) Currency() *Currency {
	return m.currency
}

// Amount returns a copy of the internal monetary value as an int64
func (m *Money) Amount() decimal.Decimal {
	return m.amount
}

// SameCurrency check if given Money is equals by currency
func (m *Money) SameCurrency(om *Money) bool {
	return m.currency.equals(om.currency)
}

func (m *Money) assertSameCurrency(om *Money) error {
	if !m.SameCurrency(om) {
		return errors.New("currencies don't match")
	}

	return nil
}

// Equals checks equality between two Money types
func (m *Money) Equals(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount.Equal(om.amount), nil
}

// GreaterThan checks whether the value of Money is greater than the other
func (m *Money) GreaterThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount.GreaterThan(om.amount), nil
}

// GreaterThanOrEqual checks whether the value of Money is greater or equal than the other
func (m *Money) GreaterThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount.GreaterThanOrEqual(om.amount), nil
}

// LessThan checks whether the value of Money is less than the other
func (m *Money) LessThan(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount.LessThan(om.amount), nil
}

// LessThanOrEqual checks whether the value of Money is less or equal than the other
func (m *Money) LessThanOrEqual(om *Money) (bool, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return false, err
	}

	return m.amount.LessThanOrEqual(om.amount), nil
}

// IsZero returns boolean of whether the value of Money is equals to zero
func (m *Money) IsZero() bool {
	return m.amount.Equal(decimal.Zero)
}

// IsPositive returns boolean of whether the value of Money is positive
func (m *Money) IsPositive() bool {
	return m.amount.Sign() == 1
}

// IsNegative returns boolean of whether the value of Money is negative
func (m *Money) IsNegative() bool {
	return m.amount.Sign() == -1
}

// Absolute returns new Money struct from given Money using absolute monetary value
func (m *Money) Absolute() *Money {
	return &Money{amount: m.amount.Abs(), currency: m.currency}
}

// Negative returns new Money struct from given Money using negative monetary value
func (m *Money) Negative() *Money {
	if m.IsNegative() {
		return &Money{amount: m.amount, currency: m.currency}
	}
	return &Money{amount: m.amount.Neg(), currency: m.currency}
}

// Add returns new Money struct with value representing sum of Self and Other Money
func (m *Money) Add(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: m.amount.Add(om.amount), currency: m.currency}, nil
}

// Subtract returns new Money struct with value representing difference of Self and Other Money
func (m *Money) Subtract(om *Money) (*Money, error) {
	if err := m.assertSameCurrency(om); err != nil {
		return nil, err
	}

	return &Money{amount: m.amount.Sub(om.amount), currency: m.currency}, nil
}

// Multiply returns new Money struct with value representing Self multiplied value by multiplier
func (m *Money) Multiply(mul int64) *Money {
	return &Money{amount: m.amount.Mul(decimal.New(mul, 0)), currency: m.currency}
}

// Divide returns new Money struct with value representing Self division value by given divider
func (m *Money) Divide(div int64) *Money {
	return &Money{amount: m.amount.Div(decimal.New(div, 0)), currency: m.currency}
}

// Round returns new Money struct with value rounded to nearest zero
func (m *Money) Round(scale int32) *Money {
	//return &Money{amount: m.amount.Round(int32(c.Fraction)), currency: m.currency}
	return &Money{amount:m.amount.Round(scale), currency: m.currency}
}

// Split returns slice of Money structs with split Self value in given number.
// After division leftover pennies will be distributed round-robin amongst the parties.
// This means that parties listed first will likely receive more pennies than ones that are listed later
func (m *Money) Split(n int) ([]*Money, error) {
	if n <= 0 {
		return nil, errors.New("split must be higher than zero")
	}

	arr := make([]*Money, n)
	quo, rem := m.amount.QuoRem(decimal.NewFromFloat(float64(n)), int32(m.currency.Fraction))

	// 1 with reminder exponent for subtraction
	remUnit := decimal.New(1, rem.Exponent())

	for i := 0; i < n; i++ {
		if !rem.Equal(decimal.Zero) {
			rem = rem.Sub(remUnit)
			arr[i] = &Money{amount: quo.Add(remUnit), currency: m.currency}
		} else {
			arr[i] = &Money{amount: quo, currency: m.currency}
		}
	}

	//var idx int
	//for !rem.Equal(decimal.Zero) {
	//	one := decimal.New(1, rem.Exponent())
	//	rem = rem.Sub(one)
	//	arr[idx].amount = arr[idx].amount.Add(one)
	//	idx++
	//}

	return arr, nil
}

// Allocate returns slice of Money structs with split Self value in given ratios.
// It lets split money by given ratios without losing pennies and as Split operations distributes
// leftover pennies amongst the parties with round-robin principle.
func (m *Money) Allocate(ratios ...int) ([]*Money, error) {
	if len(ratios) == 0 {
		return nil, errors.New("no ratios specified")
	}

	// Calculate sum of ratios
	var sum int
	for _, r := range ratios {
		sum += r
	}

	var total decimal.Decimal
	var resultMoneys []*Money
	for _, ratio := range ratios {
		party := &Money{
			amount:  m.amount.Mul(decimal.New(int64(ratio), 0)).DivRound(decimal.New(int64(sum), 0), int32(m.currency.Fraction)),
			currency: m.currency,
		}

		resultMoneys = append(resultMoneys, party)
		total = total.Add(party.amount)
	}

	// Calculate leftover value and divide to first parties
	left := m.amount.Sub(total)

	unit := decimal.New(1, left.Exponent())
	if left.LessThan(decimal.Zero) {
		unit = unit.Neg()
	}

	// 1 with currency fraction


	for i := 0; !left.Equal(decimal.Zero); i++ {
		resultMoneys[i].amount = resultMoneys[i].amount.Add(unit) //mutate.calc.add(resultMoneys[i].amount, &Amount{sub})
		left = left.Sub(unit)//-= sub
	}

	return resultMoneys, nil
}

// Display lets represent Money struct as string in given Currency value
func (m *Money) Display() string {
	c := m.currency.get()

	str := m.amount.Abs().StringFixed(int32(c.Fraction))

	str = strings.Replace(c.Template, "1", str, 1)
	str = strings.Replace(str, "$", c.Grapheme, 1)

	if m.IsNegative() {
		str = "-" + str
	}

	return str
}
