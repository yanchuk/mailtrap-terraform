package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://mailtrap.io"
	sendingAPIURL    = "https://send.api.mailtrap.io"
	sandboxAPIURL    = "https://sandbox.api.mailtrap.io"
)

// Client represents a Mailtrap API client
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// NewClient creates a new Mailtrap API client
func NewClient(apiToken string) *Client {
	return &Client{
		baseURL:  defaultBaseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL sets a custom base URL for the client
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// doRequest performs an HTTP request with proper authentication
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	// Determine the base URL based on the endpoint
	baseURL := c.baseURL
	if endpoint == "/api/send" || endpoint == "/api/batch" {
		baseURL = sendingAPIURL
	} else if endpoint[:10] == "/api/send/" {
		baseURL = sandboxAPIURL
	}

	fullURL, err := url.JoinPath(baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to construct URL: %w", err)
	}

	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Api-Token", c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// handleResponse processes the HTTP response and unmarshals JSON if needed
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp struct {
			Error   string `json:"error"`
			Errors  interface{} `json:"errors"`
			Message string `json:"message"`
		}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if errorResp.Error != "" {
				return fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp.Error)
			}
			if errorResp.Message != "" {
				return fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp.Message)
			}
			if errorResp.Errors != "" {
				return fmt.Errorf("API error (%d): %v", resp.StatusCode, errorResp.Errors)
			}
		}
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	if result != nil && len(body) > 0 {
		if err := json.Unmarshal(body, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Get performs a GET request
func (c *Client) Get(endpoint string, result interface{}) error {
	resp, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	return c.handleResponse(resp, result)
}

// Post performs a POST request
func (c *Client) Post(endpoint string, body, result interface{}) error {
	resp, err := c.doRequest("POST", endpoint, body)
	if err != nil {
		return err
	}
	return c.handleResponse(resp, result)
}

// Patch performs a PATCH request
func (c *Client) Patch(endpoint string, body, result interface{}) error {
	resp, err := c.doRequest("PATCH", endpoint, body)
	if err != nil {
		return err
	}
	return c.handleResponse(resp, result)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string, result interface{}) error {
	resp, err := c.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	return c.handleResponse(resp, result)
}
