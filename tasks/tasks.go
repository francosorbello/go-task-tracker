package tasks

import (
	"backend/jsondatabase"
	"fmt"
	"slices"
	"time"
)

type TaskStatus int

const (
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
	ID          int
	Description string
	Status      TaskStatus
	CreatedAt   string
	UpdatedAt   string
}

func (task Task) GetID() int {
	return task.ID
}

func (task Task) SetID(id int) jsondatabase.Storable {
	task.ID = id
	return task
}

func (task Task) String() string {
	var statusName string
	switch task.Status {
	case Todo:
		statusName = "TO DO"
	case InProgress:
		statusName = "IN PROGRESS"
	case Done:
		statusName = "DONE"
	}
	return fmt.Sprintf("{ID: %d | description: %s | status: %s}", task.ID, task.Description, statusName)
}

func AddTask(description string) (Task, error) {
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

	for i, t := range tasks {
		if t.GetID() == id {
			index = i
			break
		}
	}

	return index
}

func UpdateTask(id int, description string) error {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return getErr
	}

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Description = description
		setUpdatedDate(&tasks[taskToUpdate])
		writeErr := db.WriteAll(tasks)
		if writeErr != nil {
			return writeErr
		}
	} else {
		return fmt.Errorf("no task with id %d available", id)
	}
	return nil
}

func UpdateTaskStatus(id int, status TaskStatus) error {
	var db jsondatabase.Database[Task]
	if openErr := db.Open(); openErr != nil {
		return openErr
	}

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return getErr
	}

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Status = status
		setUpdatedDate(&tasks[taskToUpdate])
		writeErr := db.WriteAll(tasks)
		if writeErr != nil {
			return writeErr
		}
	} else {
		return fmt.Errorf("no task with id %d available", id)
	}

	return db.Close()
}

func setUpdatedDate(task *Task) {
	task.UpdatedAt = time.Now().String()
}

func DeleteTask(id int) error {
	var db jsondatabase.Database[Task]
	if openErr := db.Open(); openErr != nil {
		return openErr
	}

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return getErr
	}

	taskToDelete := findTaskIndex(id, tasks)

	if taskToDelete != -1 {
		tasks = slices.Delete(tasks, taskToDelete, taskToDelete+1)
		if len(tasks) == 0 {
			clearErr := db.Clear()
			if clearErr != nil {
				return clearErr
			}
		} else {
			writeErr := db.WriteAll(tasks)
			if writeErr != nil {
				return writeErr
			}
		}
	} else {
		return fmt.Errorf("no task with id %d available", id)
	}
	return db.Close()
}

func ChangeTaskStatus(id int, newStatus TaskStatus) error {
	var db jsondatabase.Database[Task]
	if openErr := db.Open(); openErr != nil {
		return openErr
	}

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return getErr
	}

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Status = newStatus
		if writErr := db.WriteAll(tasks); writErr != nil {
			return writErr
		}
	} else {
		return fmt.Errorf("no task with id %d available", id)
	}

	return db.Close()
}

func ListTasks(status int) error {
	var db jsondatabase.Database[Task]
	db.Open()

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return getErr
	}

	for _, task := range tasks {
		if status == -1 {
			fmt.Println(task)
		} else {
			if task.Status == TaskStatus(status) {
				fmt.Println(task)
			}
		}
	}
	return db.Close()
}
