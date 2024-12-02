package db

const (
	DBNAME     = "hotel-reservation-db"
	DBURI      = "mongodb://localhost:27017"
	TESTDBName = "hotel-reservation-test"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
