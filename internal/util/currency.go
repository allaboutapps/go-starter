// nolint:revive
package util

import "github.com/go-openapi/swag"

const (
	centFactor = 100
)

func Int64PtrWithCentsToFloat64Ptr(c *int64) *float64 {
	if c == nil {
		return nil
	}

	return Int64WithCentsToFloat64Ptr(*c)
}

func Int64WithCentsToFloat64Ptr(c int64) *float64 {
	return swag.Float64(float64(c) / centFactor)
}

func IntPtrWithCentsToFloat64Ptr(c *int) *float64 {
	if c == nil {
		return nil
	}

	return IntWithCentsToFloat64Ptr(*c)
}

func IntWithCentsToFloat64Ptr(c int) *float64 {
	return swag.Float64(float64(c) / centFactor)
}

func Float64PtrToInt64PtrWithCents(f *float64) *int64 {
	if f == nil {
		return nil
	}

	return swag.Int64(Float64PtrToInt64WithCents(f))
}

func Float64PtrToInt64WithCents(f *float64) int64 {
	return int64(swag.Float64Value(f) * centFactor)
}

func Float64ToInt64WithCents(f float64) int64 {
	return int64(f * centFactor)
}

func Float64PtrToIntPtrWithCents(f *float64) *int {
	if f == nil {
		return nil
	}

	return swag.Int(Float64PtrToIntWithCents(f))
}

func Float64PtrToIntWithCents(f *float64) int {
	return int(swag.Float64Value(f) * centFactor)
}
