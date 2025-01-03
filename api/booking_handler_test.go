package api

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"hotelReservetion/api/middleware"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup()

	defer db.tearddown(t)

	admin := fixtures.AddUser(db.Store, "admin@o2.pl", "Admin", "Admin", types.AdminRole)
	hotel := fixtures.AddHotel(db.Store, "Bellucia", "France", 2, nil)
	room := fixtures.AddRoom(db.Store, types.Normal, false, 300, hotel.ID)

	booking := fixtures.AddBooking(db.Store, room.ID, admin.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))

	_ = booking

	app := fiber.New()
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentications)
	adminGroup := apiv1.Group("/admin", middleware.AuthorizeRole(types.AdminRole))

	bookingHandler := NewBookingHandler(db.Store)
	adminGroup.Get("/users", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	var bookings []*types.Booking
	if err := json.NewDecoder(response.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	fmt.Println(bookings)

}
