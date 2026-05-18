package main

import (
	"log"
	"vk-walking/pkg/cmd"
	"vk-walking/pkg/config"
	"vk-walking/pkg/db"
	"vk-walking/pkg/util"
)

func main() {
	if err := util.InitLocalStorage(); err != nil {
		log.Fatalf("Error creating files/folders: %v", err)

	}

	store := db.Store{}
	err := store.LoadFromFile(config.LocalFile)
	if err != nil {
		log.Fatalf("Fatal error: failed to load walkings database: %v", err)
	}

	cmd.Run(&store)
}
