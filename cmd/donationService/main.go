package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/donation-service/internal/routes"
	"github.com/sshaparenko/donation-service/internal/utils"
)

/*
DEFAULT_PORT stores the default port value
*/
const DEFAULT_PORT = "8081"

/*
NewFilberApp creates a new fiber App and sets up the routes for it
*/
func NewFilberApp() *fiber.App {
	var app *fiber.App = fiber.New()
	routes.SetupRoutes(app)
	return app
}

func main() {
	var app *fiber.App = NewFilberApp()

	var PORT string = utils.GetValue("PORT")
	if PORT == "" {
		PORT = DEFAULT_PORT
	}

	if err := app.Listen(fmt.Sprintf(":%s", PORT)); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
