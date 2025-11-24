// ref: https://github.com/hashicorp-demoapp/hashicups-client-go/blob/main/client.go
package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const testapiUrl string = "https://testapi.io/"

type Client struct {
	URL        string
	HTTPClient *http.Client
	Token      string
}

type MeApiResponse struct {
	Name        string `json:"name"`
	Year        int64  `json:"year"`
	HomepageUrl string `json:"homepage"`
	ApiPath     string `json:"path"`
}

func NewClient(uname string, token string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/api/%s", testapiUrl, uname)
	c := Client{
		URL:        endpoint,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
	// If username not provided, return empty client
	if uname == "" {
		return &c, nil
	}

	if token != "" {
		c.Token = token
	}

	fmt.Printf(">>> Configured endpoint: %s\n", endpoint)
	return &c, nil
}

func (c *Client) GetResponse() (*MeApiResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/me", c.URL), nil)
	fmt.Printf("Trying to call API %s\n", fmt.Sprintf("%s/me", c.URL))

	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate HTTP client.\n")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke API requests.\n")
	}
	resp, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	r := &MeApiResponse{}
	if json.Unmarshal(resp, r) != nil {
		return nil, fmt.Errorf("Failed to unmarshal response to JSON.\n")
	}

	return r, nil
}
