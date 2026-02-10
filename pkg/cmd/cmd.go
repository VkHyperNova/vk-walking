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
			err := w.Add()
			if err != nil {
				fmt.Println("Error:", err)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "u", "update":
			if err := w.Update(id); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Updated!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()

		case "d", "delete":
			if err := w.Delete(id); err != nil {
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
