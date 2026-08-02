package catalogue

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func (c *Client) Download(url, filename string) error {
	resp, err := c.http.Get(url)
	if err != nil {
		return fmt.Errorf("download page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return fmt.Errorf("save page: %w", err)
	}

	return nil
}
