package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
	Italic = "\033[3m"
)

const (
	driveLabel      = "VK DATA"
	driveMountPoint = "/media/veikko/VK DATA"
	file            = "walkings.json"
	baseDB          = "WALKINGS"
	baseLocal       = "DATABASES"
	// baseLocal = "DEVELOPMENT"
)

var LocalFile = filepath.Join(baseLocal, baseDB, file)
var BackupFile = filepath.Join(driveMountPoint, baseLocal, baseDB, file)

const usage = Bold + "walks" + Reset + ` - walking tracker

` + Bold + "Usage:" + Reset + `

  walks [command]

` + Bold + "Commands:" + Reset + `

  ` + Cyan + "add, a" + Reset + `            add a new walk
  ` + Cyan + "update, u <id>" + Reset + `    edit the walk with this ID
  ` + Cyan + "delete, d <id>" + Reset + `    delete the walk with this ID
  ` + Cyan + "undo" + Reset + `              remove the last walk
  ` + Cyan + "list, l" + Reset + `           show the last 10 walks
  ` + Cyan + "help" + Reset + `              show this message

Run without a command to show the report.

`

func main() {

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			fmt.Print(usage)
			return
		}
	}

	w := Walks{}
	if err := w.Load(LocalFile); err != nil {
		log.Fatalf("load from file: %v", err)
	}

	var err error
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "a", "add":
			err = w.Add()
		case "u", "update":
			err = w.Update(os.Args[2])
		case "d", "delete":
			err = w.Delete(os.Args[2])
		case "undo":
			err = w.Undo()
		default:
			fmt.Printf("%serror: unknown command: %s%s\n\n", Red, os.Args[1], Reset)
			fmt.Print(usage)
			os.Exit(1)
		}
	}

	w.Report()

	if err != nil {
		fmt.Printf("%serror: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

}
