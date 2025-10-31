package tasks

import (
	"backend/jsondatabase"
)

type TaskStatus int

const(
	Todo TaskStatus = iota
	InProgress
	Done
)

type Task struct {
	ID int
	Description string
	Status TaskStatus
	CreatedAt string
	UpdatedAt string
}

func TestWrite() {
	test_task := Task{0,"hello",0,"",""}
	
}