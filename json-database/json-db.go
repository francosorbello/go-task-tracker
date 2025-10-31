package jsondatabase

import (
	"errors"
	"os"
	"strings"
)

var ErrNotAJson = errors.New("provided path is not a json path")

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

func CreateDb(customPath string) error{
	if DbExists(customPath) {
		return nil
	}
	
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
	_, err := os.Create(pth)
	if err != nil {
		return err
	}

	return nil
}