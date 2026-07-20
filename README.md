# vk-walking

A command-line walking tracker written in Go. Log and manage your walks through an interactive terminal shell, with data persisted locally.

## Tech Stack

- [Go](https://go.dev/)
- Local file-based storage (no external database required)

## Project Structure

```
├── DATABASES/     # Local data storage for walking records
├── pkg/
│   ├── cmd/         # Command handling / interactive shell logic
│   ├── config/       # Configuration (e.g. local storage file path)
│   ├── db/            # Walking data model (Store) and load/save logic
│   └── util/           # Utility helpers (local storage init, etc.)
├── main.go            # Entry point
├── go.mod
└── go.sum
```

## How It Works

On startup, `vk-walking`:

1. Initializes local storage (`util.InitLocalStorage`).
2. Loads existing walking data from a local file (`db.Store.LoadFromFile`).
3. Starts an interactive command shell (`cmd.Run`) where you can log and manage your walks.

## Getting Started

### Prerequisites

- Go installed

### Installation

```bash
git clone https://github.com/VkHyperNova/vk-walking.git
cd vk-walking
go build -o vk-walking .
```

### Usage

Run the built binary:

```bash
./vk-walking
```

Or run directly with Go:

```bash
go run main.go
```

This launches an interactive shell where you can log walks and view your walking history, persisted locally between runs.

## License

No license specified — all rights reserved by the author unless stated otherwise.
