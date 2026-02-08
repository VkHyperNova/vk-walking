package main

import (
	"log"
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/config"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func main() {

	// Create necessary files
	err := util.CreateNecessaryFiles()
	if err != nil {
		log.Fatalf("Fatal error: failed to create necessary files: %v", err)
	}

	// Initialize Walkings database
	books := db.Walkings{}

	// Reload Database
	err = books.ReadFromFile(config.LocalPath)
	if err != nil {
		log.Fatalf("Fatal error: failed to load books database: %v", err)
	}

	// Start
	cmd.CommandLine(&books)
}




