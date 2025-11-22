// ref: https://github.com/hashicorp-demoapp/hashicups-client-go/blob/main/client.go
package apiclient

import (
	"fmt"
	"net/http"
	"time"
)

const testapiUrl string = "https://testapi.io/"

type Client struct {
	URL        string
	HTTPClient *http.Client
	Token      string
}

func NewClient(uname string, token string) (*Client, error) {
	endpoint := fmt.Sprintf("%s/api/%s/resource/", testapiUrl, uname)
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
