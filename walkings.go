package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Walk struct {
	Id       int    `json:"id"` // keep as int
	Distance string `json:"distance"`
	Duration string `json:"duration"`
	Calories string `json:"calories"`
	Date     string `json:"date"`
}

type Walks struct {
	Walks []Walk `json:"walks"`
}

/* Main Functions */

func (w *Walks) Add() error {

	walk, err := promptWalkInput(Walk{})
	if err != nil {
		return err
	}
	walk.Id = w.nextID()
	walk.Date = time.Now().Format("02-01-2006")
	w.Walks = append(w.Walks, walk)
	return w.Save()
}

func (w *Walks) Update(idStr string) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}

	index, err := w.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := promptWalkInput(w.Walks[index])
	if err != nil {
		return err
	}

	w.Walks[index] = updated

	return w.Save()
}

func (w *Walks) Delete(idStr string) error {
	id, err := parseID(idStr)
	if err != nil {
		return err
	}

	index, err := w.indexOf(id)
	if err != nil {
		return err
	}

	// Show what will be deleted and ask for confirmation
	e := w.Walks[index]
	fmt.Printf("Delete walk #%d (%s, %s, %s)? (y/n): ",
		e.Id, e.Date, e.Distance, e.Duration)

	answer, err := stdin.ReadString('\n')
	if err != nil || strings.ToLower(strings.TrimSpace(answer)) != "y" {
		return errors.New("delete canceled")
	}

	w.Walks = append(w.Walks[:index], w.Walks[index+1:]...)

	return w.Save()
}

func (w *Walks) Undo() error {
	if len(w.Walks) == 0 {
		return errors.New("no walks to undo")
	}

	// Show what will be removed and ask for confirmation
	last := w.Walks[len(w.Walks)-1]
	fmt.Printf("Remove last walk #%d (%s, %s, %s)? (y/n): ",
		last.Id, last.Date, last.Distance, last.Duration)

	answer, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || strings.ToLower(strings.TrimSpace(answer)) != "y" {
		return errors.New("undo canceled")
	}

	w.Walks = w.Walks[:len(w.Walks)-1]

	return w.Save()
}

/* Helper Functions */

func (w *Walks) nextID() int {

	maxID := 0

	for _, book := range w.Walks {
		if book.Id > maxID {
			maxID = book.Id
		}
	}

	return maxID + 1
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID: %q", s)
	}
	return id, nil
}

func (w *Walks) indexOf(id int) (int, error) {
	for i, walk := range w.Walks {
		if walk.Id == id {
			fmt.Println(walk)
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

var stdin = bufio.NewReader(os.Stdin)

func prompt(label string) (string, error) {
	fmt.Print(label + " ")
	answer, err := stdin.ReadString('\n')
	if err != nil {
		return "", err
	}
	// Accept "4,8" as "4.8"
	return strings.ReplaceAll(strings.TrimSpace(answer), ",", "."), nil
}

func promptValid(label, current string, validate func(string) error) (string, error) {
	for {
		if current != "" {
			fmt.Printf("%s [%s]: ", label, current)
		} else {
			fmt.Printf("%s ", label)
		}

		answer, err := stdin.ReadString('\n')
		if err != nil {
			return "", err
		}

		val := strings.ReplaceAll(strings.TrimSpace(answer), ",", ".")
		if val == "" {
			val = current
		}

		if err := validate(val); err != nil {
			fmt.Println(Red + err.Error() + Reset)
			continue
		}
		return val, nil
	}
}

func promptWalkInput(w Walk) (Walk, error) {
	var err error

	if w.Distance, err = promptValid("Distance (mi):", w.Distance, validateNumber); err != nil {
		return Walk{}, err
	}
	if w.Duration, err = promptValid("Duration (H:MM:SS):", w.Duration, validateDuration); err != nil {
		return Walk{}, err
	}
	if w.Calories, err = promptValid("Calories:", w.Calories, validateInt); err != nil {
		return Walk{}, err
	}

	return w, nil
}

func validateNumber(s string) error {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil || n <= 0 {
		return errors.New("enter a positive number")
	}
	return nil
}

func validateInt(s string) error {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return errors.New("enter a whole number")
	}
	return nil
}

func validateDuration(s string) error {
	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return errors.New("use format H:MM:SS or H:MM")
	}

	total := 0
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return errors.New("use format H:MM:SS or H:MM")
		}
		// Every part except the first must be a valid minute/second value
		if i > 0 && n > 59 {
			return errors.New("minutes and seconds must be 0-59")
		}
		total += n
	}

	if total == 0 {
		return errors.New("duration can't be zero")
	}
	return nil
}
