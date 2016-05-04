package witty

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	baseURL    = "https://api.wit.ai"
	apiVersion = "20141022"
)

// Client manages communication with the wit.ai API.
type Client struct {
	HTTPClient *http.Client

	BaseURL     string
	AccessToken string
	APIVersion  string

	*chatService
}

// NewClient returns a new wit.ai API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		HTTPClient:  httpClient,
		BaseURL:     baseURL,
		AccessToken: token,
		APIVersion:  apiVersion,
	}
	c.chatService = &chatService{client: c}
	return c
}

// NewRequest creates a request to the wit.ai API.
func (c *Client) NewRequest(method, path string, params *url.Values) (*http.Request, error) {
	params.Set("v", c.APIVersion)
	urlStr := fmt.Sprintf("%s/%s?%v", c.BaseURL, path, params.Encode())
	log.Print(urlStr)

	req, err := http.NewRequest(method, urlStr, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	return req, err
}

// Do uses the wit.ai API client's HTTP client to execute the request
// and unmarshals the response into v.
// It also handles unmarshaling errors returned by the API.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}
