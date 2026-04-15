package cmd

import (
	"fmt"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine(walkings *db.Walkings) {
	for {
		walkings.PrintCLI()

		var cmd string
		var id int

		fmt.Print("=> ")

		fmt.Scanln(&cmd, &id)

		cmd = strings.ToLower(cmd)

		switch cmd {
		case "a", "add":
			if err := walkings.Add(); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Added!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "u", "update":
			if err := walkings.Update(id); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Println(color.Yellow + "\nItem Updated!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "d", "delete":
			if err := walkings.Delete(id); err != nil {
				fmt.Println(color.Red+"Error:"+color.Reset, err)
			} else {
				fmt.Printf(color.Yellow + "\n Item Removed!" + color.Reset)
			}
			util.PressAnyKey()
			util.ClearScreen()
		case "undo":
			walkings.Undo()
			util.ClearScreen()
		case "showall":
			util.ClearScreen()
			walkings.PrintAllWalks()
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
