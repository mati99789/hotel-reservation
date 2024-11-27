package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelReservetion/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (s *HotelHandler) HandleGetHotelByID(c *fiber.Ctx) error {
	hotelId := c.Params("id")

	if hotelId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Hotel id is required",
		})
	}

	hotel, err := s.store.Hotel.GetHotelByID(c.Context(), hotelId)

	if err != nil {
		return err
	}

	return c.JSON(hotel)
}

func (s *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
	}
	filter := bson.M{"hotelId": oid}
	rooms, err := s.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}

	return c.JSON(hotels)
}
