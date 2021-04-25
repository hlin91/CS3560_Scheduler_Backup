// Package schedule provides functionality for creating and managing a schedule of tasks
package schedule

import (
	"fmt"
)

// Menuer is the minimum interface required for a menu option
type menuer interface {
	Title() string
	Exec() error
}

type Menu struct {
	Options []menuer
}

func (m Menu) Display() {
	for i, o := range m.Options {
		fmt.Printf("%d. %s\n", i+1, o.Title())
	}
}

func (m Menu) Process(input int) error {
	if input < 1 || input > len(m.Options) {
		return fmt.Errorf("Process: input out of range")
	}
	return m.Options[input-1].Exec()
}
