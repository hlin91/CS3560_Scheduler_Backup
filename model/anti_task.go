// Package model provides functionality for creating and managing a schedule of tasks
// anti_task.go provides an implementation for anti tasks
package model

import "math"

// AntiTask implements an anti task in the schedule
type AntiTask struct {
	Task
}

func NewAntiTask(name, taskType string, date int, startTime, duration float32) (AntiTask, error) {
	t, err := NewTask(name, taskType, date, startTime, duration)
	if err != nil {
		return AntiTask{}, err
	}
	return AntiTask{t}, nil
}

// Cancels determines if this anti task cancels out another task
func (a AntiTask) Cancels(t Task) bool {
	date1, _ := a.GetStartDate()
	date2, _ := t.GetStartDate()
	if date2.Before(date1) {
		// Cannot cancel
		return false
	}
	// Difference in start time in hours
	timeDelta := float32(math.Abs(float64(date1.Unix()-date2.Unix())) / (60 * 60))
	return a.Duration >= timeDelta+t.Duration
}

// GetCancelledSubtask returns the subtask this anti task cancels and a bool to indicate
// if such a task was found
func (a AntiTask) GetCancelledSubtask(r RecurringTask) (Task, bool) {
	aStart, _ := a.GetStartDate()
	rStart, _ := r.GetStartDate()
	rEnd, _ := r.GetEndDate()
	if aStart.Before(rStart) || rEnd.Before(aStart) {
		// This anti task is outside of the recurring range
		return Task{}, false
	}
	if a.StartTime != r.StartTime || a.Duration != r.Duration {
		// Start time or duration does not match up
		return Task{}, false
	}
	// Determine if the anti task lines up with the recurrence cycle
	if dayDelta := int64(aStart.Sub(rStart).Hours()); dayDelta%int64(r.Frequency) == 0 {
		t, _ := NewTask(r.Name, r.Type, a.Date, a.StartTime, a.Duration)
		return t, true
	}
	return Task{}, false
}

//!--
