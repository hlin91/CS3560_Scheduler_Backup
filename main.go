// main.go is the "View" of the program
package main

import (
	"github.com/hlin91/CS3560_Scheduler_Backup/controller"
	"github.com/hlin91/CS3560_Scheduler_Backup/model"
)

func main() {
	s := model.NewSchedule()       // Create the schedule "Model"
	menu := controller.MakeMenu(s) // Create the menu "Controller"
	menu.Run()
}

//!--
