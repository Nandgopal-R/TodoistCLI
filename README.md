# TodoistCLI

A simple CLI-based task management app built in Go using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for interactive terminal applications.

## Features

- Add, list, toggle completion, and delete tasks.
- Tasks stored in a CSV file (`tasks.csv`).
- Interactive terminal interface with cursor navigation.

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) (v1.18 or higher)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [Lip Gloss](https://github.com/charmbracelet/lipgloss)

### Install Dependencies
Run `go mod tidy` in the project directory.

### Running the App
Run `go run main.go` to start the application.

## Usage

### Commands

- **Add a Task**: Press `+` to enter task addition mode.
- **View Tasks**: Press `l` to list all tasks.
- **Toggle Task Completion**: While viewing tasks, press `enter` to toggle the completion status of the selected task.
- **Delete Tasks**:
  - Press `d` to delete completed tasks.
  - Press `d` while in the task list to delete tasks.
- **Quit**: Press `ctrl+c` or `q` to quit the application.

### Key Bindings

- `+` to add a task.
- `l` to list tasks.
- `d` to delete tasks.
- `esc` to go back or cancel the current operation.
- `enter` to confirm actions (e.g., save, toggle completion, delete).
