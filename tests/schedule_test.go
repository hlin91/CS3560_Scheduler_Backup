// Package tests contains unit tests
// tests.go contains unit tests for the scheduler based on test scenarios from the professor
package tests

import (
	"testing"

	"github.com/hlin91/CS3560_Scheduler_Backup/model"
)

func TestScenario1(t *testing.T) {
	s := model.NewSchedule()
	if err := s.LoadFile("../data/Set1.json"); err != nil {
		t.Errorf("Failed to load Set1: %v", err)
		return
	}
	if err := s.DeleteTask("Intern Interview"); err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}
	if err := s.AddTransientTask("Intern Interview", "Appointment", 20200427, 17, 2.5); err != nil {
		t.Errorf("Failed to add task: %v", err)
	}
	if err := s.AddTransientTask("Watch a movie", "Movie", 20200429, 21.5, 2); err == nil {
		t.Errorf("Added invalid task with type %q", "Movie")
	}
	if err := s.AddTransientTask("Watch a movie", "Visit", 20200430, 18.5, 2); err == nil {
		t.Errorf("Added conflicting task")
	}
	if err := s.LoadFile("../data/Set2.json"); err == nil {
		t.Errorf("Loaded json file with conflicting tasks")
	}
}

func TestScenario2(t *testing.T) {
	s := model.NewSchedule()
	if err := s.LoadFile("../data/Set2.json"); err != nil {
		t.Errorf("Failed to load Set2: %v", err)
		return
	}
	if err := s.AddAntiTask("Skip-out", "Cancellation", 20200430, 19.25, .75); err == nil {
		t.Errorf("Added invalid anti task")
	}
	if err := s.AddAntiTask("Skip a meal", "Cancellation", 20200428, 17, 1); err != nil {
		t.Errorf("Failed to add valid anti task: %v", err)
	}
	if err := s.LoadFile("../data/Set1.json"); err != nil {
		t.Errorf("Failed to load Set1")
	}
}

func TestReversion(t *testing.T) {
	s := model.NewSchedule()
	s.AddRecurringTask("CS3560-Tu", "Class", 20200414, 19, 1.25, 20200505, 7)
	if err := s.LoadFile("../data/ReversionTest1.json"); err == nil {
		t.Errorf("Loaded ReversionTest1.json file with conflicting task")
	}
	if _, ok := s.TransientTasks["No conflict 1"]; ok {
		t.Errorf("Schedule failed to revert after conflict in ReversionTest1.json")
	}
	if err := s.LoadFile("../data/Set1.json"); err == nil {
		t.Errorf("Loaded Set1.json file with conflicting task")
	}
	if _, ok := s.AntiTasks["Skip For Visit"]; ok {
		t.Errorf("Schedule failed to revert after conflict in Set1.json")
	}
}

func TestAnti(t *testing.T) {
	s := model.NewSchedule()
	s.AddRecurringTask("CS3560-Tu", "Class", 20200414, 19, 1.25, 20200505, 7)
	if err := s.AddAntiTask("Bad Anti Task", "Cancellation", 20200415, 19, 1.25); err == nil {
		t.Errorf("Added bad anti task to schedule")
	}
	if err := s.AddAntiTask("Holiday", "Cancellation", 20200421, 19, 1.25); err != nil {
		t.Errorf("Failed to add good anti task")
	}
	if err := s.AddTransientTask("Pooping", "Appointment", 20200421, 19, 1); err != nil {
		t.Errorf("Anti task failed to allow transient task to be scheduled")
	}
}
