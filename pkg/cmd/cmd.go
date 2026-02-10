package cmd

import (
	"fmt"
	"log"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine() {
	for {
		// Initialize Walkings database
		w := db.Walkings{}

		// Reload Database
		err := w.ReadFromFile(config.LocalPath)
		if err != nil {
			log.Fatalf("Fatal error: failed to load walkings database: %v", err)
		}

		w.PrintCLI()

		util.PrintArrow()

		command, id, ok := util.ReadCommand()
		if !ok {
			continue
		}

		switch command {
		case "a", "add":
			if err := w.Add(); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Added!" + color.Reset)
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
