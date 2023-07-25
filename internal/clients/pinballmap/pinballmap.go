package pinballmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pinman/internal/utils"
)

type Client struct {
	apiHost   string
	apiScheme string
	client    utils.HttpDoer
}

type ClientInterface interface {
	GetLocations(nameFilter string) ([]Location, error)
	GetLocation(id int) (*Location, error)
}

func NewClient() *Client {
	return &Client{
		apiScheme: "https",
		apiHost:   "pinballmap.com",
		client:    http.DefaultClient,
	}
}

func NewClientWithHttpClient(client utils.HttpDoer) *Client {
	return &Client{
		apiScheme: "https",
		apiHost:   "pinballmap.com",
		client:    client,
	}
}

func (c *Client) NewUrl(path string, params url.Values) *url.URL {
	urrl := &url.URL{
		Scheme: c.apiScheme,
		Host:   c.apiHost,
		Path:   path,
	}

	if params != nil {
		urrl.RawQuery = params.Encode()
	}

	return urrl
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

type LocationResponse struct {
	Locations []Location `json:"locations"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

// GetLocations retrieves all locations from the Pinball Map API matching the given name filter.
// https://pinballmap.com/api/v1/docs/1.0/locations/index.html
func (c *Client) GetLocations(nameFilter string) ([]Location, error) {
	params := url.Values{}
	params.Add("by_location_name", nameFilter)
	params.Add("no_details", "true")

	req, err := http.NewRequest("GET", c.NewUrl("/api/v1/locations.json", params).String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}

		return nil, fmt.Errorf("unexpected response: %s", errResp.Errors)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var locations LocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("retrieving locations: %w", err)
	}

	return locations.Locations, nil
}

func (c *Client) GetLocation(id int) (*Location, error) {
	req, err := http.NewRequest("GET", c.NewUrl(fmt.Sprintf("/api/v1/locations/%d.json", id), nil).String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errResp := ErrorResponse{}
		err := json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("decoding error response: %w", err)
		}

		return nil, fmt.Errorf("unexpected response: %s", errResp.Errors)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	var location Location
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &location, nil
}
