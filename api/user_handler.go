package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/db"
	"hotelReservetion/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		ctx = context.Background()
		id  = c.Params("id")
	)

	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "James",
		LastName:  "",
	}
	return c.JSON(user)
}
