
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

### Task Management

```bash
# Build the application
go build -o task-cli

# Add a new task
./task-cli add "Buy groceries"

# Update an existing task
./task-cli update 1 "Buy groceries and cook dinner"

# Delete a task
./task-cli delete 1
```

### Status Management

```bash
# Mark a task as in progress
task-cli mark-in-progress 1

# Mark a task as done
task-cli mark-done 1
```

### Task Listing

```bash
# List all tasks
task-cli list

# List tasks by status
task-cli list done        # Completed tasks
task-cli list todo        # Pending tasks
task-cli list in-progress # Tasks in progress
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