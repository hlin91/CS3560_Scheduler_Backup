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
	options = append(options, NewScheduleMenuItem("Load from file", s, loadFile))
	options = append(options, NewScheduleMenuItem("Write tasks to file", s, writeTasks))
	options = append(options, NewScheduleMenuItem("Write tasks by month", s, writeTasksByMonth))
	options = append(options, NewScheduleMenuItem("Write tasks by week", s, writeTasksByWeek))
	options = append(options, NewScheduleMenuItem("Write tasks by day", s, writeTasksByDay))
	m := []Menuer{}
	for _, o := range options {
		temp := o
		m = append(m, Menuer(&temp))
	}
	return NewMenu(m)
}

// The following functions implement each option in the menu

// createTask allows the user to create and add a task to the schedule
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

// viewTask allows the user to view the details of a task by name
func viewTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a task name: ")
	input.Scan()
	if t, ok := s.TransientTasks[input.Text()]; ok {
		fmt.Println(SEP_STRING)
		fmt.Println(t)
		fmt.Println(SEP_STRING)
		return nil
	}
	if t, ok := s.RecurringTasks[input.Text()]; ok {
		fmt.Println(SEP_STRING)
		fmt.Println(t)
		fmt.Println(SEP_STRING)
		return nil
	}
	if t, ok := s.AntiTasks[input.Text()]; ok {
		fmt.Println(SEP_STRING)
		fmt.Println(t)
		fmt.Println(SEP_STRING)
		return nil
	}
	return fmt.Errorf("task name does not exist in schedule")
}

// viewTaskByMonth allows the user to view all tasks for a specified month
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
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	fmt.Println(SEP_STRING)
	for _, t := range tasks {
		fmt.Println(t)
		fmt.Println(SEP_STRING)
	}
	return nil
}

// viewTaskByWeek allows the user to view all tasks for the week of a specified day
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
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	fmt.Println(SEP_STRING)
	for _, t := range tasks {
		fmt.Println(t)
		fmt.Println(SEP_STRING)
	}
	return nil
}

// viewTaskByDay allows the user to view all tasks for a specified day
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
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	fmt.Println(SEP_STRING)
	for _, t := range tasks {
		fmt.Println(t)
		fmt.Println(SEP_STRING)
	}
	return nil
}

// deleteTask allows the user to delete a task by name
func deleteTask(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	return s.DeleteTask(input.Text())
}

// editTask allows the user to edit the details of a task by name
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

// loadFile allows the user to load the tasks from a specified json file
func loadFile(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the file to load: ")
	input.Scan()
	filePath := input.Text()
	return s.LoadFile(filePath)
}

// writeTasks allows the user to write all tasks to a specified json file
func writeTasks(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the file to write to: ")
	input.Scan()
	filePath := input.Text()
	return s.WriteTasks(filePath)
}

// writeTasksByMonth allows the user to write all tasks for a specified month to a json file
func writeTasksByMonth(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the file to write to: ")
	input.Scan()
	filePath := input.Text()
	fmt.Print("Enter a month (1-12): ")
	input.Scan()
	month, err := strconv.Atoi(input.Text())
	if err != nil {
		return fmt.Errorf("bad month entered")
	}
	tasks, err := s.GetTasksByMonth(month)
	if err != nil {
		return fmt.Errorf("error fetching tasks: %v", err)
	}
	return s.WriteTaskList(filePath, tasks)
}

// writeTasksByWeek allows the user to write all tasks for a specified week to a json file
func writeTasksByWeek(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the file to write to: ")
	input.Scan()
	filePath := input.Text()
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
		return fmt.Errorf("error fetching tasks: %v", err)
	}
	return s.WriteTaskList(filePath, tasks)
}

// writeTasksByDay allows the user to write all tasks for a specified day to a json file
func writeTasksByDay(s *model.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the path of the file to write to: ")
	input.Scan()
	filePath := input.Text()
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
		return fmt.Errorf("error fetching tasks: %v", err)
	}
	return s.WriteTaskList(filePath, tasks)
}

//!--
