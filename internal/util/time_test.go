package util_test

import (
	"testing"
	"time"

	"allaboutapps.dev/aw/go-starter/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartOfMonth(t *testing.T) {
	t.Parallel()

	d := util.Date(2020, 3, 12)
	expected := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfMonth(d))

	d = util.Date(2020, 12, 35)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfMonth(d))
}

func TestTimeFromString(t *testing.T) {
	t.Parallel()

	expected := time.Date(2020, 3, 29, 12, 34, 54, 0, time.UTC)

	d, err := util.TimeFromString("2020-03-29T12:34:54Z")
	require.NoError(t, err)

	assert.Equal(t, expected, d)
}

func TestStartOfQuarter(t *testing.T) {
	t.Parallel()

	d := util.Date(2020, 3, 31)
	expected := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 1, 1)
	expected = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 12, 1)
	expected = time.Date(2020, 10, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 12, 35)
	expected = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))

	d = util.Date(2020, 4, 1)
	expected = time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfQuarter(d))
}

func TestStartOfWeek(t *testing.T) {
	t.Parallel()

	d := util.Date(2020, 3, 12)
	expected := time.Date(2020, 3, 9, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))

	d = util.Date(2020, 6, 15)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))

	d = util.Date(2020, 6, 21)
	expected = time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expected, util.StartOfWeek(d))
}
