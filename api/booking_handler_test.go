package api

import (
	"fmt"
	"hotelReservetion/db/fixtures"
	"hotelReservetion/types"
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

	fmt.Println(booking)
}
