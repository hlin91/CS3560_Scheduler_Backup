// Package model provides functionality for creating and managing a schedule of tasks
// schedule.go provides an implementation for a schedule to store and manage tracks
package model

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Schedule struct {
	TransientTasks map[string]Task
	AntiTasks      map[string]AntiTask
	RecurringTasks map[string]RecurringTask
}

// NewSchedule creates and returns a schedule
func NewSchedule() *Schedule {
	return &Schedule{
		TransientTasks: map[string]Task{},
		AntiTasks:      map[string]AntiTask{},
		RecurringTasks: map[string]RecurringTask{},
	}
}

// AddTransientTask creates and adds a transient task to the schedule
func (s *Schedule) AddTransientTask(name, taskType string, date int, startTime, duration float32) error {
	if len(name) == 0 {
		return fmt.Errorf("AddTransientTask: name cannot be empty")
	}
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

// AddSubtask creates and adds a recurring subtask to the schedule
func (s *Schedule) AddSubtask(name, taskType string, date int, startTime, duration float32) error {
	if len(name) == 0 {
		return fmt.Errorf("AddSubtask: name cannot be empty")
	}
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddSubtask: task name already exists")
	}
	if !isRecurringType(taskType) {
		return fmt.Errorf("AddSubtask: %q is not a recurring type", taskType)
	}
	t, err := NewTask(name, taskType, date, startTime, duration)
	if err != nil {
		return fmt.Errorf("AddSubtask: error creating task: %v", err)
	}
	if s.hasAddConflict(t) {
		return fmt.Errorf("AddSubtask: task creates scheduling conflict")
	}
	s.TransientTasks[name] = t
	return nil
}

