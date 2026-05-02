package server

import (
	"net/http"
	"net/http/httptest"
	// "strings"
	"testing"

	// "github.com/pak-app/gosuper/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDaemonStatusController_Success(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/daemon/stop", nil)
	w := httptest.NewRecorder()

	daemonStatusController(w, req)

	// 5. Check HTTP response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := w.Body.String()
	assert.Contains(t, body, "status")
	assert.Contains(t, body, "up_time")
	assert.Contains(t, body, "start_date")

}
