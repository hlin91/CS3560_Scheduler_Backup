// Package model provides functionality for creating and managing a schedule of tasks
// json_containers.go provides structs that define the json marshaling formats for tasks
package model

// taskContainer is a container for the fields of Tasks and AntiTasks
// This class is not needed but provided for the sake of consistency due to the unfortunate
// necessity of recurContainer
type taskContainer struct {
	Name      string
	Type      string
	Date      int
	StartTime float32
	Duration  float32
}

// recurContainer is a container for the fields of RecurringTask for the sole purpose of complying
// with the convention of importing/exporting the "Date" field as "StartDate" for recurring tasks
// We cannot simply embed taskContainer due to the discrepency in "Date" and "StartDate" naming conventions
type recurContainer struct {
	Name      string
	Type      string
	StartDate int
	StartTime float32
	Duration  float32
	EndDate   int
	Frequency int
}

// taskToContainer populates a taskContainer with the fields of a Task
func taskToContainer(t Task) taskContainer {
	return taskContainer{
		Name:      t.Name,
		Type:      t.Type,
		Date:      t.Date,
		StartTime: t.StartTime,
		Duration:  t.Duration,
	}
}

// recurToContainer populates a recurContainer with the fields of a RecurringTask
func recurToContainer(r RecurringTask) recurContainer {
	return recurContainer{
		Name:      r.Name,
		Type:      r.Type,
		StartDate: r.Date,
		StartTime: r.StartTime,
		Duration:  r.Duration,
		EndDate:   r.EndDate,
		Frequency: r.Frequency,
	}
}

//!--
