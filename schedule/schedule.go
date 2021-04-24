// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

type Schedule struct {
	TransientTasks map[string]Task
	AntiTasks      map[string]AntiTask
	RecurringTask  map[string]RecurringTask
}

func NewSchedule() Schedule {
	return Schedule{}
}

// hasNameConflict checks if a task of the same name already exists in the schedule
func (s Schedule) hasNameConflict(name string) bool {
	if _, ok := s.TransientTasks[name]; ok {
		return true
	}
	if _, ok := s.AntiTasks[name]; ok {
		return true
	}
	if _, ok := s.RecurringTask[name]; ok {
		return true
	}
	return false
}

//!--
