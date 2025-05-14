package main

import (
	engine "newtodoapp/engine"

	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	traceID := uuid.New().String()
	type contextKey string
	key := "TraceID"

	ctx := context.WithValue(context.Background(), contextKey(key), traceID)

	args := os.Args[1:]

	err := engine.ExecuteCommand(args)
	if err != nil {
		logger.With(key, ctx.Value(key)).Error(err.Error())
		return
	}

	// Handle interrupt signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Application exiting on interrupt signal")
}
