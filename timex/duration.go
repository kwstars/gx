package timex

import "time"

// RoundDurationToMinutes rounds a duration to the nearest whole minute.
func RoundDurationToMinutes(duration time.Duration) time.Duration {
	minutes := duration.Minutes()
	roundedMinutes := float64(int(minutes + 0.5))
	return time.Duration(roundedMinutes) * time.Minute
}

// ClampDuration restricts a duration to be within the specified min and max bounds.
func ClampDuration(duration, min, max time.Duration) time.Duration {
	if duration < min {
		return min
	}
	if duration > max {
		return max
	}
	return duration
}

// SplitDuration splits a duration into hours, minutes, and seconds.
func SplitDuration(duration time.Duration) (hours int, minutes int, seconds int) {
	totalSeconds := int(duration.Seconds())
	hours = totalSeconds / 3600
	minutes = (totalSeconds % 3600) / 60
	seconds = totalSeconds % 60
	return hours, minutes, seconds
}
