package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pak-app/gosuper/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestServiceStopController_Success(t *testing.T) {
	// 1. Setup mock
	mockSup := new(MockSupervisor)
	mockSup.On("StopAllServices").Return()

	// 2. Prepare global daemonServer
	daemonServer = &DaemonServer{
		Supervisors: map[string]core.SupervisorInterface{
			"test-supervisor": mockSup,
		},
	}

	// 3. Create request (POST, with query param)
	req := httptest.NewRequest(http.MethodPost, "/service/stop?supervisor_name=test-supervisor", nil)
	w := httptest.NewRecorder()

	// 4. Call handler
	serviceStopController(w, req)

	// 5. Check HTTP response
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := w.Body.String()
	assert.Contains(t, body, "Services stopped")

	// 6. Verify mock call
	mockSup.AssertCalled(t, "StopAllServices")
}

func TestServiceStopController_MissingSupervisor(t *testing.T) {

	daemonServer = &DaemonServer{
		Supervisors: map[string]core.SupervisorInterface{},
	}

	req := httptest.NewRequest(http.MethodPost, "/service/stop?supervisor_name=missing", nil)
	w := httptest.NewRecorder()
	serviceStopController(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "supervisor with name missing doesn't exist")
}

func TestServiceStopController_MissingSupervisorName(t *testing.T) {

	daemonServer = &DaemonServer{
		Supervisors: map[string]core.SupervisorInterface{},
	}

	req := httptest.NewRequest(http.MethodPost, "/service/stop", nil)
	w := httptest.NewRecorder()
	serviceStopController(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "supervisor name doesn't define in query body")
}

func TestServiceStartController_Success(t *testing.T) {
	mockSup := new(MockSupervisor)
	mockSup.On("LoadServices").Return()
	mockSup.On("RunServices").Return(nil)

	// Override factory
    // Mocking core.NewSupervisor via package scope variable
	oldFactory := newSupervisor
	newSupervisor = func() core.SupervisorInterface { return mockSup }
	defer func() { newSupervisor = oldFactory }()

	daemonServer = &DaemonServer{
		Supervisors: map[string]core.SupervisorInterface{
			"existing-sup": mockSup,
		},
	}

	// JSON body with a config specifying the same supervisor name
	body := `{
        "supervisor": {"name": "existing-sup"},
        "services": {}
    }`
	req := httptest.NewRequest(http.MethodPost, "/service/start", strings.NewReader(body))
	w := httptest.NewRecorder()

	serviceStartController(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "services run successfully")
	mockSup.AssertCalled(t, "RunServices")
}

func TestServicesStartController_MissingSupName(t *testing.T) {

	// JSON body with a config specifying the same supervisor name
	body := `{
        "supervisor": {},
        "services": {}
    }`
	req := httptest.NewRequest(http.MethodPost, "/service/start", strings.NewReader(body))
	w := httptest.NewRecorder()

	serviceStartController(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "supervisor name doesn't set in config file")
}
