package pinballmap

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"pinman/internal/clients/generic"
)

const apiHost = "https://pinballmap.com"

type Client struct {
	apiHost       string
	genericClient generic.ClientInterface
}

type ClientInterface interface {
	GetLocations(nameFilter string) ([]Location, error)
	GetLocation(id int) (*Location, error)
}

func NewClient() *Client {
	return &Client{
		apiHost:       apiHost,
		genericClient: generic.NewClient(),
	}
}

func NewClientWithGenericClient(genericClient generic.ClientInterface) *Client {
	return &Client{
		apiHost:       apiHost,
		genericClient: genericClient,
	}
}

type Location struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	NumMachines int    `json:"num_machines"`
}

type LocationsResponse struct {
	Locations []Location `json:"locations"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

// GetLocations retrieves all locations from the Pinball Map API matching the given name filter.
// https://pinballmap.com/api/v1/docs/1.0/locations/index.html
func (c *Client) GetLocations(nameFilter string) ([]Location, error) {
	params := &url.Values{}
	params.Add("by_location_name", nameFilter)
	params.Add("no_details", "true")

	req, err := http.NewRequest("GET", generic.FormatRequestUrl(c.apiHost, "/api/v1/locations.json", params), nil)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	var response LocationsResponse
	var errorResponse ErrorResponse
	code, err := c.genericClient.Do(req, &response, &errorResponse)
	if err != nil {
		log.Error().Err(err).Msg("failed to get locations")
		return nil, fmt.Errorf("performing request: %w", err)
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("request failed with error: %s", errorResponse.Errors)
	}

	return response.Locations, nil
}

func (c *Client) GetLocation(id int) (*Location, error) {
	req, err := http.NewRequest(
		"GET",
		generic.FormatRequestUrl(c.apiHost, fmt.Sprintf("/api/v1/locations/%d.json", id)),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	var location Location
	var errorResponse ErrorResponse
	code, err := c.genericClient.Do(req, &location, &errorResponse)
	if err != nil {
		log.Error().Err(err).Int("pinballMapLocationId", id).Msg("failed to get location")
		return nil, fmt.Errorf("performing request: %w", err)
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("request failed with error: %s", errorResponse.Errors)
	}

	return &location, nil
}