// AddAntiTask creates and adds an anti task to the schedule
func (s *Schedule) AddAntiTask(name, taskType string, date int, startTime, duration float32) error {
	var cancelledExists bool
	if len(name) == 0 {
		return fmt.Errorf("AddAntiTask: name cannot be empty")
	}
	if s.hasNameConflict(name) {
		return fmt.Errorf("AddAntiTask: task name already exists")
	}
	if !isAntiType(taskType) {
		return fmt.Errorf("AddAntiTask: %q is not an anti type", taskType)
	}
	a, err := NewAntiTask(name, taskType, date, startTime, duration)
	if err != nil {
		return fmt.Errorf("AddAntiTask: error creating task: %v", err)
	}
	for _, t := range s.AntiTasks {
		if t.Overlaps(a.Task) {
			return fmt.Errorf("AddAntiTask: task overlaps with another anti task")
		}
	}
	for _, t := range s.RecurringTasks {
		if _, ok := a.GetCancelledSubtask(t); ok {
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
	if len(name) == 0 {
		return fmt.Errorf("AddRecurringTask: name cannot be empty")
	}
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
	if s.hasAddConflictRecurring(t) {
		return fmt.Errorf("AddRecurringTask: task creates scheduling conflict")
	}
	s.RecurringTasks[name] = t
	return nil
}

// DeleteTask deletes a task in the schedule by name
func (s *Schedule) DeleteTask(name string) error {
	if _, ok := s.TransientTasks[name]; ok {
		delete(s.TransientTasks, name)
		return nil
	}
	if r, ok := s.RecurringTasks[name]; ok {
		delete(s.RecurringTasks, name)
		// Delete all corresponding anti tasks for the recurring task
		for _, a := range s.AntiTasks {
			if _, ok := a.GetCancelledSubtask(r); ok {
				delete(s.AntiTasks, a.Name)
			}
		}
		return nil
	}
	if a, ok := s.AntiTasks[name]; ok {
		// For anti tasks, we have to check if deleting will not create a conflict
		if s.hasDeleteConflict(a.Task) {
			return fmt.Errorf("DeleteTask: deletion creates a schedule conflict")
		}
		delete(s.AntiTasks, name)
		return nil
	}
	return fmt.Errorf("DeleteTask: task name does not exist in schedule")
}

// EditTransienTask edits the details of an existing transient task in the schedule
func (s *Schedule) EditTransientTask(taskName, newName, newType string, newDate int, newStartTime, newDuration float32) error {
	t, ok := s.TransientTasks[taskName]
	if !ok {
		return fmt.Errorf("EditTransientTask: task name does not exist in schedule")
	}
	if newName != taskName && s.hasNameConflict(newName) {
		return fmt.Errorf("EditTransientTask: new name already exists in schedule")
	}
	newTask, err := NewTask(newName, newType, newDate, newStartTime, newDuration)
	if err != nil {
		return fmt.Errorf("EditTransientTask: %v", err)
	}
	if t.Date == newDate && t.StartTime == newStartTime && t.Duration == newDuration {
		// Only the name changed and type changed
		delete(s.TransientTasks, taskName)
		s.TransientTasks[newName] = newTask
		return nil
	}
	s.DeleteTask(taskName)
	if s.hasAddConflict(newTask) {
		// Add back the old task
		s.TransientTasks[taskName] = t
		return fmt.Errorf("EditTransientTask: new details create a schedule conflict")
	}
	s.TransientTasks[newName] = newTask
	return nil
}

// EditAntiTask edits the details of an existing anti task in the schedule
func (s *Schedule) EditAntiTask(taskName, newName string, newDate int, newStartTime, newDuration float32) error {
	a, ok := s.AntiTasks[taskName]
	if !ok {
		return fmt.Errorf("EditAntiTask: task name does not exist in schedule")
	}
	if newName != taskName && s.hasNameConflict(newName) {
		return fmt.Errorf("EditAntiTask: new name already exists in schedule")
	}
	newTask, err := NewAntiTask(newName, a.Type, newDate, newStartTime, newDuration)
	if err != nil {
		return fmt.Errorf("EditAntiTask: %v", err)
	}
	if a.Date == newDate && a.StartTime == newStartTime && a.Duration == newDuration {
		// Only name changed
		delete(s.AntiTasks, taskName)
		s.AntiTasks[newName] = newTask
		return nil
	}
	if err := s.DeleteTask(taskName); err != nil {
		return fmt.Errorf("EditAntiTask: %v", err)
	}
	// Find a corresponding recurring task
	var foundCancelledTask bool
	for _, r := range s.RecurringTasks {
		if _, ok := newTask.GetCancelledSubtask(r); ok {
			foundCancelledTask = true
			break
		}
	}
	if !foundCancelledTask {
		s.AntiTasks[taskName] = a
		return fmt.Errorf("EditAntiTask: new anti task does not correspond with any recurring task")
	}
	s.AntiTasks[newName] = newTask
	return nil
}

// EditRecurringTask edits the details of an existing recurring task in the schedule
func (s *Schedule) EditRecurringTask(taskName, newName, newType string, newDate int, newStartTime, newDuration float32, newEndDate, newFrequency int) error {
	r, ok := s.RecurringTasks[taskName]
	if !ok {
		return fmt.Errorf("EditRecurringTask: task name does not exist in schedule")
	}
	if newName != taskName && s.hasNameConflict(newName) {
		return fmt.Errorf("EditRecurringTask: new name already exists in schedule")
	}
	newTask, err := NewRecurringTask(newName, newType, newDate, newStartTime, newDuration, newEndDate, newFrequency)
	if err != nil {
		return fmt.Errorf("EditRecurringTask: %v", err)
	}
	if r.Date == newDate && r.StartTime == newStartTime && r.Duration == newDuration && r.EndDate == newEndDate && r.Frequency == newFrequency {
		// Only name changed and type changed
		delete(s.RecurringTasks, taskName)
		s.RecurringTasks[newName] = newTask
		return nil
	}
	delete(s.RecurringTasks, taskName)
	if s.hasAddConflictRecurring(newTask) {
		// Add back old task
		s.RecurringTasks[taskName] = r
		return fmt.Errorf("EditRecurringTask: new details create a schedule conflict")
	}
	// Delete all anti tasks of the old recurring task that do not match up with the new task
	for _, a := range s.AntiTasks {
		if _, ok := a.GetCancelledSubtask(r); ok {
			if _, ok := a.GetCancelledSubtask(newTask); !ok {
				delete(s.AntiTasks, a.Name)
			}
		}
	}
	return nil
}

// Project specifications are vagues so we'll consider all years in get by date range functions

// GetTasksByMonth gets all tasks/subtasks within a specified month
func (s Schedule) GetTasksByMonth(month int) ([]Task, error) {
	result := []Task{}
	// Get the transient tasks
	for _, t := range s.TransientTasks {
		if t.GetStartMonth() == month {
			result = append(result, t)
		}
	}
	// Get the recurring subtasks
	for _, r := range s.RecurringTasks {
		subtasks, err := r.GetSubtasks()
		if err != nil {
			return []Task{}, fmt.Errorf("GetTasksByMonth: error getting subtasks: %v", err)
		}
		for _, sub := range subtasks {
			if sub.GetStartMonth() == month && !s.hasAnti(sub) {
				result = append(result, sub)
			}
		}
	}
	return result, nil
}

// GetTasksByDay gets all tasks/subtasks starting at a specified month and day
func (s Schedule) GetTasksByDay(month, day int) ([]Task, error) {
	result := []Task{}
	byMonth, err := s.GetTasksByMonth(month)
	if err != nil {
		return result, fmt.Errorf("GetTasksByDay: %v", err)
	}
	for _, t := range byMonth {
		if t.GetStartDay() == day {
			result = append(result, t)
		}
	}
	return result, nil
}

// GetTasksByWeek gets all tasks/subtasks occuring in the week of the specified month and day
func (s Schedule) GetTasksByWeek(month, day int) ([]Task, error) {
	result := []Task{}
	byMonth, err := s.GetTasksByMonth(month)
	if err != nil {
		return result, fmt.Errorf("GetTasksByWeek: %v", err)
	}
	for _, t := range byMonth {
		date, err := t.GetStartDate()
		if err != nil {
			return result, fmt.Errorf("GetTasksByWeek: %v", err)
		}
		targetDate := time.Date(t.GetStartYear(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
		_, week := date.ISOWeek()
		_, targetWeek := targetDate.ISOWeek()
		if week == targetWeek {
			result = append(result, t)
		}
	}
	return result, nil
}

/// Tasks should be added to the schedule in this order
// 1. Recurring tasks
// 2. Anti tasks
// 3. Transient tasks/Subtasks
// This will prevent scheduling conflicts due to insertion order

// LoadFile loads the contents of the json file at the specified path into the schedule
// We expect the json file to contain a single list of tasks
func (s *Schedule) LoadFile(path string) error {
	transientTaskBuff := []map[string]interface{}{} // Buffer to hold transient tasks read from file
	antiTaskBuff := []map[string]interface{}{}      // Buffer to hold anti tasks read from file
	recurTaskBuff := []map[string]interface{}{}     // Buffer to hold recurring tasks read from file
	subTaskBuff := []map[string]interface{}{}       // Buffer to hold subtasks read from file
	content, err := os.ReadFile(path)               // Load contents of file as a byte slice
	if err != nil {
		return fmt.Errorf("LoadFile: error reading file %q: %v", path, err)
	}
	var tasksRead []interface{}
	err = json.Unmarshal(content, &tasksRead)
	if err != nil {
		return fmt.Errorf("LoadFile: error unmarshaling json: %v", err)
	}
	for _, i := range tasksRead {
		// Because the json format is pre-determined, we have to discriminate based on number of keys and type field
		t := i.(map[string]interface{})
		if len(t) != NUM_TASK_KEYS && len(t) != NUM_RECUR_KEYS {
			// Wrong number of keys given (ie. bad data)
			return fmt.Errorf("LoadFile: error parsing tasks: wrong number of keys")
		}
		if len(t) == NUM_TASK_KEYS {
			// Either a recurring subtask, an anti task, or a transient task
			if _, ok := t[TYPE_KEY]; !ok {
				return fmt.Errorf("LoadFile: error parsing tasks: missing %q field", TYPE_KEY)
			}
			taskType, ok := t[TYPE_KEY].(string)
			if !ok {
				return fmt.Errorf("LoadFile: error parsing tasks: could not assert type field to string")
			}
			if isTransientType(taskType) {
				transientTaskBuff = append(transientTaskBuff, t)
				continue
			}
			if isAntiType(taskType) {
				antiTaskBuff = append(antiTaskBuff, t)
				continue
			}
			if isRecurringType(taskType) {
				subTaskBuff = append(subTaskBuff, t)
				continue
			}
			return fmt.Errorf("LoadFile: error parsing tasks: bad type found: %q", taskType)
		}
		if len(t) == NUM_RECUR_KEYS {
			// A potential recurring task
			recurTaskBuff = append(recurTaskBuff, t)
			continue
		}
		return fmt.Errorf("LoadFile: error parsing tasks: bad number of keys found")
	}
	// Add the recurring tasks, then the anti tasks, then the transient and subtasks
	for _, m := range recurTaskBuff {
		if err := recurKeysPresent(m); err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: task values missing: %v", err)
		}
		name, taskType, date, startTime, duration, endDate, frequency, err := mapToRecurInfo(m)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
		err = s.AddRecurringTask(name, taskType, date, startTime, duration, endDate, frequency)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
	}
	for _, m := range antiTaskBuff {
		if err := taskKeysPresent(m); err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: task values missing: %v", err)
		}
		name, taskType, date, startTime, duration, err := mapToTaskInfo(m)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
		err = s.AddAntiTask(name, taskType, date, startTime, duration)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
	}
	for _, m := range transientTaskBuff {
		if err := taskKeysPresent(m); err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: task values missing: %v", err)
		}
		name, taskType, date, startTime, duration, err := mapToTaskInfo(m)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
		err = s.AddTransientTask(name, taskType, date, startTime, duration)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
	}
	for _, m := range subTaskBuff {
		if err := taskKeysPresent(m); err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: task values missing: %v", err)
		}
		name, taskType, date, startTime, duration, err := mapToTaskInfo(m)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
		// Append date to disambiguate subtask name
		name += fmt.Sprintf(" (%4d-%02d-%02d)", date/10000, (date/100)%100, date%100)
		err = s.AddSubtask(name, taskType, date, startTime, duration)
		if err != nil {
			return fmt.Errorf("LoadFile: error loading tasks: %v", err)
		}
	}
	return nil
}

