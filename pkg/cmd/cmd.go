package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine(w *db.Walkings) {
	for {

		w.PrintCLI()

		var userInput string = ""
		var userInputID int = 0

		fmt.Print("=> ")

		fmt.Scanln(&userInput, &userInputID)

		userInput = strings.ToLower(userInput)

		switch userInput {
		case "a", "add":
			newWalk, err := w.UserInput(db.Walk{})
			if err != nil {
				log.Fatal(err)
			}

			err = w.Add(newWalk)
			if err != nil {
				log.Fatal(err)
			}

		case "u", "update":
			index, foundWalk, err := w.FindWalk(userInputID)
			if err != nil {
				log.Fatal(err)
			}

			updatedWalk, err := w.UserInput(foundWalk)
			if err != nil {
				log.Fatal(err)
			}

			err = w.Update(index, updatedWalk)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Updated!")

		case "d", "delete":
			index, _, err := w.FindWalk(userInputID)
			if err != nil {
				log.Fatal(err)
			}

			err = w.Delete(index)
			if err != nil {
				log.Fatal(err)
			}

			w.ResetIDs()
			
			fmt.Println("Walk removed!")

		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
		}
	}
}
