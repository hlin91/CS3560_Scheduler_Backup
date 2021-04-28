// Package controller provides functions to edit the schedule or view (ie. command line)
// menu.go provides classes to create a generic menu
package controller

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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

// Run the menu until the user passes in the escape string
func (m *Menu) Run() {
	input := bufio.NewScanner(os.Stdin)
	m.Clear()
	displayHeader()
	m.Display()
	for input.Scan() {
		if input.Text() == ESCAPE {
			return
		}
		option, err := strconv.Atoi(input.Text())
		if err != nil {
			fmt.Println("Error: bad option")
			displayHeader()
			m.Display()
			continue
		}
		m.Clear()
		err = m.Process(option)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Success!")
		}
		fmt.Print("Press enter to continue...")
		input.Scan()
		m.Clear()
		displayHeader()
		m.Display()
	}
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

//!--
