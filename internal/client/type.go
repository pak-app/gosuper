package client

import (
	"encoding/json"
	"log"
	"net/http"
	"errors"
)

type Daemon struct {
	UpTime 		int		`json:"up_time"`
	StartDate 	string	`json:"start_date"`
	Status 		string	`json:"status"`
}

// Client represents the gosuper API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

func (c *Client) SartDaemonRequest() error {

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/daemon/start", nil)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)
	
	if res.StatusCode != http.StatusOK {
		return errors.New("failed to start daemon")
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
		return errors.New("failed to stop daemon")
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
		return nil, errors.New("daemon is dead")
	}
	defer res.Body.Close()

	var result Daemon

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Println("failed to get status of daemon")
		return nil, errors.New("daemon is dead")
	}

	log.Printf("Daemon status %v and it started from %v.\nUp Time: %v", result.Status, result.StartDate, result.UpTime)

	return &result, nil
}
