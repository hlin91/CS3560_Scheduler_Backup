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

// hasAnti checks if an anti task that cancels the specified task exists in the schedule
func (s Schedule) hasAnti(task Schedulable) bool {
	for _, anti := range s.AntiTasks {
		if anti.Cancels(task) {
			return true
		}
	}
	return false
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
	// Check against all transient tasks
	for n, t := range s.TransientTasks {
		if n == task.Name() {
			// Don't check against itself
			continue
		}
		if task.Overlaps(t) {
			return true
		}
	}
	// Check against all recurring tasks
	for n, t := range s.RecurringTasks {
		if n == task.Name() {
			continue
		}
		overlaps, _ := t.GetOverlappingSubtasks(task)
		for _, o := range overlaps {
			if !s.hasAnti(o) {
				return true
			}
		}
	}
	return false
}

// hasDeleteConflict checks if a task will produce a scheduling conflict if deleted
func (s Schedule) hasDeleteConflict(task Schedulable) bool {
	// Only have to check deletion conflicts if task is an anti task
	a, ok := task.(AntiSchedulable)
	if !ok {
		return false
	}
	for _, t := range s.RecurringTasks {
		// For every cancelled subtask, check if there is an overlap in the schedule with that
		// subtask
		if cancelled, ok := a.GetCancelledSubtask(t.(RecurringSchedulable)); ok {
			if s.hasAddConflict(cancelled) {
				return true
			}
		}
	}
	return false
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
	var cancelledExists bool
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddAntiTask: task name already exists")
	}
	if !isAntiType(taskType) {
		return fmt.Errorf("AddAntiTask: %q is not a transient type", taskType)
	}
	a, err := NewAntiTask(name, taskType, date, startTime, duration)
	if err != nil {
		return fmt.Errorf("AddAntiTask: error creating task: %v", err)
	}
	for _, t := range s.RecurringTasks {
		if _, ok := a.GetCancelledSubtask(t.(RecurringSchedulable)); ok {
			cancelledExists = true
			break
		}
	}
	if !cancelledExists {
		return fmt.Errorf("AddAntiTask: no corresponding recurring task exists")
	}
	s.AntiTasks[name] = a
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
func (s Schedule) GetTask(name string) (Schedulable, error) {
	if t, ok := s.TransientTasks[name]; ok {
		return t, nil
	}
	if t, ok := s.AntiTasks[name]; ok {
		return t, nil
	}
	if t, ok := s.RecurringTasks[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("GetTask: task name does not exist in schedule")
}

// DeleteTask deletes a task in the schedule by name
func (s *Schedule) DeleteTask(name string) error {
	if _, ok := s.TransientTasks[name]; ok {
		delete(s.TransientTasks, name)
		return nil
	}
	if _, ok := s.RecurringTasks[name]; ok {
		delete(s.RecurringTasks, name)
		return nil
	}
	if a, ok := s.AntiTasks[name]; ok {
		if s.hasDeleteConflict(a) {
			return fmt.Errorf("DeleteTask: deletion creates a schedule conflict")
		}
		delete(s.AntiTasks, name)
		return nil
	}
	return fmt.Errorf("DeleteTask: task name does not exist in schedule")
}

// TODO: Implement EditTask

//!--
