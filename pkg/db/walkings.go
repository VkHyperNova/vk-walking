package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
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
	DATE     int     `json:"date"`
}

type Walkings struct {
	WALKINGS []Walk `json:"walkings"`
}

func (w *Walkings) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + "VK-WALKING 1.0" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)
	w.PrintTopTen()
	w.PrintOverallStats()
	w.PrintStatsByYear()
	fmt.Println(color.Cyan + "\n< Add Update Delete showall undo Quit >" + color.Reset)

}

func (w *Walkings) PrintTopTen() {
	// Sort descending by distance
	sort.Slice(w.WALKINGS, func(i, j int) bool {
		return w.WALKINGS[i].DISTANCE > w.WALKINGS[j].DISTANCE
	})

	// Determine how many to print (up to 10)
	n := 10
	if len(w.WALKINGS) < 10 {
		n = len(w.WALKINGS)
	}

	// Print top n
	for i := 0; i < n; i++ {
		walk := w.WALKINGS[i]
		number := fmt.Sprintf("%d. (ID: %d) ", i+1, walk.ID)
		distance := fmt.Sprintf("%s%s%.2f miles(%.2f km)%s | ", color.Blue, color.Bold, walk.DISTANCE, walk.DISTANCE*1.60934, color.Reset)
		name := fmt.Sprintf(" %s%s%s | ", color.Green, walk.NAME, color.Reset)
		steps := fmt.Sprintf("%s%d%s steps | ", color.Yellow, walk.STEPS, color.Reset)
		calories := fmt.Sprintf("%s%d%s calories ", color.Yellow, walk.CALORIES, color.Reset)
		pace := fmt.Sprintf(" %s ", walk.PACE)
		duration := fmt.Sprintf(" %s ", walk.DURATION)
		fmt.Println(number + distance + steps + calories + pace + duration + name)
	}
}

func (w *Walkings) PrintAllWalks() {
	for i, walk := range w.WALKINGS {
		fmt.Printf(
			"%d. ID:%d  Distance:%.2f  Duration:%s  Pace:%s  Steps:%d  Calories:%d  Date:%d  Name:%s\n",
			i+1,
			walk.ID,
			walk.DISTANCE,
			walk.DURATION,
			walk.PACE,
			walk.STEPS,
			walk.CALORIES,
			walk.DATE,
			walk.NAME,
		)
	}
}

func (w *Walkings) PrintOverallStats() {

	totalmiles := 0.0
	totalsteps := 0
	totalcalories := 0

	for _, walk := range w.WALKINGS {
		totalmiles += walk.DISTANCE
		totalsteps += walk.STEPS
		totalcalories += walk.CALORIES
	}
	fmt.Printf(color.Blue+color.Bold+color.Italic+"\nOVERALL: Total Miles: %.2f (%.2fkm) | %d steps | %d calories\n"+color.Reset, totalmiles, totalmiles*1.60934, totalsteps, totalcalories)

}

func (w *Walkings) PrintStatsByYear() {

	var years []int

	for _, walk := range w.WALKINGS {
		years = util.AppendIfMissing(years, walk.DATE)
	}

	for i := 0; i < len(years); i++ {
		year := years[i]
		totalmiles := 0.0
		totalsteps := 0
		totalcalories := 0
		for _, walk := range w.WALKINGS {
			if walk.DATE == years[i] {
				totalmiles += walk.DISTANCE
				totalsteps += walk.STEPS
				totalcalories += walk.CALORIES
			}
		}

		fmt.Printf(color.Blue+color.Bold+color.Italic+"\n%d Total Miles: %.2f (%.2fkm) | %d steps | %d calories\n"+color.Reset, year, totalmiles, totalmiles*1.60934, totalsteps, totalcalories)
	}
}

func (w *Walkings) Add() error {

	// Get new walk data
	newWalk, err := w.UserInput(Walk{})
	if err != nil {
		return err
	}

	// Add unique ID
	newWalk.ID = w.NewID()

	// Add
	w.WALKINGS = append(w.WALKINGS, newWalk)

	// Save
	err = w.Save()
	if err != nil {
		return err
	}

	fmt.Println(color.Green + "\nItem Added!" + color.Reset)

	return nil
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
	walks, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	// Save
	err = os.WriteFile(config.LocalPath, walks, 0644)
	if err != nil {
		return err
	}

	// Save Backup
	err = os.WriteFile(config.BackupPath, walks, 0644)
	if err != nil {
		return err
	}

	// Save Backup with Date
	err = os.WriteFile(config.BackupPathWithDate, walks, 0644)
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
	dateStr := util.PromptWithSuggestion("Date", strconv.Itoa(oldWalk.DATE))

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

	date, err := strconv.Atoi(dateStr)
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

func (w *Walkings) Update(id int) error {

	// Invalid IDs Guard
	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Find and Update
	for index, walk := range w.WALKINGS {

		// Find walk
		if walk.ID == id {

			// Get updated fields
			updatedWalk, err := w.UserInput(walk)
			if err != nil {
				return err
			}

			// Update
			w.WALKINGS[index] = updatedWalk

			// Save
			return w.Save()
		}
	}
	return fmt.Errorf("item with ID %d not found", id)
}

func (w *Walkings) Delete(id int) error {

	// Invalid IDs Guard
	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	// Find and Delete
	for index, walk := range w.WALKINGS {
		if walk.ID == id {

			// Delete
			w.WALKINGS = append((w.WALKINGS)[:index], (w.WALKINGS)[index+1:]...)

			w.ResetIDs()

			return w.Save()
		}
	}

	return fmt.Errorf("item with ID %d not found", id)
}

func (w *Walkings) ResetIDs() {
	for key := range w.WALKINGS {
		w.WALKINGS[key].ID = key + 1
	}
}

func (w *Walkings) Undo() bool {

	// Remove the last item
	w.WALKINGS = w.WALKINGS[:len(w.WALKINGS)-1]

	w.Save()

	return true
}
