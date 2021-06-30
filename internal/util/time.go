package util

import "time"

const (
	DateFormat = "2006-01-02"
)

// TimeFromString returns an instance of time.Time from a given string asuming RFC3339 format
func TimeFromString(timeString string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeString)
}

func DateFromString(dateString string) (time.Time, error) {
	return time.Parse(DateFormat, dateString)
}

func EndOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month()+1, 1, 0, 0, 0, -1, d.Location())
}

func EndOfDay(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day()+1, 0, 0, 0, -1, d.Location())
}

func StartOfMonth(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
}

func StartOfQuarter(d time.Time) time.Time {
	quarter := (int(d.Month()) - 1) / 3
	m := quarter*3 + 1
	return time.Date(d.Year(), time.Month(m), 1, 0, 0, 0, 0, d.Location())
}

// StartOfWeek returns the monday (assuming week starts at monday) of the week of the date
func StartOfWeek(d time.Time) time.Time {
	dayOffset := int(d.Weekday()) - 1

	// go time is starting weeks at sunday
	if dayOffset < 0 {
		dayOffset = 6
	}
	return time.Date(d.Year(), d.Month(), d.Day()-dayOffset, 0, 0, 0, 0, d.Location())
}

func Date(year int, month int, day int, loc *time.Location) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, loc)
}

func AddWeeks(d time.Time, weeks int) time.Time {
	return d.AddDate(0, 0, 7*weeks)
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
