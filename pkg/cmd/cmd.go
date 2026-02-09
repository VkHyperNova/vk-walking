package cmd

import (
	"fmt"
	"vk-walking/pkg/color"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine(w *db.Walkings) {
	for {

		w.PrintCLI()

		util.PrintArrow()

		command, id, ok := util.ReadCommand()
		if !ok {
			continue
		}

		switch command {
		case "a", "add":
			newWalk, err := w.UserInput(db.Walk{})
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = w.Add(newWalk)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println(color.Green + "\nItem Added!" + color.Reset)
			util.PressAnyKey()
			util.ClearScreen()

		case "u", "update":

			if id == 0 {
				fmt.Println("Please provide an ID")
				continue
			}

			index, foundWalk, err := w.FindWalk(id)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			updatedWalk, err := w.UserInput(foundWalk)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = w.Update(index, updatedWalk)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Println(color.Yellow + "\nItem Updated!" + color.Reset)
			util.PressAnyKey()
			util.ClearScreen()

		case "d", "delete":

			if id == 0 {
				fmt.Println("Please provide an ID")
				continue
			}
			index, _, err := w.FindWalk(id)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			err = w.Delete(index)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			w.ResetIDs()

			fmt.Println(color.Red + "\nItem Removed!" + color.Reset)
			util.PressAnyKey()
			util.ClearScreen()
		case "undo":
			w.Undo()
			util.ClearScreen()

		case "q", "quit":
			util.ClearScreen()
			return
		default:
			fmt.Println("Unknown command. Try: add, update, delete, quit")
			util.PressAnyKey()
			util.ClearScreen()
		}
	}
}
