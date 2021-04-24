// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
)

type Schedule struct {
	TransientTasks map[string]Task
	AntiTasks      map[string]AntiTask
	RecurringTasks map[string]RecurringTask
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
	if _, ok := s.RecurringTasks[name]; ok {
		return true
	}
	return false
}

// hasAddConflict checks if a task will produce scheduling conflicts if added
func (s Schedule) hasAddConflict(task Schedulable) bool {
	// TODO: Implement this
	return true
}

// hasDeleteConflict checks if a task will produce a scheduling conflict if deleted
func (s Schedule) hasDeleteConflict(task Schedulable) bool {
	// TODO: Implement this
	return true
}

// AddTransientTask creates and adds a transient task to the schedule
func (s *Schedule) AddTransientTask(name, taskType string, date int, startTime, duration float32) error {
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddTransientTask: task name already exists")
	}
	if !isTransientType(taskType) {
		return fmt.Errorf("AddTransientTask: %q is not a transient type", taskType)
	}
	t, err := NewTask(name, taskType, date, startTime, duration)
	if err != nil {
		return fmt.Errorf("AddTransientTask: error creating task: %v", err)
	}
	if s.hasAddConflict(t) {
		return fmt.Errorf("AddTransientTask: task creates scheduling conflict")
	}
	s.TransientTasks[name] = t
	return nil
}

// AddAntiTask creates and adds an anti task to the schedule
func (s *Schedule) AddAntiTask(name, taskType string, date int, startTime, duration float32) error {
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddAntiTask: task name already exists")
	}
	if !isAntiType(taskType) {
		return fmt.Errorf("AddAntiTask: %q is not a transient type", taskType)
	}
	t, err := NewAntiTask(name, taskType, date, startTime, duration)
	if err != nil {
		return fmt.Errorf("AddAntiTask: error creating task: %v", err)
	}
	s.AntiTasks[name] = t
	return nil
}

// AddRecurringTask creates and adds a recurring task to the schedule
func (s *Schedule) AddRecurringTask(name, taskType string, date int, startTime, duration float32, endDate, frequency int) error {
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddRecurringTask: task name already exists")
	}
	if !isRecurringType(taskType) {
		return fmt.Errorf("AddRecurringTask: %q is not a recurring type", taskType)
	}
	t, err := NewRecurringTask(name, taskType, date, startTime, duration, endDate, frequency)
	if err != nil {
		return fmt.Errorf("AddRecurringTask: error creating task: %v", err)
	}
	s.RecurringTasks[name] = t
	return nil
}

// GetTask gets a task in the schedule by name
func (s Schedule) GetTask(name string) (Schedulable, bool) {
	if t, ok := s.TransientTasks[name]; ok {
		return t, ok
	}
	if t, ok := s.AntiTasks[name]; ok {
		return t, ok
	}
	if t, ok := s.RecurringTasks[name]; ok {
		return t, ok
	}
	return nil, false
}

//!--
