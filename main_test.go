package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func TestCreateEndpoint(t *testing.T) {

	router := initialize()
	ts := httptest.NewServer(router.Handler)
	defer ts.Close()

	go actor()

	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			name := "test" + strconv.Itoa(index)
			description := "des" + strconv.Itoa(index)
			item := map[string]string{
				"name":        name,
				"description": description,
			}

			jsonData, err := json.Marshal(item)
			if err != nil {
				t.Errorf("Failed to marshal JSON: %v", err)
				return
			}

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Errorf("Failed to create request: %v", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			ctx := context.WithValue(req.Context(), "traceID", "1")
			createHandler(w, req.WithContext(ctx))
		}(i)
	}

	wg.Wait()

	wg.Add(1)
	items := make(map[string]string)
	go func() {
		defer wg.Done()
		w := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, "/get", nil)
		if err != nil {
			t.Errorf("Failed to create request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), "traceID", "1")
		getHandler(w, req.WithContext(ctx))

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", w.Code)
			return
		}
		if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
			t.Errorf("Failed to decode response: %v", err)
			return
		}
	}()

	wg.Wait()

	if len(items) != 10000 {
		t.Errorf("Expected 10000 items, got %d", len(items))
		return
	}

	close(requestChan)
}
