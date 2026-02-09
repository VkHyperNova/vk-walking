package main

import (
	"log"
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/config"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func main() {
	util.CreateDDriveFiles()

	util.CreateLocalFiles()

	// Initialize Walkings database
	w := db.Walkings{}

	// Reload Database
	err := w.ReadFromFile(config.LocalPath)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	// Start
	cmd.CommandLine(&w)
}
