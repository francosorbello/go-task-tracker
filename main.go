package main

import (
	"backend/tasks"
	"fmt"
	"os"
)

func main() {
	argsWithProg := os.Args
    fmt.Println(argsWithProg)
	// tasks.TestWrite()
	tasks.AddTask("this is a task")
}