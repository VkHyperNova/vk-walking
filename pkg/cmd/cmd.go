package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func CommandLine(walkings *db.Walkings) {
	for {

		walkings.PrintCLI()

		var userInput string = ""
		var userInputID int = 0

		fmt.Print("=> ")

		fmt.Scanln(&userInput, &userInputID)

		userInput = strings.ToLower(userInput)

		switch userInput {
		case "a", "add":
			newWalk, err := walkings.UserInput()
			if err != nil {
				log.Fatal(err)
			}

			err = walkings.Add(newWalk)
			if err != nil {
				log.Fatal(err)
			}

		// case "u", "update":
		// 	index, updateSuggestions := books.FindBook(userInputID)
		// 	if index != -1 {
		// 		updatedBook := books.UserInput(updateSuggestions)
		// 		err := books.Update(index, updatedBook)
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 	} else {
		// 		fmt.Printf("Book %d not found\n", userInputID)
		// 	}

		// case "d", "delete":
		// 	index, _ := books.FindBook(userInputID)
		// 	if index != -1 {
		// 		err := books.Delete(index)
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		} else {
		// 			fmt.Println("Item removed!")
		// 		}
		// 	}

		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		default:
			util.ClearScreen()
		}
	}
}