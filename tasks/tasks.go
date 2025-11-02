package tasks

import (
	"backend/jsondatabase"
	"fmt"
	"time"
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

func (task Task) SetID(id int) {
	task.ID = id
}


func TestRead() {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()

	fmt.Println(tasks)
}

func AddTask(description string) {
	var newTask Task
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	newTask.Description = description
	newTask.CreatedAt = time.Now().String()

	db.Append(newTask)
}