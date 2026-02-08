package db

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"
	"vk-walking/pkg/util"
)

type Walk struct {
	ID       int     `json:"id"`
	NAME     string  `json:"name"`
	DISTANCE float64 `json:"distance"`
	DURATION string  `json:"duration"`
	PACE     string  `json:"pace"`
	STEPS    int     `json:"steps"`
	CALORIES int     `json:"calories"`
	DATE     string  `json:"date"`
}

type Walkings struct {
	WALKINGS []Walk `json:"walkings"`
}

func (w *Walkings) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + "VK-WALKING" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)

}

func (w *Walkings) Add(newWalk Walk) error {

	// Add unique ID
	newWalk.ID = w.NewID()

	// Add
	w.WALKINGS = append(w.WALKINGS, newWalk)

	return w.Save()
}

func (w *Walkings) NewID() int {

	maxID := 0

	for _, book := range w.WALKINGS {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	return maxID + 1
}

func (w *Walkings) Save() error {

	// Format JSON
	books, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	// Save
	err = os.WriteFile(config.LocalPath, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup
	err = os.WriteFile(config.BackupPath, books, 0644)
	if err != nil {
		return err
	}

	// Save Backup with Date
	err = os.WriteFile(config.BackupPathWithDate, books, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (w *Walkings) UserInput() (Walk, error) {

	// Get Data
	var answers []string
	for _, question := range config.Questions {
		input := util.Input(question)
		answers = append(answers, input)
	}

	// Convert string to float64
	distance, err := strconv.ParseFloat(answers[1], 64)
	if err != nil {
		return Walk{}, err
	}

	// Convert string to int
	steps, err := strconv.Atoi(answers[4])
	if err != nil {
		return Walk{}, err
	}

	calories, err := strconv.Atoi(answers[5])
	if err != nil {
		return Walk{}, err
	}


	return Walk{
		ID:       0,
		NAME:     answers[0],
		DISTANCE: distance,
		DURATION: answers[2],
		PACE:     answers[3],
		STEPS:    steps,
		CALORIES: calories,
		DATE:     answers[6],
	}, nil

}
