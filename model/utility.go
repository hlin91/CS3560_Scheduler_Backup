// Package model provides functionality for creating and managing a schedule of tasks
package model

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
	year := date / 10000
	month := (date / 100) % 100
	day := date % 100
	t := time.Date(date/10000, time.Month((date/100)%100), date%100, 0, 0, 0, 0, time.UTC)
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return t, fmt.Errorf("bad date")
	}
	return t, nil
}

// dateToInt converts a time.Time struct to an integer date format
func dateToInt(date time.Time) int {
	return (date.Year() * 10000) + (int(date.Month()) * 100) + date.Day()
}

// dateIntToString converts a integer date format to a more readable string
func dateIntToString(date int) string {
	return fmt.Sprintf("%04d-%02d-%02d", date/10000, (date/100)%100, date%100)
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

// indexNames adds index numbers to the names of a list of subtasks
func indexSubtasks(tasks []Task) {
	for i, t := range tasks {
		t.Name = fmt.Sprintf("%s %d", t.Name, i)
	}
}

//!--
