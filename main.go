package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

const (
	StatusTodo       = "todo"
	StatusInProgress = "in-progress"
	StatusDone       = "done"
)

const dataFile = "tasks.json"

func processArgs(args []string) []string {
	if len(args) < 2 {
		return args
	}

	subcommand := args[1]
	subcommandArgs := args[2:]

	if len(subcommandArgs) == 0 {
		return args
	}

	if strings.HasPrefix(subcommandArgs[0], "-") {
		return args
	}

	switch subcommand {
	case "add":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-desc", subcommandArgs[0]}
		}
	case "update":
		if len(subcommandArgs) >= 2 {
			return []string{args[0], subcommand, "-id", subcommandArgs[0], "-desc", subcommandArgs[1]}
		}
	case "delete", "mark-in-progress", "mark-done":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-id", subcommandArgs[0]}
		}
	case "list":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-status", subcommandArgs[0]}
		}
	case "help":
		if len(subcommandArgs) > 0 {
			return []string{args[0], subcommand, "-flag", subcommandArgs[0]}
		}
	}

	return args
}

func main() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "Task description")
	addCmd.StringVar(addDescription, "d", "", "Task description (shorthand)")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "Task ID")
	updateDescription := updateCmd.String("desc", "", "New task description")
	updateCmd.IntVar(updateID, "i", 0, "Task ID (shorthand)")
	updateCmd.StringVar(updateDescription, "d", "", "New task description (shorthand)")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Task ID")
	deleteCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	markInProgressCmd := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	markInProgressID := markInProgressCmd.Int("id", 0, "Task ID")
	markInProgressCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	markDoneCmd := flag.NewFlagSet("mark-done", flag.ExitOnError)
	markDoneID := markDoneCmd.Int("id", 0, "Task ID")
	markDoneCmd.IntVar(deleteID, "i", 0, "Task ID (shorthand)")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listStatus := listCmd.String("status", "", "Filter tasks by status (todo, in-progress, done)")
	listCmd.StringVar(listStatus, "s", "", "Filter tasks by status (shorthand)")

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)
	helpFlag := helpCmd.String("flag", "", "Command to display help for")
	helpCmd.StringVar(helpFlag, "f", "", "Command to display help for (shorthand)")

	os.Args = processArgs(os.Args)

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommands: add, update, delete, mark-in-progress, mark-done, list, help")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDescription == "" {
			fmt.Println("Error: Task description is required")
			addCmd.PrintDefaults()
			os.Exit(1)
		}
		addTask(*addDescription)

	case "update":
		updateCmd.Parse(os.Args[2:])
		if *updateID == 0 || *updateDescription == "" {
			fmt.Println("Error: Task ID and description are required")
			updateCmd.PrintDefaults()
			os.Exit(1)
		}
		updateTask(fmt.Sprintf("%d", *updateID), *updateDescription)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Println("Error: Task ID is required")
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		deleteTask(fmt.Sprintf("%d", *deleteID))

	case "mark-in-progress":
		markInProgressCmd.Parse(os.Args[2:])
		if *markInProgressID == 0 {
			fmt.Println("Error: Task ID is required")
			markInProgressCmd.PrintDefaults()
			os.Exit(1)
		}
		updateTaskStatus(fmt.Sprintf("%d", *markInProgressID), StatusInProgress)

	case "mark-done":
		markDoneCmd.Parse(os.Args[2:])
		if *markDoneID == 0 {
			fmt.Println("Error: Task ID is required")
			markDoneCmd.PrintDefaults()
			os.Exit(1)
		}
		updateTaskStatus(fmt.Sprintf("%d", *markDoneID), StatusDone)

	case "list":
		listCmd.Parse(os.Args[2:])
		listTasks(*listStatus)
	case "help":
		helpCmd.Parse(os.Args[2:])
		switch *helpFlag {
		case "add":
			fmt.Println("Add a new task to the tracking system")
			fmt.Println("\nUsage:")
			fmt.Println("  task add -desc <description>     # Using full flag")
			fmt.Println("  task add -d <description>        # Using short flag")
			fmt.Println("  task add <description>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task add -desc \"Complete monthly report\"")
			fmt.Println("  task add -d \"Complete monthly report\"")
			fmt.Println("  task add \"Complete monthly report\"")
			addCmd.PrintDefaults()
		case "update":
			fmt.Println("Update an existing task's description")
			fmt.Println("\nUsage:")
			fmt.Println("  task update -id <task_id> -desc <new_description>      # Using full flags")
			fmt.Println("  task update -i <task_id> -d <new_description>         # Using short flags")
			fmt.Println("  task update <task_id> <new_description>             # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task update -id 1 -desc \"Review monthly report\"")
			fmt.Println("  task update -i 1 -d \"Review monthly report\"")
			fmt.Println("  task update 1 \"Review monthly report\"")
			updateCmd.PrintDefaults()
		case "delete":
			fmt.Println("Delete a task from the system")
			fmt.Println("\nUsage:")
			fmt.Println("  task delete -id <task_id>      # Using full flag")
			fmt.Println("  task delete -i <task_id>       # Using short flag")
			fmt.Println("  task delete <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task delete -id 1")
			fmt.Println("  task delete -i 1")
			fmt.Println("  task delete 1")
			deleteCmd.PrintDefaults()
		case "mark-in-progress":
			fmt.Println("Mark a task as in progress")
			fmt.Println("\nUsage:")
			fmt.Println("  task mark-in-progress -id <task_id>      # Using full flag")
			fmt.Println("  task mark-in-progress -i <task_id>       # Using short flag")
			fmt.Println("  task mark-in-progress <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task mark-in-progress -id 1")
			fmt.Println("  task mark-in-progress -i 1")
			fmt.Println("  task mark-in-progress 1")
			markInProgressCmd.PrintDefaults()
		case "mark-done":
			fmt.Println("Mark a task as completed")
			fmt.Println("\nUsage:")
			fmt.Println("  task mark-done -id <task_id>      # Using full flag")
			fmt.Println("  task mark-done -i <task_id>       # Using short flag")
			fmt.Println("  task mark-done <task_id>          # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task mark-done -id 1")
			fmt.Println("  task mark-done -i 1")
			fmt.Println("  task mark-done 1")
			markDoneCmd.PrintDefaults()
		case "list":
			fmt.Println("List all tasks or filter by status")
			fmt.Println("\nUsage:")
			fmt.Println("  task list -status <status>      # Using full flag")
			fmt.Println("  task list -s <status>           # Using short flag")
			fmt.Println("  task list <status>              # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task list                       # Show all tasks")
			fmt.Println("  task list -status todo          # Using full flag")
			fmt.Println("  task list -s in-progress        # Using short flag")
			fmt.Println("  task list done                  # Direct usage")
			listCmd.PrintDefaults()
		case "help":
			fmt.Println("Display help for specific commands")
			fmt.Println("\nUsage:")
			fmt.Println("  task help -flag <command>      # Using full flag")
			fmt.Println("  task help -f <command>         # Using short flag")
			fmt.Println("  task help <command>            # Direct usage")
			fmt.Println("\nExamples:")
			fmt.Println("  task help -flag add")
			fmt.Println("  task help -f add")
			fmt.Println("  task help add")
			helpCmd.PrintDefaults()
		default:
			fmt.Println("Task Tracking System - Available Commands:")
			fmt.Println("\nCommands:")
			fmt.Println("  add             Add a new task")
			fmt.Println("  update          Update an existing task")
			fmt.Println("  delete          Delete a task")
			fmt.Println("  mark-in-progress Mark a task as in progress")
			fmt.Println("  mark-done       Mark a task as completed")
			fmt.Println("  list            List tasks")
			fmt.Println("  help            Show this help or specific command help")
			fmt.Println("\nTo see detailed help for a specific command:")
			fmt.Println("  task help -flag <command>")
			fmt.Println("\nExample: task help -flag add")
		}

	default:
		fmt.Println("Error: Unknown command")
		os.Exit(1)
	}
}

