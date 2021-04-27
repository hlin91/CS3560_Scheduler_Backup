package main

import (
	"bufio"
	"fmt"
	"github.com/hlin91/CS3560_Scheduler_Backup/schedule"
	"os"
	"strconv"
	"strings"
)

const (
	ESCAPE = "quit"
)

func main() {
	s := schedule.NewSchedule()
	menu := makeMenu(&s)
	input := bufio.NewScanner(os.Stdin)
	displayHeader()
	menu.Display()
	for input.Scan() {
		if input.Text() == ESCAPE {
			return
		}
		option, err := strconv.Atoi(input.Text())
		if err != nil {
			fmt.Println("Error: bad option")
			displayHeader()
			menu.Display()
			continue
		}
		err = menu.Process(option)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Success!")
		}
		fmt.Print("Press enter to continue...")
		input.Scan()
		displayHeader()
		menu.Display()
	}
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

func displayHeader() {
	fmt.Println("==================================================")
	fmt.Println("Welcome to PSS!")
	fmt.Printf("Enter %q to quit\n", ESCAPE)
	fmt.Println("==================================================")
}

// Make the menu and populate with options to interact with schedule
func makeMenu(s *schedule.Schedule) schedule.Menu {
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
	options := []schedule.ScheduleMenuItem{}
	// Add create task option
	options = append(options, schedule.NewScheduleMenuItem("Create a task", s, createTask))
	options = append(options, schedule.NewScheduleMenuItem("View a task", s, viewTask))
	options = append(options, schedule.NewScheduleMenuItem("View by month", s, viewTaskByMonth))
	options = append(options, schedule.NewScheduleMenuItem("View by week", s, viewTaskByWeek))
	options = append(options, schedule.NewScheduleMenuItem("View by day", s, viewTaskByDay))
	options = append(options, schedule.NewScheduleMenuItem("Delete a task", s, deleteTask))
	// TODO: Add edit task option
	// TODO: Add file IO options
	m := []schedule.Menuer{}
	for _, o := range options {
		temp := o
		m = append(m, schedule.Menuer(&temp))
	}
	return schedule.NewMenu(m)
}

func displayTransientTypes() {
	fmt.Println("Available types...")
	fmt.Println(schedule.VISIT)
	fmt.Println(schedule.SHOPPING)
	fmt.Println(schedule.APPOINTMENT)
}

func displayRecurringTypes() {
	fmt.Println("Available types...")
	fmt.Println(schedule.CLASS)
	fmt.Println(schedule.STUDY)
	fmt.Println(schedule.SLEEP)
	fmt.Println(schedule.EXERCISE)
	fmt.Println(schedule.WORK)
	fmt.Println(schedule.MEAL)
}

func createTask(s *schedule.Schedule) error {
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
				return fmt.Errorf("bad date entered")
			}
			fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			startTime, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad start time entered")
			}
			fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			duration, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad duration entered")
			}
			return s.AddTransientTask(name, taskType, date, float32(startTime), float32(duration))
		case "2":
			valid = true
			fmt.Print("Enter task name: ")
			input.Scan()
			name := input.Text()
			fmt.Print("Enter date (eg. 2020-11-14): ")
			input.Scan()
			date, err := stringToDateInt(input.Text())
			if err != nil {
				return fmt.Errorf("bad date entered")
			}
			fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			startTime, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad start time entered")
			}
			fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			duration, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad duration entered")
			}
			return s.AddAntiTask(name, schedule.CANCEL, date, float32(startTime), float32(duration))
		case "3":
			valid = true
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
				return fmt.Errorf("bad date entered")
			}
			fmt.Print("Enter start time (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			startTime, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad start time entered")
			}
			fmt.Print("Enter duration (eg. '8.5' for 8 hours 30 min): ")
			input.Scan()
			duration, err := strconv.ParseFloat(input.Text(), 32)
			if err != nil {
				return fmt.Errorf("bad duration entered")
			}
			fmt.Print("Enter end date (eg. 2020-11-14): ")
			input.Scan()
			endDate, err := stringToDateInt(input.Text())
			if err != nil {
				return fmt.Errorf("bad date entered")
			}
			fmt.Print("Enter frequency (1 or 7): ")
			input.Scan()
			frequency, err := strconv.Atoi(input.Text())
			if err != nil {
				return fmt.Errorf("bad frequency entered")
			}
			return s.AddRecurringTask(name, taskType, date, float32(startTime), float32(duration), endDate, frequency)
		default:
			fmt.Print("Invalid option. Try again: ")
		}
	}
	return nil
}

func viewTask(s *schedule.Schedule) error {
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

func viewTaskByMonth(s *schedule.Schedule) error {
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

func viewTaskByWeek(s *schedule.Schedule) error {
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

func viewTaskByDay(s *schedule.Schedule) error {
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

func deleteTask(s *schedule.Schedule) error {
	input := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter task name: ")
	input.Scan()
	return s.DeleteTask(input.Text())
}

// TODO: Make file IO menu options

//!--
