package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"
)

/* Report */

func (w *Walks) Report() {

	clearScreen()

	fmt.Println(Bold + Blue + "VK-WALKING 2" + Reset)
	fmt.Println(("------------------------"))

	w.lastTenEntries()
}

/* Helper Functions */

func (w *Walks) lastTenEntries() {
	if len(w.Walks) == 0 {
		fmt.Println("No walks yet.")
		return
	}

	// Slice the last 10 (or fewer if there aren't 10 yet)
	start := 0
	if len(w.Walks) > 10 {
		start = len(w.Walks) - 10
	}
	entries := w.Walks[start:]

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "ID\tDATE\tDISTANCE\tDURATION\tCALORIES")
	fmt.Fprintln(tw, "--\t----\t--------\t--------\t--------")
	for _, e := range entries {
		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\t%s\n",
			e.Id, e.Date, e.Distance, e.Duration, e.Calories)
	}
	tw.Flush()
}

func clearScreen() {

	var c *exec.Cmd

	if runtime.GOOS == "windows" {
		c = exec.Command("cmd", "/c", "cls")
	} else {
		c = exec.Command("clear")
	}

	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}