func loadTasks() TaskList {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return TaskList{Tasks: []Task{}}
	}

	var tasks TaskList
	if err := json.Unmarshal(data, &tasks); err != nil {
		fmt.Printf("Error reading tasks: %v\n", err)
		os.Exit(1)
	}

	return tasks
}

func saveTasks(tasks TaskList) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}
}

func addTask(description string) {
	tasks := loadTasks()

	newID := 1
	if len(tasks.Tasks) > 0 {
		newID = tasks.Tasks[len(tasks.Tasks)-1].ID + 1
	}

	newTask := Task{
		ID:          newID,
		Description: description,
		Status:      StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks.Tasks = append(tasks.Tasks, newTask)
	saveTasks(tasks)
	fmt.Printf("Added task %d: %s\n", newID, description)
}

func updateTask(idStr, newDescription string) {
	id := parseID(idStr)
	tasks := loadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].Description = newDescription
			tasks.Tasks[i].UpdatedAt = time.Now()
			saveTasks(tasks)
			fmt.Printf("Updated task %d\n", id)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func deleteTask(idStr string) {
	id := parseID(idStr)
	tasks := loadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks = append(tasks.Tasks[:i], tasks.Tasks[i+1:]...)
			saveTasks(tasks)
			fmt.Printf("Deleted task %d\n", id)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func updateTaskStatus(idStr, newStatus string) {
	id := parseID(idStr)
	tasks := loadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].Status = newStatus
			tasks.Tasks[i].UpdatedAt = time.Now()
			saveTasks(tasks)
			fmt.Printf("Updated task %d status to %s\n", id, newStatus)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func listTasks(status string) {
	tasks := loadTasks()

	if len(tasks.Tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	for _, task := range tasks.Tasks {
		if status == "" || task.Status == status {
			fmt.Printf("[%d] %s (Status: %s)\n", task.ID, task.Description, task.Status)
		}
	}
}

func parseID(idStr string) int {
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		fmt.Printf("Invalid task ID: %s\n", idStr)
		os.Exit(1)
	}
	return id
}
