
# Task Tracker CLI

Task Tracker is a simple yet powerful command-line interface (CLI) application for managing and tracking your tasks. It allows you to keep track of what you need to do, what you're currently working on, and what you've already completed.

## Features

- Complete task management (Add, Update, Delete)
- Task status tracking (Todo, In Progress, Completed)
- Task listing with different filters
- Persistent storage in JSON format
- Intuitive command-line interface

## Requirements

- Go 1.21 or higher
- No external dependencies required
- Uses only standard Go packages
- Data is stored locally in a JSON file

## Data Structure

Each task contains the following properties:

```json
{
  "id": 1,
  "description": "Task description",
  "status": "todo",
  "createdAt": "2024-01-01T10:00:00Z",
  "updatedAt": "2024-01-01T10:00:00Z"
}
```

## Available Commands

Each command can be used in three different ways:
1. Using full flags
2. Using short flags
3. Direct usage without flags

### Task Management

```bash
# Build the application
go build -o task-cli

# Add a new task
task-cli add -desc "Buy groceries"      # Using full flag
task-cli add -d "Buy groceries"         # Using short flag
task-cli add "Buy groceries"            # Direct usage

# Update an existing task
task-cli update -id 1 -desc "Buy groceries"      # Using full flags
task-cli update -i 1 -d "Buy groceries"         # Using short flags
task-cli update 1 "Buy groceries"               # Direct usage

# Delete a task
task-cli delete -id 1      # Using full flag
task-cli delete -i 1       # Using short flag
task-cli delete 1          # Direct usage
```

### Status Management

```bash
# Mark a task as in progress
task-cli mark-in-progress -id 1      # Using full flag
task-cli mark-in-progress -i 1       # Using short flag
task-cli mark-in-progress 1          # Direct usage

# Mark a task as done
task-cli mark-done -id 1      # Using full flag
task-cli mark-done -i 1       # Using short flag
task-cli mark-done 1          # Direct usage
```

### Task Listing

```bash
# List all tasks
task-cli list                  # Show all tasks

# List tasks by status
task-cli list -status done     # Using full flag
task-cli list -s done          # Using short flag
task-cli list done             # Direct usage

# More status examples
task-cli list todo             # List pending tasks
task-cli list in-progress      # List tasks in progress
```

### Help Command

```bash
# Display general help
task-cli help

# Get help for a specific command
task-cli help -flag add          # Using full flag
task-cli help -f add            # Using short flag
task-cli help add               # Direct usage
```

## Storage

Tasks are stored in a JSON file in the current directory. The file is created automatically if it doesn't exist. The JSON file structure is as follows:

```json
{
  "tasks": [
    {
      "id": 1,
      "description": "Buy groceries",
      "status": "todo",
      "createdAt": "2024-01-01T10:00:00Z",
      "updatedAt": "2024-01-01T10:00:00Z"
    }
  ]
}
```

## Error Handling

The application includes error handling for the following cases:

- Non-existent task ID
- Invalid commands
- File read/write errors
- Invalid JSON format

## Project Constraints

- Written in Go (Golang)
- No external libraries or frameworks used
- Only uses Go standard library packages
- Arguments are passed via command line
- Local storage in JSON file