package util_test

import (
	"math"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestTypeConversions(t *testing.T) {
	i := 19
	if i > math.MaxInt32 || i < math.MinInt32 {
		t.Error("Int value is out of range")
	}

	i64 := int64(i)
	i32 := int32(i)
	res := util.IntPtrToInt64Ptr(&i)
	assert.Equal(t, &i64, res)

	res2 := util.Int64PtrToIntPtr(&i64)
	assert.Equal(t, &i, res2)

	res3 := util.IntToInt32Ptr(i)
	assert.Equal(t, &i32, res3)
}
