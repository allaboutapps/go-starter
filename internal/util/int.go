package util

import "github.com/go-openapi/swag"

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
	return swag.Int32(int32(num))
}
