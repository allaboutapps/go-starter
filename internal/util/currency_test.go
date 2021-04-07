package util_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyConversion(t *testing.T) {
	tests := []struct {
		name string
		val  int
	}{
		{
			name: "1",
			val:  999999999999999, // this is the max size with accurate precision
		},
		{
			name: "2",
			val:  0,
		},
		{
			name: "3",
			val:  -999999999999999,
		},
		{
			name: "4",
			val:  3333333333333333,
		},
		{
			name: "5",
			val:  1111111111111111,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			in := int64(tt.val)
			res := util.Int64PtrWithCentsToFloat64Ptr(&in)
			out := util.Float64PtrToInt64PtrWithCents(res)
			assert.Equal(t, in, *out)

			inInt := int(in)
			res = util.IntPtrWithCentsToFloat64Ptr(&inInt)
			outInt := util.Float64PtrToIntPtrWithCents(res)
			assert.Equal(t, inInt, *outInt)
		})
	}
}

func TestCurrencyConversionNil(t *testing.T) {
	res := util.Int64PtrWithCentsToFloat64Ptr(nil)
	assert.Nil(t, res)

	res = util.IntPtrWithCentsToFloat64Ptr(nil)
	assert.Nil(t, res)

	res2 := util.Float64PtrToInt64PtrWithCents(nil)
	assert.Nil(t, res2)

	res3 := util.Float64PtrToIntPtrWithCents(nil)
	assert.Nil(t, res3)
}
