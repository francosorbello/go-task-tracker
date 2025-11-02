package main

import (
	"backend/tasks"
	"errors"
	"fmt"
	"os"
	"strconv"
)

type InvalidArgNumberError struct {
	commandName string
	expectedArgs int
	receivedArgs int
}

func (e *InvalidArgNumberError) Error() string{
	return fmt.Sprintf("%s command expects %d args, but received %d",e.commandName,e.expectedArgs,e.receivedArgs)
}

func validateArgs(commandName string, expected int, received int) error {
	if (received == 0) || (expected != received) {
		return &InvalidArgNumberError{commandName,expected, received}
	}
	return nil
}

type Command struct {
	argNumber int
}

type Executable interface {
	Verify(args []string) error
	Execute(args []string)
}

type AddCommand Command

func (com *AddCommand) Verify(args []string) error{
	err := validateArgs("add",com.argNumber,len(args))
	
	return err	
}

func (com *AddCommand) Execute(args []string) {
	task := tasks.AddTask(args[0])
	fmt.Printf("Task added sucessfully (ID: %d)",task.ID)
}

type UpdateCommand Command

func (com *UpdateCommand) Verify(args []string) error{
	err := validateArgs("add",com.argNumber,len(args))
	
	return err		
}

func (com *UpdateCommand) Execute(args []string) {
	taskId,_ := strconv.Atoi(args[0])
	desc := args[1]
	tasks.UpdateTask(taskId,desc)
}


func main() {
	receivedArgs := os.Args[1:]
	if len(receivedArgs) == 0 {
		panic(errors.New("no command"))
	}
	commName := receivedArgs[0]
	var comm Executable
	
	switch commName {
	case "add":
		comm = &AddCommand{1}
	case "update":
		comm = &UpdateCommand{2}
	default:
		panic(errors.New("No command named "+commName))
	}
	err := comm.Verify(receivedArgs[1:])
	if err != nil {
		panic(err)
	}
	comm.Execute(receivedArgs[1:])
}