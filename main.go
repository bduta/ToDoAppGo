package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"

	"newtodoapp/server"

	"github.com/google/uuid"
)

func main() {
	server := initialize()
	go actor()
	log.Fatal(http.ListenAndServe(":5000", server))
}

func initialize() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/create", traceMiddleware(http.HandlerFunc(createHandler)))
	router.Handle("/get", traceMiddleware(http.HandlerFunc(getHandler)))
	router.Handle("/update", traceMiddleware(http.HandlerFunc(updateHandler)))
	router.Handle("/delete", traceMiddleware(http.HandlerFunc(deleteHandler)))
	router.Handle("/list", traceMiddleware(http.HandlerFunc(listHandler)))

	router.Handle("/about", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("static", "about.html")
		http.ServeFile(w, r, filePath)
	}))

	return router
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

type request struct {
	request     *http.Request
	response    http.ResponseWriter
	signalChan  chan bool
	requestType string
}

const (
	getRequest    = "get"
	createRequest = "create"
	updateRequest = "update"
	deleteRequest = "delete"
	listRequest   = "list"
)

var requestChan = make(chan request, 10)

func getHandler(w http.ResponseWriter, r *http.Request) {
	sendRequestToActor(w, r, getRequest)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	sendRequestToActor(w, r, createRequest)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	sendRequestToActor(w, r, updateRequest)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	sendRequestToActor(w, r, deleteRequest)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	sendRequestToActor(w, r, listRequest)
}

func sendRequestToActor(w http.ResponseWriter, r *http.Request, requestType string) {
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: requestType,
	}
	<-signalChan
}

func actor() {
	for {
		req, ok := <-requestChan
		if !ok {
			return // Exit the loop when the channel is closed
		}
		switch req.requestType {
		case getRequest:
			server.GetItems(req.response, req.request)
		case createRequest:
			server.CreateItem(req.response, req.request)
		case updateRequest:
			server.UpdateItem(req.response, req.request)
		case deleteRequest:
			server.DeleteItem(req.response, req.request)
		case listRequest:
			server.ListItems(req.response, req.request)
		}
		req.signalChan <- true
		close(req.signalChan)
	}
}
