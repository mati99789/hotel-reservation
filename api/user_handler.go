package api

import (
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "James",
		LastName:  "",
	}
	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("James")
}
