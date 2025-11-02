package jsondatabase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type NotAJsonError struct {
    filename string
}

func (e *NotAJsonError) Error() string {
    return fmt.Sprintf("File %s is not a json file",e.filename)
}

var ErrNotAJsonOld = errors.New("provided path is not a json path")

type Storable interface {
    GetID() int
    SetID(id int) Storable
}

type Database[T Storable] struct {
    path string
    file *os.File
}

func (db *Database[T]) Open() {
    fmt.Println("Opening db")
    db.file = getDb(db.path)
}

func (db *Database[T]) Close() {
    fmt.Println("Closing database")
    db.file.Close()
}

func (db *Database[T]) Append(item T) {
    items := db.GetAll()
    fmt.Println("Appending item",item, "to list",items)
    lastId := 0
    for _,other := range items {
        fmt.Println(other.GetID())
        if other.GetID() > lastId {
            lastId = other.GetID()
        }
    }

    modifiedItem := item.SetID(lastId + 1)
    items = append(items, modifiedItem.(T))
    db.WriteAll(items)
}

func (db *Database[T]) GetAll() []T {
    decoder := json.NewDecoder(db.file)
    // decoder.Token()

    entries := []T{}

    for decoder.More() {
        err := decoder.Decode(&entries)
        if err != nil {
            panic(err)
        }
        fmt.Println("entry:",entries) 
    }

    return entries
}

func (db *Database[T]) WriteAll(items []T) {
    fmt.Println("Writing items",items)
    data,err := json.MarshalIndent(items,"","    ")
    if err != nil {
        panic(err)
    }
    _,err = db.file.WriteAt(data,0)
    if err != nil {
        panic(err)
    }
}


func isValidPath(pth string) error {
    if strings.Contains(pth, ".json") {
        return nil
    }
    return &NotAJsonError{pth}
}


func getDb(customPath string) *os.File {
    if customPath == "" {
        customPath = "db.json"
    } else {
       err := isValidPath(customPath)
        if err != nil {
            panic(err)
        } 
    }
    file, err := os.OpenFile(customPath,os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        panic(err)
    }
    return file
}
