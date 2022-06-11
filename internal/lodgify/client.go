package lodgify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type LodgifyClient struct {
	BaseURL    string
	APIKey     string
	HttpClient *http.Client
}

type LodgifyClientArgs struct {
	BaseURL string
	APIKey  string
}

var (
	getBookingsURL string = "reservations/bookings"
)

func NewClient(args LodgifyClientArgs) (LodgifyConnector, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(5) * time.Second,
	}

	lodgifyClient := LodgifyClient{
		HttpClient: httpClient,
		BaseURL:    args.BaseURL,
		APIKey:     args.APIKey,
	}

	return &lodgifyClient, nil
}

func (lodgify *LodgifyClient) GetBookings() (*BookingsResponse, error) {
	req, err := http.NewRequest("GET", lodgify.BaseURL+"/"+getBookingsURL, nil)
	bookings := BookingsResponse{}

	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("page", "1")
	req.URL.RawQuery = params.Encode()

	resp, err := lodgify.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(bookings)

	if decodeErr != nil {
		return nil, nil
	}

	return &bookings, nil
}

func (lodgify *LodgifyClient) addDefaultHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKey", lodgify.APIKey)
}
