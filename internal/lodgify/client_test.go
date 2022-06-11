package lodgify_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBookings(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

	}))
	// Close the server when test finishes
	defer server.Close()

}
