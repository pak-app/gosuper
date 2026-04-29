package core

import (
	"sync"
	"github.com/pak-app/gosuper/internal/config"
)

type Supervisor struct {
	Services         map[string]*Service
	SupervisorConfig *config.SupervisorConfig
	mu               sync.RWMutex
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		Services: make(map[string]*Service),
	}
}

func (sup *Supervisor) LoadServices(cfg *config.Config) {

	sup.SupervisorConfig = &cfg.Supervisor

	// Add programs
	for name, serviceCfg := range cfg.Services {
		sup.addService(name, &serviceCfg)
	}
}

func (sup *Supervisor) RunServices() {

	for _, service := range sup.Services {
		service.start()
	}
}

func (sup *Supervisor) addService(name string, serviceCfg *config.ServiceConfig) {

	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.Services[name] = &Service{
		OriginalConfig: serviceCfg,
		CurrentState:   Stopped,
		Name: name,
		stopSignal: make(chan struct{}, 1),
	}
}
