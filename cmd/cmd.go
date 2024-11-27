package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/basdf/tast-tracker-cli/internal"
	"github.com/basdf/tast-tracker-cli/storage"
	jsonStorage "github.com/basdf/tast-tracker-cli/storage/json"
)

const (
	MinArgsForList = 2
)

type CMD struct {
	storage storage.Storage
	regex   *regexp.Regexp
}

func New(filename string) *CMD {
	re := regexp.MustCompile(`"[^"]+"|[a-zA-Z0-9-]+`)
	if !strings.Contains(filename, ".json") {
		filename += ".json"
	}
	storage, _ := jsonStorage.NewStorage(filename)
	return &CMD{regex: re, storage: storage}
}

func (c *CMD) analizedInput(input string) []string {
	return c.regex.FindAllString(input, -1)
}

func (c *CMD) ExecuteCommand(input string, sigChan chan os.Signal) bool {
	args := c.analizedInput(input)
	if len(args) == 0 {
		fmt.Println("Comando no encontrado")
		return false
	}
	switch args[1] {
	case "add":
		c.add(args[2])
	case "update":
		ID, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		c.update(ID, args[3])
	case "delete":
		ID, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		c.delete(ID)
	case "list":
		if len(args) == MinArgsForList {
			args = append(args, "")
		}
		c.list(args[2])
	case "mark-in-progress":
		ID, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		c.markInProgress(ID)
	case "mark-done":
		ID, err := strconv.Atoi(args[2])
		if err != nil {
			panic(err)
		}
		c.markDone(ID)
	case "exit":
		sigChan <- syscall.SIGINT
		return true
	case "save":
		c.save()
	}
	return false
}

func (c *CMD) add(description string) {
	LastID, err := c.storage.GetLastID()
	if err != nil {
		panic(err)
	}
	task := internal.Task{
		ID:          LastID + 1,
		Description: description,
		Status:      internal.TaskStatusTODO,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err = c.storage.Add(task)
	if err != nil {
		panic(err)
	}
}

func (c *CMD) update(id int, data string) {
	err := c.storage.Update(id, data)
	if err != nil {
		panic(err)
	}
}

func (c *CMD) delete(id int) {
	err := c.storage.Delete(id)
	if err != nil {
		panic(err)
	}
}

func (c *CMD) list(filter string) {
	tasks, err := c.storage.Get(filter)
	if err != nil {
		panic(err)
	}
	for _, task := range tasks {
		fmt.Println(task)
	}
}

func (c *CMD) markInProgress(id int) {
	err := c.storage.UpdateStatus(id, internal.TaskStatusINPROGRESS)
	if err != nil {
		panic(err)
	}
}

func (c *CMD) markDone(id int) {
	err := c.storage.UpdateStatus(id, internal.TaskStatusDONE)
	if err != nil {
		panic(err)
	}
}

func (c *CMD) Close() {
	err := c.storage.Close()
	if err != nil {
		panic(err)
	}
}

func (c *CMD) save() {
	err := c.storage.Save()
	if err != nil {
		panic(err)
	}
}
