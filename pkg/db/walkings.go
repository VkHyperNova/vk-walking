package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"
	"vk-walking/pkg/util"
)

type Walk struct {
	Id       int    `json:"id"` // keep as int
	Distance string `json:"distance"`
	Duration string `json:"duration"`
	Steps    string `json:"steps"`
	Calories string `json:"calories"`
	Date     string `json:"date"`
}

type WalkData struct {
	Data []Walk `json:"data"`
}

/* Main Functions */

func (w *WalkData) PrintCLI() {

	// Program information
	fmt.Println(color.PrintBoldBlue("------------------------"))
	fmt.Println(color.PrintBoldBlue("VK-WALKING 1.0"))
	fmt.Println(color.PrintBoldBlue("------------------------"))
	w.PrintLatest()
}

func (w *WalkData) Add() error {

	newWalk, err := w.GetUserInput(Walk{})
	if err != nil {
		return err
	}

	newWalk.Id = w.NewID()

	w.Data = append(w.Data, newWalk)

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

	updated, err := w.GetUserInput((w.Data)[index])
	if err != nil {
		return err
	}

	(w.Data)[index] = updated

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

	w.Data = append((w.Data)[:index], (w.Data)[index+1:]...)

	return w.Save()
}

func (w *WalkData) PrintAllWalks() {
	for i, walk := range w.Data {
		fmt.Printf(
			"%d. ID:%d  Distance:%s  Duration:%s  Steps:%s  Calories:%s  Date:%s\n",
			i+1,
			walk.Id,
			walk.Distance,
			walk.Duration,
			walk.Steps,
			walk.Calories,
			walk.Date,
		)
	}
}

func(w *WalkData) PrintLatest() {
	for i := 1; i < 4; i++ {

		fmt.Printf("(ID:%s) Miles: %s | Steps: %s | Calories: %s | Time: %s\n",
			color.PrintBoldYellow(strconv.Itoa((w.Data)[len(w.Data) - i].Id)),
			(w.Data)[len(w.Data) - i].Distance,
			(w.Data)[len(w.Data) - i].Steps,
			(w.Data)[len(w.Data) - i].Calories,
			(w.Data)[len(w.Data) - i].Duration,
		)
	}
}

/* Dir Functions */

func (w *WalkData) NewID() int {

	maxID := 0

	for _, book := range w.Data {
		if book.Id > maxID {
			maxID = book.Id
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
		{"Distance:", &suggestion.Distance},
		{"Duration:", &suggestion.Duration},
		{"Steps:", &suggestion.Steps},
		{"Calories:", &suggestion.Calories},
		{"Date:", &suggestion.Date},
	}

	for _, p := range prompts {
		val, err := util.PromptWithSuggestion(p.label, *p.target)
		if err != nil {
			return Walk{}, err
		}
		*p.target = val
	}

	return Walk{
		Id:       suggestion.Id,
		Distance: suggestion.Distance,
		Duration: suggestion.Duration,
		Steps:    suggestion.Steps,
		Calories: suggestion.Calories,
		Date:     suggestion.Date,
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
	for i, walk := range w.Data {
		if walk.Id == id {
			fmt.Println(walk)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (w *WalkData) Undo() bool {
	if len(w.Data) == 0 {
		fmt.Println("No walks to undo.")
		return false
	}

	lastWalk := w.Data[len(w.Data)-1]
	fmt.Println(lastWalk)

	answer, err := util.PromptWithSuggestion("Are you sure you want to delete?", "No")
	if err != nil {
		fmt.Print(err)
	}

	if strings.ToLower(answer) == "y" || strings.ToLower(answer) == "yes" {
		w.Data = w.Data[:len(w.Data)-1]

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