// WriteTasks writes all tasks in the schedule to a specified file in JSON format
func (s Schedule) WriteTasks(path string) error {
	allTasks := []interface{}{}
	// Compile all the tasks
	for _, t := range s.TransientTasks {
		allTasks = append(allTasks, taskToContainer(t))
	}
	for _, t := range s.AntiTasks {
		allTasks = append(allTasks, taskToContainer(t.Task))
	}
	for _, t := range s.RecurringTasks {
		allTasks = append(allTasks, recurToContainer(t))
	}
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("WriteTasks: error creating file: %v", err)
	}
	content, err := json.MarshalIndent(allTasks, "", "\t")
	if err != nil {
		return fmt.Errorf("WriteTasks: error marshaling json: %v", err)
	}
	_, err = outFile.Write(content)
	if err != nil {
		return fmt.Errorf("WriteTasks: error writing to file: %v", err)
	}
	err = outFile.Close()
	if err != nil {
		return fmt.Errorf("WriteTasks: error closing file: %v", err)
	}
	return nil
}

// WriteTaskList writes a list of tasks into a specified file in JSON format
func (s Schedule) WriteTaskList(path string, tasks []Task) error {
	l := []interface{}{}
	for _, t := range tasks {
		l = append(l, taskToContainer(t))
	}
	outFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("WriteTaskList: error creating file: %v", err)
	}
	content, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		return fmt.Errorf("WriteTaskList: error marshaling json: %v", err)
	}
	_, err = outFile.Write(content)
	if err != nil {
		return fmt.Errorf("WriteTaskList: error writing to file: %v", err)
	}
	err = outFile.Close()
	if err != nil {
		return fmt.Errorf("WriteTaskList: error closing file: %v", err)
	}
	return nil
}

