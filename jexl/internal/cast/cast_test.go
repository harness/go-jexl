// Copyright (c) 2014 Steve Francia
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cast

import (
	"encoding/json"
	"testing"
	"time"
)

func TestToBool(t *testing.T) {
	cases := []struct {
		in   any
		want bool
	}{
		{true, true},
		{false, false},
		{int(1), true},
		{int(0), false},
		{int64(1), true},
		{int64(0), false},
		{float64(1.0), true},
		{float64(0.0), false},
		{uint(1), true},
		{uint(0), false},
		{"true", true},
		{"false", false},
		{"1", true},
		{"0", false},
		{json.Number("1"), true},
		{json.Number("0"), false},
		{nil, false},
		{"bad", false},
	}
	for _, c := range cases {
		got := ToBool(c.in)
		if got != c.want {
			t.Fatalf("ToBool(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToFloat64(t *testing.T) {
	cases := []struct {
		in   any
		want float64
	}{
		{float64(3.14), 3.14},
		{float32(1.5), float64(float32(1.5))},
		{int(42), 42},
		{int8(8), 8},
		{int16(16), 16},
		{int32(32), 32},
		{int64(64), 64},
		{uint(10), 10},
		{uint8(8), 8},
		{uint16(16), 16},
		{uint32(32), 32},
		{uint64(64), 64},
		{"3.14", 3.14},
		{json.Number("2.71"), 2.71},
		{true, 1},
		{false, 0},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToFloat64(c.in)
		if got != c.want {
			t.Fatalf("ToFloat64(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToFloat32(t *testing.T) {
	cases := []struct {
		in   any
		want float32
	}{
		{float32(1.5), 1.5},
		{float64(2.0), 2.0},
		{int(5), 5},
		{int64(10), 10},
		{"1.5", 1.5},
		{json.Number("2.5"), 2.5},
		{true, 1},
		{false, 0},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToFloat32(c.in)
		if got != c.want {
			t.Fatalf("ToFloat32(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToInt64(t *testing.T) {
	cases := []struct {
		in   any
		want int64
	}{
		{int(42), 42},
		{int64(64), 64},
		{int32(32), 32},
		{int16(16), 16},
		{int8(8), 8},
		{uint(10), 10},
		{uint64(20), 20},
		{float64(3.9), 3},
		{float32(2.1), 2},
		{"42", 42},
		{"42.0", 42},
		{json.Number("100"), 100},
		{true, 1},
		{false, 0},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToInt64(c.in)
		if got != c.want {
			t.Fatalf("ToInt64(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToInt32(t *testing.T) {
	cases := []struct {
		in   any
		want int32
	}{
		{int32(32), 32},
		{int64(64), 64},
		{float64(1.9), 1},
		{"7", 7},
		{json.Number("8"), 8},
		{true, 1},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToInt32(c.in)
		if got != c.want {
			t.Fatalf("ToInt32(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToInt16(t *testing.T) {
	cases := []struct {
		in   any
		want int16
	}{
		{int16(16), 16},
		{int32(100), 100},
		{float64(3.0), 3},
		{"15", 15},
		{json.Number("20"), 20},
		{true, 1},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToInt16(c.in)
		if got != c.want {
			t.Fatalf("ToInt16(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToInt8(t *testing.T) {
	cases := []struct {
		in   any
		want int8
	}{
		{int8(8), 8},
		{int32(50), 50},
		{float64(2.0), 2},
		{"5", 5},
		{json.Number("6"), 6},
		{true, 1},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToInt8(c.in)
		if got != c.want {
			t.Fatalf("ToInt8(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToInt(t *testing.T) {
	cases := []struct {
		in   any
		want int
	}{
		{int(42), 42},
		{int64(10), 10},
		{float64(7.9), 7},
		{"99", 99},
		{json.Number("50"), 50},
		{true, 1},
		{false, 0},
		{nil, 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToInt(c.in)
		if got != c.want {
			t.Fatalf("ToInt(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToUint(t *testing.T) {
	cases := []struct {
		in   any
		want uint
	}{
		{uint(5), 5},
		{int(10), 10},
		{float64(4.0), 4},
		{"8", 8},
		{json.Number("9"), 9},
		{true, 1},
		{false, 0},
		{nil, 0},
		{int(-1), 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToUint(c.in)
		if got != c.want {
			t.Fatalf("ToUint(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToUint64(t *testing.T) {
	cases := []struct {
		in   any
		want uint64
	}{
		{uint64(100), 100},
		{int(50), 50},
		{float64(3.0), 3},
		{"20", 20},
		{json.Number("30"), 30},
		{true, 1},
		{nil, 0},
		{int64(-1), 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToUint64(c.in)
		if got != c.want {
			t.Fatalf("ToUint64(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToUint32(t *testing.T) {
	cases := []struct {
		in   any
		want uint32
	}{
		{uint32(32), 32},
		{int(16), 16},
		{float64(5.0), 5},
		{"12", 12},
		{json.Number("13"), 13},
		{true, 1},
		{nil, 0},
		{int(-1), 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToUint32(c.in)
		if got != c.want {
			t.Fatalf("ToUint32(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToUint16(t *testing.T) {
	cases := []struct {
		in   any
		want uint16
	}{
		{uint16(16), 16},
		{int(8), 8},
		{float64(6.0), 6},
		{"10", 10},
		{json.Number("11"), 11},
		{true, 1},
		{nil, 0},
		{int(-1), 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToUint16(c.in)
		if got != c.want {
			t.Fatalf("ToUint16(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToUint8(t *testing.T) {
	cases := []struct {
		in   any
		want uint8
	}{
		{uint8(8), 8},
		{int(4), 4},
		{float64(7.0), 7},
		{"9", 9},
		{json.Number("15"), 15},
		{true, 1},
		{nil, 0},
		{int(-1), 0},
		{"bad", 0},
	}
	for _, c := range cases {
		got := ToUint8(c.in)
		if got != c.want {
			t.Fatalf("ToUint8(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToString(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
		{int(42), "42"},
		{int8(8), "8"},
		{int16(16), "16"},
		{int32(32), "32"},
		{int64(64), "64"},
		{uint(10), "10"},
		{uint8(8), "8"},
		{uint16(16), "16"},
		{uint32(32), "32"},
		{uint64(64), "64"},
		{float64(3.14), "3.14"},
		{float32(1.5), "1.5"},
		{json.Number("99"), "99"},
		{[]byte("bytes"), "bytes"},
		{nil, ""},
	}
	for _, c := range cases {
		got := ToString(c.in)
		if got != c.want {
			t.Fatalf("ToString(%#v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToDuration(t *testing.T) {
	cases := []struct {
		in   any
		want time.Duration
	}{
		{time.Second, time.Second},
		{int64(1000), time.Duration(1000)},
		{int(500), time.Duration(500)},
		{float64(100), time.Duration(100)},
		{"1s", time.Second},
		{"500ms", 500 * time.Millisecond},
		{nil, 0},
	}
	for _, c := range cases {
		got := ToDuration(c.in)
		if got != c.want {
			t.Fatalf("ToDuration(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToStringSlice(t *testing.T) {
	cases := []struct {
		in   any
		want []string
	}{
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]any{"x", "y"}, []string{"x", "y"}},
		{[]int{1, 2}, []string{"1", "2"}},
	}
	for _, c := range cases {
		got := ToStringSlice(c.in)
		if len(got) != len(c.want) {
			t.Fatalf("ToStringSlice(%#v) len=%d, want %d", c.in, len(got), len(c.want))
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("ToStringSlice(%#v)[%d] = %q, want %q", c.in, i, got[i], c.want[i])
			}
		}
	}
}

func TestToIntSlice(t *testing.T) {
	cases := []struct {
		in   any
		want []int
	}{
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]any{1, 2, 3}, []int{1, 2, 3}},
	}
	for _, c := range cases {
		got := ToIntSlice(c.in)
		if len(got) != len(c.want) {
			t.Fatalf("ToIntSlice(%#v) len=%d, want %d", c.in, len(got), len(c.want))
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("ToIntSlice(%#v)[%d] = %d, want %d", c.in, i, got[i], c.want[i])
			}
		}
	}
}

func TestToStringMap(t *testing.T) {
	in := map[string]any{"a": 1, "b": "x"}
	got := ToStringMap(in)
	if len(got) != 2 {
		t.Fatalf("ToStringMap len=%d, want 2", len(got))
	}
}

func TestToStringMapString(t *testing.T) {
	in := map[string]any{"k": "v"}
	got := ToStringMapString(in)
	if got["k"] != "v" {
		t.Fatalf("ToStringMapString got %q, want %q", got["k"], "v")
	}
}

func TestToStringMapBool(t *testing.T) {
	in := map[string]any{"yes": true, "no": false}
	got := ToStringMapBool(in)
	if !got["yes"] || got["no"] {
		t.Fatalf("ToStringMapBool unexpected result %v", got)
	}
}

func TestToStringMapInt(t *testing.T) {
	in := map[string]any{"n": 42}
	got := ToStringMapInt(in)
	if got["n"] != 42 {
		t.Fatalf("ToStringMapInt got %d, want 42", got["n"])
	}
}

func TestToStringMapInt64(t *testing.T) {
	in := map[string]any{"n": int64(99)}
	got := ToStringMapInt64(in)
	if got["n"] != 99 {
		t.Fatalf("ToStringMapInt64 got %d, want 99", got["n"])
	}
}

func TestToBoolSlice(t *testing.T) {
	in := []any{true, false, true}
	got := ToBoolSlice(in)
	want := []bool{true, false, true}
	for i, v := range got {
		if v != want[i] {
			t.Fatalf("ToBoolSlice[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestToFloat64Slice(t *testing.T) {
	in := []any{1.1, 2.2, 3.3}
	got := ToFloat64Slice(in)
	if len(got) != 3 {
		t.Fatalf("ToFloat64Slice len=%d, want 3", len(got))
	}
}

func TestToInt64Slice(t *testing.T) {
	in := []any{int64(1), int64(2)}
	got := ToInt64Slice(in)
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("ToInt64Slice unexpected result %v", got)
	}
}

func TestToUintSlice(t *testing.T) {
	in := []any{uint(1), uint(2)}
	got := ToUintSlice(in)
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("ToUintSlice unexpected result %v", got)
	}
}

func TestToSlice(t *testing.T) {
	got := ToSlice([]any{1, "two"})
	if len(got) != 2 {
		t.Fatalf("ToSlice len=%d, want 2", len(got))
	}
}

func TestToStringMapStringSlice(t *testing.T) {
	in := map[string][]string{"k": {"a", "b"}}
	got := ToStringMapStringSlice(in)
	if len(got["k"]) != 2 {
		t.Fatalf("ToStringMapStringSlice unexpected %v", got)
	}
}

func TestToTime(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	cases := []struct {
		in   any
		want time.Time
	}{
		{now, now},
		{int64(1700000000), now},
		{int(1700000000), now},
	}
	for _, c := range cases {
		got := ToTime(c.in)
		if !got.Equal(c.want) {
			t.Fatalf("ToTime(%#v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToString_pointer(t *testing.T) {
	s := "ptr"
	got := ToString(&s)
	if got != "ptr" {
		t.Fatalf("ToString(ptr) = %q, want %q", got, "ptr")
	}
}

func TestToBool_pointer(t *testing.T) {
	v := true
	if !ToBool(&v) {
		t.Fatal("ToBool(ptr true) expected true")
	}
}

func TestToInt_weekday(t *testing.T) {
	got := ToInt(time.Wednesday)
	if got != int(time.Wednesday) {
		t.Fatalf("ToInt(Wednesday) = %d, want %d", got, int(time.Wednesday))
	}
}

func TestToInt_month(t *testing.T) {
	got := ToInt(time.March)
	if got != int(time.March) {
		t.Fatalf("ToInt(March) = %d, want %d", got, int(time.March))
	}
}

func TestToFloat64_pointer(t *testing.T) {
	v := float64(9.9)
	got := ToFloat64(&v)
	if got != 9.9 {
		t.Fatalf("ToFloat64(ptr) = %v, want 9.9", got)
	}
}

func TestToStringSlice_string(t *testing.T) {
	got := ToStringSlice("hello world")
	if len(got) != 2 || got[0] != "hello" || got[1] != "world" {
		t.Fatalf("ToStringSlice(string) unexpected %v", got)
	}
}

func TestToStringSlice_int64s(t *testing.T) {
	got := ToStringSlice([]int64{1, 2})
	if len(got) != 2 || got[0] != "1" || got[1] != "2" {
		t.Fatalf("ToStringSlice([]int64) unexpected %v", got)
	}
}

func TestToStringSlice_uint8s(t *testing.T) {
	got := ToStringSlice([]uint8{65, 66})
	if len(got) != 2 {
		t.Fatalf("ToStringSlice([]uint8) unexpected %v", got)
	}
}

func TestToUint_negativeFloat(t *testing.T) {
	got := ToUint(float64(-1))
	if got != 0 {
		t.Fatalf("ToUint(negative float) = %d, want 0", got)
	}
}

func TestToStringMap_fromJSON(t *testing.T) {
	got := ToStringMap(`{"a":1}`)
	if _, ok := got["a"]; !ok {
		t.Fatalf("ToStringMap(json string) missing key 'a': %v", got)
	}
}

func TestToStringMapString_fromJSON(t *testing.T) {
	got := ToStringMapString(`{"k":"v"}`)
	if got["k"] != "v" {
		t.Fatalf("ToStringMapString(json string) got %q, want %q", got["k"], "v")
	}
}

// Cover all numeric input types for each ToIntXX / ToUintXX function.

func TestToInt32_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), uint(3), uint8(4), uint16(5), uint32(6), uint64(7), float32(8)}
	for _, in := range cases {
		ToInt32(in)
	}
}

func TestToInt16_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int32(2), uint(3), uint8(4), uint16(5), uint32(6), uint64(7), float32(8)}
	for _, in := range cases {
		ToInt16(in)
	}
}

func TestToInt8_allInputTypes(t *testing.T) {
	cases := []any{int16(1), int32(2), uint(3), uint8(4), uint16(5), uint32(6), uint64(7), float32(8)}
	for _, in := range cases {
		ToInt8(in)
	}
}

func TestToInt_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), uint(4), uint8(5), uint16(6), uint32(7), uint64(8), float32(9)}
	for _, in := range cases {
		ToInt(in)
	}
}

func TestToUint_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint8(5), uint16(6), uint32(7), uint64(8), float32(9), float64(10)}
	for _, in := range cases {
		ToUint(in)
	}
}

func TestToUint64_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint8(5), uint16(6), uint32(7), float32(9), bool(true)}
	for _, in := range cases {
		ToUint64(in)
	}
}

func TestToUint32_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint8(5), uint16(6), uint64(7), float32(8), float64(9), bool(false)}
	for _, in := range cases {
		ToUint32(in)
	}
}

func TestToUint16_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint8(5), uint32(6), uint64(7), float32(8), float64(9)}
	for _, in := range cases {
		ToUint16(in)
	}
}

func TestToUint8_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint16(5), uint32(6), uint64(7), float32(8), float64(9)}
	for _, in := range cases {
		ToUint8(in)
	}
}

func TestToFloat32_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), int64(4), uint(5), uint8(6), uint16(7), uint32(8), uint64(9)}
	for _, in := range cases {
		ToFloat32(in)
	}
}

func TestToFloat64_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), uint(4), uint8(5), uint16(6), uint32(7), uint64(8)}
	for _, in := range cases {
		ToFloat64(in)
	}
}

func TestToBool_allInputTypes(t *testing.T) {
	cases := []any{int8(1), int16(0), int32(1), uint(0), uint8(1), uint16(0), uint32(1), uint64(0), float32(1), time.Duration(1)}
	for _, in := range cases {
		ToBool(in)
	}
}

func TestToString_allNumericTypes(t *testing.T) {
	cases := []any{int8(1), int16(2), int32(3), uint(4), uint8(5), uint16(6), uint32(7), uint64(8), float32(9)}
	for _, in := range cases {
		ToString(in)
	}
}

func TestToStringSlice_allTypes(t *testing.T) {
	ToStringSlice([]int8{1, 2})
	ToStringSlice([]int{1, 2})
	ToStringSlice([]int32{1, 2})
	ToStringSlice([]uint{1, 2})
	ToStringSlice([]uint32{1, 2})
	ToStringSlice([]uint64{1, 2})
	ToStringSlice([]float32{1.0, 2.0})
	ToStringSlice([]float64{1.0, 2.0})
}

func TestToStringMapStringSlice_allInputTypes(t *testing.T) {
	ToStringMapStringSlice(map[string][]any{"k": {"a", "b"}})
	ToStringMapStringSlice(map[string]string{"k": "v"})
	ToStringMapStringSlice(map[string]any{"k": []any{"a"}})
	ToStringMapStringSlice(map[string]any{"k": []string{"a"}})
	ToStringMapStringSlice(map[string]any{"k": "v"})
	ToStringMapStringSlice(map[any][]string{"k": {"a"}})
	ToStringMapStringSlice(map[any]string{"k": "v"})
	ToStringMapStringSlice(`{"k":["a","b"]}`)
}

func TestToSlice_mapSlice(t *testing.T) {
	in := []map[string]any{{"a": 1}, {"b": 2}}
	got := ToSlice(in)
	if len(got) != 2 {
		t.Fatalf("ToSlice([]map) len=%d, want 2", len(got))
	}
}

func TestToTime_string(t *testing.T) {
	got := ToTime("2024-01-15")
	if got.IsZero() {
		t.Fatal("ToTime(string) returned zero time")
	}
}

func TestToTimeInDefaultLocation(t *testing.T) {
	got := ToTimeInDefaultLocation(int64(1700000000), time.UTC)
	want := time.Unix(1700000000, 0).UTC()
	if !got.Equal(want) {
		t.Fatalf("ToTimeInDefaultLocation = %v, want %v", got, want)
	}
}

func TestToDuration_string_noUnit(t *testing.T) {
	got := ToDuration("100")
	if got != time.Duration(100) {
		t.Fatalf("ToDuration(\"100\") = %v, want 100ns", got)
	}
}

func TestToDurationSlice(t *testing.T) {
	in := []any{time.Second, time.Millisecond}
	got := ToDurationSlice(in)
	if len(got) != 2 || got[0] != time.Second || got[1] != time.Millisecond {
		t.Fatalf("ToDurationSlice unexpected result %v", got)
	}
}

func TestStringToDate(t *testing.T) {
	cases := []string{
		"2024-01-15",
		"2024-01-15T10:00:00",
		"2024-01-15T10:00:00Z",
		"Mon, 15 Jan 2024 10:00:00 +0000",
	}
	for _, s := range cases {
		got, err := StringToDate(s)
		if err != nil {
			t.Fatalf("StringToDate(%q) error: %v", s, err)
		}
		if got.IsZero() {
			t.Fatalf("StringToDate(%q) returned zero", s)
		}
	}
}

func TestToTime_uint(t *testing.T) {
	ToTime(uint(1700000000))
	ToTime(uint32(1700000000))
	ToTime(uint64(1700000000))
	ToTime(int32(1000))
	ToTime(json.Number("1700000000"))
}

func TestToStringMap_fromJSONString(t *testing.T) {
	got := ToStringMap(`{"x":1,"y":2}`)
	if len(got) != 2 {
		t.Fatalf("ToStringMap json len=%d, want 2", len(got))
	}
}

func TestToStringMapInt_mapAnyAny(t *testing.T) {
	in := map[any]any{"k": 42}
	got := ToStringMapInt(in)
	if got["k"] != 42 {
		t.Fatalf("ToStringMapInt(map[any]any) got %d, want 42", got["k"])
	}
}

func TestToStringMapInt_mapStringInt(t *testing.T) {
	in := map[string]int{"k": 7}
	got := ToStringMapInt(in)
	if got["k"] != 7 {
		t.Fatalf("ToStringMapInt(map[string]int) got %d, want 7", got["k"])
	}
}

func TestToStringMapInt64_mapAnyAny(t *testing.T) {
	in := map[any]any{"k": int64(99)}
	got := ToStringMapInt64(in)
	if got["k"] != 99 {
		t.Fatalf("ToStringMapInt64(map[any]any) got %d, want 99", got["k"])
	}
}

func TestToStringMapStringSlice_mapInterface(t *testing.T) {
	in := map[any]any{"k": []string{"a", "b"}}
	got := ToStringMapStringSlice(in)
	if len(got["k"]) != 2 {
		t.Fatalf("ToStringMapStringSlice(map[any]any) unexpected %v", got)
	}
}

func TestToStringMapStringSlice_mapInterfaceInterface(t *testing.T) {
	in := map[any][]any{"k": {"a", "b"}}
	got := ToStringMapStringSlice(in)
	if len(got["k"]) != 2 {
		t.Fatalf("ToStringMapStringSlice(map[any][]any) unexpected %v", got)
	}
}

func TestToUint_negativeIntTypes(t *testing.T) {
	cases := []any{int8(-1), int16(-1), int32(-1), int64(-1)}
	for _, in := range cases {
		got := ToUint(in)
		if got != 0 {
			t.Fatalf("ToUint(%#v) = %d, want 0", in, got)
		}
	}
}

func TestToUint64_negativeIntTypes(t *testing.T) {
	cases := []any{int8(-1), int16(-1), int32(-1), int64(-1), float32(-1), float64(-1)}
	for _, in := range cases {
		got := ToUint64(in)
		if got != 0 {
			t.Fatalf("ToUint64(%#v) = %d, want 0", in, got)
		}
	}
}

func TestToUint32_negativeIntTypes(t *testing.T) {
	cases := []any{int8(-1), int16(-1), int32(-1), int64(-1), float32(-1), float64(-1)}
	for _, in := range cases {
		got := ToUint32(in)
		if got != 0 {
			t.Fatalf("ToUint32(%#v) = %d, want 0", in, got)
		}
	}
}

func TestToUint16_negativeIntTypes(t *testing.T) {
	cases := []any{int8(-1), int16(-1), int32(-1), int64(-1), float32(-1), float64(-1)}
	for _, in := range cases {
		got := ToUint16(in)
		if got != 0 {
			t.Fatalf("ToUint16(%#v) = %d, want 0", in, got)
		}
	}
}

func TestToUint8_negativeIntTypes(t *testing.T) {
	cases := []any{int8(-1), int16(-1), int32(-1), int64(-1), float32(-1), float64(-1)}
	for _, in := range cases {
		got := ToUint8(in)
		if got != 0 {
			t.Fatalf("ToUint8(%#v) = %d, want 0", in, got)
		}
	}
}

func TestToStringMapString_mapAnyV(t *testing.T) {
	in := map[any]string{"k": "v"}
	got := ToStringMapString(in)
	if got["k"] != "v" {
		t.Fatalf("ToStringMapString(map[any]string) got %q", got["k"])
	}
}

func TestToStringMapString_mapAnyAny(t *testing.T) {
	in := map[any]any{"k": "v"}
	got := ToStringMapString(in)
	if got["k"] != "v" {
		t.Fatalf("ToStringMapString(map[any]any) got %q", got["k"])
	}
}

func TestToStringMap_mapAnyV(t *testing.T) {
	in := map[any]any{"k": 1}
	got := ToStringMap(in)
	if _, ok := got["k"]; !ok {
		t.Fatalf("ToStringMap(map[any]any) missing key: %v", got)
	}
}

func TestToStringMapInt_mapAnyInt(t *testing.T) {
	in := map[any]int{"k": 5}
	got := ToStringMapInt(in)
	if got["k"] != 5 {
		t.Fatalf("ToStringMapInt(map[any]int) got %d, want 5", got["k"])
	}
}

func TestToStringMapInt_fromJSON(t *testing.T) {
	got := ToStringMapInt(`{"n":7}`)
	if got["n"] != 7 {
		t.Fatalf("ToStringMapInt(json) got %d, want 7", got["n"])
	}
}

func TestToStringMapInt64_fromJSON(t *testing.T) {
	got := ToStringMapInt64(`{"n":8}`)
	if got["n"] != 8 {
		t.Fatalf("ToStringMapInt64(json) got %d, want 8", got["n"])
	}
}

func TestToStringMapBool_mapAnyAny(t *testing.T) {
	in := map[any]any{"yes": true}
	got := ToStringMapBool(in)
	if !got["yes"] {
		t.Fatalf("ToStringMapBool(map[any]any) unexpected %v", got)
	}
}

func TestToDuration_floatTypes(t *testing.T) {
	got := ToDuration(float32(200))
	if got != time.Duration(200) {
		t.Fatalf("ToDuration(float32) = %v, want 200", got)
	}
}

func TestToTime_nilAndBadInput(t *testing.T) {
	got := ToTime(nil)
	if !got.IsZero() {
		t.Fatalf("ToTime(nil) expected zero, got %v", got)
	}
}

func TestToTime_timezoneFormats(t *testing.T) {
	got := ToTime("2024-01-15T10:00:00+05:30")
	if got.IsZero() {
		t.Fatal("ToTime(RFC3339 with timezone) returned zero")
	}
	got2 := ToTime("Mon, 15 Jan 2024 10:00:00 UTC")
	if got2.IsZero() {
		t.Fatal("ToTime(RFC1123 with named timezone) returned zero")
	}
}

func TestToStringSlice_errorsAndFallback(t *testing.T) {
	got := ToStringSlice(42)
	if len(got) != 1 || got[0] != "42" {
		t.Fatalf("ToStringSlice(int) = %v, want [\"42\"]", got)
	}
}

func TestResolveAlias_unsupportedKind(t *testing.T) {
	type myStruct struct{ v int }
	_, ok := resolveAlias(myStruct{1})
	if ok {
		t.Fatal("resolveAlias(struct) should return false")
	}
}
