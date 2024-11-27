package internal

import (
	"errors"
)

var (
	ErrInvalidInput = errors.New("invalid input command")
	ErrTaskNotFound = errors.New("task not found")
)
