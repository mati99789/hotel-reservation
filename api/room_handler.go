package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
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

func (p BookRoomParams) Validate() error {
	now := time.Now()
	if now.After(p.From) || now.After(p.To) {
		return fmt.Errorf("can't book a room in the past")
	}

	if p.To.Before(p.From) || p.To.Equal(p.From) {
		return fmt.Errorf("check-out date must be after check-in date")
	}

	if p.NumPersons <= 0 {
		return fmt.Errorf("number of persons must be positive")
	}

	return nil
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

	if err := params.Validate(); err != nil {
		return err
	}

	roomId, err := primitive.ObjectIDFromHex(c.Params("id"))
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

	ok, err = h.isRoomAvailable(c.Context(), roomId, params)

	if err != nil {
		return err
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("Room %s is already booked for the selected dates", roomId.Hex()),
		})
	}

	booking := types.Booking{
		RoomID:     roomId,
		UserID:     userID,
		NumPersons: params.NumPersons,
		From:       params.From,
		To:         params.To,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}

	return c.Status(fiber.StatusOK).JSON(inserted)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{"room_id": roomID, "from": bson.M{"$gte": params.From}, "to": bson.M{"$lte": params.To}}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
