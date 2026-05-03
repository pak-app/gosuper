package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

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
	StartedAt      time.Time
}

// later CPU, and Memory usage
type ServiceStatus struct {
	Name      string       `json:"name"`
	State     ServiceState `json:"state"`
	PID       int          `json:"pid"`
	StartedAt time.Time    `jsonn:"started_at"`
	UptimeMS  int64        `json:"up_time_ms"`
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
	service.StartedAt = time.Now()
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

	service.mu.Lock()
	defer service.mu.Unlock()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			// non‑zero exit code: exitErr.ExitCode()
			// treat as crash
			log.Println("Exit error happened: ", exitErr.ExitCode())
		} else {
			// some other error (e.g., I/O error waiting) — rarely happens
			log.Println("Error happened")
		}
	} else {
		// clean exit
		service.CurrentState = Stopped
		service.PID = 0
		log.Printf("Service %s ended and stopped\n", service.Name)
		return
	}

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

func (service *Service) status() ServiceStatus {
	service.mu.RLock()
	defer service.mu.RUnlock()

	var delta int64

	if service.StartedAt.IsZero() {
		delta = 0
	} else {
		delta = time.Since(service.StartedAt).Milliseconds()
	}

	return ServiceStatus{
		Name:      service.Name,
		State:     service.CurrentState,
		PID:       service.PID,
		StartedAt: service.StartedAt,
		UptimeMS:  delta,
	}
}
