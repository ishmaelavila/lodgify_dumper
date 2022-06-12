package lodgify

type LodgifyConnector interface {
	GetBookings() ([]Booking, error)
	GetPropertyByID(id int) (*Property, error)
}
