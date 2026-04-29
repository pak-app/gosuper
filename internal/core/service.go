package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/pak-app/gosuper/internal/config"
)

type ServiceState string

const (
	Running  ServiceState = "Running"
	Stopped  ServiceState = "Stopped"
	Starting ServiceState = "Starting"
	Failed   ServiceState = "Failed"
	Fatal    ServiceState = "Fatal"
	Stopping ServiceState = "Stopping"
)

type Service struct {
	OriginalConfig *config.ServiceConfig
	Command        *exec.Cmd
	CurrentState   ServiceState
	PID            int
	RestartCount   int
	Name           string
	mu             sync.RWMutex
	stopSignal     chan struct{}
}

func (service *Service) start() error {

	if len(service.OriginalConfig.Command) == 0 {
		return fmt.Errorf(`Command is empty for %v service`, service.Name)
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
		return fmt.Errorf("failed to start %v: %w", service.Name, err)
	}

	service.mu.Lock()
	service.CurrentState = Running
	service.PID = cmd.Process.Pid
	service.mu.Unlock()

	go service.monitor(cmd) // monitor service process

	return nil
}

func (service *Service) stop() {
	
	service.mu.Lock()
	defer service.mu.Unlock()
	
	if service.CurrentState != Running || service.Command == nil || service.Command.Process == nil {
		return
	}

	service.CurrentState = Stopping

	// Signal the monitor goroutine that this is intentional
	select {
	case service.stopSignal <- struct{}{}:
	default:
	}

	// Kill the process
	err := service.Command.Process.Kill()
	if err != nil {
		log.Printf("failed to kill service %s: %w\n", service.Name, err)
	}

}

func (service *Service) monitor(cmd *exec.Cmd) {

	err := cmd.Wait()

	if err != nil {
		log.Panicln("Error in waiting goroutine: ", err)
	}

	service.mu.Lock()
	defer service.mu.Unlock()

	select {
	case <-service.stopSignal:
		// Intentionally stopped
		service.CurrentState = Stopped
		service.PID = 0
		log.Printf("Service %s stopped intentionally.\n", service.Name)
	default:
		// Crashed or exited naturally
		service.CurrentState = Failed
		service.PID = 0
		log.Printf("Service %s exited unexpectedly (err: %v)\n", service.Name, err)

		// TODO: Add auto-restart logic here
	}
}
