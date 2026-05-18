package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
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

type Store struct {
	Walks []Walk `json:"walks"`
}

/* Main Functions */

func (w *Store) PrintDashboard() {

	// Program information
	fmt.Println(color.PrintBoldBlue("------------------------"))
	fmt.Println(color.PrintBoldBlue("VK-WALKING 1.0"))
	fmt.Println(color.PrintBoldBlue("------------------------"))
	w.PrintLatest()
}

func (w *Store) Add() error {

	newWalk, err := w.promptWalkInput(Walk{})
	if err != nil {
		return err
	}

	newWalk.Id = w.nextID()

	w.Walks = append(w.Walks, newWalk)

	return w.saveToFile()
}

func (w *Store) Update(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := w.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := w.promptWalkInput((w.Walks)[index])
	if err != nil {
		return err
	}

	(w.Walks)[index] = updated

	return w.saveToFile()
}

func (w *Store) Delete(id int) error {

	if id <= 0 {
		return fmt.Errorf("invalid ID: %d", id)
	}

	index, err := w.indexOf(id)
	if err != nil {
		return err
	}

	confirm := util.Confirm()
	if !confirm {
		return fmt.Errorf("Abort")
	}

	w.Walks = append((w.Walks)[:index], (w.Walks)[index+1:]...)

	return w.saveToFile()
}

func (w *Store) PrintAll() {
	for i, walk := range w.Walks {
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

func (w *Store) PrintLatest() {
	for i := 3; i > 0; i-- {

		fmt.Printf("(ID:%s) Miles: %s | Steps: %s | Calories: %s | Time: %s\n",
			color.PrintBoldYellow(strconv.Itoa((w.Walks)[len(w.Walks)-i].Id)),
			(w.Walks)[len(w.Walks)-i].Distance,
			(w.Walks)[len(w.Walks)-i].Steps,
			(w.Walks)[len(w.Walks)-i].Calories,
			(w.Walks)[len(w.Walks)-i].Duration,
		)
	}
}

func (w *Store) Export() error {

	input, err := util.PromptWithSuggestion("Export db to d drive? (y/n) ", "n")
	if err != nil {
		return err
	}

	if input == "y" || input == "yes" {

		if err := util.InitBackupStorage(); err != nil {
			return err
		}

		if err := w.LoadFromFile(config.LocalFile); err != nil {
			return fmt.Errorf("load from file: %w", err)
		}

		finance, err := json.MarshalIndent(w, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(config.BackupFile, finance, 0644); err != nil {
			return err
		}

		fmt.Printf("Database exported to %s\nPress Enter!", config.BackupFile)
		return nil
	}

	fmt.Println("Export canceled!")
	return nil
}

func (w *Store) Import() error {

	input, err := util.PromptWithSuggestion("Import db from d drive? (y/n) ", "n")
	if err != nil {
		return err
	}

	if input == "y" || input == "yes" {

		if err := util.InitBackupStorage(); err != nil {
			return err
		}

		if err := w.LoadFromFile(config.BackupFile); err != nil {
			return fmt.Errorf("load from file: %w", err)
		}

		finance, err := json.MarshalIndent(w, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(config.LocalFile, finance, 0644); err != nil {
			return err
		}

		fmt.Printf("Database imported from %s\nPress Enter!", config.BackupFile)
		return nil
	}

	fmt.Println("Import canceled!")

	return nil
}

/* Dir Functions */

func (w *Store) nextID() int {

	maxID := 0

	for _, book := range w.Walks {
		if book.Id > maxID {
			maxID = book.Id
		}
	}

	return maxID + 1
}

func (w *Store) saveToFile() error {

	copySlice := make([]Walk, len(w.Walks))
	copy(copySlice, w.Walks)
	copyQuotes := Store{Walks: copySlice}

	// Format JSON
	byteValue, err := json.MarshalIndent(copyQuotes, "", "  ")
	if err != nil {
		return err
	}

	// Save local
		if err := os.WriteFile(config.LocalFile, byteValue, 0644); err != nil {
		return err
	}

	// Save Backup
	if err := util.InitBackupStorage(); err != nil {
		fmt.Println(color.Yellow + "Backup init failed: " + err.Error() + color.Reset)
		return nil // or return err, depending on your needs
	}

	if err := os.WriteFile(config.BackupFile, byteValue, 0644); err != nil {
		fmt.Println(color.Yellow + "Backup write failed: " + err.Error() + color.Reset)
		return nil // same decision here
	}

	return nil
}

func (w *Store) promptWalkInput(suggestion Walk) (Walk, error) {
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

func (w *Store) LoadFromFile(path string) error {

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

func (w *Store) indexOf(id int) (int, error) {
	for i, walk := range w.Walks {
		if walk.Id == id {
			fmt.Println(walk)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func (w *Store) Undo() error {
	if len(w.Walks) == 0 {
		return fmt.Errorf("no walks to undo.")
	}

	lastWalk := w.Walks[len(w.Walks)-1]
	fmt.Println(lastWalk)

	confirm := util.Confirm()
	if !confirm {
		return fmt.Errorf("Abort")
	}

	w.Walks = w.Walks[:len(w.Walks)-1]

	if err := w.saveToFile(); err != nil {
		return err
	}
	return nil
}
