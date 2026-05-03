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

func (m *MockSupervisor) Status() core.SupervisorStatus{
	args := m.Called()
	return args.Get(0).(core.SupervisorStatus)
}
