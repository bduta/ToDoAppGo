package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"newtodoapp/engine"
)

func TestCreateEndpoint(t *testing.T) {

	engine.ToDoListFileName = "ToDoList_test.txt"
	defer os.Remove(engine.ToDoListFileName)

	router := initialize()
	ts := httptest.NewServer(router.Handler)
	defer ts.Close()

	go actor()

	var wg sync.WaitGroup
	client := &http.Client{}

	// Send 100 POST requests concurrently
	for i := 0; i < 200; i++ {
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

			req, err := http.NewRequest(http.MethodPost, ts.URL+"/create", bytes.NewBuffer(jsonData))
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
	close(requestChan)

	data, err := os.ReadFile(engine.ToDoListFileName)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	lines := bytes.Split(data, []byte("\n"))
	if len(lines)-1 != 200 { // Subtract 1 for the trailing newline
		t.Errorf("Unexpected number of entries: got %v, want %v", len(lines)-1, 200)
	}
}
