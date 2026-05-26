// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"regexp"
	"testing"
	"time"
)

// Ensure CurrentDate returns a string matching YYYY-MM-DD.
func TestCurrentDate_format(t *testing.T) {
	got := CurrentDate()
	if !regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(got) {
		t.Fatalf("unexpected format: %s", got)
	}
}

// Ensure CurrentTime returns a string matching HH:MM:SS.
func TestCurrentTime_format(t *testing.T) {
	got := CurrentTime()
	if !regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`).MatchString(got) {
		t.Fatalf("unexpected format: %s", got)
	}
}

// Ensure DateFormat formats a time.Time value with a Java-style pattern.
func TestDateFormat_dateOnly(t *testing.T) {
	d := time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)
	if DateFormat(d, "yyyy-MM-dd") != "2026-05-24" {
		t.Fatalf("unexpected result: %s", DateFormat(d, "yyyy-MM-dd"))
	}
}

// Ensure DateFormat includes time components when pattern contains HH:mm:ss.
func TestDateFormat_dateTime(t *testing.T) {
	d := time.Date(2026, 5, 24, 13, 45, 30, 0, time.UTC)
	if DateFormat(d, "yyyy-MM-dd HH:mm:ss") != "2026-05-24 13:45:30" {
		t.Fatalf("unexpected result: %s", DateFormat(d, "yyyy-MM-dd HH:mm:ss"))
	}
}

// Ensure DateFormat handles US-style month/day/year pattern.
func TestDateFormat_usStyle(t *testing.T) {
	d := time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)
	if DateFormat(d, "MM/dd/yyyy") != "05/24/2026" {
		t.Fatalf("unexpected result: %s", DateFormat(d, "MM/dd/yyyy"))
	}
}

// Ensure PlusMinutes adds the correct duration.
func TestPlusMinutes_basic(t *testing.T) {
	d := time.Date(2026, 5, 24, 10, 0, 0, 0, time.UTC)
	got := PlusMinutes(d, 30)
	if got.Hour() != 10 || got.Minute() != 30 {
		t.Fatalf("expected 10:30, got %v", got)
	}
}

// Ensure PlusHours adds the correct duration.
func TestPlusHours_basic(t *testing.T) {
	d := time.Date(2026, 5, 24, 10, 0, 0, 0, time.UTC)
	got := PlusHours(d, 3)
	if got.Hour() != 13 {
		t.Fatalf("expected hour 13, got %d", got.Hour())
	}
}

// Ensure PlusDays advances the date by the given number of days.
func TestPlusDays_basic(t *testing.T) {
	d := time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)
	got := PlusDays(d, 7)
	if got.Day() != 31 || got.Month() != time.May {
		t.Fatalf("expected 2026-05-31, got %v", got)
	}
}
