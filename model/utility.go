// Package model provides functionality for creating and managing a schedule of tasks
// utility.go provides various helpful functions private to the package
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
	CANCEL = "Cancellation"
	// Recurring Types
	CLASS    = "Class"
	STUDY    = "Study"
	SLEEP    = "Sleep"
	EXERCISE = "Exercise"
	WORK     = "Work"
	MEAL     = "Meal"
	// Number of keys in transient/anti tasks
	NUM_TASK_KEYS  = 5
	NUM_RECUR_KEYS = 7
	// Key names for JSON marshaling
	NAME_KEY       = "Name"
	TYPE_KEY       = "Type"
	DATE_KEY       = "StartDate"
	START_TIME_KEY = "StartTime"
	DURATION_KEY   = "Duration"
	END_DATE_KEY   = "EndDate"
	FREQUENCY_KEY  = "Frequency"
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

// indexSubtasks adds index numbers to the names of a list of subtasks
func indexSubtasks(tasks []Task) {
	for i := 0; i < len(tasks); i++ {
		tasks[i].Name = fmt.Sprintf("%s (%d)", tasks[i].Name, i+1)
	}
}

// taskKeysPresent checks if the necessary keys for a task are present in a generic string map
func taskKeysPresent(m map[string]interface{}) bool {
	if _, ok := m[NAME_KEY]; !ok {
		return false
	}
	if _, ok := m[TYPE_KEY]; !ok {
		return false
	}
	if _, ok := m[DATE_KEY]; !ok {
		return false
	}
	if _, ok := m[START_TIME_KEY]; !ok {
		return false
	}
	if _, ok := m[DURATION_KEY]; !ok {
		return false
	}
	return true
}

// recurKeysPresent checks if the necessary keys for a recurring task are present in a generic map
func recurKeysPresent(m map[string]interface{}) bool {
	if !taskKeysPresent(m) {
		return false
	}
	if _, ok := m[END_DATE_KEY]; !ok {
		return false
	}
	if _, ok := m[FREQUENCY_KEY]; !ok {
		return false
	}
	return true
}

// mapToTaskInfo extracts task information from a generic map
func mapToTaskInfo(m map[string]interface{}) (string, string, int, float32, float32, error) {
	name, ok := m[NAME_KEY].(string)
	if !ok {
		return "", "", 0, 0, 0, fmt.Errorf("bad name value")
	}
	taskType, ok := m[TYPE_KEY].(string)
	if !ok {
		return "", "", 0, 0, 0, fmt.Errorf("bad type value")
	}
	date, ok := m[DATE_KEY].(float64)
	if !ok {
		return "", "", 0, 0, 0, fmt.Errorf("bad date value")
	}
	startTime, ok := m[START_TIME_KEY].(float64)
	if !ok {
		return "", "", 0, 0, 0, fmt.Errorf("bad start time value")
	}
	duration, ok := m[DURATION_KEY].(float64)
	if !ok {
		return "", "", 0, 0, 0, fmt.Errorf("bad duration value")
	}
	return name, taskType, int(date), float32(startTime), float32(duration), nil
}

// mapToRecurInfo extracts recurring task information from a generic map
func mapToRecurInfo(m map[string]interface{}) (string, string, int, float32, float32, int, int, error) {
	name, taskType, date, startTime, duration, err := mapToTaskInfo(m)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, err
	}
	endDate, ok := m[END_DATE_KEY].(float64)
	if !ok {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad end date value")
	}
	frequency, ok := m[FREQUENCY_KEY].(float64)
	if !ok {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad frequency value")
	}
	return name, taskType, date, startTime, duration, int(endDate), int(frequency), nil
}

//!--
