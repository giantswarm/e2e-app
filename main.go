package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/giantswarm/micrologger"
	yaml "gopkg.in/yaml.v2"
)

const (
	Port = "8000"
)

var (
	description = "Application run for integration tests within a e2e test environment on Kubernetes."
	gitCommit   = "n/a"
	name        = "e2e-app"
	source      = "https://github.com/giantswarm/e2e-app"
)

func main() {
	if (len(os.Args) > 1) && (os.Args[1] == "version") {
		d, err := yaml.Marshal(newVersionResponse())
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", d)

		return
	}

	var lHelp bool
	var sHelp bool
	flag.BoolVar(&lHelp, "help", false, "Print help usage.")
	flag.BoolVar(&sHelp, "h", false, "Print help usage.")
	flag.Parse()
	if lHelp || sHelp {
		flag.Usage()
		return
	}

	ctx := context.Background()
	mux := http.NewServeMux()
	logger, err := micrologger.New(micrologger.Config{})
	if err != nil {
		panic(err)
	}

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newVersionResponse())
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
		logger.LogCtx(ctx, "level", "debug", "message", fmt.Sprintf("running server at http://0.0.0.0:%s", Port))
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

	logger.LogCtx(ctx, "level", "debug", "message", "received termination signal")
	logger.LogCtx(ctx, "level", "debug", "message", "draining server connections")

	ctx, _ = context.WithTimeout(ctx, 5*time.Second)
	server.Shutdown(ctx)

	logger.LogCtx(ctx, "level", "debug", "message", "shutting down")
}
