package storage

import "github.com/basdf/tast-tracker-cli/internal"

type Storage interface {
	Add(data internal.Task) error
	Update(ID int, data string) error
	UpdateStatus(ID int, status internal.TaskStatus) error
	Delete(ID int) error
	Get(filter string) ([]*internal.Task, error)
	GetLastID() (int, error)
	Close() error
	Save() error
}