// hasAnti checks if an anti task that cancels the specified task exists in the schedule
func (s Schedule) hasAnti(task Task) bool {
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
func (s Schedule) hasAddConflict(task Task) bool {
	// Check against all transient tasks
	for n, t := range s.TransientTasks {
		if n == task.Name {
			// Don't check against itself
			continue
		}
		if task.Overlaps(t) {
			return true
		}
	}
	// Check against all recurring tasks
	for _, t := range s.RecurringTasks {
		overlaps, _ := t.GetOverlappingSubtasks(task)
		for _, o := range overlaps {
			if !s.hasAnti(o) {
				return true
			}
		}
	}
	return false
}

// hasAddConflictRecurring checks if a recurring task will produce scheduling conflicts if added
func (s Schedule) hasAddConflictRecurring(task RecurringTask) bool {
	for _, t := range s.TransientTasks {
		overlaps, _ := task.GetOverlappingSubtasks(t)
		for _, o := range overlaps {
			if !s.hasAnti(o) {
				return true
			}
		}
	}
	// Check against all recurring tasks
	// As per the project specs, anti tasks cannot be applied to overlaps between 2 recurring tasks
	for n, t := range s.RecurringTasks {
		if n == task.Name {
			continue
		}
		if task.OverlapsRecurring(t) {
			return true
		}
	}
	return false
}

// hasDeleteConflict checks if a task will produce a scheduling conflict if deleted
func (s Schedule) hasDeleteConflict(task Task) bool {
	// Only have to check deletion conflicts if task is an anti task
	if !isAntiType(task.Type) {
		return false
	}
	a := AntiTask{task}
	for _, t := range s.RecurringTasks {
		// For every cancelled subtask, check if there is an overlap in the schedule with that
		// subtask
		if cancelled, ok := a.GetCancelledSubtask(t); ok {
			if s.hasAddConflict(cancelled) {
				return true
			}
		}
	}
	return false
}

//!--
