package cmd

import (
	"fmt"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine(w *db.WalkData) {
	for {
		w.PrintCLI()

		var cmd string
		var n int

		fmt.Print("=> ")

		fmt.Scanln(&cmd, &n)

		cmd = strings.ToLower(cmd)

		switch cmd {
		case "a", "add":
			if err := w.Add(); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Added!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "u", "update":
			if err := w.Update(n); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Updated!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "d", "delete":
			if err := w.Delete(n); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Printf(color.Yellow + "\n Item Removed!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "undo":
			w.Undo()
			util.ClearScreen()
		case "showall":
			util.ClearScreen()
			w.PrintAllWalks()
			util.PressAnyKey()
		case "stats":
			w.PrintStats(n)
		case "distance":
			w.PrintDistance(n)
		case "steps":
			w.PrintSteps(n)
		case "calories":
			w.PrintCalories(n)
		case "duration":
			w.PrintDuration(n)
		case "q", "quit":
			util.ClearScreen()
			return
		default:
			fmt.Println("Unknown command. Try: add, update, delete, undo, stats, showall, distance, steps, calories, duration, quit")
			util.PressAnyKey()
			util.ClearScreen()
		}
	}
}
