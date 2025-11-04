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
	return fmt.Sprintf("File %s is not a json file", e.filename)
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

func (db *Database[T]) Open() error {
	file,err := getDb(db.path)
    if err != nil {
        return err
    }
    db.file = file   
    return nil
}

func (db *Database[T]) Close() error {
	return db.file.Close()
}


func (db *Database[T]) Append(item T) (T,error) {
    var zero T
	items, getErr := db.GetAll()
    if getErr != nil {
        return zero,getErr
    }
	lastId := item.GetID()
	for _, other := range items {
		if other.GetID() > lastId {
			lastId = other.GetID()
		}
	}

	// TODO: this creates a new copy of the task.
	// Ideally, we would set the id using pointers.
	modifiedItem := item.SetID(lastId + 1)
	items = append(items, modifiedItem.(T))
	db.WriteAll(items)
	return modifiedItem.(T),nil
}

func (db *Database[T]) GetAll() ([]T,error) {
	decoder := json.NewDecoder(db.file)
	// decoder.Token()

	entries := []T{}

	for decoder.More() {
		err := decoder.Decode(&entries)
		if err != nil {
			return nil, err
		}
	}

	return entries,nil
}

func (db *Database[T]) WriteAll(items []T) error {
	data, marshallErr := json.MarshalIndent(items, "", "    ")
	if marshallErr != nil {
        return marshallErr
	}

	clearErr := db.Clear()
    if clearErr != nil {
        return clearErr
    }
	
    _, writingErr := db.file.WriteAt(data, 0)
	if writingErr != nil {
		return writingErr
	}

	return nil
}

func (db *Database[T]) Clear() error{
	if db.path == "" {
		db.path = "db.json"
	}
	return os.Truncate(db.path, 0)
}

func isValidPath(pth string) error {
	if strings.Contains(pth, ".json") {
		return nil
	}
	return &NotAJsonError{pth}
}

func getDb(customPath string) (*os.File,error) {
	if customPath == "" {
		customPath = "db.json"
	} else {
		pathErr := isValidPath(customPath)
		if pathErr != nil {
			return nil, pathErr
		}
	}
	file, openFileErr := os.OpenFile(customPath, os.O_RDWR|os.O_CREATE, 0644)
	if openFileErr != nil {
		return nil, openFileErr
	}
	return file,nil
}
