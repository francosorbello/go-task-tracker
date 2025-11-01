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
	fmt.Println("set id",id,"for task",task.Description)
	task.ID = id
}

func TestWrite() {
	test_task := Task{0,"hello",0,"",""}
	jsondatabase.Insert(test_task)
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