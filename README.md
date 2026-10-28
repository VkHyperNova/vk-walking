# VK-Walking

A small command-line walking tracker written in Go. Log each walk's distance, duration and calories, see your latest walks in a clean table, and keep everything in a single local JSON file with a backup copy on a USB drive.

## Features

- **Quick entry**: answer three prompts to log a walk; the date is filled in automatically
- **Input validation**: bad values are rejected and you are asked again
- **Edit and delete by ID**: fix a typo in an old walk or remove one, with a confirmation before anything is deleted
- **Undo**: remove the last walk if you added it by mistake
- **Latest walks table**: the last 10 walks are shown every time you run the program
- **Plain JSON storage**: no database required
- **Backup**: data is also copied to the `VK DATA` drive when it is connected

## Requirements

- [Go](https://go.dev/dl/) 1.25 or newer
- Linux or Windows (the screen is cleared with `clear` / `cls`)
- Optional: a drive labelled `VK DATA` for backups

## Installation

```bash
git clone https://github.com/VkHyperNova/vk-walking.git
cd vk-walking
go build
```

This produces a `vk-walking` binary. You can also run it without building:

```bash
go run . add
```

## Usage

```
vk-walking [command]
```

| Command | Description |
|---|---|
| *(none)* | Show the last 10 walks |
| `add`, `a` | Add a new walk |
| `update <id>`, `u <id>` | Edit the walk with this ID |
| `delete <id>`, `d <id>` | Delete the walk with this ID (asks for confirmation) |
| `undo` | Remove the last walk (asks for confirmation) |
| `help`, `-h`, `--help` | Show usage |

### Adding a walk

```
$ vk-walking add
Distance (mi): 2.4
Duration (H:MM:SS): 0:48:30
Calories: 210
```

The ID is assigned automatically and the date is set to today.

| Field | Format | Example |
|---|---|---|
| Distance | positive number, in miles (`,` is accepted as a decimal separator) | `2.4` or `2,4` |
| Duration | `H:MM:SS` or `H:MM` | `1:12:23`, `1:23` |
| Calories | positive whole number | `210` |

If a value is invalid, an error is shown and you are asked again.

### Editing a walk

```bash
vk-walking update 42
```

Each prompt shows the current value in brackets. Press Enter to keep it, or type a new value. The date is not changed when you edit a walk.

### Deleting a walk

```bash
vk-walking delete 42
```

The walk is shown and you must answer `y` to confirm. Anything else cancels. The IDs are the first column of the table.

### The table

```
VK-WALKING 2
------------------------
ID   DATE         DISTANCE   DURATION   CALORIES
--   ----         --------   --------   --------
41   20-10-2026   2.4        0:48:30    210
42   22-10-2026   3.1        1:02:10    265
```

Dates use the `DD-MM-YYYY` format.

## Data storage

Walks are stored as JSON in `<baseLocal>/WALKINGS/walkings.json`, relative to the directory you run the program from. The file is created on the first save, and a missing file simply means an empty list.

```json
{
  "walks": [
    {
      "id": 42,
      "distance": "3.1",
      "duration": "1:02:10",
      "calories": "265",
      "date": "22-10-2026"
    }
  ]
}
```

`baseLocal` is a constant in `main.go`. It is currently set to `DEVELOPMENT`; switch it to `DATABASES` for real data.

## Backup

On each save, the data is also written to `/media/veikko/VK DATA/<baseLocal>/WALKINGS/walkings.json` if the `VK DATA` drive can be mounted. The drive label and mount point are constants at the top of `main.go`; change them to match your setup.

## Project structure

```
├── main.go         # Entry point, command dispatch, paths and constants
├── walkings.go     # Data model and commands (add, update, delete, undo, input validation)
├── report.go       # Table output
├── storage.go      # Loading and saving JSON, drive backup
├── go.mod
└── go.sum
```

## Dependencies

- [`peterh/liner`](https://github.com/peterh/liner)
- [`VkHyperNova/vkmount`](https://github.com/VkHyperNova/vkmount) for mounting the backup drive

## License

No license specified. All rights reserved by the author unless stated otherwise.