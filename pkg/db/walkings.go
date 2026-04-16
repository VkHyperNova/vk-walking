package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"
	"vk-walking/pkg/util"
)

type Walk struct {
	ID       int    `json:"id"` // keep as int
	DISTANCE string `json:"distance"`
	DURATION string `json:"duration"`
	STEPS    string `json:"steps"`
	CALORIES string `json:"calories"`
	DATE     string `json:"date"`
}

type WalkData struct {
	WALKINGS []Walk `json:"walkings"`
}

/* Main Functions */

func (w *WalkData) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + "VK-WALKING 1.0" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)

	w.PrintTopDistance()
	w.PrintTopSteps()
	w.PrintTopCalories()
	w.PrintTopDuration()

}

func (w *WalkData) Add() error {

	newWalk, err := w.GetUserInput(Walk{})
	if err != nil {
		return err
	}

	newWalk.ID = w.NewID()

	w.WALKINGS = append(w.WALKINGS, newWalk)

	return w.Save()
}

func (w *WalkData) Update(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := w.findIndex(id)
	if err != nil {
		return err
	}

	updated, err := w.GetUserInput((w.WALKINGS)[index])
	if err != nil {
		return err
	}

	(w.WALKINGS)[index] = updated

	return w.Save()
}

func (w *WalkData) Delete(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := w.findIndex(id)
	if err != nil {
		return err
	}

	confirm := util.Confirm()
	if !confirm {
		return fmt.Errorf("Abort")
	}

	w.WALKINGS = append((w.WALKINGS)[:index], (w.WALKINGS)[index+1:]...)

	return w.Save()
}

/* Top Stats */

func (w *WalkData) PrintTopDistance() {
	// Copy the slice
	sorted := make([]Walk, len(w.WALKINGS))
	copy(sorted, w.WALKINGS)
	// Sort copy descending by distance
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DISTANCE > sorted[j].DISTANCE
	})
	// Determine how many to print (up to 5)
	n := 5
	if len(sorted) < n {
		n = len(sorted)
	}
	// Print top n
	fmt.Print(color.Blue + color.Bold + "\nTop Distance \n" + color.Reset)
	for i := 0; i < n; i++ {
		walk := sorted[i]
		number := fmt.Sprintf("ID: %d ", walk.ID)
		floatDistance, _ := strconv.ParseFloat("3.14", 64)
		distance := fmt.Sprintf("%s%s%s miles(%.2f km)%s | ", color.Yellow, color.Bold, walk.DISTANCE, floatDistance*1.60934, color.Reset)
		steps := fmt.Sprintf("%s steps | ", walk.STEPS)
		calories := fmt.Sprintf("%s calories ", walk.CALORIES)
		duration := fmt.Sprintf(" %s ", walk.DURATION)
		fmt.Println(number + distance + steps + calories + duration)
	}
}

func (w *WalkData) PrintTopSteps() {
	// Copy the slice
	sorted := make([]Walk, len(w.WALKINGS))
	copy(sorted, w.WALKINGS)
	// Sort copy descending by distance
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].STEPS > sorted[j].STEPS
	})
	// Determine how many to print (up to 5)
	n := 5
	if len(sorted) < n {
		n = len(sorted)
	}
	// Print top n
	fmt.Print(color.Blue + color.Bold + "\nTop Steps \n" + color.Reset)
	for i := 0; i < n; i++ {
		walk := sorted[i]
		number := fmt.Sprintf("ID: %d ", walk.ID)
		floatDistance, _ := strconv.ParseFloat("3.14", 64)
		distance := fmt.Sprintf("%s miles(%.2f km) | ", walk.DISTANCE, floatDistance*1.60934)
		steps := fmt.Sprintf("%s%s%s%s steps | ", color.Yellow, color.Bold, walk.STEPS, color.Reset)
		calories := fmt.Sprintf("%s calories ", walk.CALORIES)
		duration := fmt.Sprintf(" %s ", walk.DURATION)
		fmt.Println(number + distance + steps + calories + duration)
	}
}

