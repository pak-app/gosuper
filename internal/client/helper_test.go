package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newFakeDaemonServer(t *testing.T) (*httptest.Server, *Client) {
    t.Helper()
    handler := http.NewServeMux()
    // Register routes matching your daemon's endpoints
    handler.HandleFunc("POST /service/start", successMessageResponse("services run successfully", 200))
    handler.HandleFunc("POST /service/stop", successMessageResponse("services stopped", 200))
    handler.HandleFunc("GET /service/status", successServiceStatusResponse_1Sup("test"))

    server := httptest.NewServer(handler)
    // Create a client pointing to the test server
    client := newClientForTesting(server.URL, server.Client().Transport)
    return server, client
}
