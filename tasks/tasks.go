package tasks

import (
	"backend/jsondatabase"
	"fmt"
	"slices"
	"time"
)

type TaskStatus int

const(
	Todo TaskStatus = iota
	InProgress
	Done
)

func StatusNameToValue(name string) TaskStatus {
	switch name {
	case "todo":
		return Todo
	case "done":
		return Done
	case "in-progress":
		return InProgress
	}
	return -1
}

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

func findTaskIndex(id int, tasks []Task) int {
	index := -1

	for i,t := range tasks {
		if t.GetID() == id {
			index = i
			break
		}
	}

	return index
}

func UpdateTask(id int, description string) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()
	
	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Description = description
		db.WriteAll(tasks)
	} else {
		panic(fmt.Sprintf("no task with id %d available",id))
	}
}

func UpdateTaskStatus(id int, status TaskStatus) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()
	
	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Status = status
		db.WriteAll(tasks)
	} else {
		panic(fmt.Sprintf("no task with id %d available",id))
	}
}

func DeleteTask(id int) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()

	taskToDelete := findTaskIndex(id, tasks)

	if taskToDelete != -1 {
		tasks = slices.Delete(tasks,taskToDelete,taskToDelete+1)
		if len(tasks) == 0 {
			db.Clear()
		}else {
			db.WriteAll(tasks)
		}
	} else {
		panic(fmt.Sprintf("no task with id %d available",id))
	}
} 

func ChangeTaskStatus(id int, newStatus TaskStatus) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks := db.GetAll()

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Status = newStatus
		db.WriteAll(tasks)
	} else {
		panic(fmt.Sprintf("no task with id %d available",id))
	}
}

func ListTasks(status int) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()
	
	tasks := db.GetAll()
	for _,task := range tasks {
		if status == -1 {
			fmt.Println(task)
		} else {
			if task.Status == TaskStatus(status) {
				fmt.Println(task)
			}
		}
	}
}