package timex

import (
	"fmt"
	"time"
)

// FormatDuration formats a time.Duration into a string in the format "HH:MM:SS".
func FormatDuration(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// FormatCoolDown formats a time.Duration into a compact string representation.
func FormatCoolDown(duration time.Duration) string {
	totalSeconds := int(duration.Seconds())
	if totalSeconds < 60 {
		return fmt.Sprintf("%ds", totalSeconds)
	} else if totalSeconds < 3600 {
		minutes := totalSeconds / 60
		seconds := totalSeconds % 60
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	} else {
		hours := totalSeconds / 3600
		minutes := (totalSeconds % 3600) / 60
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
}
