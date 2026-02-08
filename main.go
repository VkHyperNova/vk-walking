package main

import (
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func main() {

	util.CreateNecessaryFiles()

	// Initialize Walkings database
	books := db.Walkings{}

	// Start
	cmd.CommandLine(&books)
}




