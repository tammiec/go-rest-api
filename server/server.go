package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/tammiec/go-rest-api/appcontext"
	"github.com/tammiec/go-rest-api/handlers"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	//nolint
	w.Write([]byte("Hello!\n"))
}

func NewRouter(handlers *handlers.Handlers) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", baseHandler)

	router.HandleFunc("/health", handlers.Health.Ping).Methods(http.MethodGet)

	router.HandleFunc("/users", handlers.Users.List).Methods(http.MethodGet)
	router.HandleFunc("/users", handlers.Users.Create).Methods(http.MethodPost)
	router.HandleFunc("/users/{id:[0-9]+}", handlers.Users.Get).Methods(http.MethodGet)
	router.HandleFunc("/users/{id:[0-9]+}", handlers.Users.Update).Methods(http.MethodPut)
	router.HandleFunc("/users/{id:[0-9]+}", handlers.Users.Delete).Methods(http.MethodDelete)

	return router
}

func NewServer(ctx appcontext.AppContext) *http.Server {
	handlersDeps := handlers.Deps{
		Health: ctx.Services.Health,
		Users:  ctx.Services.Users,
		Render: ctx.Clients.Render,
	}

	_handlers := handlers.New(&handlersDeps)

	router := NewRouter(_handlers)

	return NewServerWithRouter(ctx.Config.HTTPPort, router)
}

func NewServerWithRouter(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		// https://en.wikipedia.org/wiki/Slowloris_(computer_security)
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
}
