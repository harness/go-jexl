// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package coerce

import (
	"time"

	"github.com/harness/go-jexl/jexl/internal/cast"
)

// ToBool converts any type to a bool.
func ToBool(v any) bool {
	return cast.ToBool(v)
}

// ToTime converts any type to a time.Time.
func ToTime(v any) time.Time {
	return cast.ToTime(v)
}

// ToDuration converts any type to a time.Duration.
func ToDuration(v any) time.Duration {
	return cast.ToDuration(v)
}

// ToFloat64 converts any type to a float64.
func ToFloat64(v any) float64 {
	return cast.ToFloat64(v)
}

// ToFloat32 converts any type to a float32.
func ToFloat32(v any) float32 {
	return cast.ToFloat32(v)
}

// ToInt64 converts any type to an int64.
func ToInt64(v any) int64 {
	return cast.ToInt64(v)
}

// ToInt32 converts any type to an int32.
func ToInt32(v any) int32 {
	return cast.ToInt32(v)
}

// ToInt16 converts any type to an int16.
func ToInt16(v any) int16 {
	return cast.ToInt16(v)
}

// ToInt8 converts any type to an int8.
func ToInt8(v any) int8 {
	return cast.ToInt8(v)
}

// ToInt converts any type to an int.
func ToInt(v any) int {
	return cast.ToInt(v)
}

// ToUint converts any type to a uint.
func ToUint(v any) uint {
	return cast.ToUint(v)
}

// ToUint64 converts any type to a uint64.
func ToUint64(v any) uint64 {
	return cast.ToUint64(v)
}

// ToUint32 converts any type to a uint32.
func ToUint32(v any) uint32 {
	return cast.ToUint32(v)
}

// ToUint16 converts any type to a uint16.
func ToUint16(v any) uint16 {
	return cast.ToUint16(v)
}

// ToUint8 converts any type to a uint8.
func ToUint8(v any) uint8 {
	return cast.ToUint8(v)
}

// ToString converts any type to a string.
func ToString(v any) string {
	return cast.ToString(v)
}

// ToStringMapString converts any type to a map[string]string.
func ToStringMapString(v any) map[string]string {
	return cast.ToStringMapString(v)
}

// ToStringMapStringSlice converts any type to a map[string][]string.
func ToStringMapStringSlice(v any) map[string][]string {
	return cast.ToStringMapStringSlice(v)
}

// ToStringMapBool converts any type to a map[string]bool.
func ToStringMapBool(v any) map[string]bool {
	return cast.ToStringMapBool(v)
}

// ToStringMapInt converts any type to a map[string]int.
func ToStringMapInt(v any) map[string]int {
	return cast.ToStringMapInt(v)
}

// ToStringMapInt64 converts any type to a map[string]int64.
func ToStringMapInt64(v any) map[string]int64 {
	return cast.ToStringMapInt64(v)
}

// ToStringMap converts any type to a map[string]interface{}.
func ToStringMap(v any) map[string]interface{} {
	return cast.ToStringMap(v)
}

// ToSlice converts any type to a []interface{}.
func ToSlice(v any) []interface{} {
	return cast.ToSlice(v)
}

// ToBoolSlice converts any type to a []bool.
func ToBoolSlice(v any) []bool {
	return cast.ToBoolSlice(v)
}

// ToStringSlice converts any type to a []string.
func ToStringSlice(v any) []string {
	return cast.ToStringSlice(v)
}

// ToIntSlice converts any type to a []int.
func ToIntSlice(v any) []int {
	return cast.ToIntSlice(v)
}

// ToInt64Slice converts any type to a []int64.
func ToInt64Slice(v any) []int64 {
	return cast.ToInt64Slice(v)
}

// ToUintSlice converts any type to a []uint.
func ToUintSlice(v any) []uint {
	return cast.ToUintSlice(v)
}

// ToFloat64Slice converts any type to a []float64.
func ToFloat64Slice(v any) []float64 {
	return cast.ToFloat64Slice(v)
}

// ToDurationSlice converts any type to a []time.Duration.
func ToDurationSlice(v any) []time.Duration {
	return cast.ToDurationSlice(v)
}
