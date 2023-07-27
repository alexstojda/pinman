package generic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func FormatRequestUrl(host string, path string, params ...*url.Values) string {
	// ensure the path starts with a slash
	if path[0] != '/' {
		path = fmt.Sprintf("/%s", path)
	}

	// ensure the host does not end with a slash
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	r := fmt.Sprintf("%s%s", host, path)

	if len(params) > 0 && params[0] != nil {
		r = fmt.Sprintf("%s?%s", r, params[0].Encode())
	}

	return r
}

type Client struct {
	httpClient HttpDoer
}

type ClientInterface interface {
	Do(request *http.Request, responseObj any, errorResponseObj any) (int, error)
}

func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

func NewClientWithHttpClient(httpClient HttpDoer) *Client {
	return &Client{
		httpClient,
	}
}

func (c *Client) Do(request *http.Request, responseObj any, errorResponseObj any) (int, error) {
	// ensure that responseObj is a non-nil pointer
	if responseObj == nil || reflect.TypeOf(responseObj).Kind() != reflect.Ptr {
		return -1, fmt.Errorf("responseObj must be a pointer")
	}

	// ensure that errorResponseObj is a non-nil pointer
	if errorResponseObj == nil || reflect.TypeOf(errorResponseObj).Kind() != reflect.Ptr {
		return -1, fmt.Errorf("errorResponseObj must be a pointer")
	}

	if request == nil {
		return -1, fmt.Errorf("request must not be nil")
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return -1, fmt.Errorf("doing request: %w", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		if err := json.NewDecoder(resp.Body).Decode(errorResponseObj); err != nil {
			return -1, fmt.Errorf("decoding error response: %w", err)
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(responseObj); err != nil {
			return -1, fmt.Errorf("decoding response: %w", err)
		}
	}

	return resp.StatusCode, nil
}
