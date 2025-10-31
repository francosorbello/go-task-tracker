package jsondatabase

import (
	"encoding/json"
	"os"
)

type DbFile = os.File

type Scanner interface {
	Scan(src any) error
}

func Insert[T any](f *DbFile, item T) error {
	data,err := json.Marshal(item)

	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil
}