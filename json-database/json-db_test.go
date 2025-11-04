package jsondatabase

import (
	"errors"
	"os"
	"testing"
)

type MockupStorable struct {
	Id          int
	StringField string
}

func (ms MockupStorable) SetID(id int) Storable {
	ms.Id = id
	return ms
}

func (ms MockupStorable) GetID() int {
	return ms.Id
}

func TestOpenEmpty(t *testing.T) {
	db := Database[MockupStorable]{}

	got := db.Open()
	defer db.file.Close()

	t.Cleanup(func() {
		err := os.Remove("./db.json")
		if err != nil {
			panic(err)
		}
	})

	if got != nil {
		t.Error(got.Error())
	}

	_, got = os.Stat("./db.json")
	if got != nil && errors.Is(got, os.ErrNotExist) {
		t.Errorf("open should create a db.json file, but file doesnt exist")
	}
}

func TestOpenWithInvalidPath(t *testing.T) {
	db := Database[MockupStorable]{path: "custompath"}
	got := db.Open()
	if got == nil {
		t.Errorf("expected to receive an error. got nil")
	}

	target := &NotAJsonError{}
	if got != nil && !(errors.As(got, &target)) {
		t.Errorf("expected NotAJsonError, got %s", got.Error())
	}
}

func TestOpenWithValidPath(t *testing.T) {
	db := Database[MockupStorable]{path: "custompath.json"}
	got := db.Open()
	defer db.file.Close()

	t.Cleanup(func() {
		err := os.Remove("./custompath.json")
		if err != nil {
			panic(err)
		}
	})

	if got != nil {
		t.Error(got.Error())
	}

	_, got = os.Stat("./custompath.json")
	if got != nil && errors.Is(got, os.ErrNotExist) {
		t.Errorf("open should create a custompath.json file, but file doesnt exist")
	}
}

func TestWriteAllSingleItem(t *testing.T) {
	db := Database[MockupStorable]{path: "writeAllDB.json"}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./writeAllDB.json")
		if err != nil {
			panic(err)
		}
	})

	itemToWrite := MockupStorable{0, "test"}
	items := []MockupStorable{itemToWrite}
	err := db.WriteAll(items)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWriteAllMultipleItems(t *testing.T) {
	db := Database[MockupStorable]{path: "writeAllDB.json"}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./writeAllDB.json")
		if err != nil {
			panic(err)
		}
	})

	itemToWrite := MockupStorable{0, "test"}
	itemToWrite2 := MockupStorable{1, "test2"}
	items := []MockupStorable{itemToWrite, itemToWrite2}
	err := db.WriteAll(items)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestWriteAllEmpty(t *testing.T) {
	db := Database[MockupStorable]{path: "writeAllDB.json"}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./writeAllDB.json")
		if err != nil {
			panic(err)
		}
	})

	items := []MockupStorable{}
	err := db.WriteAll(items)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestClear(t *testing.T) {
	db := Database[MockupStorable]{}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./db.json")
		if err != nil {
			panic(err)
		}
	})

	itemToWrite := MockupStorable{0, "test"}
	itemToWrite2 := MockupStorable{1, "test2"}
	items := []MockupStorable{itemToWrite, itemToWrite2}
	err := db.WriteAll(items)
	if err != nil {
		t.Error(err.Error())
	}

	db.Clear()

	fileStat, err := db.file.Stat()
	if err != nil {
		t.Error(err.Error())
	}
	filledSize := fileStat.Size()

	if filledSize != 0 {
		t.Errorf("cleaning database failed. Expected 0 bytes inside the file, but received %d", filledSize)
	}
}

func TestAppend(t *testing.T) {
	db := Database[MockupStorable]{}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./db.json")
		if err != nil {
			panic(err)
		}
	})

	itemToAppend := MockupStorable{1, "appended item"}
	result, err := db.Append(itemToAppend)
	if err != nil {
		t.Error(err.Error())
	}

	if itemToAppend.Id == result.Id {
		t.Errorf("append method should override id to ensure there are no repeated elements")
	}
}

func TestGetAll(t *testing.T) {
	db := Database[MockupStorable]{}
	db.Open()
	defer db.Close()

	t.Cleanup(func() {
		err := os.Remove("./db.json")
		if err != nil {
			panic(err)
		}
	})

	items := []MockupStorable{
		{1, "Test item 1"},
		{2, "Test item 2"},
	}
	db.WriteAll(items)

	result, err := db.GetAll()
	if err != nil {
		t.Error(err.Error())
	}
	
	if len(result) != len(items) {
		t.Errorf("expected %d items but received %d",len(items),len(result))
	}
}
