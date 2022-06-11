package lodgify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBookings(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), "https://api.lodgify.com/v2/reservations/bookings")
		// Send response to be tested
		rw.Write([]byte(`OK`))
	}))
	// Close the server when test finishes
	defer server.Close()

}
