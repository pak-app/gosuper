package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/types"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Client represents the gosuper API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

type Daemon struct {
	UpTime    int    `json:"up_time"`
	StartDate string `json:"start_date"`
	Status    string `json:"status"`
}

// New creates a new client connected to the daemon's Unix socket
func New(socketPath string) (*Client, error) {

	if socketPath == "" {
		return nil, fmt.Errorf("Socket path doesn't provided or exist")
	}

	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				// Override the dialer to use the Unix domain socket
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.Dial("unix", socketPath)
				},
			},
		},
		// The hostname ("localhost") is ignored by the Unix socket dialer,
		// but http.NewRequest requires a valid URL format.
		baseURL: "http://localhost",
	}, nil
}

func (c *Client) SartDaemonRequest() error {

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/daemon/start", nil)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start daemon")
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) StopDaemonRequest() error {

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/daemon/stop", nil)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to stop daemon")
	}
	defer res.Body.Close()

	return nil
}

func (c *Client) DaemonStatusRequest() (*Daemon, error) {

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/daemon/status", nil)

	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("daemon is dead")
	}
	defer res.Body.Close()

	var result Daemon

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("daemon is dead")
	}

	log.Printf("Daemon status %v and it started from %v.\nUp Time: %v", result.Status, result.StartDate, result.UpTime)

	return &result, nil
}

func (c *Client) ServiceStartRequest(serviceConfig *config.Config) error {

	jsonData, err := json.Marshal(serviceConfig)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/service/start", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)

	var responseData types.SimpleResponse

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start service/services: %v", res.StatusCode)
	}

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("service started successfully: %s", responseData.Message)
	defer res.Body.Close()

	return nil
}

func (c *Client) ServiceStopRequest(name string) error {
    url := fmt.Sprintf("%s/service/stop?group_name=%s", c.baseURL, url.QueryEscape(name))
    
    req, err := http.NewRequest(http.MethodPost, url, nil)
    if err != nil {
        return err
    }
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}