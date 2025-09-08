package errors

import "errors"

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrTodoAlreadyExists = errors.New("todo already exists")
	ErrorCantChangeTodo  = errors.New("i cant chage todo by ID")
)
