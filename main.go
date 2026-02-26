package main

import (
	"fmt"
	"log"
	"os"
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/config"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func main() {

	if err := util.CreateFilesAndFolders(); err != nil {
		fmt.Println("Error creating files/folders:", err)
		os.Exit(1)
	}

	// Initialize Walkings database
	w := db.Walkings{}

	// Reload Database
	err := w.ReadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	// Start
	cmd.CommandLine(&w)
}
