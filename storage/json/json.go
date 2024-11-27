package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/basdf/tast-tracker-cli/internal"
)

const (
	DefaultFilePermissions = 0660
)

type Storage struct {
	file *os.File
	data []*internal.Task
}

func NewStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, DefaultFilePermissions)
	if err != nil {
		fmt.Println("Error al crear o abrir el archivo:", err)
		return nil, nil
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}
	data := []*internal.Task{}
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &data)
		if err != nil {
			fmt.Println("Error al leer el archivo:", err)
		}
	}
	return &Storage{file: file, data: data}, nil
}
func (js *Storage) Add(data internal.Task) error {
	js.data = append(js.data, &data)
	return nil
}
func (js *Storage) Update(id int, data string) error {
	for index, task := range js.data {
		if task.ID == id {
			js.data[index].Description = data
			return nil
		}
	}
	return errors.New("Task not found")
}
func (js *Storage) UpdateStatus(id int, status internal.TaskStatus) error {
	for index, task := range js.data {
		if task.ID == id {
			js.data[index].Status = status
			return nil
		}
	}
	return errors.New("Task not found")
}
func (js *Storage) Delete(id int) error {
	for index, task := range js.data {
		if task.ID == id {
			js.data = append(js.data[:index], js.data[index+1:]...)
			return nil
		}
	}
	return errors.New("Task not found")
}
func (js *Storage) Get(filter string) ([]*internal.Task, error) {
	if filter == "" {
		return js.data, nil
	}
	tasks := make([]*internal.Task, 0)
	for _, task := range js.data {
		if string(task.Status) == filter {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
func (js *Storage) GetLastID() (int, error) {
	if len(js.data) == 0 {
		return 0, nil
	}
	lastElement := js.data[len(js.data)-1]
	return lastElement.ID, nil

}
func (js *Storage) Close() error {
	err := js.Save()
	if err != nil {
		return fmt.Errorf("Error al cerrar el archivo: %w", err)
	}
	err = js.file.Close()
	if err != nil {
		return fmt.Errorf("Error al cerrar el archivo: %w", err)
	}
	return nil
}
func (js *Storage) Save() error {
	file, err := json.MarshalIndent(js.data, "", "  ")
	if err != nil {
		return fmt.Errorf("Error al transformar los datos: %w", err)
	}
	err = os.WriteFile(js.file.Name(), file, DefaultFilePermissions)
	if err != nil {
		return fmt.Errorf("Error al guardar los datos: %w", err)
	}
	return nil
}
