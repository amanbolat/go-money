package money_test

import (
	"testing"
	"github.com/shopspring/decimal"
	"github.com/amanbolat/go-money"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestNew(t *testing.T) {
	m := money.New(1, "EUR")
	expect := decimal.New(1, -int32(m.Currency().Fraction))

	assert.Truef(t, m.Amount().Equal(expect), "Expected %s got %s", expect.String(), m.Amount().String())

	assert.Equalf(t, "EUR", m.Currency().Code, "Expected currency %s got %s", "EUR", m.Currency().Code)

	m = money.New(-100, "EUR")
	expect = decimal.New(-100, -int32(m.Currency().Fraction))

	assert.Truef(t, m.Amount().Equal(expect), "Expected %s got %s", expect.String(), m.Amount().String())
}

func TestCurrency(t *testing.T) {
	code := "MOCK"
	decimals := 5
	money.AddCurrency(code, "M$", "1 $", ".", ",", decimals)
	m := money.New(1, code)
	c := m.Currency().Code
	if c != code {
		t.Errorf("Expected %d got %d", code, c)
	}
	f := m.Currency().Fraction
	if f != decimals {
		t.Errorf("Expected %d got %d", decimals, f)
	}
}

func TestMoney_SameCurrency(t *testing.T) {
	m := money.New(0, "EUR")
	om := money.New(0, "USD")

	if m.SameCurrency(om) {
		t.Errorf("Expected %s not to be same as %s", m.Currency().Code, om.Currency().Code)
	}

	om = money.New(0, "EUR")

	if !m.SameCurrency(om) {
		t.Errorf("Expected %s to be same as %s", m.Currency().Code, om.Currency().Code)
	}
}

func TestMoney_Equals(t *testing.T) {
	m := money.New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := money.New(tc.amount, "EUR")
		r, err := m.Equals(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals %d == %t got %t", m.Amount(),
				om.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_GreaterThan(t *testing.T) {
	m := money.New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		om := money.New(tc.amount, "EUR")
		r, err := m.GreaterThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Greater Than %d == %t got %t", m.Amount(),
				om.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_GreaterThanOrEqual(t *testing.T) {
	m := money.New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		om := money.New(tc.amount, "EUR")
		r, err := m.GreaterThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equals Or Greater Than %d == %t got %t", m.Amount(),
				om.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_LessThan(t *testing.T) {
	m := money.New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		om := money.New(tc.amount, "EUR")
		r, err := m.LessThan(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Less Than %d == %t got %t", m.Amount(),
				om.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_LessThanOrEqual(t *testing.T) {
	m := money.New(0, "EUR")
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, true},
	}

	for _, tc := range tcs {
		om := money.New(tc.amount, "EUR")
		r, err := m.LessThanOrEqual(om)

		if err != nil || r != tc.expected {
			t.Errorf("Expected %d Equal Or Less Than %d == %t got %t", m.Amount(),
				om.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_IsZero(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, true},
		{1, false},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.IsZero()

		if r != tc.expected {
			t.Errorf("Expected %d to be zero == %t got %t", m.Amount(), tc.expected, r)
		}
	}
}

func TestMoney_IsNegative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.IsNegative()

		if r != tc.expected {
			t.Errorf("Expected %d to be negative == %t got %t", m.Amount(),
				tc.expected, r)
		}
	}
}

func TestMoney_IsPositive(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected bool
	}{
		{-1, false},
		{0, false},
		{1, true},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.IsPositive()

		if r != tc.expected {
			t.Errorf("Expected %d to be positive == %t got %t", m.Amount(),
				tc.expected, r)
		}
	}
}

func TestMoney_Absolute(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, 1},
		{0, 0},
		{1, 1},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.Absolute().Amount()

		assert.Truef(t, r.Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected absolute %s to be %d got %d", m.Amount().String(), tc.expected, r.String())
	}
}

func TestMoney_Negative(t *testing.T) {
	tcs := []struct {
		amount   int64
		expected int64
	}{
		{-1, -1},
		{0, -0},
		{1, -1},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.Negative().Amount()

		assert.Truef(t, r.Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected negative %s to be %d got %s", m.Amount().String(), tc.expected, r.String())
	}
}

func TestMoney_Add(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 10},
		{10, 5, 15},
		{1, -1, 0},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount1, "EUR")
		om := money.New(tc.amount2, "EUR")
		r, err := m.Add(om)

		if err != nil {
			t.Error(err)
		}

		assert.Truef(t, r.Amount().Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected %d + %d = %d got %d", tc.amount1, tc.amount2, tc.expected, r.Amount())
	}

}

func TestMoney_Add2(t *testing.T) {
	m := money.New(100, "EUR")
	dm := money.New(100, "GBP")
	r, err := m.Add(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Subtract(t *testing.T) {
	tcs := []struct {
		amount1  int64
		amount2  int64
		expected int64
	}{
		{5, 5, 0},
		{10, 5, 5},
		{1, -1, 2},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount1, "EUR")
		om := money.New(tc.amount2, "EUR")
		r, err := m.Subtract(om)

		if err != nil {
			t.Error(err)
		}

		assert.Truef(t, r.Amount().Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected %d - %d = %d got %d", tc.amount1, tc.amount2, tc.expected, r.Amount())
	}
}

func TestMoney_Subtract2(t *testing.T) {
	m := money.New(100, "EUR")
	dm := money.New(100, "GBP")
	r, err := m.Subtract(dm)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Multiply(t *testing.T) {
	tcs := []struct {
		amount     int64
		multiplier int64
		expected   int64
	}{
		{5, 5, 25},
		{10, 5, 50},
		{1, -1, -1},
		{1, 0, 0},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.Multiply(tc.multiplier).Amount()

		assert.Truef(t, r.Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected %d * %d = %d got %d", tc.amount, tc.multiplier, tc.expected, r)
	}
}

func TestMoney_Divide(t *testing.T) {
	tcs := []struct {
		amount   int64
		divisor  int64
		expected int64
	}{
		{5, 5, 1},
		{10, 5, 2},
		{1, -1, -1},
		{10, 3, 3},
		{11, 3, 4},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.Divide(tc.divisor).Round(int32(m.Currency().Fraction)).Amount()

		assert.Truef(t, r.Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected %d / %d = %d got %s", tc.amount, tc.divisor, tc.expected, r.String())
	}
}

func TestMoney_Round(t *testing.T) {
	tcs := []struct {
		amount   int64
		scale    int32
		expected int64
	}{
		{125, 0, 100},
		{175, 0, 200},
		{349, 0, 300},
		{351, 0, 400},
		{0, 0, 0},
		{-1, 0, 0},
		{-75, 0, -100},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		r := m.Round(tc.scale).Amount()

		assert.Truef(t, r.Equal(decimal.New(tc.expected, -int32(m.Currency().Fraction))), "Expected rounded %d to be %d got %s", tc.amount, tc.expected, r.String())
	}
}

func TestMoney_Split(t *testing.T) {
	tcs := []struct {
		amount   int64
		split    int
		expected []decimal.Decimal
	}{
		{100, 3, []decimal.Decimal{decimal.New(34, -2), decimal.New(33, -2), decimal.New(33, -2)}},
		{100, 4, []decimal.Decimal{decimal.New(25, -2), decimal.New(25, -2), decimal.New(25, -2), decimal.New(25, -2)}},
		{5, 3, []decimal.Decimal{decimal.New(2, -2), decimal.New(2, -2), decimal.New(1, -2)}},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		var rs []decimal.Decimal
		split, _ := m.Split(tc.split)

		for _, party := range split {
			rs = append(rs, party.Amount())
		}

		for i, actual := range rs {
			assert.Truef(t, tc.expected[i].Equal(actual), "[%d] Expected %s, but got %s", i+1, tc.expected[i].String(), actual.String())
		}

		//if !reflect.DeepEqual(tc.expected, rs) {
		//	t.Errorf("Expected split of %d to be %v got %v", tc.amount, tc.expected, rs)
		//}
	}
}

func TestMoney_Split2(t *testing.T) {
	m := money.New(100, "EUR")
	r, err := m.Split(-10)

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Allocate(t *testing.T) {
	fn := decimal.New

	tcs := []struct {
		amount   int64
		ratios   []int
		expected []decimal.Decimal
	}{
		{100, []int{50, 50}, []decimal.Decimal{fn(50, -2), fn(50, -2)}},
		{100, []int{30, 30, 30}, []decimal.Decimal{fn(34, -2), fn(33, -2), fn(33, -2)}},
		{200, []int{25, 25, 50}, []decimal.Decimal{fn(50, -2), fn(50, -2), fn(100, -2)}},
		{5, []int{50, 25, 25}, []decimal.Decimal{fn(3, -2), fn(1, -2), fn(1, -2)}},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, "EUR")
		var rs []decimal.Decimal
		split, _ := m.Allocate(tc.ratios...)

		for _, party := range split {
			rs = append(rs, party.Amount())
		}

		for i, expected := range tc.expected {
			assert.Truef(t, expected.Equal(rs[i]), "Expected %s, got %s", expected.String(), rs[i].String())
		}
	}
}

func TestMoney_Allocate2(t *testing.T) {
	m := money.New(100, "EUR")
	r, err := m.Allocate()

	if r != nil || err == nil {
		t.Error("Expected err")
	}
}

func TestMoney_Chain(t *testing.T) {
	m := money.New(10, "EUR")
	om := money.New(5, "EUR")
	// 10 + 5 = 15 / 5 = 3 * 4 = 12 - 5 = 7
	e := int64(7)

	m, err := m.Add(om)

	if err != nil {
		t.Error(err)
	}

	m = m.Divide(5).Multiply(4)
	m, err = m.Subtract(om)


	assert.NoError(t, err)
	assert.Truef(t, m.Amount().Equal(decimal.New(7, -int32(m.Currency().Fraction))), "Expected %d got %s", e, m.Amount().String())
}

func TestMoney_Format(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, "GBP", "£1.00"},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}

}

func TestMoney_Display(t *testing.T) {
	tcs := []struct {
		amount   int64
		code     string
		expected string
	}{
		{100, "AED", "1.00 .\u062f.\u0625"},
		{1, "USD", "$0.01"},
	}

	for _, tc := range tcs {
		m := money.New(tc.amount, tc.code)
		r := m.Display()

		if r != tc.expected {
			t.Errorf("Expected formatted %d to be %s got %s", tc.amount, tc.expected, r)
		}
	}
}

//func TestMoney_Allocate3(t *testing.T) {
//	pound := money.New(100, "GBP")
//	parties, err := pound.Allocate(33, 33, 33)
//
//	if err != nil {
//		t.Error(err)
//	}
//
//	if parties[0].Display() != "£0.34" {
//		t.Errorf("Expected %s got %s", "£0.34", parties[0].Display())
//	}
//
//	if parties[1].Display() != "£0.33" {
//		t.Errorf("Expected %s got %s", "£0.33", parties[1].Display())
//	}
//
//	if parties[2].Display() != "£0.33" {
//		t.Errorf("Expected %s got %s", "£0.33", parties[2].Display())
//	}
//}

func TestMoney_Comparison(t *testing.T) {
	pound := money.New(100, "GBP")
	twoPounds := money.New(200, "GBP")
	twoEuros := money.New(200, "EUR")

	if r, err := pound.GreaterThan(twoPounds); err != nil || r {
		t.Errorf("Expected %d Greater Than %d == %t got %t", pound.Amount(),
			twoPounds.Amount(), false, r)
	}

	if r, err := pound.LessThan(twoPounds); err != nil || !r {
		t.Errorf("Expected %d Less Than %d == %t got %t", pound.Amount(),
			twoPounds.Amount(), true, r)
	}

	if r, err := pound.LessThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThan(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.Equals(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.LessThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}

	if r, err := pound.GreaterThanOrEqual(twoEuros); err == nil || r {
		t.Error("Expected err")
	}
}

func TestMoney_Currency(t *testing.T) {
	pound := money.New(100, "GBP")

	if pound.Currency().Code != "GBP" {
		t.Errorf("Expected %s got %s", "GBP", pound.Currency().Code)
	}
}

func TestMoney_Amount(t *testing.T) {
	pound := money.New(100, "GBP")

	if !pound.Amount().Equal(decimal.New(100, -int32(pound.Currency().Fraction))) {
		t.Errorf("Expected %d got %d", 100, pound.Amount())
	}
}

func TestMoney_MarshalJSON(t *testing.T) {
	usd := money.New(125, "USD")
	b, err := usd.MarshalJSON()
	assert.NoError(t, err)
	t.Log(string(b))
}

func TestMoney_UnmarshalJSON(t *testing.T) {
	jsonB := []byte(`{"amount":"125.22","currency":"usd"}`)
	m := &money.Money{}
	err := json.Unmarshal(jsonB, m)
	assert.NoError(t, err)
	assert.Equal(t, "125.22", m.Amount().String())
	assert.Equal(t, "USD", m.Currency().Code)
}