package datex

import "time"

// IsWithinDateRange checks if the target date is within the specified start and end dates.
func IsWithinDateRange(startDate, endDate, targetDate time.Time) bool {
	return !targetDate.Before(startDate) && !targetDate.After(endDate)
}

// GetDateRangeDays returns the number of days between the start and end dates.
func GetDateRangeDays(startDate, endDate time.Time) int {
	y1, m1, d1 := startDate.Date()
	y2, m2, d2 := endDate.Date()
	date1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)
	diff := date2.Sub(date1)
	return int(diff.Hours() / 24)
}
