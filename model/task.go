// Package model provides functionality for creating and managing a schedule of tasks
// task.go provides an implementation for transient tasks
package model

import (
	"fmt"
	"math"
	"time"
)

// Task is the base class and the transient task
type Task struct {
	Name      string
	Type      string
	Date      int
	StartTime float32
	Duration  float32
}

func NewTask(name, taskType string, date int, startTime, duration float32) (Task, error) {
	var result Task
	if startTime < 0 || startTime > 23.75 {
		return result, fmt.Errorf("bad start time")
	}
	if duration < 0 || duration > 23.75 {
		return result, fmt.Errorf("bad duration")
	}
	if _, err := intToDate(date); err != nil {
		return result, fmt.Errorf("bad date")
	}
	result.Name = name
	result.Date = date
	result.StartTime = startTime
	// Round duration to nearest .25
	duration = float32(math.Round(float64(duration)/.25) * .25)
	result.Duration = duration
	result.Type = taskType
	return result, nil
}

func (t Task) String() string {
	return fmt.Sprintf("--------------------------------\nName: %v\nType: %v\nStart Date: %v\nStart Time: %v\nDuration: %v\n--------------------------------", t.Name, t.Type, dateIntToString(t.Date), t.StartTime, t.Duration)
}

func (t Task) GetStartYear() int {
	return t.Date / 10000
}

func (t Task) GetStartMonth() int {
	return (t.Date / 100) % 100
}

func (t Task) GetStartDay() int {
	return t.Date % 100
}

// GetStartDate gets the start date of the task as a Time struct
func (t Task) GetStartDate() (time.Time, error) {
	date, err := intToDate(t.Date)
	if err != nil {
		return date, fmt.Errorf("GetStartDate: %v", err)
	}
	// Account for the start time
	date = date.Add((time.Hour * time.Duration(t.StartTime)) +
		(time.Minute * 60 * time.Duration(float64(t.StartTime)-math.Floor(float64(t.StartTime)))))
	return date, nil
}

// GetStartDateWithouttime gets the start date as a Time struct without accounting for start time
func (t Task) GetStartDateWithoutTime() (time.Time, error) {
	return intToDate(t.Date)
}

// Overlaps returns true if this task overlaps with the duration of another task
func (t Task) Overlaps(op Task) bool {
	time1, _ := t.GetStartDate()
	time2, _ := op.GetStartDate()
	// Difference in start date in hours
	timeDelta := math.Abs(float64(time1.Unix()-time2.Unix())) / (60 * 60)
	var earlierTask Task
	if time1.Before(time2) {
		earlierTask = t
	} else {
		earlierTask = op
	}
	return timeDelta < float64(earlierTask.Duration)
}

// Before returns true if this task occurs strictly before another task
func (t Task) Before(op Task) bool {
	if t.Date == op.Date {
		return t.StartTime < op.StartTime
	}
	return t.Date < op.Date
}

//!--
