package lodgify

type Booking struct {
	ID           int     `json:"id"`
	Arrival      string  `json:"arrival"`
	Departure    string  `json:"departure"`
	PropertyID   int     `json:"property_id"`
	TotalAmount  float64 `json:"total_amount"`
	PropertyName string  `json:"property_name",omitempty`
}

type BookingsResponse struct {
	Count int       `json:"count"`
	Items []Booking `json:"items"`
}

type Property struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
