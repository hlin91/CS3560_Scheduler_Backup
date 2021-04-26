// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
)

// Menuer is the minimum interface required for a menu option
type Menuer interface {
	Title() string
	Exec() error
}

type Menu struct {
	Options []Menuer
}

func NewMenu(options []Menuer) Menu {
	return Menu{options}
}

func (m Menu) Display() {
	for i, o := range m.Options {
		fmt.Printf("%d. %s\n", i+1, o.Title())
	}
	fmt.Print("\nEnter an option: ")
}

func (m Menu) Process(input int) error {
	if input < 1 || input > len(m.Options) {
		return fmt.Errorf("input out of range")
	}
	return m.Options[input-1].Exec()
}

type ScheduleMenuItem struct {
	title    string
	schedule *Schedule
	hook     func(*Schedule) error
}

func (s ScheduleMenuItem) Title() string {
	return s.title
}

func (s *ScheduleMenuItem) Exec() error {
	return s.hook(s.schedule)
}

func NewScheduleMenuItem(title string, schedule *Schedule, hook func(*Schedule) error) ScheduleMenuItem {
	return ScheduleMenuItem{
		title:    title,
		schedule: schedule,
		hook:     hook,
	}
}
