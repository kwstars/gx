package timex

import "time"

// IsWithinTimeRange checks if the target time is within the specified start and end times.
func IsWithinTimeRange(start, end, target time.Time) bool {
	return !target.Before(start) && !target.After(end)
}

// GetTimeRangeDuration returns the duration between the start and end times.
func GetTimeRangeDuration(start, end time.Time) time.Duration {
	return end.Sub(start)
}
