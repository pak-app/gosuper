package server

import (
	// "github.com/pak-app/gosuper/internal/core"
	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/core"
	"github.com/stretchr/testify/mock"
)

// MockSupervisor implements core.SupervisorInterface
type MockSupervisor struct {
	mock.Mock
}

func (m *MockSupervisor) RunServices() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockSupervisor) StopAllServices() {
	m.Called()
}

func (m *MockSupervisor) LoadServices(cfg *config.Config) {
	m.Called()
}

func (m *MockSupervisor) Status() core.SupervisorStatus {
	args := m.Called()
	return args.Get(0).(core.SupervisorStatus)
}

type MockDaemonServer struct {
	mock.Mock
}

func (m *MockDaemonServer) GetAllStatus() map[string]core.SupervisorStatus {
	args := m.Called()
	return args.Get(0).(map[string]core.SupervisorStatus)
}

func (m *MockDaemonServer) GetSupervisor(name string) (core.SupervisorInterface, bool) {
	args := m.Called()
	return args.Get(0).(core.SupervisorInterface), args.Bool(1)
}

func (m *MockDaemonServer) StoreSupervisor(name string, sup core.SupervisorInterface) {
	m.Called()
}

func (m *MockDaemonServer) SupervisorCount() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockDaemonServer) RemoveSupervisor(name string) {
	m.Called()
}
