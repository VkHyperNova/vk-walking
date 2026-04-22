package cmd

import (
	"fmt"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func Run(store *db.Store) {
	for {
		store.PrintDashboard()

		var command string
		var id int

		fmt.Print("=> ")

		fmt.Scanln(&command, &id)

		command = strings.ToLower(command)

		switch command {
		case "a", "add":
			if err := store.Add(); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nAdded!" + color.Reset)
			}
		case "u", "update":
			if err := store.Update(id); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nUpdated!" + color.Reset)
			}
		case "d", "delete":
			if err := store.Delete(id); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nRemoved!" + color.Reset)
			}
		case "undo":
			if err := store.Undo(); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(color.Yellow + "\nUndone!" + color.Reset)
			}
		case "l", "list":
			store.PrintAll()
		case "stats":
			store.PrintStats(id)
		case "distance":
			store.PrintDistance(id)
		case "steps":
			store.PrintSteps(id)
		case "calories":
			store.PrintCalories(id)
		case "duration":
			store.PrintDuration(id)
		case "q", "quit":
			util.ClearScreen()
			return
		default:
			fmt.Println("Unknown command. Try: add, update, delete, undo, stats, list, distance, steps, calories, duration, quit")
		}
	}
}
