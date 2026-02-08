package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"vk-walking/pkg/config"
	"github.com/peterh/liner"
)

func CreateNecessaryFiles() error {

	// Make local Folder
	if err := os.MkdirAll(config.FolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", config.FolderName, err)
	}

	// Make local File
	if _, err := os.Stat(config.LocalPath); os.IsNotExist(err) {
		if err := os.WriteFile(config.LocalPath, []byte(`{"walkings": []}`), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", config.FileName, err)
		}
	}

	// Make backup folder in another drive
	if err := os.MkdirAll(config.BackupFolder + config.FolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", config.BackupFolder, err)
	}

	return nil
}

func ClearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

func PromptWithSuggestion(name string, suggestion string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", suggestion, -1)
	if err != nil {
		panic(err)
	}

	return input
}