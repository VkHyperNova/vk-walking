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
	DATE     int     `json:"date"`
}

type Walkings struct {
	WALKINGS []Walk `json:"walkings"`
}

func (w *Walkings) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + "VK-WALKING 1.0" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)
	w.PrintAllWalks()
	w.PrintOverallStats()
	w.PrintStatsByYear()
	fmt.Println(color.Cyan + "\n< Add Update Delete Quit >" + color.Reset)

}

func (w *Walkings) PrintAllWalks() {
	for _, walk := range w.WALKINGS {
		id := fmt.Sprint(color.Yellow, walk.ID, color.Reset)
		name := fmt.Sprint(color.Green + "\"" + walk.NAME + "\"" + color.Reset)
		distance := fmt.Sprint(strconv.FormatFloat(walk.DISTANCE, 'f', 2, 64) + " miles")
		duration := fmt.Sprint("DURATION: " + walk.DURATION)
		pace := fmt.Sprint("PACE: " + walk.PACE)
		steps := fmt.Sprint("STEPS: " + strconv.Itoa(walk.STEPS))
		calories := fmt.Sprint("CALORIES: " + strconv.Itoa(walk.CALORIES))
		date := fmt.Sprint(walk.DATE)
		fmt.Println(id, name, color.Cyan+color.Italic+distance+" | ", duration+" | ", pace+" | ", steps+" | ", calories+" | ", date+color.Reset)
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

	fmt.Printf(color.Blue+color.Bold+color.Italic+"\n%d Total Miles: %.2f (%.2fkm) | %d steps | %d calories\n"+color.Reset,year, totalmiles, totalmiles*1.60934, totalsteps, totalcalories)
	}
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

func (w *Walkings) ResetIDs() {

	for key := range w.WALKINGS {
		w.WALKINGS[key].ID = key + 1
	}

	w.Save()
}

func (w *Walkings) Undo() bool {

	// Remove the last item
	w.WALKINGS = w.WALKINGS[:len(w.WALKINGS)-1]

	w.Save()

	return true
}
