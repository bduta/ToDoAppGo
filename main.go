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
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: getRequest,
	}
	<-signalChan
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: createRequest,
	}
	<-signalChan
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: updateRequest,
	}
	<-signalChan
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: deleteRequest,
	}
	<-signalChan
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	signalChan := make(chan bool)
	requestChan <- request{
		request:     r,
		response:    w,
		signalChan:  signalChan,
		requestType: listRequest,
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
