package tasks

import (
	"errors"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()

	err := os.Remove("./db.json")
	if err != nil {
		panic(err)
	}

	os.Exit(exitVal)
}

func TestAddTask(t *testing.T) {
	expectedDesc := "task 1"
	result, err := AddTask(expectedDesc)
	if err != nil {
		t.Error(err.Error())
	}
	if result.Description != expectedDesc {
		t.Errorf("description should be %s, but received %s", expectedDesc, result.Description)
	}

	if result.CreatedAt == "" {
		t.Errorf("Adding a task should populate CreatedAt field")
	}
}

func TestAddTaskEmpty(t *testing.T) {
	task, err := AddTask("")
	if err == nil || !errors.Is(err, ErrNoDescriptionProvided) {
		t.Errorf("expected ErrNoDescriptionProvided error")
	}
	if task != ZeroTask {
		t.Errorf("expected ZeroTask, received %s", task)
	}
}

func TestUpdateTask(t *testing.T) {
	updatedText := "Updated task 1"
	task, _ := AddTask("Test 1")
	updatedTask, err := UpdateTask(task.ID, updatedText)
	if err != nil {
		t.Error(err.Error())
	}

	if updatedTask.Description != updatedText {
		t.Errorf("description doesnt match. expected %s, received %s", updatedText, updatedTask.Description)
	}

	taskToCompare, err := GetTask(task.ID)
	if err != nil {
		t.Error(err.Error())
	}
	if taskToCompare.Description != updatedText {
		t.Errorf("description doesnt match. expected %s, received %s", updatedText, taskToCompare.Description)
	}
}

func TestUpdateTaskEmpty(t *testing.T) {
	AddTask("Test task")

	task, err := UpdateTask(1, "")
	if err == nil || !errors.Is(err, ErrNoDescriptionProvided) {
		t.Errorf("expected ErrNoDescriptionProvided error")
	}
	if task != ZeroTask {
		t.Errorf("expected ZeroTask, received %s", task)
	}
}

func TestUpdateNonExistingTask(t *testing.T) {
	_, err := UpdateTask(10, "")
	if err == nil {
		t.Errorf("expected error when updating non existing task")
	}
}

func TestUpdateTaskStatus(t *testing.T) {
	updatedStatus := InProgress
	task, _ := AddTask("Test 1")
	updatedTask, err := UpdateTaskStatus(task.ID, updatedStatus)
	if err != nil {
		t.Error(err.Error())
	}

	if updatedTask.Status != updatedStatus {
		t.Errorf("status doesnt match. expected %d, received %d", updatedStatus, updatedTask.Status)
	}

	taskToCompare, err := GetTask(task.ID)
	if err != nil {
		t.Error(err.Error())
	}
	if taskToCompare.Status != updatedStatus {
		t.Errorf("status doesnt match. expected %d, received %d", updatedStatus, taskToCompare.Status)
	}
}

func TestUpdateTaskStatusWrongStatus(t *testing.T) {
	expectedErrType := &ErrInvalidStatus{}
	task, _ := AddTask("test task")
	wrongStatus := TaskStatus(10)
	_,err := UpdateTaskStatus(task.ID,wrongStatus)
	
	if err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected ErrInvalidStatus, received %s",err)
	}
}

func TestListTasks(t *testing.T) {
	err := ListTasks(AllTasks)
	if err != nil {
		t.Error(err)
	}
	
	err = ListTasks(int(Todo))
	if err != nil {
		t.Error(err)
	}

	err = ListTasks(int(InProgress))
	if err != nil {
		t.Error(err)
	}

	err = ListTasks(int(Done))
	if err != nil {
		t.Error(err)
	}
}

func TestListTaskInvalidStatus(t *testing.T) {
	expectedErrType := &ErrInvalidStatus{}
	err := ListTasks(10)
	if err == nil || !errors.As(err, &expectedErrType) {
		t.Errorf("expected ErrInvalidStatus, received %s",err)
	}
}