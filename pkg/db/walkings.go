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
	ID       int     `json:"id"`
	DISTANCE float64 `json:"distance"`
	DURATION string  `json:"duration"`
	STEPS    int     `json:"steps"`
	CALORIES int     `json:"calories"`
	DATE     int     `json:"date"`
}

type Walkings struct {
	WALKINGS []Walk `json:"walkings"`
}

/* Main Functions */

func (w *Walkings) PrintCLI() {

	// Program information
	fmt.Println(color.Cyan + "VK-WALKING 1.0" + color.Reset)
	fmt.Println(color.Cyan + "------------------------" + color.Reset)

	w.PrintTopDistance()
	w.PrintTopSteps()
	w.PrintTopCalories()
	w.PrintTopDuration()

	w.PrintOverallStats()
	w.PrintStatsByYear()
}

func (w *Walkings) Add() error {

	newWalk, err := w.GetUserInput(Walk{})
	if err != nil {
		return err
	}

	newWalk.ID = w.NewID()

	w.WALKINGS = append(w.WALKINGS, newWalk)

	return w.Save()
}

func (w *Walkings) Update(id int) error {

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

func (w *Walkings) Delete(id int) error {

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

// printTop is a helper function that handles the common logic for all PrintTop* functions.
// It takes a label for the header, a less function to determine sort order,
// and a highlight function to format each walk's output line.
func (w *Walkings) printTop(label string, less func(a, b Walk) bool, highlight func(walk Walk) string) {
    // Copy the slice to avoid mutating the original data
    sorted := make([]Walk, len(w.WALKINGS))
    copy(sorted, w.WALKINGS)

    // Sort the copy using the provided comparison function
    sort.Slice(sorted, func(i, j int) bool {
        return less(sorted[i], sorted[j])
    })

    // Cap at 5 results, or less if the slice is smaller
    n := 5
    if len(sorted) < n {
        n = len(sorted)
    }

    // Print the section header
    fmt.Print(color.Blue + color.Bold + "\n" + label + "\n" + color.Reset)

    // Print each walk using the provided highlight/format function
    for i := 0; i < n; i++ {
        fmt.Println(highlight(sorted[i]))
    }
}

// PrintTopDistance prints the top 5 walks sorted by distance (descending).
// Distance is highlighted in the output.
func (w *Walkings) PrintTopDistance() {
    w.printTop("Top Distance", func(a, b Walk) bool {
        return a.DISTANCE > b.DISTANCE
    }, func(walk Walk) string {
        return fmt.Sprintf("ID: %d %s%.2f miles(%.2f km)%s | %d steps | %d calories  %s ",
            walk.ID, color.Yellow+color.Bold, walk.DISTANCE, walk.DISTANCE*1.60934, color.Reset,
            walk.STEPS, walk.CALORIES, walk.DURATION)
    })
}

// PrintTopSteps prints the top 5 walks sorted by step count (descending).
// Step count is highlighted in the output.
func (w *Walkings) PrintTopSteps() {
    w.printTop("Top Steps", func(a, b Walk) bool {
        return a.STEPS > b.STEPS
    }, func(walk Walk) string {
        return fmt.Sprintf("ID: %d %.2f miles(%.2f km) | %s%s%d%s steps | %d calories  %s ",
            walk.ID, walk.DISTANCE, walk.DISTANCE*1.60934,
            color.Yellow, color.Bold, walk.STEPS, color.Reset,
            walk.CALORIES, walk.DURATION)
    })
}

// PrintTopCalories prints the top 5 walks sorted by calories burned (descending).
// Calorie count is highlighted in the output.
func (w *Walkings) PrintTopCalories() {
    w.printTop("Top Calories", func(a, b Walk) bool {
        return a.CALORIES > b.CALORIES
    }, func(walk Walk) string {
        return fmt.Sprintf("ID: %d %.2f miles(%.2f km) | %d steps | %s%s%d%s calories  %s ",
            walk.ID, walk.DISTANCE, walk.DISTANCE*1.60934,
            walk.STEPS, color.Yellow, color.Bold, walk.CALORIES, color.Reset,
            walk.DURATION)
    })
}

// PrintTopDuration prints the top 5 walks sorted by duration (descending).
// Duration is highlighted in the output.
func (w *Walkings) PrintTopDuration() {
    w.printTop("Top Duration", func(a, b Walk) bool {
        return a.DURATION > b.DURATION
    }, func(walk Walk) string {
        return fmt.Sprintf("ID: %d %.2f miles(%.2f km) | %d steps | %d calories  %s%s%s%s ",
            walk.ID, walk.DISTANCE, walk.DISTANCE*1.60934,
            walk.STEPS, walk.CALORIES,
            color.Yellow, color.Bold, walk.DURATION, color.Reset)
    })
}

/* Other Stats */

func (w *Walkings) PrintAllWalks() {
	for i, walk := range w.WALKINGS {
		fmt.Printf(
			"%d. ID:%d  Distance:%.2f  Duration:%s  Pace:%s  Steps:%d  Calories:%d  Date:%d\n",
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

func (w *Walkings) PrintOverallStats() {

	walksCount := len(w.WALKINGS)

	totalmiles := 0.0
	totalsteps := 0
	totalcalories := 0

	for _, walk := range w.WALKINGS {
		totalmiles += walk.DISTANCE
		totalsteps += walk.STEPS
		totalcalories += walk.CALORIES
	}
	fmt.Printf("\nTotal walks count: %d", walksCount)
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

/* Dir Functions */

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

func (w *Walkings) GetUserInput(oldWalk Walk) (Walk, error) {

	// Get Data (strings)
	distanceStr := util.PromptWithSuggestion("Distance", strconv.FormatFloat(oldWalk.DISTANCE, 'f', 2, 64))
	duration := util.PromptWithSuggestion("Duration", oldWalk.DURATION)
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
		DISTANCE: distance,
		DURATION: duration,
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

func (w *Walkings) findIndex(id int) (int, error) {
	for i, walk := range w.WALKINGS{
		if walk.ID == id {
			fmt.Println(walk)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (w *Walkings) Undo() bool {
	if len(w.WALKINGS) == 0 {
		fmt.Println("No walks to undo.")
		return false
	}

	lastWalk := w.WALKINGS[len(w.WALKINGS)-1]
	fmt.Println(lastWalk)

	answer := strings.ToLower(util.PromptWithSuggestion("Are you sure you want to delete?", "No"))

	if answer == "y" || answer == "yes" {
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

