package client

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
)

// New creates a new client connected to the daemon's Unix socket
func New(socketPath string) (*Client, error) {

	if socketPath == "" {
		log.Printf("Socket path doesn't provided for client side")
		return nil, errors.New("Socket path doesn't provided or exist")
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

