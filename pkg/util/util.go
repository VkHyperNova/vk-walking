package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"

	"github.com/peterh/liner"
)

func CreateLocalFiles() error {
	// Create necessary files
	err := CreateLocalFolder(config.FolderName)
	if err != nil {
		log.Fatalf("Fatal error: failed to create necessary files: %v", err)
	}

	// Create necessary files
	err = CreateLocalJSON(config.LocalPath)
	if err != nil {
		log.Fatalf("Fatal error: failed to create necessary files: %v", err)
	}

	config.LocalSave = true
	return nil
}

func CreateDDriveFiles() {
	mounted := HardDriveMountCheck()
	if !mounted {
		input := Input("Do you want to continue? (y/n)")
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			os.Exit(0)
		}
	} else {
		CreateDDriveFolder()
	}
}

func CreateLocalJSON(fileName string) error {

	// Make local json File
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		if err := os.WriteFile(fileName, []byte(`{"walkings": []}`), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", fileName, err)
		}
	}

	return nil
}

func CreateLocalFolder(folderName string) error {
	// Make local Folder
	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", folderName, err)
	}

	return nil
}

func CreateDDriveFolder() error {
	// Make backup folder in another drive
	if err := os.MkdirAll(config.BackupFolder+config.FolderName, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", config.BackupFolder, err)
	}
	config.DDriveSave = true
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

func HardDriveMountCheck() bool {
	if runtime.GOOS != "linux" {
		fmt.Println("This program only works on Linux.")
		return false
	}

	mountPoint := "/media/veikko/VK\\040DATA" // match /proc/mounts format

	file, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Println("Cannot open /proc/mounts:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[1] == mountPoint {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning /proc/mounts:", err)
		return false
	}

	fmt.Println(color.Red + "\nVK DATA is NOT mounted" + color.Reset)
	return false
}

func Input(prompt string) string {

	line := liner.NewLiner()
	defer line.Close()

	userInput, err := line.Prompt(prompt)
	if err != nil {
		panic(err)
	}
	return userInput
}

func PrintArrow() {
	arrow := color.Yellow + "< Local Only => " + color.Reset

	if config.DDriveSave {
		arrow = color.Green + "< Local/DDrive => " + color.Reset
	}

	fmt.Print(arrow)
}

func ReadCommand() (string, int, bool) {
	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Input error:", err)
		return "", 0, false
	}

	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	if len(parts) == 0 {
		return "", 0, false
	}

	command := strings.ToLower(parts[0])
	id := 0

	if len(parts) > 1 {
		if _, err := fmt.Sscan(parts[1], &id); err != nil {
			fmt.Println("Invalid ID")
			return "", 0, false
		}
	}

	return command, id, true
}

func PressAnyKey() {
	fmt.Print()
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
