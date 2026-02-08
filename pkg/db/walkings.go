package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func (w *Walkings) UserInput(oldWalk Walk) (Walk, error) {

	// Get Data (strings)
	name := util.PromptWithSuggestion("Name", oldWalk.NAME)
	distanceStr := util.PromptWithSuggestion("Distance", strconv.FormatFloat(oldWalk.DISTANCE, 'f', 2, 64))
	duration := util.PromptWithSuggestion("Duration", oldWalk.DURATION)
	pace := util.PromptWithSuggestion("Pace", oldWalk.PACE)
	stepsStr := util.PromptWithSuggestion("Steps", strconv.Itoa(oldWalk.STEPS))
	caloriesStr := util.PromptWithSuggestion("Calories", strconv.Itoa(oldWalk.CALORIES))
	date := util.PromptWithSuggestion("Date", oldWalk.DATE)

	// Convert string to float64
	distance, err := strconv.ParseFloat(distanceStr, 64)
	if err != nil {
		return Walk{}, err
	}

	// Convert string to int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return Walk{}, err
	}

	calories, err := strconv.Atoi(caloriesStr)
	if err != nil {
		return Walk{}, err
	}

	return Walk{
		ID:       oldWalk.ID,
		NAME:     name,
		DISTANCE: distance,
		DURATION: duration,
		PACE:     pace,
		STEPS:    steps,
		CALORIES: calories,
		DATE:     date,
	}, nil
}

func (w *Walkings) ReadFromFile(path string) error {

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer file.Close()

	// Read entire file contents
	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(byteValue, w); err != nil {
		return fmt.Errorf("error parsing JSON from file %s: %w", path, err)
	}

	return nil
}

func (w *Walkings) FindWalk(id int) (int, Walk, error) {
	for index, foundWalk := range w.WALKINGS {
		if foundWalk.ID == id {
			return index, foundWalk, nil
		}
	}

	return -1, Walk{}, errors.New("walk not found")
}

func (w *Walkings) Update(index int, updatedWalk Walk) error {

	// Set correct ID
	updatedWalk.ID = w.WALKINGS[index].ID

	// Update
	w.WALKINGS[index] = updatedWalk

	return w.Save()
}

func (w *Walkings) Delete(index int) error {
	w.WALKINGS = append((w.WALKINGS)[:index], (w.WALKINGS)[index+1:]...)
	return w.Save()
}
