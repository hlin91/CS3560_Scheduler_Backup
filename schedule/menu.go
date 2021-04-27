// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// Menuer is the minimum interface required for a menu option
type Menuer interface {
	Title() string
	Exec() error
}

// Menu is an object that manages a list of menu options
type Menu struct {
	Options []Menuer
}

// NewMenu creates and returns a new menu
func NewMenu(options []Menuer) Menu {
	return Menu{options}
}

// Display prints all the options in the menu and an input prompt
func (m Menu) Display() {
	for i, o := range m.Options {
		fmt.Printf("%d. %s\n", i+1, o.Title())
	}
	fmt.Print("\nEnter an option: ")
}

// Process calls the appropriate menu option for the given input
func (m Menu) Process(input int) error {
	if input < 1 || input > len(m.Options) {
		return fmt.Errorf("input out of range")
	}
	return m.Options[input-1].Exec()
}

// Clear clears the screen
func (m Menu) Clear() {
	switch runtime.GOOS {
	case "linux":
		fallthrough
	case "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// ScheduleMenuItem defines a menu option that manages the schedule
type ScheduleMenuItem struct {
	title    string
	schedule *Schedule
	hook     func(*Schedule) error
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
func NewScheduleMenuItem(title string, schedule *Schedule, hook func(*Schedule) error) ScheduleMenuItem {
	return ScheduleMenuItem{
		title:    title,
		schedule: schedule,
		hook:     hook,
	}
}
