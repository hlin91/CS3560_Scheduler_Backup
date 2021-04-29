// Package model provides functionality for creating and managing a schedule of tasks
// recurring_task.go provides an implementation for recurring tasks
package model

import (
	"fmt"
	"math"
	"time"
)

// RecurringTask implements a recurrint task in the schedule
type RecurringTask struct {
	Task
	EndDate   int
	Frequency int
}

func NewRecurringTask(name, taskType string, date int, startTime, duration float32, endDate, frequency int) (RecurringTask, error) {
	t, err := NewTask(name, taskType, date, startTime, duration)
	if err != nil {
		return RecurringTask{}, err
	}
	start, _ := intToDate(date)
	end, err := intToDate(endDate)
	if err != nil {
		return RecurringTask{}, fmt.Errorf("bad end date")
	}
	if end.Before(start) {
		return RecurringTask{}, fmt.Errorf("end date before start date")
	}
	if frequency < 1 || frequency > 7 {
		return RecurringTask{}, fmt.Errorf("bad frequency")
	}
	// if frequency != 1 && frequency != 7 {
	// 	return RecurringTask{}, fmt.Errorf("bad frequency")
	// }
	result := RecurringTask{
		Task:      t,
		EndDate:   endDate,
		Frequency: frequency,
	}
	return result, nil
}

func (r RecurringTask) String() string {
	return r.Task.String() + fmt.Sprintf("\nEnd Date: %v\nFrequency: %v", dateIntToString(r.EndDate), r.Frequency)
}

func (r RecurringTask) GetEndYear() int {
	return r.EndDate / 10000
}

func (r RecurringTask) GetEndMonth() int {
	return (r.EndDate / 100) % 100
}

func (r RecurringTask) GetEndDay() int {
	return r.EndDate % 100
}

func (r RecurringTask) GetEndDate() (time.Time, error) {
	date, err := intToDate(r.EndDate)
	if err != nil {
		return date, fmt.Errorf("GetEndDate: %v", err)
	}
	// Account for start time
	date = date.Add((time.Hour * time.Duration(r.StartTime)) +
		(time.Minute * 60 * time.Duration(float64(r.StartTime)-math.Floor(float64(r.StartTime)))))
	// Account for duration
	date = date.Add((time.Hour * time.Duration(r.Duration)) +
		(time.Minute * 60 * time.Duration(float64(r.Duration)-math.Floor(float64(r.Duration)))))
	return date, nil
}

// GetSubtasks expands the recurring tasks into a series of subtasks
func (r RecurringTask) GetSubtasks() ([]Task, error) {
	result := []Task{}
	startDate, err := r.GetStartDate()
	if err != nil {
		return result, fmt.Errorf("GetSubtasks: %v", err)
	}
	endDate, err := r.GetEndDate()
	if err != nil {
		return result, fmt.Errorf("GetSubtasks: %v", err)
	}
	for startDate.Before(endDate) {
		t, err := NewTask(r.Name, r.Type, dateToInt(startDate), r.StartTime, r.Duration)
		if err != nil {
			return []Task{}, err
		}
		result = append(result, t)
		startDate = startDate.Add(24 * time.Hour * time.Duration(r.Frequency))
	}
	return result, nil
}

// GetOverlappingSubtasks returns the set of subtasks that overlap a given task
func (r RecurringTask) GetOverlappingSubtasks(task Task) ([]Task, error) {
	result := []Task{}
	rStartDate, err := r.GetStartDateWithoutTime()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	rEndDate, err := r.GetEndDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	taskDate, err := task.GetStartDateWithoutTime()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	if rEndDate.Before(taskDate) {
		return result, nil
	}
	// Difference in start days in days
	startDayDelta := int(taskDate.Unix()-rStartDate.Unix()) / (60 * 60 * 24)
	if startDayDelta < 0 {
		// Only have to check the very first recurring subtask
		t, err := NewTask(r.Name, r.Type, r.Date, r.StartTime, r.Duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
		}
		if t.Overlaps(task) {
			result = append(result, t)
		}
		return result, nil
	}
	// Reintroduce start times into start dates
	rStartDate, _ = r.GetStartDate()
	taskDate, _ = task.GetStartDate()
	if startDayDelta%r.Frequency == 0 {
		// Get the subtask for this day
		t, err := NewTask(r.Name, r.Type, task.Date, r.StartTime, r.Duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: error creating today subtask: %v", err)
		}
		if t.Overlaps(task) {
			result = append(result, t)
		}
	}
	// Have to check the cycle before and the cycle after for potential overlaps
	prevCycleDistance := startDayDelta % r.Frequency
	if startDayDelta == 0 {
		prevCycleDistance = r.Frequency
	}
	nextCycleDistance := r.Frequency - (startDayDelta % r.Frequency)
	yesterday := time.Unix(taskDate.Unix()-((24*60*60)*int64(prevCycleDistance)), 0)
	tomorrow := time.Unix(taskDate.Unix()+((24*60*60)*int64(nextCycleDistance)), 0)
	yt, err := NewTask(r.Name, r.Type, dateToInt(yesterday), r.StartTime, r.Duration)
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: error creating yesterday task: %v", err)
	}
	tt, err := NewTask(r.Name, r.Type, dateToInt(tomorrow), r.StartTime, r.Duration)
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: error creating tomorrow task: %v", err)
	}
	if rStartDate.Before(yesterday) && yt.Overlaps(task) {
		result = append(result, yt)
	}
	if tomorrow.Before(rEndDate) && tt.Overlaps(task) {
		result = append(result, tt)
	}
	indexSubtasks(result) // Add index numbers to subtask names
	return result, nil
}

func (r RecurringTask) Overlaps(task Task) bool {
	l, err := r.GetOverlappingSubtasks(task)
	if err != nil {
		fmt.Printf("Warning: RecurringTask.Overlaps: %v\n", err)
	}
	return len(l) > 0
}

// GetOverlappingSubtasksRecurring returns the set of subtasks that overlap a recurring task
// This can also find the overlap with a non recurring task but is less optimal
func (r RecurringTask) GetOverlappingSubtasksRecurring(task RecurringTask) ([]Task, error) {
	result := []Task{}
	startDate, err := r.GetStartDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasksRecurring: %v", err)
	}
	endDate, err := r.GetEndDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasksRecurring: %v", err)
	}
	for startDate.Before(endDate) {
		t, err := NewTask(r.Name, r.Type, dateToInt(startDate), r.StartTime, r.Duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasksRecurring: %v", err)
		}
		if task.Overlaps(t) {
			result = append(result, t)
		}
		startDate = startDate.Add(24 * time.Hour * time.Duration(r.Frequency))
	}
	indexSubtasks(result) // Add index numbers to subtask names
	return result, nil
}

func (r RecurringTask) OverlapsRecurring(task RecurringTask) bool {
	l, err := r.GetOverlappingSubtasksRecurring(task)
	if err != nil {
		fmt.Printf("Warning: RecurringTask.OverlapsRecurring: %v\n", err)
	}
	return len(l) > 0
}

//!--
