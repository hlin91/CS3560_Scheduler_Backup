// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

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
	return fmt.Sprintf("--------------------------------\nName: %v\nType: %v\nStart Date: %v\nStart Time: %v\nDuration: %v\n--------------------------------", t.Name, t.Type, t.Date, t.StartTime, t.Duration)
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
	if t.Before(op) {
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
	result := RecurringTask{
		Task:      t,
		EndDate:   endDate,
		Frequency: frequency,
	}
	return result, nil
}

func (r RecurringTask) String() string {
	return r.Task.String() + fmt.Sprintf("--------------------------------\nEnd Date: %v\nFrequency: %v\n--------------------------------", r.EndDate, r.Frequency)
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
		startDate = startDate.Add(24 * time.Hour)
	}
	return result, nil
}

// GetOverlappingSubtasks returns the set of subtasks that overlap a given task
func (r RecurringTask) GetOverlappingSubtasks(task Task) ([]Task, error) {
	result := []Task{}
	rStartDate, err := r.GetStartDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	rEndDate, err := r.GetEndDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	taskDate, err := task.GetStartDate()
	if err != nil {
		return result, fmt.Errorf("GetOverlappingSubtasks: %v", err)
	}
	if taskDate.Before(rStartDate) || rEndDate.Before(taskDate) {
		return result, nil
	}
	startDayDelta := int(math.Abs(float64(rStartDate.Unix()-taskDate.Unix())) / (60 * 60 * 24))
	// Check the cycle before, the cycle of, and the cycle after the task for overlaps
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
	if r.Frequency == 1 {
		// Have to check the day before and the day after
		yesterday := time.Unix(taskDate.Unix()-(24*60*60), 0)
		tomorrow := time.Unix(taskDate.Unix()+(24*60*60), 0)
		yt, err := NewTask(r.Name, r.Type, dateToInt(yesterday), r.StartTime, r.Duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: error creating yesterday task: %v", err)
		}
		tt, err := NewTask(r.Name, r.Type, dateToInt(tomorrow), r.StartTime, r.Duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: error creating tomorrow task: %v", err)
		}
		if yt.Overlaps(task) {
			result = append(result, yt)
		}
		if tt.Overlaps(task) {
			result = append(result, tt)
		}
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
