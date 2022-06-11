package lodgify

type Booking struct {
	ID          int     `json:"id"`
	Arrival     string  `json:"arrival"`
	Departure   string  `json:"departure"`
	PropertyID  int     `json:"property_id"`
	TotalAmount float64 `json:"total_amount"`
}

type BookingsResponse struct {
	Count int       `json:"count"`
	Items []Booking `json:"items"`
}
