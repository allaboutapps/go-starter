package util_test

import (
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartOfMonth(t *testing.T) {
	d := util.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfMonth(d))

	d = util.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfMonth(d))
}

func TestTimeFromString(t *testing.T) {
	expected := time.Date(2020, 3, 29, 12, 34, 54, 0, time.UTC)

	d, err := util.TimeFromString("2020-03-29T12:34:54Z")
	require.NoError(t, err)

	assert.Equal(t, expected, d)
}

func TestStartOfQuarter(t *testing.T) {
	d := util.Date(2020, 3, 31, time.UTC)
	expected := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 1, 1, time.UTC)
	expected = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 12, 1, time.UTC)
	expected = time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 4, 1, time.UTC)
	expected = time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))
}

func TestStartOfWeek(t *testing.T) {
	d := util.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))

	d = util.Date(2020, 6, 15, time.UTC)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))

	d = util.Date(2020, 6, 21, time.UTC)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))
}

func TestDateFromString(t *testing.T) {
	res, err := util.DateFromString("2020-01-03")
	require.NoError(t, err)

	require.True(t, res.Equal(time.Date(2020, 1, 3, 0, 0, 0, 0, time.UTC)))

	res, err = util.DateFromString("2020-xx-03")
	require.Error(t, err)
	assert.Empty(t, res)
}

func TestEndOfMonth(t *testing.T) {
	d := util.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 31, 23, 59, 59, 999999999, time.UTC)
	assert.True(t, expected.Equal(util.EndOfMonth(d)))

	d = util.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 31, 23, 59, 59, 999999999, time.UTC)
	res := util.EndOfMonth(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2021, 1, 31, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.TruncateTime(res))

}

func TestEndOfDay(t *testing.T) {
	d := util.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2020, 3, 12, 23, 59, 59, 999999999, time.UTC)
	assert.True(t, expected.Equal(util.EndOfDay(d)))

	d = util.Date(2020, 12, 35, time.UTC)
	expected = time.Date(2021, 1, 4, 23, 59, 59, 999999999, time.UTC)
	res := util.EndOfDay(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.TruncateTime(res))
}

func TestDateAdds(t *testing.T) {
	d := util.Date(2020, 3, 12, time.UTC)
	expected := time.Date(2022, 4, 12, 0, 0, 0, 0, time.UTC)
	res := util.AddMonths(d, 25)
	assert.True(t, expected.Equal(res))

	d = util.Date(2020, 1, 30, time.UTC)
	expected = time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	res = util.AddMonths(d, 1)
	assert.True(t, expected.Equal(res))

	d = util.Date(2020, 1, 30, time.UTC)
	expected = time.Date(2020, 3, 5, 0, 0, 0, 0, time.UTC)
	res = util.AddWeeks(d, 5)
	assert.True(t, expected.Equal(res))
}

func TestDayBefore(t *testing.T) {
	d := util.Date(2020, 3, 1, time.UTC)
	expected := time.Date(2020, 2, 29, 23, 59, 59, 999999999, time.UTC)
	res := util.DayBefore(d)
	assert.True(t, expected.Equal(res))

	expected = time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.TruncateTime(res))

}
