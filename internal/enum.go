package internal

type Enum interface {
	Values() []Enum
}
type TaskStatus string

func (t TaskStatus) Values() []Enum {
	return TaskStatusValues
}

const (
	TaskStatusTODO       TaskStatus = "todo"
	TaskStatusDONE       TaskStatus = "done"
	TaskStatusINPROGRESS TaskStatus = "in-progress"
)

var TaskStatusValues []Enum

func init() {
	TaskStatusValues = []Enum{TaskStatusTODO, TaskStatusDONE, TaskStatusINPROGRESS}
}
