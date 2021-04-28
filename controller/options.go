// Package controller provides functions to edit the schedule or view (ie. command line)
// options.go provides the options the user can choose from for interacting with the schedule
package controller

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/hlin91/CS3560_Scheduler_Backup/model"
)

const (
	ESCAPE = "quit"
)

// Make the menu and populate with options to interact with schedule
func MakeMenu(s *model.Schedule) Menu {
	/**********************************************
		 * Options are...
		 * Create a task
		 * View a task
		 * View by month
		 * View by week
		 * View by day
		 * Delete a task
		 * Read schedule from file
		 * Write schedule to file
		 * Write by day
		 * Write by week
		 * Write by month
	     **********************************************/
	options := []ScheduleMenuItem{}
	// Add create task option
	options = append(options, NewScheduleMenuItem("Create a task", s, createTask))
	options = append(options, NewScheduleMenuItem("Delete a task", s, deleteTask))
	options = append(options, NewScheduleMenuItem("Edit a task", s, editTask))
	options = append(options, NewScheduleMenuItem("View a task", s, viewTask))
	options = append(options, NewScheduleMenuItem("View by month", s, viewTaskByMonth))
	options = append(options, NewScheduleMenuItem("View by week", s, viewTaskByWeek))
	options = append(options, NewScheduleMenuItem("View by day", s, viewTaskByDay))
	// TODO: Add edit task option
	// TODO: Add file IO options
	m := []Menuer{}
	for _, o := range options {
		temp := o
		m = append(m, Menuer(&temp))
	}
	return NewMenu(m)
}

func createTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	valid := false
	fmt.Println("Select the type of task to add")
	fmt.Println("1. Transient task")
	fmt.Println("2. Anti task")
	fmt.Println("3. Recurring task")
	fmt.Print("Enter an option: ")
	for !valid {
		switch input.Scan(); input.Text() {
		case "1":
			valid = true
			name, taskType, date, startTime, duration, err := requestTaskInfo()
			if err != nil {
				return err
			}
			return s.AddTransientTask(name, taskType, date, startTime, duration)
		case "2":
			valid = true
			name, date, startTime, duration, err := requestAntiInfo()
			if err != nil {
				return err
			}
			return s.AddAntiTask(name, model.CANCEL, date, float32(startTime), float32(duration))
		case "3":
			valid = true
			name, taskType, date, startTime, duration, endDate, frequency, err := requestRecurringInfo()
			if err != nil {
				return err
			}
			return s.AddRecurringTask(name, taskType, date, startTime, duration, endDate, frequency)
		default:
			fmt.Print("Invalid option. Try again: ")
		}
	}
	return nil
}

func viewTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a task name: ")
	input.Scan()
	if t, ok := s.TransientTasks[input.Text()]; ok {
		fmt.Println(t)
		return nil
	}
	if t, ok := s.RecurringTasks[input.Text()]; ok {
		fmt.Println(t)
		return nil
	}
	if t, ok := s.AntiTasks[input.Text()]; ok {
		fmt.Println(t)
		return nil
	}
	return fmt.Errorf("task name does not exist in schedule")
}

func viewTaskByMonth(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a month (1-12): ")
	input.Scan()
	month, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad month entered")
	}
	tasks, err := s.GetTasksByMonth(month)
	if err != nil {
		return err
	}
	for i, t := range tasks {
		fmt.Println(t)
		if i < len(tasks)-1 {
			fmt.Println("-----------------------------------------")
		}
	}
	return nil
}

func viewTaskByWeek(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a month (1-12): ")
	input.Scan()
	month, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad month entered")
	}
	fmt.Print("Enter a day (1-31): ")
	input.Scan()
	day, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad day entered")
	}
	tasks, err := s.GetTasksByWeek(month, day)
	if err != nil {
		return err
	}
	for i, t := range tasks {
		fmt.Println(t)
		if i < len(tasks)-1 {
			fmt.Println("-----------------------------------------")
		}
	}
	return nil
}

func viewTaskByDay(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a month (1-12): ")
	input.Scan()
	month, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad month entered")
	}
	fmt.Print("Enter a day (1-31): ")
	input.Scan()
	day, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad day entered")
	}
	tasks, err := s.GetTasksByDay(month, day)
	if err != nil {
		return err
	}
	for i, t := range tasks {
		fmt.Println(t)
		if i < len(tasks)-1 {
			fmt.Println("-----------------------------------------")
		}
	}
	return nil
}

func deleteTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	return s.DeleteTask(input.Text())
}

func editTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the name of the task to edit: ")
	input.Scan()
	taskName := input.Text()
	if _, ok := s.TransientTasks[taskName]; ok {
		// Edit a transient task
		newName, newType, newDate, newStartTime, newDuration, err := requestTaskInfo()
		if err != nil {
			return err
		}
		err = s.EditTransientTask(taskName, newName, newType, newDate, newStartTime, newDuration)
		return err
	}
	if _, ok := s.AntiTasks[taskName]; ok {
		// Edit an anti task
		newName, newDate, newStartTime, newDuration, err := requestAntiInfo()
		if err != nil {
			return err
		}
		err = s.EditAntiTask(taskName, newName, newDate, newStartTime, newDuration)
		return err
	}
	if _, ok := s.RecurringTasks[taskName]; ok {
		// Edit a recurring task
		newName, newType, newDate, newStartTime, newDuration, newEndDate, newFrequency, err := requestRecurringInfo()
		if err != nil {
			return err
		}
		err = s.EditRecurringTask(taskName, newName, newType, newDate, newStartTime, newDuration, newEndDate, newFrequency)
		return err
	}
	return fmt.Errorf("could not find task with name %q", taskName)
}

// TODO: Make file IO menu options

//!--
