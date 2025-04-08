package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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

func main() {
	// Subcommands
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDescription := addCmd.String("desc", "", "Task description")

	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	updateID := updateCmd.Int("id", 0, "Task ID")
	updateDescription := updateCmd.String("desc", "", "New task description")

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "Task ID")

	markInProgressCmd := flag.NewFlagSet("mark-in-progress", flag.ExitOnError)
	markInProgressID := markInProgressCmd.Int("id", 0, "Task ID")

	markDoneCmd := flag.NewFlagSet("mark-done", flag.ExitOnError)
	markDoneID := markDoneCmd.Int("id", 0, "Task ID")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listStatus := listCmd.String("status", "", "Filter tasks by status (todo, in-progress, done)")

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommands: add, update, delete, mark-in-progress, mark-done, list")
		fmt.Println("\nUsage:")
		fmt.Println("  add -desc \"Task description\"")
		fmt.Println("  update -id <task-id> -desc \"New description\"")
		fmt.Println("  delete -id <task-id>")
		fmt.Println("  mark-in-progress -id <task-id>")
		fmt.Println("  mark-done -id <task-id>")
		fmt.Println("  list [-status <todo|in-progress|done>]")
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
