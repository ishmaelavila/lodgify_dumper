package lodgify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
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
	getBookingsURL     string = "reservations/bookings"
	getPropertyByIDURL string = "properties"
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

//GetBookings will attempt to retrieve ALL bookings for the account, this may result in multiple requests. Additionally, it resolves the related Property which results in a extra request.
func (lodgify *LodgifyClient) GetBookings() ([]Booking, error) {
	req, err := http.NewRequest("GET", lodgify.BaseURL+"/"+getBookingsURL, nil)
	lodgify.addDefaultHeaders(req)

	bookings := []Booking{}

	if err != nil {
		return nil, err
	}

	properties := map[int]Property{}

	for i := 1; ; i++ {
		params := url.Values{}
		params.Set("page", strconv.Itoa(i))
		req.URL.RawQuery = params.Encode()

		lodgify.Logger.Infof("getbookings fetching page %d", i)
		resp, err := lodgify.HttpClient.Do(req)

		if err != nil {
			return nil, fmt.Errorf("lodgify getbookings error creating request on page %d", i)
		}

		if resp.StatusCode > 201 {
			return nil, lodgify.handleErrorResponse(resp)
		}

		decodedResponse := BookingsResponse{}
		decodeErr := json.NewDecoder(resp.Body).Decode(&decodedResponse)

		if decodeErr != nil {
			return nil, fmt.Errorf("getbookings error decoding response: %w", decodeErr)
		}

		if len(decodedResponse.Items) < 1 {
			break
		}

		for j := 0; j < len(decodedResponse.Items); j++ {

			booking := &decodedResponse.Items[j]
			propID := booking.PropertyID
			property, ok := properties[propID]

			if !ok {
				propertyLookedUp, err := lodgify.GetPropertyByID(propID)

				if err != nil {
					lodgify.Logger.Errorf("error retrieving property with id %d for booking %d: %w", propID, booking.ID, err)
				}

				decodedResponse.Items[j].PropertyName = propertyLookedUp.Name
				properties[propID] = *propertyLookedUp

			} else {
				decodedResponse.Items[j].PropertyName = property.Name
			}

		}

		bookings = append(bookings, decodedResponse.Items...)
		time.Sleep(300 * time.Millisecond)
	}

	lodgify.Logger.Infof("getbookings finished with %d bookings", len(bookings))

	return bookings, nil
}

func (lodgify *LodgifyClient) GetPropertyByID(id int) (*Property, error) {
	req, err := http.NewRequest("GET", lodgify.BaseURL+"/"+getPropertyByIDURL+"/"+strconv.Itoa(id), nil)
	lodgify.addDefaultHeaders(req)

	if err != nil {
		return nil, fmt.Errorf("lodgify getPropertyByID error creating request %w", err)
	}

	resp, err := lodgify.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("lodgify getPropertyByID error making request: %w", err)
	}

	if resp.StatusCode > 201 {
		return nil, lodgify.handleErrorResponse(resp)
	}

	property := Property{}
	decodeErr := json.NewDecoder(resp.Body).Decode(&property)

	if decodeErr != nil {
		return nil, fmt.Errorf("getPropertyByID error decoding response: %w", decodeErr)
	}

	return &property, nil
}

func (lodgify *LodgifyClient) handleErrorResponse(resp *http.Response) error {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("lodgify getbookings could not decode response body: %w", err)
	}
	return fmt.Errorf("lodgify getbookings returned an error response %d: %s", resp.StatusCode, string(bodyBytes))

}

func (lodgify *LodgifyClient) addDefaultHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-ApiKey", lodgify.APIKey)
}
