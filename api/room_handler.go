package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservetion/db"
	"hotelReservetion/types"
	"time"
)

type BookRoomParams struct {
	From       time.Time `json:"from" time_format:"2006-01-02T15:04:05"`
	To         time.Time `json:"to" time_format:"2006-01-02T15:04:05"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	roomId := c.Params("id")

	ooid, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user claims",
		})
	}
	userIdInterface, ok := claims["id"]
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user id"})
	}

	userId, ok := userIdInterface.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user id"})
	}

	userID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	booking := types.Booking{
		RoomID:     ooid,
		UserID:     userID,
		NumPersons: params.NumPersons,
		From:       params.From,
		To:         params.To,
	}

	fmt.Printf("%+v\n", booking)
	return nil
}
