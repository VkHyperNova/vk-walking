package main

import "fmt"

func main() {
	fmt.Println("walking")
}

type Walk struct {
	ID       int     `json:"id"`
	NAME     string  `json:"name"`
	DISTANCE float64 `json:"distance"`
	DURATION string  `json:"duration"`
	PACE     string  `json:"pace"`
	STEPS    int     `json:"steps"`
	CALORIES int     `json:"calories"`

	DATE     string  `json:"date"`
}

type Walkings struct {
	WALKINGS []Walk `json:"walkings"`
}
