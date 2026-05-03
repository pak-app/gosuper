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
	err := client.ServiceStartRequest(cfg)
	assert.NoError(t, err)
}

func TestServiceStopRequest_Success(t *testing.T) {
	server, client := newFakeDaemonServer(t)
	defer server.Close()

	err := client.ServiceStopRequest("test")
	assert.NoError(t, err)
}

func TestServiceStatusRequest_Success(t *testing.T) {
	server, client := newFakeDaemonServer(t)
	defer server.Close()

	status, err := client.ServiceStatusRequest("sup-test")

	assert.NoError(t, err)
	assert.NotEmpty(t, status)
}

