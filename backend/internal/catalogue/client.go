package catalogue

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download page: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	return resp, nil
}
