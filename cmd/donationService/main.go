package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/sshaparenko/donation-service/pkg/database"
	"github.com/sshaparenko/donation-service/pkg/routes"
)

const DEFAULT_PORT string = "8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Second * 3))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://play.google.com/*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           120,
	}))

	routes.SetupRotes(r)
	printAscii()
	printStartupMessages()

	database.Init()
	defer database.DB.Close()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", DEFAULT_PORT), r))
}

func printStartupMessages() {
	log.Printf("Starting service at port %s", DEFAULT_PORT)
	log.Printf("Handlers: ")
	log.Printf("Threads: %d", runtime.GOMAXPROCS(0))
	log.Printf("PID: %d", os.Getpid())
}

func printAscii() {
	b, err := os.ReadFile("../../common/ascii.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
