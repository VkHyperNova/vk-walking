package db

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/util"
)

/* Top Stats */

type sortConfig struct {
	name string
	less func(a, b Walk) bool
}

func (w *WalkData) printSorted(number int, cfg sortConfig) {
	if number == 0 {
		number = 5
	}
	fmt.Println(color.PrintBoldYellow("\n" + strings.ToLower(cfg.name) + " " + strconv.Itoa(number)))
	data := w.sorted(cfg.less)
	printTopTen(data, cfg.name, number)
}

func printTopTen(sortedData []Walk, name string, number int) {
	fmt.Print(color.PrintBoldBlue("\n" + name + "\n"))
	for i := 0; i < number; i++ {
		w := sortedData[i]
		distanceToFloat, _ := strconv.ParseFloat(w.Distance, 64)
		distanceInKilometer := distanceToFloat * 1.60934
		highlight := func(field, value string) string {
			if name == field {
				return color.PrintBoldYellow(value)
			}
			return value
		}
		fmt.Printf("(ID:%d) Miles: %s (%.2f km) | Steps: %s | Calories: %s | Time: %s\n",
			w.Id,
			highlight("Distance", w.Distance),
			distanceInKilometer,
			highlight("Steps", w.Steps),
			highlight("Calories", w.Calories),
			highlight("Duration", w.Duration),
		)
	}
}

func (w *WalkData) sorted(less func(a, b Walk) bool) []Walk {
	sortedData := make([]Walk, len(w.Data))
	copy(sortedData, w.Data)
	sort.Slice(sortedData, func(i, j int) bool {
		return less(sortedData[i], sortedData[j])
	})
	return sortedData
}

func (w *WalkData) PrintDistance(number int) {
	w.printSorted(number, sortConfig{
		name: "Distance",
		less: func(a, b Walk) bool {
			af, _ := strconv.ParseFloat(a.Distance, 64)
			bf, _ := strconv.ParseFloat(b.Distance, 64)
			return af > bf
		},
	})
}

func (w *WalkData) PrintSteps(number int) {
	w.printSorted(number, sortConfig{
		name: "Steps",
		less: func(a, b Walk) bool {
			ai, _ := strconv.Atoi(a.Steps)
			bi, _ := strconv.Atoi(b.Steps)
			return ai > bi
		},
	})
}

func (w *WalkData) PrintCalories(number int) {
	w.printSorted(number, sortConfig{
		name: "Calories",
		less: func(a, b Walk) bool {
			ai, _ := strconv.Atoi(a.Calories)
			bi, _ := strconv.Atoi(b.Calories)
			return ai > bi
		},
	})
}

func (w *WalkData) PrintDuration(number int) {
	w.printSorted(number, sortConfig{
		name: "Duration",
		less: func(a, b Walk) bool {
			return util.TimeToSeconds(a.Duration) > util.TimeToSeconds(b.Duration)
		},
	})
}

func (w *WalkData) PrintStats(number int) {
	if number == 0 {
		number = 5
	}
	fmt.Println(color.PrintBoldYellow("\nstats " + strconv.Itoa(number)))
	w.PrintDistance(number)
	w.PrintSteps(number)
	w.PrintCalories(number)
	w.PrintDuration(number)
}