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

	// supervisor mock functions
	mockSup := new(MockSupervisor)
	mockSup.On("StopAllServices").Return()

	// daemon server mock functions
	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetSupervisor", "test-supervisor").Return(mockSup, true)
	mockDmn.On("RemoveSupervisor", "test-supervisor").Return()

	// 2. Prepare global daemonServer
	daemonServer = mockDmn

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

	mockSup := new(MockSupervisor)

	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetSupervisor", "missing").Return(mockSup, false)

	daemonServer = mockDmn

	req := httptest.NewRequest(http.MethodPost, "/service/stop?supervisor_name=missing", nil)
	w := httptest.NewRecorder()
	serviceStopController(w, req)

	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "supervisor with name missing doesn't exist")
}

func TestServiceStopController_MissingSupervisorName(t *testing.T) {

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

	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetSupervisor", "existing-sup").Return(mockSup, true)
	mockDmn.On("StoreSupervisor", "existing-sup", mockSup).Return()
	
	// Override factory
	// Mocking core.NewSupervisor via package scope variable
	oldFactory := newSupervisor
	newSupervisor = func() core.SupervisorInterface { return mockSup }
	defer func() { newSupervisor = oldFactory }()


	daemonServer = mockDmn

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
	mockDmn.AssertCalled(t, "GetSupervisor", "existing-sup")
	mockDmn.AssertCalled(t, "StoreSupervisor", "existing-sup", mockSup)
}

func TestServicesStartController_MissingSupName(t *testing.T) {

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

func TestServiceStatusContorller_Success(t *testing.T) {

	mockSup := new(MockSupervisor)
	mockSup.On("Status").Return(core.SupervisorStatus{
		Name:           "sup_1",
		ServicesStatus: make(map[string]core.ServiceStatus),
	})

	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetSupervisor", "sup_1").Return(mockSup, true)

	daemonServer = mockDmn

	req := httptest.NewRequest(http.MethodPost, "/service/status?supervisor_name=sup_1", nil)
	w := httptest.NewRecorder()

	serviceStatusController(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "sup_1")
	mockSup.AssertCalled(t, "Status")
}

func TestServiceStatusContorller_MissingSupName(t *testing.T) {

	mockSup := new(MockSupervisor)

	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetSupervisor", "sup_2").Return(mockSup, false)

	daemonServer = mockDmn

	req := httptest.NewRequest(http.MethodPost, "/service/status?supervisor_name=sup_2", nil)
	w := httptest.NewRecorder()

	serviceStatusController(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "supervisor doesn't exist")
}

func TestServiceStatusController_AllSupStatus(t *testing.T) {

	fakeDaemonStatus := map[string]core.SupervisorStatus{
		"sup_1": {Name: "sup_1"},
		"sup_2": {Name: "sup_2"},
	}

	mockDmn := new(MockDaemonServer)
	mockDmn.On("GetAllStatus").Return(fakeDaemonStatus)
	mockDmn.On("SupervisorCount").Return(2)

	daemonServer = mockDmn

	req := httptest.NewRequest(http.MethodPost, "/service/status", nil)
	w := httptest.NewRecorder()

	serviceStatusController(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "sup_1")
	assert.Contains(t, w.Body.String(), "sup_2")
	mockDmn.AssertCalled(t, "GetAllStatus")
}
