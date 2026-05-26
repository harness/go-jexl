package decimal

import (
	"math/big"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		value int64
		exp   int32
		want  string
	}{
		{123, -2, "1.23"},
		{-456, -1, "-45.6"},
		{0, 0, "0"},
		{1, 3, "1000"},
		{100, -2, "1"},
	}
	for _, c := range cases {
		got := New(c.value, c.exp).String()
		if got != c.want {
			t.Fatalf("New(%d, %d) = %q, want %q", c.value, c.exp, got, c.want)
		}
	}
}

func TestNewFromInt(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-42, "-42"},
		{9999999999, "9999999999"},
	}
	for _, c := range cases {
		got := NewFromInt(c.in).String()
		if got != c.want {
			t.Fatalf("NewFromInt(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNewFromBigInt(t *testing.T) {
	cases := []struct {
		value *big.Int
		exp   int32
		want  string
	}{
		{big.NewInt(100), -2, "1"},
		{big.NewInt(-500), -1, "-50"},
		{big.NewInt(0), 0, "0"},
	}
	for _, c := range cases {
		got := NewFromBigInt(c.value, c.exp).String()
		if got != c.want {
			t.Fatalf("NewFromBigInt(%s, %d) = %q, want %q", c.value, c.exp, got, c.want)
		}
	}
}

func TestNewFromString(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"1.23", "1.23", false},
		{"-4.5", "-4.5", false},
		{"0", "0", false},
		{"1e2", "100", false},
		{"1.5e3", "1500", false},
		{"1.5e-2", "0.015", false},
		{"not-a-number", "", true},
		{"1..2", "", true},
	}
	for _, c := range cases {
		got, err := NewFromString(c.in)
		if c.wantErr {
			if err == nil {
				t.Fatalf("NewFromString(%q) expected error, got nil", c.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("NewFromString(%q) unexpected error: %v", c.in, err)
		}
		if got.String() != c.want {
			t.Fatalf("NewFromString(%q) = %q, want %q", c.in, got.String(), c.want)
		}
	}
}

func TestNewFromFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{1.5, "1.5"},
		{-3.14, "-3.14"},
		{100.0, "100"},
		{0.001, "0.001"},
	}
	for _, c := range cases {
		got := NewFromFloat(c.in).String()
		if got != c.want {
			t.Fatalf("NewFromFloat(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNewFromFloat_panicsOnNaN(t *testing.T) {
	// Ensure NaN input panics rather than silently producing a bad value
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewFromFloat(NaN) expected panic, got none")
		}
	}()
	nan := float64(0)
	nan /= float64(0)
	NewFromFloat(nan)
}

func TestAdd(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"1.5", "2.5", "4"},
		{"0.1", "0.2", "0.3"},
		{"-1", "1", "0"},
		{"100", "-50.5", "49.5"},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Add(b).String()
		if got != c.want {
			t.Fatalf("(%s).Add(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestSub(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"5", "3", "2"},
		{"0.3", "0.1", "0.2"},
		{"1", "1", "0"},
		{"-1", "-1", "0"},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Sub(b).String()
		if got != c.want {
			t.Fatalf("(%s).Sub(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"2", "3", "6"},
		{"1.5", "2", "3"},
		{"-4", "2.5", "-10"},
		{"0", "999", "0"},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Mul(b).String()
		if got != c.want {
			t.Fatalf("(%s).Mul(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestDiv(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"10", "2", "5"},
		{"1", "3", "0.3333333333333333"},
		{"-6", "2", "-3"},
		{"7", "2", "3.5"},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Div(b).String()
		if got != c.want {
			t.Fatalf("(%s).Div(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestMod(t *testing.T) {
	cases := []struct {
		a, b string
		want string
	}{
		{"10", "3", "1"},
		{"7", "2", "1"},
		{"6", "2", "0"},
		{"5.5", "2", "1.5"},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Mod(b).String()
		if got != c.want {
			t.Fatalf("(%s).Mod(%s) = %q, want %q", c.a, c.b, got, c.want)
		}
	}
}

func TestNeg(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"1.5", "-1.5"},
		{"-3", "3"},
		{"0", "0"},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Neg().String()
		if got != c.want {
			t.Fatalf("(%s).Neg() = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAbs(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"-5.5", "5.5"},
		{"3", "3"},
		{"0", "0"},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Abs().String()
		if got != c.want {
			t.Fatalf("(%s).Abs() = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCmp(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"1", "2", -1},
		{"2", "2", 0},
		{"3", "2", 1},
		{"1.5", "1.50", 0},
		{"-1", "0", -1},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Cmp(b)
		if got != c.want {
			t.Fatalf("(%s).Cmp(%s) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1", "1", true},
		{"1.0", "1", true},
		{"1", "2", false},
		{"-1", "-1", true},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		got := a.Equal(b)
		if got != c.want {
			t.Fatalf("(%s).Equal(%s) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"0", true},
		{"0.0", true},
		{"1", false},
		{"-1", false},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.IsZero()
		if got != c.want {
			t.Fatalf("(%s).IsZero() = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestIsNegative(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"-1", true},
		{"0", false},
		{"1", false},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.IsNegative()
		if got != c.want {
			t.Fatalf("(%s).IsNegative() = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSign(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"-5", -1},
		{"0", 0},
		{"5", 1},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Sign()
		if got != c.want {
			t.Fatalf("(%s).Sign() = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestIntPart(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"3.7", 3},
		{"-2.9", -2},
		{"100", 100},
		{"0.5", 0},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.IntPart()
		if got != c.want {
			t.Fatalf("(%s).IntPart() = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestFloat64(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"1.5", 1.5},
		{"-3.14", -3.14},
		{"0", 0},
		{"100", 100},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got, _ := d.Float64()
		if got != c.want {
			t.Fatalf("(%s).Float64() = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		in     string
		places int32
		want   string
	}{
		{"1.456", 2, "1.46"},
		{"1.454", 2, "1.45"},
		{"2.5", 0, "3"},
		{"-2.5", 0, "-3"},
		{"1.005", 2, "1.01"},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Round(c.places).String()
		if got != c.want {
			t.Fatalf("(%s).Round(%d) = %q, want %q", c.in, c.places, got, c.want)
		}
	}
}

func TestPow(t *testing.T) {
	cases := []struct {
		base string
		exp  string
		want string
	}{
		{"2", "10", "1024"},
		{"3", "3", "27"},
		{"2", "0", "1"},
		{"0", "5", "0"},
		{"4", "0.5", "2"},
	}
	for _, c := range cases {
		base, _ := NewFromString(c.base)
		exp, _ := NewFromString(c.exp)
		got := base.Pow(exp).String()
		if got != c.want {
			t.Fatalf("(%s).Pow(%s) = %q, want %q", c.base, c.exp, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	// Ensure String trims trailing zeros
	cases := []struct {
		in   string
		want string
	}{
		{"1.500", "1.5"},
		{"100.00", "100"},
		{"-0.10", "-0.1"},
		{"0.0", "0"},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.String()
		if got != c.want {
			t.Fatalf("String of %q = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExponent(t *testing.T) {
	cases := []struct {
		in   string
		want int32
	}{
		{"1.23", -2},
		{"100", 0},
		{"1e3", 3},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Exponent()
		if got != c.want {
			t.Fatalf("(%s).Exponent() = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestRescalePair(t *testing.T) {
	cases := []struct {
		a, b     string
		wantAExp int32
		wantBExp int32
	}{
		{"1.5", "2.50", -2, -2},
		{"1", "0.001", -3, -3},
		{"1.0", "1.0", -1, -1},
	}
	for _, c := range cases {
		a, _ := NewFromString(c.a)
		b, _ := NewFromString(c.b)
		ra, rb := RescalePair(a, b)
		if ra.exp != c.wantAExp || rb.exp != c.wantBExp {
			t.Fatalf("RescalePair(%s, %s) exps = (%d, %d), want (%d, %d)",
				c.a, c.b, ra.exp, rb.exp, c.wantAExp, c.wantBExp)
		}
	}
}

func TestTruncate(t *testing.T) {
	cases := []struct {
		in        string
		precision int32
		want      string
	}{
		{"1.456", 2, "1.45"},
		{"1.9", 0, "1"},
		{"-1.9", 0, "-1"},
		{"100.123", 1, "100.1"},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.Truncate(c.precision).String()
		if got != c.want {
			t.Fatalf("(%s).Truncate(%d) = %q, want %q", c.in, c.precision, got, c.want)
		}
	}
}

func TestNumDigits(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"0", 1},
		{"9", 1},
		{"10", 2},
		{"999", 3},
		{"-123", 3},
		{"1.23", 3},
	}
	for _, c := range cases {
		d, _ := NewFromString(c.in)
		got := d.NumDigits()
		if got != c.want {
			t.Fatalf("(%s).NumDigits() = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestZeroValue(t *testing.T) {
	// Ensure zero-value Decimal is safe to use without initialization
	var d Decimal
	if !d.IsZero() {
		t.Fatalf("zero-value Decimal.IsZero() = false, want true")
	}
	if d.Sign() != 0 {
		t.Fatalf("zero-value Decimal.Sign() = %d, want 0", d.Sign())
	}
	if d.String() != "0" {
		t.Fatalf("zero-value Decimal.String() = %q, want \"0\"", d.String())
	}
}
