package tasks

import (
	"backend/jsondatabase"
	"fmt"
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

func (task Task) GetID() int {
	return task.ID
}

func TestWrite() {
	test_task := Task{0,"hello",0,"",""}
	jsondatabase.Insert(test_task)
}

func TestRead() {
	task := jsondatabase.FindByID2[Task](0)

	fmt.Println(task)
}
