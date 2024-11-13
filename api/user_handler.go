package api

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"hotelReservetion/utils"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		update bson.M
		userID = c.Params("id")
	)

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	if err := c.BodyParser(&update); err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUserById(c.Context(), filter, update); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"Message": "User updated!"})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	_, err := h.userStore.DeleteUserById(c.Context(), userId)

	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIErrors{
			Errors: errors,
		})
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)

	if err != nil {
		return err
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)

	user, err := h.userStore.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, fiber.ErrNotFound) {

			return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{
				"Status":  "Error",
				"Message": "User not found",
			})
		}
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(users)
}
