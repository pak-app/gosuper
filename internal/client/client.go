package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/core"
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

func newClientForTesting(baseURL string, transport http.RoundTripper) *Client {
	return &Client{
		httpClient: &http.Client{Transport: transport},
		baseURL:    baseURL,
	}
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

	if err != nil {
		return fmt.Errorf("request failed: %w", err) // response is nil here, don't touch it
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to start service/services: %v", res.StatusCode)
	}

	var responseData types.SimpleResponse

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	log.Printf("service started successfully: %s", responseData.Message)
	defer res.Body.Close()

	return nil
}

func (c *Client) ServiceStopRequest(supervisorName string) error {
    // Parse the base URL so we can safely modify it
    base, err := url.Parse(c.baseURL)
    if err != nil {
        return fmt.Errorf("invalid base URL: %w", err)
    }

    // Set the specific path
    base.Path = "/service/stop"

    // Add query parameters
    q := base.Query()
    q.Set("supervisor_name", supervisorName)
    base.RawQuery = q.Encode()

    req, err := http.NewRequest(http.MethodPost, base.String(), nil)
    if err != nil {
        return err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    var responseData types.SimpleResponse
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    log.Println("Stop request response: ", responseData.Message)
    return nil
}

func (c *Client) ServiceStatusRequest(supervisorName string) (map[string]*core.SupervisorStatus, error) {
    // Parse the client's base URL
    base, err := url.Parse(c.baseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid base URL: %w", err)
    }

    // Set the endpoint path
    base.Path = "/service/status"

    // Add query parameter if supervisor name is provided
    if supervisorName != "" {
        q := base.Query()
        q.Set("supervisor_name", supervisorName)
        base.RawQuery = q.Encode()
    }

    req, err := http.NewRequest(http.MethodGet, base.String(), nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    responseData := make(map[string]*core.SupervisorStatus)
    if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return responseData, nil
}
