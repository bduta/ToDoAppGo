package engine

import (
	"newtodoapp/models"
	"os"
	"testing"
)

func TestCreateTheToDoListFileIfNeeded(t *testing.T) {
	ToDoListFileName = "ToDoList_test_create.txt"
	defer os.Remove(ToDoListFileName)

	os.Remove(ToDoListFileName)
	created, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !created {
		t.Fatalf("Expected file to be created, but it wasn't")
	}

	file, err := os.Open(ToDoListFileName)
	if err != nil {
		t.Fatalf("Expected file to exist, but got error: %v", err)
	}
	file.Close()
}

func TestGenerateItemId(t *testing.T) {
	ToDoListFileName = "ToDoList_test_generate_id.txt"
	defer os.Remove(ToDoListFileName)

	_, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	id, err := generateItemId(false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("Expected ID to be 1, got %d", id)
	}
}

func TestWriteAndReadExistingList(t *testing.T) {
	ToDoListFileName = "ToDoList_test_write_read.txt"
	defer os.Remove(ToDoListFileName)

	_, err := createTheToDoListFileIfNeeded()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	item := models.ToDoItem{Id: 1, Name: "Test", Description: "Test description"}
	err = writeItemToFile(item)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	list, err := readExistingList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("Expected list length to be 1, got %d", len(list))
	}

	if list[0].Name != "Test" || list[0].Description != "Test description" {
		t.Fatalf("Expected item to match, got %+v", list[0])
	}
}
