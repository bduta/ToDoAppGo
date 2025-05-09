package server

import (
	"encoding/json"
	"net/http"

	"newtodoapp/engine"
	"newtodoapp/models"
)

const jsonContentType = "application/json"

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	error := engine.CreateItem(input.Name, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to create ToDo item: ` + error.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ToDo item created successfully", "traceID": "` + traceID + `"}`))
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	toDosJson, err := engine.GetItems()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to get ToDo items: ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(toDosJson); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	error := engine.UpdateItem(input.Id, input.Description)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to update ToDo item: ` + error.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item updated successfully", "traceID": "` + traceID + `"}`))
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	var input models.ToDoItem

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request payload", "traceID": "` + traceID + `"}`))
		return
	}

	err := engine.DeleteItem(input.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to delete ToDo item: ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item deleted successfully", "traceID": "` + traceID + `"}`))
}
