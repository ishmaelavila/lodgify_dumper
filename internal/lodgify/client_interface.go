package lodgify

type LodgifyConnector interface {
	GetBookings() (*BookingsResponse, error)
}
