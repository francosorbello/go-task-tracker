package jsondatabase

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

var ErrNotAJson = errors.New("provided path is not a json path")

type Storable interface {
    GetID() int
    SetID(id int)
}

type Database[T Storable] struct {
    path string
    file *os.File
}

func (db *Database[T]) Open() {
    db.file = GetOrCreate(db.path)
}

func (db *Database[T]) Close() {
    db.file.Close()
}

func (db *Database[T]) Append(item T) {
    items := db.GetAll()
    lastId := 1
    for _,other := range items {
        if other.GetID() > lastId {
            lastId = other.GetID()
        }
    }

    item.SetID(lastId + 1)
    items = append(items, item)
    db.WriteAll(items)
}

func (db *Database[T]) GetAll() []T {
    decoder := json.NewDecoder(db.file)
    // decoder.Token()

    // data := T
    var data T
    entries := []T{}

    for decoder.More() {
        err := decoder.Decode(&data)
        if err != nil {
            panic(err)
        }
        
        entries = append(entries, data)
    }
    return entries
}

func (db *Database[T]) WriteAll(items []T) {
    for _,item := range items {
        data, err := json.Marshal(item)
        if err != nil {
            panic(err)
        }

        _, err = db.file.Write(data)
        if err != nil {
            panic(err)
        }
    }
}


func IsValidPath(pth string) error {
    if strings.Contains(pth, ".json") {
        return nil
    }
    return ErrNotAJson
}

func DbExists(dbName string) bool {
    _, err := os.Stat(dbName)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return false
        }
        panic(err)
    }

    return true
}

func CreateDb(customPath string) (*os.File, error){	
    pth := ""
    if customPath == "" {
        pth = "db.json"
    } else {
        err := IsValidPath(customPath)
        if err != nil && errors.Is(err, ErrNotAJson) {
            customPath += ".json"
        }
        pth = customPath
    }
    file, err := os.Create(pth)
    if err != nil {
        return nil,err
    }

    return file,nil
}

func GetDb(customPath string) *os.File {
    file, err := os.OpenFile(customPath,os.O_CREATE| os.O_RDWR, 0644)
    if err != nil {
        panic(err)
    }
    return file
}

func GetOrCreate(customPath string) *os.File {
    if customPath == "" {
        customPath = "db.json"
    }
    
    var file *os.File

    if !DbExists(customPath) {
        file,_ = CreateDb(customPath)
    } else {
        file = GetDb(customPath)
    }
    return file
}

func Insert(item any) {
    var data,err = json.Marshal(item)

    if err != nil {
        panic(err)
    } 

    file := GetOrCreate("") 
    defer file.Close()

    file.Write(data)
}

func FindByID(id int) map[string]any {
    var file = GetOrCreate("")
    defer file.Close()

    decoder := json.NewDecoder(file)
    // decoder.Token()

    data := map[string]any{}

    for decoder.More() {
        err := decoder.Decode(&data)
        if err != nil {
            panic(err)
        }
        
        value, ok := data["ID"]
        if ok && value == id {
            return data   
        }
    }
    return nil
}

func FindByID2[T Storable] (id int) Storable {
    var file = GetOrCreate("")
    defer file.Close()

    decoder := json.NewDecoder(file)
    // decoder.Token()

    // data := T
    var data T

    for decoder.More() {
        err := decoder.Decode(&data)
        if err != nil {
            panic(err)
        }
        
        value := data.GetID()
        if value == id {
            return data   
        }
    }
    return nil
}