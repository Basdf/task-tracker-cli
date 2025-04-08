package services

import (
	"fmt"
	"time"

	"task-tracker-cli/models"
	"task-tracker-cli/storage"
)

func AddTask(description string) {
	tasks := storage.LoadTasks()

	newID := 1
	if len(tasks.Tasks) > 0 {
		newID = tasks.Tasks[len(tasks.Tasks)-1].ID + 1
	}

	newTask := models.Task{
		ID:          newID,
		Description: description,
		Status:      models.StatusTodo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks.Tasks = append(tasks.Tasks, newTask)
	storage.SaveTasks(tasks)
	fmt.Printf("Added task %d: %s\n", newID, description)
}

func UpdateTask(idStr, newDescription string) {
	id := storage.ParseID(idStr)
	tasks := storage.LoadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].Description = newDescription
			tasks.Tasks[i].UpdatedAt = time.Now()
			storage.SaveTasks(tasks)
			fmt.Printf("Updated task %d\n", id)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func DeleteTask(idStr string) {
	id := storage.ParseID(idStr)
	tasks := storage.LoadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks = append(tasks.Tasks[:i], tasks.Tasks[i+1:]...)
			storage.SaveTasks(tasks)
			fmt.Printf("Deleted task %d\n", id)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func UpdateTaskStatus(idStr, newStatus string) {
	id := storage.ParseID(idStr)
	tasks := storage.LoadTasks()

	for i := range tasks.Tasks {
		if tasks.Tasks[i].ID == id {
			tasks.Tasks[i].Status = newStatus
			tasks.Tasks[i].UpdatedAt = time.Now()
			storage.SaveTasks(tasks)
			fmt.Printf("Updated task %d status to %s\n", id, newStatus)
			return
		}
	}

	fmt.Printf("Task %d not found\n", id)
}

func ListTasks(status string) {
	tasks := storage.LoadTasks()

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
