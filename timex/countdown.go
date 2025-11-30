package timex

import "time"

// CountdownStatus represents the status of a countdown timer.
type CountdownStatus int

const (
	StatusNotStarted CountdownStatus = iota
	StatusOngoing
	StatusEnded
)

// String returns the string representation of CountdownStatus.
func (s CountdownStatus) String() string {
	switch s {
	case StatusNotStarted:
		return "NotStarted"
	case StatusOngoing:
		return "Ongoing"
	case StatusEnded:
		return "Ended"
	default:
		return "Unknown"
	}
}

// CalculateRemainingTime returns the duration remaining until the specified end time.
func CalculateRemainingTime(endTime time.Time) time.Duration {
	now := time.Now()
	if endTime.Before(now) {
		return 0
	}
	return endTime.Sub(now)
}

// GetCountdownStatus returns the status and relevant duration for a countdown period.
// For NotStarted: returns time until start.
// For Ongoing: returns time until end.
// For Ended: returns zero duration.
func GetCountdownStatus(start, end time.Time) (CountdownStatus, time.Duration) {
	now := time.Now()
	if now.Before(start) {
		return StatusNotStarted, start.Sub(now)
	}
	if now.After(end) {
		return StatusEnded, 0
	}
	return StatusOngoing, end.Sub(now)
}

// IsTimeActive checks if the current time is within the specified start and end times.
func IsTimeActive(start, end time.Time) bool {
	now := time.Now()
	return !now.Before(start) && !now.After(end)
}
