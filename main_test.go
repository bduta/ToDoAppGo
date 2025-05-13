package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"newtodoapp/engine"
)

func TestCreateEndpoint(t *testing.T) {
	// Change the ToDoListFileName to use a test file
	engine.ToDoListFileName = "ToDoList_test.txt"
	defer os.Remove(engine.ToDoListFileName) // Clean up test file after test

	server := initialize()

	// Start the server in a separate goroutine
	srv := &http.Server{Addr: ":5000", Handler: server}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Ensure the server is stopped after the test
	defer func() {
		if err := srv.Close(); err != nil {
			t.Fatalf("Failed to stop server: %v", err)
		}
	}()

	// Wait briefly to ensure the server is running
	time.Sleep(100 * time.Millisecond)

	var wg sync.WaitGroup
	client := &http.Client{}

	// Send 1000 POST requests concurrently
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			name := "test"
			description := "des"
			item := map[string]string{
				"name":        name,
				"description": description,
			}

			jsonData, err := json.Marshal(item)
			if err != nil {
				t.Errorf("Failed to marshal JSON: %v", err)
				return
			}

			req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/create", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusCreated {
				t.Errorf("Unexpected status code: got %v, want %v", resp.StatusCode, http.StatusCreated)
			}
		}(i)
	}

	wg.Wait()

	// Verify the test file contains 1000 entries
	data, err := ioutil.ReadFile(engine.ToDoListFileName)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	lines := bytes.Split(data, []byte("\n"))
	if len(lines)-1 != 10 { // Subtract 1 for the trailing newline
		t.Errorf("Unexpected number of entries: got %v, want %v", len(lines)-1, 10)
	}
}
