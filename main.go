package main

import (
	"backend/tasks"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type InvalidArgNumberError struct {
	commandName  string
	expectedArgs int
	receivedArgs int
}

func (e *InvalidArgNumberError) Error() string {
	return fmt.Sprintf("%s command expects %d args, but received %d", e.commandName, e.expectedArgs, e.receivedArgs)
}

func validateArgs(commandName string, expected int, received int) error {
	if expected != received {
		return &InvalidArgNumberError{commandName, expected, received}
	}
	return nil
}

type Command struct {
	argNumber int
	data      map[string]any
}

type Executable interface {
	Verify(args []string) error
	Execute(args []string) error
}

type AddCommand Command

func (com *AddCommand) Verify(args []string) error {
	err := validateArgs("add", com.argNumber, len(args))

	return err
}

func (com *AddCommand) Execute(args []string) error {
	task, addTaskErr := tasks.AddTask(args[0])
	if addTaskErr != nil {
		return addTaskErr
	}
	fmt.Printf("Task added sucessfully (ID: %d)", task.ID)
	return nil
}

type UpdateCommand Command

func (com *UpdateCommand) Verify(args []string) error {
	err := validateArgs("update", com.argNumber, len(args))

	return err
}

func (com *UpdateCommand) Execute(args []string) error {
	taskId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	desc := args[1]
	return tasks.UpdateTask(taskId, desc)
}

type DeleteCommand Command

func (com *DeleteCommand) Verify(args []string) error {
	err := validateArgs("delete", com.argNumber, len(args))

	return err
}

func (com *DeleteCommand) Execute(args []string) error {
	taskId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	return tasks.DeleteTask(taskId)
}

type UpdateStatusCommand Command

func (com *UpdateStatusCommand) Verify(args []string) error {
	err := validateArgs("update status", com.argNumber, len(args))
	if err != nil {
		return err
	}
	status, ok := com.data["status"].(tasks.TaskStatus)
	if !ok {
		return errors.New("command data doesnt have a status entry")
	}
	if status > tasks.Done {
		return fmt.Errorf("invalid status. Received %d, expected %s or %s", status, tasks.InProgress, tasks.Done)
	}

	return nil
}

func (com *UpdateStatusCommand) Execute(args []string) error {
	taskId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	status := com.data["status"].(tasks.TaskStatus)

	return tasks.UpdateTaskStatus(taskId, status)
}

type ListCommand Command

func (com *ListCommand) Verify(args []string) error {
	err := validateArgs("list", com.argNumber, len(args))
	if err != nil {
		return err
	}

	if com.argNumber == 1 {
		statusArg := args[0]
		if statusArg != "todo" && statusArg != "in-progress" && statusArg != "done" {
			return fmt.Errorf("invalid action for list command. Expected, todo, in-progress or done. Received %s", statusArg)
		}
	}

	return nil
}

func (com *ListCommand) Execute(args []string) error {
	status := -1
	if com.argNumber == 1 {
		statusName := args[0]
		status = int(tasks.StatusNameToValue(statusName))
		fmt.Println("==", statusName, "==")
	} else {
		fmt.Println("== All tasks ==")
	}

	err := tasks.ListTasks(status)
	if err != nil {
		return err
	}
	fmt.Println("====")
	
	return nil
}

func main() {
	receivedArgs := os.Args[1:]
	if len(receivedArgs) == 0 {
		fmt.Println("No commands received")
		return
	}

	commName := receivedArgs[0]
	var comm Executable
	var err error
	commandData := map[string]any{}

	switch commName {
	case "add":
		comm = &AddCommand{1, commandData}
	case "update":
		comm = &UpdateCommand{2, commandData}
	case "delete":
		comm = &DeleteCommand{1, commandData}
	case "mark-in-progress":
		commandData["status"] = tasks.InProgress
		comm = &UpdateStatusCommand{1, commandData}
	case "mark-done":
		commandData["status"] = tasks.Done
		comm = &UpdateStatusCommand{1, commandData}
	case "list":
		if len(receivedArgs) > 1 {
			comm = &ListCommand{1, commandData}
		} else {
			comm = &ListCommand{0, commandData}
		}
	default:
		fmt.Println("no command named " + commName)
		return
	}

	err = comm.Verify(receivedArgs[1:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = comm.Execute(receivedArgs[1:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
