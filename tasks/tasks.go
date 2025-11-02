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

func (task Task) SetID(id int) jsondatabase.Storable{
	task.ID = id
	return task
}


func TestRead() {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()

	fmt.Println(tasks)
}

func AddTask(description string) Task {
	var newTask Task
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	newTask.Description = description
	newTask.CreatedAt = time.Now().String()

	return db.Append(newTask)
}

func UpdateTask(id int, description string) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()
	
	taskToUpdate := -1 
	for i,t := range tasks {
		if t.GetID() == id {
			taskToUpdate = i
		}
	}

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Description = description
		db.WriteAll(tasks)
	} else {
		panic(fmt.Sprintf("no task with id %d available",id))
	}
}