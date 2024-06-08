package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/donation-service/internal/handlers"
)

/*
SetupRoutes is setting up the routes of an API
*/
func SetupRoutes(app *fiber.App) {
	var publicRoutes fiber.Router = app.Group("api/v1")

	publicRoutes.Get("/items", handlers.GetAllItems)
}
