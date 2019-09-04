package redir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client is the redirector client.
type Client struct {
	URL   string
	Token string
}

// SetRequest is the JSON struct for a request to POST /set.
type SetRequest struct {
	// URL is the URL that the user should be redirected to.
	URL string `json:"url"`
}

// NewClient returns a new client.
func NewClient(url, token string) *Client {
	return &Client{
		URL:   url,
		Token: token,
	}
}

// Set sets the redirector using the given URL, token, and request.
func (c *Client) Set(reqBody *SetRequest) error {
	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
