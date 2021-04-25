// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
	"math"
	"time"
)

// Task is the base class and the transient task
type Task struct {
	name      string
	taskType  string
	date      int
	startTime float32
	duration  float32
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Type() string {
	return t.taskType
}

func (t Task) Date() int {
	return t.date
}

func (t Task) StartTime() float32 {
	return t.startTime
}

func (t Task) Duration() float32 {
	return t.duration
}

func NewTask(name, taskType string, date int, startTime, duration float32) (Task, error) {
	var result Task
	if startTime < 0 || startTime > 23.75 {
		return result, fmt.Errorf("NewTask: bad start time")
	}
	if duration < 0 || duration > 23.75 {
		return result, fmt.Errorf("NewTask: bad duration")
	}
	if _, err := intToDate(date); err != nil {
		return result, fmt.Errorf("NewTask: bad date")
	}
	result.name = name
	result.date = date
	result.startTime = startTime
	// Round duration to nearest .25
	duration = float32(math.Round(float64(duration)/.25) * .25)
	result.duration = duration
	result.taskType = taskType
	return result, nil
}

func (t Task) String() string {
	return fmt.Sprintf("Name: %v\nType: %v\nStart Date: %v\nStart Time: %v\nDuration: %v", t.name, t.taskType, t.date, t.startTime, t.duration)
}

func (t Task) GetStartYear() int {
	return t.date / 10000
}

func (t Task) GetStartMonth() int {
	return (t.date / 100) % 100
}

func (t Task) GetStartDay() int {
	return t.date % 100
}

// GetStartDate gets the start date of the task as a Time struct
func (t Task) GetStartDate() (time.Time, error) {
	date, err := intToDate(t.date)
	if err != nil {
		return date, fmt.Errorf("GetStartDate: %v", err)
	}
	// Account for the start time
	date = date.Add((time.Hour * time.Duration(t.startTime)) +
		(time.Minute * 60 * time.Duration(float64(t.startTime)-math.Floor(float64(t.startTime)))))
	return date, nil
}

// GetStartDateWithouttime gets the start date as a Time struct without accounting for start time
func (t Task) GetStartDateWithoutTime() (time.Time, error) {
	return intToDate(t.date)
}

// Overlaps returns true if this task overlaps with the duration of another task
func (t Task) Overlaps(op Task) bool {
	time1, _ := t.GetStartDate()
	time2, _ := op.GetStartDate()
	timeDelta := math.Abs(float64(time1.Unix() - time2.Unix()))
	var earlierTask Task
	if t.Before(op) {
		earlierTask = t
	} else {
		earlierTask = op
	}
	return timeDelta < float64(earlierTask.Duration())
}

// Before returns true if this task occurs strictly before another task
func (t Task) Before(op Task) bool {
	if t.Date() == op.Date() {
		return t.StartTime() < op.StartTime()
	}
	return t.Date() < op.Date()
}

// RecurringTask implements a recurrint task in the schedule
type RecurringTask struct {
	Task
	endDate   int
	frequency int
}

func NewRecurringTask(name, taskType string, date int, startTime, duration float32, endDate, frequency int) (RecurringTask, error) {
	t, err := NewTask(name, taskType, date, startTime, duration)
	if err != nil {
		return RecurringTask{}, fmt.Errorf("NewRecurringTask: %v", err)
	}
	start, _ := intToDate(date)
	end, err := intToDate(endDate)
	if err != nil {
		return RecurringTask{}, fmt.Errorf("NewRecurringTask: bad end date")
	}
	if end.Before(start) {
		return RecurringTask{}, fmt.Errorf("NewRecurringTask: end date before start date")
	}
	if frequency < 1 || frequency > 7 {
		return RecurringTask{}, fmt.Errorf("NewRecurringTask: bad frequency")
	}
	result := RecurringTask{
		Task:      t,
		endDate:   endDate,
		frequency: frequency,
	}
	return result, nil
}

func (r RecurringTask) EndDate() int {
	return r.endDate
}

func (r RecurringTask) Frequency() int {
	return r.frequency
}

func (r RecurringTask) String() string {
	return r.Task.String() + fmt.Sprintf("\nEnd Date: %v\nFrequency: %v", r.endDate, r.frequency)
}

func (r RecurringTask) GetEndYear() int {
	return r.endDate / 10000
}

func (r RecurringTask) GetEndMonth() int {
	return (r.endDate / 100) % 100
}

func (r RecurringTask) GetEndDay() int {
	return r.endDate % 100
}

func (r RecurringTask) GetEndDate() (time.Time, error) {
	date, err := intToDate(r.endDate)
	if err != nil {
		return date, fmt.Errorf("GetEndDate: %v", err)
	}
	// Account for start time
	date = date.Add((time.Hour * time.Duration(r.startTime)) +
		(time.Minute * 60 * time.Duration(float64(r.startTime)-math.Floor(float64(r.startTime)))))
	// Account for duration
	date = date.Add((time.Hour * time.Duration(r.duration)) +
		(time.Minute * 60 * time.Duration(float64(r.duration)-math.Floor(float64(r.duration)))))
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
		t, err := NewTask(r.name, r.taskType, dateToInt(startDate), r.startTime, r.duration)
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
	if _, ok := task.(RecurringTask); ok {
		// Recurring tasks require a less optimal check
		return r.GetOverlappingSubtasksRecurring(task)
	}
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
	if startDayDelta%r.frequency == 0 {
		// Get the subtask for this day
		t, err := NewTask(r.name, r.taskType, task.Date(), r.startTime, r.duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: error creating today subtask: %v", err)
		}
		if t.Overlaps(task) {
			result = append(result, t)
		}
	}
	if r.frequency == 1 {
		// Have to check the day before and the day after
		yesterday := time.Unix(taskDate.Unix()-(24*60*60), 0)
		tomorrow := time.Unix(taskDate.Unix()+(24*60*60), 0)
		yt, err := NewTask(r.name, r.taskType, dateToInt(yesterday), r.startTime, r.duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasks: error creating yesterday task: %v", err)
		}
		tt, err := NewTask(r.name, r.taskType, dateToInt(tomorrow), r.startTime, r.duration)
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
	return result, nil
}

func (r RecurringTask) Overlaps(task Task) bool {
	if recur, ok := task.(RecurringTask); ok {
		return r.OverlapsRecurring(recur)
	}
	l, err := r.GetOverlappingSubtasks(task)
	if err != nil {
		fmt.Printf("Warning: RecurringTask.Overlaps: %v\n", err)
	}
	return len(l) > 0
}

// GetOverlappingSubtasksRecurring returns the set of subtasks that overlap a recurring task
// This can also find the overlap with a non recurring task but is less optimal
func (r RecurringTask) GetOverlappingSubtasksRecurring(task Task) ([]Task, error) {
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
		t, err := NewTask(r.name, r.taskType, dateToInt(startDate), r.startTime, r.duration)
		if err != nil {
			return result, fmt.Errorf("GetOverlappingSubtasksRecurring: %v", err)
		}
		if task.Overlaps(t) {
			result = append(result, t)
		}
		startDate = startDate.Add(24 * time.Hour * time.Duration(r.frequency))
	}
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
		return AntiTask{}, fmt.Errorf("NewAntiTask: %v", err)
	}
	return AntiTask{t}, nil
}

// Cancels determines if this anti task cancels out another task
func (a AntiTask) Cancels(op Task) bool {
	return a.Date() == op.Date() && a.StartTime() == op.StartTime() && a.Duration() == op.Duration()
}

// GetCancelledSubtask returns the subtask this anti task cancels and a bool to indicate
// if such a task was found
func (a AntiTask) GetCancelledSubtask(r RecurringTask) (Task, bool) {
	aStart, _ := a.GetStartDate()
	rStart, _ := r.GetStartDate()
	rEnd, _ := r.GetEndDate()
	if aStart.Before(rStart) || rEnd.Before(aStart) {
		// This anti task is outside of the recurring range
		return nil, false
	}
	if a.startTime != r.StartTime() || a.duration != r.Duration() {
		// Start time or duration does not match up
		return nil, false
	}
	// Determine if the anti task lines up with the recurrence cycle
	if dayDelta := int64(aStart.Sub(rStart).Hours()); dayDelta%int64(r.Frequency()) == 0 {
		t, _ := NewTask(r.Name(), r.Type(), a.date, a.startTime, a.duration)
		return t, true
	}
	return nil, false
}

//!--
