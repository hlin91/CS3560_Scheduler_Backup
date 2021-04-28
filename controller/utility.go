// Package controller provides functions to edit the schedule or view (ie. command line)
// utility.go provides various utility functions private to the package
package controller

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hlin91/CS3560_Scheduler_Backup/model"
)

func displayHeader() {
	fmt.Println("==================================================")
	fmt.Println("Welcome to PSS!")
	fmt.Printf("Enter %q to quit\n", ESCAPE)
	fmt.Println("==================================================")
}

// Convert a string of format YYYY-MM-DD to a date integer for the scheduler
func stringToDateInt(s string) (int, error) {
	tok := strings.Split(s, "-")
	if len(tok) != 3 {
		return 0, fmt.Errorf("stringToDateInt: string does not match expected format")
	}
	year, _ := strconv.Atoi(tok[0])
	month, _ := strconv.Atoi(tok[1])
	day, _ := strconv.Atoi(tok[2])
	date, err := strconv.Atoi(fmt.Sprintf("%04d%02d%02d", year, month, day))
	if err != nil {
		return 0, fmt.Errorf("stringToDateInt: string is non-numeric")
	}
	return date, nil
}

// requestTaskInfo asks the user to enter transient task information
func requestTaskInfo() (string, string, int, float32, float32, error) {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	name := input.Text()
	displayTransientTypes()
	fmt.Print("Enter task type: ")
	input.Scan()
	taskType := input.Text()
	fmt.Print("Enter date (eg. 2020-11-14): ")
	input.Scan()
	date, err := stringToDateInt(input.Text())
	if err != nil {
		return "", "", 0, 0, 0, fmt.Errorf("bad date entered")
	}
	fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	startTime, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", "", 0, 0, 0, fmt.Errorf("bad start time entered")
	}
	fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	duration, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", "", 0, 0, 0, fmt.Errorf("bad duration entered")
	}
	return name, taskType, date, float32(startTime), float32(duration), nil
}

// requestAntiInfo asks the user to enter anti task information
func requestAntiInfo() (string, int, float32, float32, error) {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	name := input.Text()
	fmt.Print("Enter date (eg. 2020-11-14): ")
	input.Scan()
	date, err := stringToDateInt(input.Text())
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("bad date entered")
	}
	fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	startTime, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("bad start time entered")
	}
	fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	duration, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("bad duration entered")
	}
	return name, date, float32(startTime), float32(duration), nil
}

// requestRecurringInfo asks the user to enter recurring task information
func requestRecurringInfo() (string, string, int, float32, float32, int, int, error) {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	name := input.Text()
	displayRecurringTypes()
	fmt.Print("Enter task type: ")
	input.Scan()
	taskType := input.Text()
	fmt.Print("Enter date (eg. 2020-11-14): ")
	input.Scan()
	date, err := stringToDateInt(input.Text())
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad date entered")
	}
	fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	startTime, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad start time entered")
	}
	fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
	input.Scan()
	duration, err := strconv.ParseFloat(input.Text(), 32)
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad duration entered")
	}
	fmt.Print("Enter end date (eg. 2020-11-14): ")
	input.Scan()
	endDate, err := stringToDateInt(input.Text())
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad date entered")
	}
	fmt.Print("Enter frequency (1 or 7): ")
	input.Scan()
	frequency, err := strconv.Atoi(input.Text())
	if err != nil {
		return "", "", 0, 0, 0, 0, 0, fmt.Errorf("bad frequency entered")
	}
	return name, taskType, date, float32(startTime), float32(duration), endDate, frequency, nil
}

func displayTransientTypes() {
	fmt.Println("Available types...")
	fmt.Println(model.VISIT)
	fmt.Println(model.SHOPPING)
	fmt.Println(model.APPOINTMENT)
}

func displayRecurringTypes() {
	fmt.Println("Available types...")
	fmt.Println(model.CLASS)
	fmt.Println(model.STUDY)
	fmt.Println(model.SLEEP)
	fmt.Println(model.EXERCISE)
	fmt.Println(model.WORK)
	fmt.Println(model.MEAL)
}

//!--
