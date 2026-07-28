package client

import (
	"github.com/pak-app/gosuper/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceStartRequest_Success(t *testing.T) {
	server, client := newFakeDaemonServer(t)
	defer server.Close()

	cfg := &config.Config{Supervisor: config.SupervisorConfig{Name: "test"}}
	_, err := client.ServiceStartRequest(cfg)
	assert.NoError(t, err)
}

func TestServiceStopRequest_Success(t *testing.T) {
	server, client := newFakeDaemonServer(t)
	defer server.Close()

	_, err := client.ServiceStopRequest("test")
	assert.NoError(t, err)
}

func TestServiceStatusRequest_Success(t *testing.T) {
	server, client := newFakeDaemonServer(t)
	defer server.Close()

	status, err := client.ServiceStatusRequest("sup-test")

	assert.NoError(t, err)
	assert.NotEmpty(t, status)
}

