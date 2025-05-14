package server

import (
	"encoding/json"
	"html/template"
	"net/http"

	"newtodoapp/models"
)

const jsonContentType = "application/json"

var items = make(map[string]string)

func CreateItem(w http.ResponseWriter, r *http.Request) {
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

	items[input.Name] = input.Description

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ToDo item created successfully", "traceID": "` + traceID + `"}`))
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)
	w.Header().Set("Content-Type", jsonContentType)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Method not allowed", "traceID": "` + traceID + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateItem(w http.ResponseWriter, r *http.Request) {
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

	items[input.Name] = input.Description

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item updated successfully", "traceID": "` + traceID + `"}`))
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
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

	delete(items, input.Name)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item deleted successfully", "traceID": "` + traceID + `"}`))
}

func ListItems(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)

	tmpl, err := template.ParseFiles("templates/todos.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error parsing template ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	var todoItems []models.ToDoItem
	for name, description := range items {
		todoItems = append(todoItems, models.ToDoItem{
			Name:        name,
			Description: description,
		})
	}

	err = tmpl.Execute(w, todoItems)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error: ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

}
