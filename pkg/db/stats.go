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
	field string
	less func(a, b Walk) bool
}

func (w *Store) printTopN(number int, cfg sortConfig) {
	if number == 0 {
		number = 5
	}
	fmt.Println(color.PrintBoldYellow("\n" + strings.ToLower(cfg.field) + " " + strconv.Itoa(number)))
	walks := w.sortBy(cfg.less)
	printRows(walks, cfg.field, number)
}

func printRows(walks []Walk, name string, n int) {
	fmt.Print(color.PrintBoldBlue("\n" + name + "\n"))
	for i := 0; i < n; i++ {
		w := walks[i]
		miles, _ := strconv.ParseFloat(w.Distance, 64)
		km := miles * 1.60934
		highlight := func(field, value string) string {
			if name == field {
				return color.PrintBoldYellow(value)
			}
			return value
		}
		fmt.Printf("(ID:%d) Miles: %s (%.2f km) | Steps: %s | Calories: %s | Time: %s\n",
			w.Id,
			highlight("Distance", w.Distance),
			km,
			highlight("Steps", w.Steps),
			highlight("Calories", w.Calories),
			highlight("Duration", w.Duration),
		)
	}
}

func (w *Store) sortBy(less func(a, b Walk) bool) []Walk {
	walks := make([]Walk, len(w.Walks))
	copy(walks, w.Walks)
	sort.Slice(walks, func(i, j int) bool {
		return less(walks[i], walks[j])
	})
	return walks
}

func (w *Store) PrintDistance(n int) {
	w.printTopN(n, sortConfig{
		field: "Distance",
		less: func(a, b Walk) bool {
			af, _ := strconv.ParseFloat(a.Distance, 64)
			bf, _ := strconv.ParseFloat(b.Distance, 64)
			return af > bf
		},
	})
}

func (w *Store) PrintSteps(n int) {
	w.printTopN(n, sortConfig{
		field: "Steps",
		less: func(a, b Walk) bool {
			ai, _ := strconv.Atoi(a.Steps)
			bi, _ := strconv.Atoi(b.Steps)
			return ai > bi
		},
	})
}

func (w *Store) PrintCalories(n int) {
	w.printTopN(n, sortConfig{
		field: "Calories",
		less: func(a, b Walk) bool {
			ai, _ := strconv.Atoi(a.Calories)
			bi, _ := strconv.Atoi(b.Calories)
			return ai > bi
		},
	})
}

func (w *Store) PrintDuration(n int) {
	w.printTopN(n, sortConfig{
		field: "Duration",
		less: func(a, b Walk) bool {
			return util.TimeToSeconds(a.Duration) > util.TimeToSeconds(b.Duration)
		},
	})
}

func (w *Store) PrintStats(n int) {
	if n == 0 {
		n = 5
	}
	fmt.Println(color.PrintBoldYellow("\nstats " + strconv.Itoa(n)))
	w.PrintDistance(n)
	w.PrintSteps(n)
	w.PrintCalories(n)
	w.PrintDuration(n)
}
