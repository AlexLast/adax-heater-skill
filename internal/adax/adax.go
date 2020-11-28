package adax

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

// Config defines all possible configuration options
type Config struct {
	ClientID     string `split_words:"true"`
	ClientSecret string `split_words:"true"`
	BaseURL      string `split_words:"true" default:"https://api-1.adax.no/client-api"`
}

// Client defines any required
// clients and configuration
type Client struct {
	HTTP   *http.Client
	Config *Config
}

// Response defines the structure of responses
// from the Adax API
type Response struct {
	Homes   []Home   `json:"homes"`
	Rooms   []Room   `json:"rooms"`
	Devices []Device `json:"devices"`
}

// Home defines the strucure of the home field
// returned by Adax
type Home struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Room defines the structure of the room field
// returned by Adax
type Room struct {
	ID             int64  `json:"id"`
	HomeID         int64  `json:"homeId"`
	Name           string `json:"name"`
	HeatingEnabled bool   `json:"heatingEnabled"`
	Temperature    int64  `json:"temperature"`
}

// Update defines the structure of a request
// to update a home, room or device target temperature
type Update struct {
	ID                int64  `json:"id"`
	HeatingEnabled    bool   `json:"heatingEnabled"`
	TargetTemperature string `json:"targetTemperature"` // For some reason this is a string
}

// Updates is the wrapper type for Update
type Updates struct {
	Homes   []Update `json:"homes,omitempty"`
	Rooms   []Update `json:"rooms,omitempty"`
	Devices []Update `json:"devices,omitempty"`
}

// Device defines the structure of the device field
// returned by Adax
type Device struct {
	ID     int64  `json:"id"`
	HomeID int64  `json:"homeId"`
	RoomID int64  `json:"roomId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

const (
	tokenPath         = "%s/auth/token"
	contentPath       = "%s/rest/v1/content"
	controlPath       = "%s/rest/v1/control"
	errGetStatus      = "Unable to get status, error: %v"
	errSetTemperature = "Unable to set temperature, error: %v"
)

// GetAccessToken returns an OAuth access token used for
// communication with the Adax API. This isn't required for
// normal lambda flow as the access token is passed in with the lambda
// context. This should only be used for testing.
func (c *Client) GetAccessToken() (*oauth2.Token, error) {
	// Ensure required config is set
	if c.Config.ClientID == "" || c.Config.ClientSecret == "" {
		return nil, errors.New("ClientID and ClientSecret are required to generate an access token")
	}

	// Build new OAuth client
	client := &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf(tokenPath, c.Config.BaseURL),
		},
	}

	// Return a new token
	return client.PasswordCredentialsToken(
		context.WithValue(
			context.Background(),
			oauth2.HTTPClient,
			c.HTTP,
		),
		c.Config.ClientID,
		c.Config.ClientSecret,
	)
}

// AdaxRequest performs any required Adax requests and returns the Response type
func (c *Client) AdaxRequest(token string, method string, url string, b io.Reader) (*Response, error) {
	// Build a new HTTP request
	request, err := http.NewRequest(method, url, b)

	if err != nil {
		return nil, err
	}

	// Add required headers
	for header, value := range map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json"} {
		request.Header.Add(header, value)
	}

	// Perform the HTTP request
	response, err := c.HTTP.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Ensure we get a 200 response
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Adax API returned an HTTP %d response", response.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	// Unmarshal into Response type
	r := new(Response)
	err = json.Unmarshal(body, r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

// GetStatus returns the status of all homes, rooms and devices
// for the supplied Adax access token
func (c *Client) GetStatus(token string) (*Response, error) {
	// Build the URL
	url := fmt.Sprintf(contentPath, c.Config.BaseURL)

	// Perform the Adax request
	response, err := c.AdaxRequest(token, "GET", url, bytes.NewBuffer(nil))

	if err != nil {
		return nil, fmt.Errorf(errGetStatus, err)
	}

	return response, nil
}

// SetTemperature sets the temperature for a specified room
func (c *Client) SetTemperature(token string, updates *Updates) error {
	// Build the URL
	url := fmt.Sprintf(controlPath, c.Config.BaseURL)

	// Marshal into a JSON payload
	payload, err := json.Marshal(updates)

	if err != nil {
		return fmt.Errorf(errSetTemperature, err)
	}

	// Perform the Adax request
	_, err = c.AdaxRequest(token, "POST", url, bytes.NewBuffer(payload))

	if err != nil {
		return fmt.Errorf(errSetTemperature, err)
	}

	return nil
}
