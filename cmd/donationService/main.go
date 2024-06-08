package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/donation-service/internal/routes"
	"github.com/sshaparenko/donation-service/internal/utils"
)

const DEFAULT_PORT = "8080"

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

	app.Listen(fmt.Sprintf(":%s", PORT))
}
