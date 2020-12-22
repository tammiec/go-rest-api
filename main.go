package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/tammiec/go-rest-api/appcontext"
	"github.com/tammiec/go-rest-api/config"
	"github.com/tammiec/go-rest-api/server"
)

func main() {
	_config := config.Loader()
	appContext := appcontext.NewAppContext(_config)
	serverUsers := server.NewServer(*appContext)
	start("users", serverUsers)
	waitForSigInt(serverUsers)
}

func waitForSigInt(httpServer *http.Server) {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	//nolint
	httpServer.Shutdown(ctx)

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

func start(name string, server *http.Server) {
	go func() {
		log.Printf("%s listening at %s\n", name, server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}
