// Package controller provides functions to edit the schedule or view (ie. command line)
package controller

import "github.com/hlin91/CS3560_Scheduler_Backup/model"

// ScheduleMenuItem defines a menu option that manages the schedule
type ScheduleMenuItem struct {
	title    string
	schedule *model.Schedule
	hook     func(*model.Schedule) error
}

// Title returns the title of the menu item
func (s ScheduleMenuItem) Title() string {
	return s.title
}

// Exec runs the hook attached to the menu item
func (s *ScheduleMenuItem) Exec() error {
	return s.hook(s.schedule)
}

// NewScheduleMenuItem creates and returns a new schedule menu item
func NewScheduleMenuItem(title string, schedule *model.Schedule, hook func(*model.Schedule) error) ScheduleMenuItem {
	return ScheduleMenuItem{
		title:    title,
		schedule: schedule,
		hook:     hook,
	}
}
