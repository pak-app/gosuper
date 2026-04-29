package core

import (
	"github.com/pak-app/gosuper/internal/config"
	"sync"
)

type Supervisor struct {
	Services         map[string]*Service
	SupervisorConfig *config.SupervisorConfig
	mu               sync.RWMutex
	Name             string
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		Services: make(map[string]*Service),
	}
}

func (sup *Supervisor) LoadServices(cfg *config.Config) {

	sup.SupervisorConfig = &cfg.Supervisor
	sup.Name = cfg.Supervisor.Name

	// Add programs
	for name, serviceCfg := range cfg.Services {
		sup.addService(name, &serviceCfg)
	}
}

func (sup *Supervisor) RunServices() error {

	for _, service := range sup.Services {
		err := service.start()

		if err != nil {
			return err
		}
	}
	return nil
}

func (sup *Supervisor) addService(name string, serviceCfg *config.ServiceConfig) {

	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.Services[name] = &Service{
		OriginalConfig: serviceCfg,
		CurrentState:   Stopped,
		Name:           name,
		stopSignal:     make(chan struct{}, 1),
	}
}

// it waits to all goroutines end
func (sup *Supervisor) StopAllServices() {

	sup.mu.RLock()
	defer sup.mu.RUnlock()

	var wg sync.WaitGroup // ADDED

	for _, service := range sup.Services {
		wg.Add(1)             // ADDED
		go func(s *Service) { // CHANGED (capture service as parameter)
			defer wg.Done() // ADDED
			s.stop()        // CHANGED (use parameter)
		}(service) // CHANGED (pass service)
	}

	wg.Wait()
}
