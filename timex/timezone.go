package timex

import "time"

// GetSystemTimeZone returns the system's local time zone location.
func GetSystemTimeZone() *time.Location {
	return time.Local
}

// ConvertTimeZone converts the given time to the specified target time zone location.
func ConvertTimeZone(t time.Time, targetLoc *time.Location) time.Time {
	return t.In(targetLoc)
}

// ConvertToLocalTime converts the given time to the system's local time zone.
func ConvertToLocalTime(t time.Time) time.Time {
	return t.In(time.Local)
}

