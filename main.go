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
	server := Init()
	log.Fatal(http.ListenAndServe(":5000", server))
}

func Init() *http.ServeMux {
	router := http.NewServeMux()
	router.Handle("/create", traceMiddleware(http.HandlerFunc(server.CreateHandler)))
	router.Handle("/get", traceMiddleware(http.HandlerFunc(server.GetHandler)))
	router.Handle("/update", traceMiddleware(http.HandlerFunc(server.UpdateHandler)))
	router.Handle("/delete", traceMiddleware(http.HandlerFunc(server.DeleteHandler)))
	router.Handle("/list", traceMiddleware(http.HandlerFunc(server.ListHandler)))

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
