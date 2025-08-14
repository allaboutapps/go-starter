// nolint:revive
package util

import (
	"fmt"
	"time"

	"github.com/go-openapi/swag"
)

const (
	DateFormat       = "2006-01-02"
	monthsPerQuarter = 3
	daysPerWeek      = 7
)

// TimeFromString returns an instance of time.Time from a given string asuming RFC3339 format
func TimeFromString(timeString string) (time.Time, error) {
	result, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time string: %w", err)
	}

	return result, nil
}

func DateFromString(dateString string) (time.Time, error) {
	result, err := time.Parse(DateFormat, dateString)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date string: %w", err)
	}

	return result, nil
}

func EndOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month()+1, 1, 0, 0, 0, -1, d.Location())
}

func EndOfPreviousMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, -1, d.Location())
}

func EndOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, -1, d.Location())
}

func StartOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func StartOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

func StartOfQuarter(d time.Time) time.Time {
	quarter := (int(d.Month()) - 1) / monthsPerQuarter
	m := quarter*monthsPerQuarter + 1
	return time.Date(d.Year(), time.Month(m), 1, 0, 0, 0, 0, d.Location())
}

// StartOfWeek returns the monday (assuming week starts at monday) of the week of the date.
func StartOfWeek(date time.Time) time.Time {
	dayOffset := int(date.Weekday()) - 1

	// go time is starting weeks at sunday
	if dayOffset < 0 {
		dayOffset = 6
	}

	return time.Date(date.Year(), date.Month(), date.Day()-dayOffset, 0, 0, 0, 0, date.Location())
}

func Date(year int, month int, day int, loc *time.Location) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func AddWeeks(d time.Time, weeks int) time.Time {
	return d.AddDate(0, 0, daysPerWeek*weeks)
}

func AddMonths(d time.Time, months int) time.Time {
	return d.AddDate(0, months, 0)
}

func DayBefore(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, -1, d.Location())
}

func TruncateTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// MaxTime returns the latest time.Time of the given params
func MaxTime(times ...time.Time) time.Time {
	var latestTime time.Time
	for _, t := range times {
		if t.After(latestTime) {
			latestTime = t
		}
	}

	return latestTime
}

// NonZeroTimeOrNil returns a pointer to passed time if it is not a zero time. Passing zero/uninitialised time returns nil instead.
func NonZeroTimeOrNil(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}

	return swag.Time(t)
}