func (w *WalkData) PrintTopCalories() {
	// Copy the slice
	sorted := make([]Walk, len(w.WALKINGS))
	copy(sorted, w.WALKINGS)
	// Sort copy descending by distance
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].CALORIES > sorted[j].CALORIES
	})
	// Determine how many to print (up to 5)
	n := 5
	if len(sorted) < n {
		n = len(sorted)
	}
	// Print top n
	fmt.Print(color.Blue + color.Bold + "\nTop Calories \n" + color.Reset)
	for i := 0; i < n; i++ {
		walk := sorted[i]
		number := fmt.Sprintf("ID: %d ", walk.ID)
		floatDistance, _ := strconv.ParseFloat("3.14", 64)
		distance := fmt.Sprintf("%s miles(%.2f km) | ", walk.DISTANCE, floatDistance*1.60934)
		steps := fmt.Sprintf("%s steps | ", walk.STEPS)
		calories := fmt.Sprintf("%s%s%s%s calories ", color.Yellow, color.Bold, walk.CALORIES, color.Reset)
		duration := fmt.Sprintf(" %s ", walk.DURATION)
		fmt.Println(number + distance + steps + calories + duration)
	}
}

func (w *WalkData) PrintTopDuration() {
	// Copy the slice
	sorted := make([]Walk, len(w.WALKINGS))
	copy(sorted, w.WALKINGS)
	// Sort copy descending by distance
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].DURATION > sorted[j].DURATION
	})
	// Determine how many to print (up to 5)
	n := 5
	if len(sorted) < n {
		n = len(sorted)
	}
	// Print top n
	fmt.Print(color.Blue + color.Bold + "\nTop Duration \n" + color.Reset)
	for i := 0; i < n; i++ {
		walk := sorted[i]
		number := fmt.Sprintf("ID: %d ", walk.ID)
		floatDistance, _ := strconv.ParseFloat("3.14", 64)
		distance := fmt.Sprintf("%s miles(%.2f km) | ", walk.DISTANCE, floatDistance*1.60934)
		steps := fmt.Sprintf("%s steps | ", walk.STEPS)
		calories := fmt.Sprintf("%s calories ", walk.CALORIES)
		duration := fmt.Sprintf(" %s%s%s%s ", color.Yellow, color.Bold, walk.DURATION, color.Reset)
		fmt.Println(number + distance + steps + calories + duration)
	}
}

/* Other Stats */

func (w *WalkData) PrintAllWalks() {
	for i, walk := range w.WALKINGS {
		fmt.Printf(
			"%d. ID:%d  Distance:%s  Duration:%s  Steps:%s  Calories:%s  Date:%s\n",
			i+1,
			walk.ID,
			walk.DISTANCE,
			walk.DURATION,
			walk.STEPS,
			walk.CALORIES,
			walk.DATE,
		)
	}
}

/* Dir Functions */

func (w *WalkData) NewID() int {

	maxID := 0

	for _, book := range w.WALKINGS {
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	return maxID + 1
}

func (w *WalkData) Save() error {

	// Format JSON
	walks, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	// Save
	err = os.WriteFile(config.LocalFile, walks, 0644)
	if err != nil {
		return err
	}

	// Save Backup
	err = os.WriteFile(config.BackupFile, walks, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (w *WalkData) GetUserInput(suggestion Walk) (Walk, error) {
    prompts := []struct {
        label  string
        target *string
    }{
        {"Distance:", &suggestion.DISTANCE},
        {"Duration:", &suggestion.DURATION},
        {"Steps:", &suggestion.STEPS},
        {"Calories:", &suggestion.CALORIES},
        {"Date:", &suggestion.DATE},
    }

    for _, p := range prompts {
        val, err := util.PromptWithSuggestion(p.label, *p.target)
        if err != nil {
            return Walk{}, err
        }
        *p.target = val
    }

    return Walk{
        ID:       suggestion.ID,
        DISTANCE: suggestion.DISTANCE,
        DURATION: suggestion.DURATION,
        STEPS:    suggestion.STEPS,
        CALORIES: suggestion.CALORIES,
        DATE:     suggestion.DATE,
    }, nil
}

func (w *WalkData) ReadFromFile(path string) error {

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

func (w *WalkData) findIndex(id int) (int, error) {
	for i, walk := range w.WALKINGS {
		if walk.ID == id {
			fmt.Println(walk)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (w *WalkData) Undo() bool {
	if len(w.WALKINGS) == 0 {
		fmt.Println("No walks to undo.")
		return false
	}

	lastWalk := w.WALKINGS[len(w.WALKINGS)-1]
	fmt.Println(lastWalk)

	answer, err := util.PromptWithSuggestion("Are you sure you want to delete?", "No")
	if err != nil {
		fmt.Print(err)
	}

	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		w.WALKINGS = w.WALKINGS[:len(w.WALKINGS)-1]

		if err := w.Save(); err != nil {
			fmt.Println(color.Red+"Error saving data:"+color.Reset, err)
			return false
		}

		fmt.Println(color.Yellow + "Last walk removed." + color.Reset)
		return true
	}

	fmt.Println("Undo cancelled.")
	return false
}
