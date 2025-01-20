package util

import (
	"math"

	"github.com/go-openapi/swag"
)

func IntPtrToInt64Ptr(num *int) *int64 {
	if num == nil {
		return nil
	}

	return swag.Int64(int64(*num))
}

func Int64PtrToIntPtr(num *int64) *int {
	if num == nil {
		return nil
	}

	return swag.Int(int(*num))
}

func IntToInt32Ptr(num int) *int32 {
	if num > math.MaxInt32 || num < math.MinInt32 {
		return nil
	}

	return swag.Int32(int32(num))
}
