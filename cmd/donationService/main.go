package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sshaparenko/donation-service/internal/routes"
)

const DEFAULT_PORT string = "8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	routes.SetupRotes(r)
	printStrtupMessages()

	http.ListenAndServe(fmt.Sprintf(":%s", DEFAULT_PORT), r)
}

func printStrtupMessages() {
	log.Printf("Starting Donation Service at port %s", DEFAULT_PORT)
	log.Printf("Handlers: ")
	log.Printf("Threads: %d", runtime.GOMAXPROCS(0))
	log.Printf("PID: %d", os.Getpid())
}
