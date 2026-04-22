package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"vk-walking/pkg/color"
	"vk-walking/pkg/config"

	"github.com/peterh/liner"
)

func TimeToSeconds(duration string) int {
	parts := strings.Split(duration, ":")
	h, m, sec := 0, 0, 0
	if len(parts) == 3 {
		h, _ = strconv.Atoi(parts[0])
		m, _ = strconv.Atoi(parts[1])
		sec, _ = strconv.Atoi(parts[2])
	} else if len(parts) == 2 {
		m, _ = strconv.Atoi(parts[0])
		sec, _ = strconv.Atoi(parts[1])
	}
	return h*3600 + m*60 + sec
}

func ClearScreen() {

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

func PromptWithSuggestion(name string, suggestion string) (string, error) {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", suggestion, -1)
	if err != nil {
		return input, err
	}

	return input, nil
}

func isMounted() bool {
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

func ensureFile(path string, content string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", path, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", path, err)
		}
	}

	return nil
}

func InitStorage() error {

	if err := ensureFile(config.LocalFile, config.DefaultContent); err != nil {
		return err
	}

	if !isMounted() {
		input := Confirm()
		if !input {
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
	} else {
		if err := ensureFile(config.BackupFile, config.DefaultContent); err != nil {
			return err
		}
	}

	return nil
}

func Confirm() bool {

	input, err := PromptWithSuggestion("Do you want to continue (y/n): ", "n")
	if err != nil {
		fmt.Print(err)
		return false
	}

	if input == "n" || input == "no" || input == "q" {
		fmt.Println(color.Red, "Aborted!", color.Reset)
		return false
	}
	return true
}
