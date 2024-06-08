package handlers

import (
	"github.com/gofiber/fiber/v2"
	models "github.com/sshaparenko/donation-service/internal/models/resoponce"
	"github.com/sshaparenko/donation-service/internal/services"
)

/*
GetAllItems returns a JSON based on Responce struct with a list of data, returned from service package
*/
func GetAllItems(c *fiber.Ctx) error {
	var result []int = services.GetAllItems()
	return c.JSON(models.Responce[[]int]{
		Success: true,
		Message: "Here is your message",
		Data:    result,
	})
}
