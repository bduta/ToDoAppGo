package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"encoding/json"
	"newtodoapp/models"
	"newtodoapp/server"

	_ "net/http/pprof"

	"github.com/google/uuid"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	server := initialize()
	go actor()
	log.Fatal(http.ListenAndServe(":5000", server.Handler))
}

func initialize() *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/create", createHandler)
	router.HandleFunc("/get", getHandler)
	router.HandleFunc("/update", updateHandler)
	router.HandleFunc("/delete", deleteHandler)
	router.HandleFunc("/list", listHandler)

	router.HandleFunc("/about", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("static", "about.html")
		http.ServeFile(w, r, filePath)
	}))

	server := &http.Server{
		Addr:    ":5000",
		Handler: traceMiddleware(router),
	}
	return server
}

func traceMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), "traceID", traceID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

type actorRequest struct {
	request      models.ToDoItem
	responseChan chan map[string]string
	signalChan   chan bool
	requestType  string
}

const (
	getRequest    = "get"
	createRequest = "create"
	updateRequest = "update"
	deleteRequest = "delete"
	listRequest   = "list"
)

var requestChan = make(chan actorRequest, 10)

const jsonContentType = "application/json"

func getHandler(w http.ResponseWriter, r *http.Request) {

	responseChan := make(chan map[string]string)
	signalChan := make(chan bool)
	requestChan <- actorRequest{
		responseChan: responseChan,
		signalChan:   signalChan,
		requestType:  getRequest,
	}
	items := <-responseChan
	<-signalChan

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

func createHandler(w http.ResponseWriter, r *http.Request) {
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

	responseChan := make(chan map[string]string)
	signalChan := make(chan bool)
	requestChan <- actorRequest{
		request:      input,
		responseChan: responseChan,
		signalChan:   signalChan,
		requestType:  createRequest,
	}
	<-signalChan

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "ToDo item created successfully", "traceID": "` + traceID + `"}`))
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
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

	signalChan := make(chan bool)
	requestChan <- actorRequest{
		request:     input,
		signalChan:  signalChan,
		requestType: updateRequest,
	}
	<-signalChan

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item updated successfully", "traceID": "` + traceID + `"}`))
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
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

	signalChan := make(chan bool)
	requestChan <- actorRequest{
		request:     input,
		signalChan:  signalChan,
		requestType: deleteRequest,
	}
	<-signalChan

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "ToDo item deleted successfully", "traceID": "` + traceID + `"}`))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Context().Value("traceID").(string)

	tmpl, err := template.ParseFiles("templates/todos.gohtml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "error parsing template ` + err.Error() + `", "traceID": "` + traceID + `"}`))
		return
	}

	responseChan := make(chan map[string]string)
	signalChan := make(chan bool)
	requestChan <- actorRequest{
		responseChan: responseChan,
		signalChan:   signalChan,
		requestType:  listRequest,
	}
	items := <-responseChan
	<-signalChan

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

func actor() {
	for {
		req, ok := <-requestChan
		if !ok {
			return // Exit the loop when the channel is closed
		}
		switch req.requestType {
		case getRequest:
			req.responseChan <- server.GetItems()
		case createRequest:
			server.UpdateItem(req.request.Name, req.request.Description)
		case updateRequest:
			server.UpdateItem(req.request.Name, req.request.Description)
		case deleteRequest:
			server.DeleteItem(req.request.Name)
		case listRequest:
			req.responseChan <- server.GetItems()
		}
		req.signalChan <- true
		close(req.responseChan)
		close(req.signalChan)
	}
}
