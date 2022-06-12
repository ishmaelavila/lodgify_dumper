package lodgify

type LodgifyConnector interface {
	GetBookings() ([]Booking, error)
}
