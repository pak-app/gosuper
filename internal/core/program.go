package core

import (
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
)

type Service struct {
	OriginalConfig *config.ServiceConfig
	Command        *exec.Cmd
	CurrentState   ServiceState
	PID            int
	RestartCount   int
	mu             sync.RWMutex
}
