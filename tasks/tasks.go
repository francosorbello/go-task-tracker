package tasks

import (
	"backend/jsondatabase"
	"errors"
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

const AllTasks = -1

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

var ZeroTask = Task{ID: -1, Description: ""}

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

var ErrNoDescriptionProvided = errors.New("no description provided")

type ErrTaskNotFound struct {
	ID int
}

func (e *ErrTaskNotFound) Error() string {
	return fmt.Sprintf("Task with id %d not found",e.ID)
}

type ErrInvalidStatus struct {
	receivedStatus TaskStatus
}

func (e *ErrInvalidStatus) Error() string {
	return fmt.Sprintf("invalid status. Value must be between %d and %d but received %d instead",Todo,Done,e.receivedStatus)
}

func validateStatus(status TaskStatus) error{
	if (status > Done) || (status < Todo) {
		return &ErrInvalidStatus{status}
	}
	return nil
}

func AddTask(description string) (Task, error) {
	if description == "" {
		return ZeroTask,ErrNoDescriptionProvided
	}

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

func UpdateTask(id int, description string) (Task,error) {
	if description == "" {
		return ZeroTask,ErrNoDescriptionProvided
	}
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks, getErr := db.GetAll()
	if getErr != nil {
		return ZeroTask,getErr
	}

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Description = description
		setUpdatedDate(&tasks[taskToUpdate])
		writeErr := db.WriteAll(tasks)
		if writeErr != nil {
			return ZeroTask,writeErr
		}
	} else {
		return ZeroTask,fmt.Errorf("no task with id %d available", id)
	}
	return tasks[taskToUpdate],nil
}

func UpdateTaskStatus(id int, status TaskStatus) (Task,error) {
	var db jsondatabase.Database[Task]
	if openErr := db.Open(); openErr != nil {
		return ZeroTask,openErr
	}

	if validateErr := validateStatus(status); validateErr != nil {
		db.Close()
		return ZeroTask,validateErr		
	}

	tasks, getErr := db.GetAll()
	if getErr != nil {
		db.Close()
		return ZeroTask,getErr
	}

	taskToUpdate := findTaskIndex(id, tasks)

	if taskToUpdate != -1 {
		tasks[taskToUpdate].Status = status
		setUpdatedDate(&tasks[taskToUpdate])
		writeErr := db.WriteAll(tasks)
		if writeErr != nil {
			return ZeroTask,writeErr
		}
	} else {
		return ZeroTask,fmt.Errorf("no task with id %d available", id)
	}

	return tasks[taskToUpdate],db.Close()
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

	if validateErr := validateStatus(TaskStatus(status)); status != AllTasks && validateErr != nil {
		db.Close()
		return validateErr		
	}

	tasks, getErr := db.GetAll()
	if getErr != nil {
		db.Close()
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

func GetTask(id int) (Task, error) {
	var db jsondatabase.Database[Task]
	db.Open()
	defer db.Close()

	tasks,err := db.GetAll()	
	if err != nil {
		return ZeroTask,err
	}

	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}	
	}

	return ZeroTask, &ErrTaskNotFound{id}
}
