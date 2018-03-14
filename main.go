package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giantswarm/micrologger"
)

const (
	Port = "8080"
)

func main() {
	ctx := context.Background()
	mux := http.NewServeMux()
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(err)
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World")
	}))
	mux.Handle("/delay/1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		io.WriteString(w, "Hello World")
	}))
	mux.Handle("/delay/5", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		io.WriteString(w, "Hello World")
	}))
	mux.Handle("/delay/10", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		io.WriteString(w, "Hello World")
	}))

	server := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}

	go func() {
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("running server at %s", server.Addr))
		err := server.ListenAndServe()
		if IsServerClosed(err) {
			// fall through
		} else if err != nil {
			panic(err)
		}
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan

	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	server.Shutdown(ctx)
}
