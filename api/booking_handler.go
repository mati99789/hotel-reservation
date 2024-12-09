package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"hotelReservetion/db"
	"hotelReservetion/utils"
	"net/http"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return nil
	}

	user, err := utils.GetUserFromClaims(claims)
	if err != nil {
		return err
	}

	if user.ID != booking.UserID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)

	if err != nil {
		return err
	}
	claims, ok := c.Locals("claims").(jwt.MapClaims)
	if !ok {
		return err
	}

	user, err := utils.GetUserFromClaims(claims)

	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	err = h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true})

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(genericResp{
		Type: "msg",
		Msg:  "Updated successfully",
	})
}
