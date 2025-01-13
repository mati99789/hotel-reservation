package api

import (
	"hotelReservetion/db"
	"hotelReservetion/shared"
	"hotelReservetion/types"
	"hotelReservetion/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := shared.Map{"hotelId": oid}
	rooms, err := s.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	pagination := utils.ExtractPaginationFromRequest(c)

	hotels, err := h.store.Hotel.GetHotels(c.Context(), nil, pagination)
	if err != nil {
		return err
	}

	response := utils.NewResourceResponse(len(hotels), int(pagination.Page), hotels)
	return c.JSON(response)
}

func (h *HotelHandler) HandleHotelUpdate(c *fiber.Ctx) error {
	hotelID := c.Params("id")

	if hotelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hotel id is required"})
	}

	var hotel types.Hotel

	if err := c.BodyParser(&hotel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	updatedHotel, err := h.store.Hotel.UpdateHotel(c.Context(), &hotel, hotelID)
	if err != nil {
		return err
	}

	return c.JSON(updatedHotel)
}
