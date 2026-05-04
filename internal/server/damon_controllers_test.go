package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"
	"time"

	// "github.com/pak-app/gosuper/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDaemonStatusController_Success(t *testing.T) {

	currentTimeMilli := time.Now().UnixMilli()

	mockDmn := new(MockDaemonServer)
	mockDmn.On("Status").Return(DaemonServerStatus{
		SupervisorsCounts: 10,
		StartedAt:         currentTimeMilli,
		State:             Alive,
	})

	daemonServer = mockDmn

	req := httptest.NewRequest(http.MethodGet, "/daemon/status", nil)
	w := httptest.NewRecorder()

	daemonStatusController(w, req)

	// 5. Check HTTP response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := w.Body.Bytes()

	var response DaemonServerStatus

	if err := json.Unmarshal(body, &response); err != nil {
		t.Log("Error occured:", err)
		return
	}

	assert.Equal(t, currentTimeMilli, response.StartedAt)
	assert.Equal(t, Alive, response.State)
	assert.Equal(t, 10, response.SupervisorsCounts)

	mockDmn.AssertCalled(t, "Status")
}
