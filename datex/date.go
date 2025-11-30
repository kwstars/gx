package datex

import "time"

// IsSameDate checks if two time.Time values represent the same calendar date.
func IsSameDate(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// GetDateStart returns a new time.Time representing the start of the day (00:00:00) for the given time.
func GetDateStart(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func CalculateDateDifference(t1, t2 time.Time) int {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	date1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	diff := date2.Sub(date1)
	return int(diff.Hours() / 24)
}

// AddDays adds the specified number of days to the given time.Time.
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths adds the specified number of months to the given time.Time.
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears adds the specified number of years to the given time.Time.
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}
