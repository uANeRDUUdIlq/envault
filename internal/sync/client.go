package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Payload represents an encrypted env file payload exchanged with the backend.
type Payload struct {
	Name      string `json:"name"`
	Ciphertext string `json:"ciphertext"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Client handles HTTP communication with the shared envault backend.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new sync Client targeting baseURL and authenticating with apiKey.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) authHeader(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")
}

// Push uploads an encrypted payload to the backend.
func (c *Client) Push(p Payload) error {
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("sync: marshal payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPut, c.baseURL+"/vaults/"+p.Name, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("sync: build request: %w", err)
	}
	c.authHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sync: push request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("sync: push failed with status %d", resp.StatusCode)
	}
	return nil
}

// Pull fetches an encrypted payload from the backend by name.
func (c *Client) Pull(name string) (Payload, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/vaults/"+name, nil)
	if err != nil {
		return Payload{}, fmt.Errorf("sync: build request: %w", err)
	}
	c.authHeader(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Payload{}, fmt.Errorf("sync: pull request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return Payload{}, fmt.Errorf("sync: vault %q not found on remote", name)
	}
	if resp.StatusCode != http.StatusOK {
		return Payload{}, fmt.Errorf("sync: pull failed with status %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Payload{}, fmt.Errorf("sync: read response: %w", err)
	}
	var p Payload
	if err := json.Unmarshal(data, &p); err != nil {
		return Payload{}, fmt.Errorf("sync: unmarshal payload: %w", err)
	}
	return p, nil
}
