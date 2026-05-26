// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package builtin

import (
	"strings"
	"time"

	"github.com/harness/go-jexl/jexl/coerce"
)

// PlusMinutes adds n minutes to t and returns the result.
func PlusMinutes(t any, n any) time.Time {
	return coerce.ToTime(t).Add(
		time.Duration(coerce.ToInt64(n)) * time.Minute,
	)
}

// PlusHours adds n hours to t and returns the result.
func PlusHours(t any, n any) time.Time {
	return coerce.ToTime(t).Add(
		time.Duration(coerce.ToInt64(n)) * time.Hour,
	)
}

// PlusDays adds n days to t and returns the result.
func PlusDays(t any, n any) time.Time {
	return coerce.ToTime(t).AddDate(0, 0, coerce.ToInt(n))
}

// CurrentDate returns today's date as YYYY-MM-DD.
func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// CurrentTime returns the current time as HH:MM:SS.
func CurrentTime() string {
	return time.Now().Format("15:04:05")
}

// DateFormat formats t using a Java-style pattern (e.g. "yyyy-MM-dd HH:mm:ss").
func DateFormat(t any, pattern any) string {
	return coerce.ToTime(t).Format(
		javaPatternToGo(coerce.ToString(pattern)),
	)
}

// javaPatternToGo converts a Java SimpleDateFormat pattern to a Go layout.
// Tokens are replaced longest-first to avoid partial substitution.
func javaPatternToGo(pattern string) string {
	tokens := []struct {
		java   string
		layout string
	}{
		{"yyyy", "2006"},
		{"yy", "06"},
		{"MM", "01"},
		{"dd", "02"},
		{"HH", "15"},
		{"mm", "04"},
		{"ss", "05"},
		{"SSS", "000"},
		{"Z", "-0700"},
		{"z", "MST"},
	}
	out := pattern
	for _, t := range tokens {
		out = strings.ReplaceAll(out, t.java, t.layout)
	}
	return out
}
