package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"task-tracker-cli/models"
)

const dataFile = "tasks.json"

func LoadTasks() models.TaskList {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return models.TaskList{Tasks: []models.Task{}}
	}

	var tasks models.TaskList
	if err := json.Unmarshal(data, &tasks); err != nil {
		fmt.Printf("Error reading tasks: %v\n", err)
		os.Exit(1)
	}

	return tasks
}

func SaveTasks(tasks models.TaskList) {
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

func ParseID(idStr string) int {
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		fmt.Printf("Invalid task ID: %s\n", idStr)
		os.Exit(1)
	}
	return id
}
