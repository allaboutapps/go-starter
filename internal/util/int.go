// nolint:revive
package util

import (
	"math"

	"github.com/go-openapi/swag/conv"
)

func IntPtrToInt64Ptr(num *int) *int64 {
	if num == nil {
		return nil
	}

	return conv.Pointer(int64(*num))
}

func Int64PtrToIntPtr(num *int64) *int {
	if num == nil {
		return nil
	}

	return conv.Pointer(int(*num))
}

func IntToInt32Ptr(num int) *int32 {
	if num > math.MaxInt32 || num < math.MinInt32 {
		return nil
	}

	return conv.Pointer(int32(num))
}
