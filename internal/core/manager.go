package core

import (
	"fmt"
	"os"
	"os/exec"
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

func (sup *Supervisor) LoadConfig(cfg *config.Config) {

	sup.SupervisorConfig = &cfg.Supervisor

	// Add programs
	for name, serviceCfg := range cfg.Services {
		sup.addService(name, &serviceCfg)
	}
}

func (sup *Supervisor) RunServices() {

	for name, service := range sup.Services {
		sup.startProgram(name, service)
	}

}

func (sup *Supervisor) addService(name string, serviceCfg *config.ServiceConfig) {

	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.Services[name] = &Service{
		OriginalConfig: serviceCfg,
		CurrentState:   Stopped,
	}
}

func (sup *Supervisor) startProgram(name string, service *Service) error {

	if len(service.OriginalConfig.Command) == 0 {
		return fmt.Errorf(`Command is empty for %v service`, name)
	}

	cmd := exec.Command(service.OriginalConfig.Command[0], service.OriginalConfig.Command[1:]...)

	if service.OriginalConfig.Dir != "" {
		cmd.Dir = service.OriginalConfig.Dir
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	service.mu.Lock()
	service.Command = cmd
	service.CurrentState = Starting
	service.mu.Unlock()

	err := cmd.Start()

	if err != nil {
		service.mu.Lock()
		service.CurrentState = Failed
		service.mu.Unlock()
		return fmt.Errorf("failed to start %s: %w", name, err)
	}

	service.mu.Lock()
	service.CurrentState = Running
	service.PID = cmd.Process.Pid
	service.mu.Unlock()

	return nil
}
