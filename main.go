package main

import (
	"backend/tasks"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

func main() {
	argsWithProg := os.Args
    fmt.Println(argsWithProg)
	// tasks.TestWrite()
	tasks.AddTask("task "+strconv.Itoa(rand.IntN(100)))
	// tasks.TestRead()
}