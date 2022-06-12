package lodgify_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ishmaelavila/lodgify_dumper/internal/lodgify"
	"go.uber.org/zap"
)

func getPropertyResponse(id int, name string) ([]byte, error) {

	property := lodgify.Property{
		ID:   id,
		Name: name,
	}

	resp, err := json.Marshal(property)

	return resp, err
}

func getBookingsResponseWithOneBooking(id int, propertyID int) ([]byte, error) {

	bookingResp := lodgify.BookingsResponse{
		Items: []lodgify.Booking{
			{
				ID:         id,
				PropertyID: propertyID,
			},
		},
		Count: 1,
	}

	respJson, err := json.Marshal(bookingResp)
	return respJson, err
}

func getEmptyBookingsResponse() ([]byte, error) {
	bookingResp := lodgify.BookingsResponse{
		Items: []lodgify.Booking{},
		Count: 0,
	}

	respJson, err := json.Marshal(bookingResp)
	return respJson, err
}
func TestGetBookings(t *testing.T) {
	// Start a local HTTP server

	bookingID := 123
	propertyID := 321
	propertyName := "Cool Place"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		params := req.URL.Query()
		if strings.Contains(req.URL.String(), "properties") {
			resp, err := getPropertyResponse(propertyID, propertyName)
			if err != nil {
				t.Fatalf("error creating property response: %v", err)
			}
			fmt.Fprint(writer, string(resp))
			return
		}

		var bookingResponse []byte
		var err error

		if params["page"][0] == "1" {
			bookingResponse, err = getBookingsResponseWithOneBooking(bookingID, propertyID)
			if err != nil {
				t.Fatalf("error creating bookings response: %v", err)
			}
		} else {
			bookingResponse, err = getEmptyBookingsResponse()
			if err != nil {
				t.Fatalf("error creating bookings response: %v", err)
			}
		}

		fmt.Fprint(writer, string(bookingResponse))

	}))

	// Close the server when test finishes
	defer server.Close()

	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugarLogger := logger.Sugar()

	lodgifyClient, err := lodgify.NewClient(lodgify.LodgifyClientArgs{
		HttpClient: server.Client(),
		Logger:     sugarLogger,
		BaseURL:    server.URL,
	})

	if err != nil {
		t.Fatalf("error constructing lodgify client: %v", err)
	}

	bookings, err := lodgifyClient.GetBookings()

	if err != nil {
		t.Fatalf("error getting bookings: %v", err)
	}

	if len(bookings) < 1 {
		t.Fatal("expected 1 booking, got 0")
	}

	if bookings[0].ID != 123 {
		t.Fatalf("expected ID to be 123, got %d", bookings[0].ID)
	}

}

func TestGetPropertyByID(t *testing.T) {

	propertyID := 321
	propertyName := "Cool Place"
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {

		resp, err := getPropertyResponse(propertyID, propertyName)
		if err != nil {
			t.Fatalf("error creating property response: %v", err)
		}
		fmt.Fprint(writer, string(resp))

	}))

	defer server.Close()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugarLogger := logger.Sugar()

	lodgifyClient, err := lodgify.NewClient(lodgify.LodgifyClientArgs{
		HttpClient: server.Client(),
		Logger:     sugarLogger,
		BaseURL:    server.URL,
	})

	if err != nil {
		t.Fatalf("error constructing lodgify client: %v", err)
	}

	property, err := lodgifyClient.GetPropertyByID(propertyID)

	if err != nil {
		t.Fatalf("error getting property: %v", err)
	}

	if property.ID != 321 {
		t.Fatalf("incorrect property id returned, got %d, want %d", property.ID, propertyID)
	}

}
