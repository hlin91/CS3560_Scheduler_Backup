// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
	"math"
	"time"
)

const (
	// Transient types
	VISIT       = "Visit"
	SHOPPING    = "Shopping"
	APPOINTMENT = "Appointment"
	// Anti types
	CANCEL = "Cancel"
	// Recurring Types
	CLASS    = "Class"
	STUDY    = "Study"
	SLEEP    = "Sleep"
	EXERCISE = "Exercise"
	WORK     = "Work"
	MEAL     = "Meal"
)

// isTransientType checks if the type is a valid transient type
func isTransientType(s string) bool {
	return s == VISIT || s == SHOPPING || s == APPOINTMENT
}

// isAntiType checks if the type is a valid anti type
func isAntiType(s string) bool {
	return s == CANCEL
}

// isRecurringType checks if the type is a valid recurring type
func isRecurringType(s string) bool {
	return s == CLASS || s == STUDY || s == SLEEP || s == EXERCISE || s == WORK || s == MEAL
}

// intToDate converts integer date format to a time.Time struct
func intToDate(date int) (time.Time, error) {
	const dateFormat = "2020-01-02"
	t, err := time.Parse(dateFormat, fmt.Sprintf("%04d-%02d-%02d", date/10000, (date/100)%100, date%100))
	return t, err
}

// dateToInt converts a time.Time struct to an integer date format
func dateToInt(date time.Time) int {
	return (date.Year() * 10000) + (int(date.Month()) * 100) + date.Day()
}

// datesOverlap determines if two dates with given durations overlap
func datesOverlap(date1 time.Time, duration1 int, date2 time.Time, duration2 int) bool {
	// Difference between start times in hours
	timeDelta := math.Abs(float64(date1.Unix() - date2.Unix()))
	var earlierDuration int
	if date1.Before(date2) {
		earlierDuration = duration1
	} else {
		earlierDuration = duration2
	}
	return timeDelta < float64(earlierDuration)
}

//!--
