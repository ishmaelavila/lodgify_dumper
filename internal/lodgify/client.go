package lodgify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

type LodgifyClient struct {
	BaseURL    string
	APIKey     string
	HttpClient *http.Client
	Logger     *zap.SugaredLogger
}

type LodgifyClientArgs struct {
	BaseURL string
	APIKey  string
	Logger  *zap.SugaredLogger
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
		Logger:     args.Logger,
	}

	return &lodgifyClient, nil
}

func (lodgify *LodgifyClient) GetBookings() (*BookingsResponse, error) {
	req, err := http.NewRequest("GET", lodgify.BaseURL+"/"+getBookingsURL, nil)
	lodgify.addDefaultHeaders(req)
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

	if resp.StatusCode > 201 {
		lodgify.Logger.Errorf("status code: %d", resp.StatusCode)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lodgify.Logger.Errorf("could not decode response body")
		}
		lodgify.Logger.Errorf("error response %s", string(bodyBytes))
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(&bookings)

	if decodeErr != nil {
		return nil, fmt.Errorf("getbookings error decoding response: %w", decodeErr)
	}

	return &bookings, nil
}

func (lodgify *LodgifyClient) addDefaultHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKey", lodgify.APIKey)
}
