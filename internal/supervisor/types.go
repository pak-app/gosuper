package supervisor

import (
	"time"
	"os/exec"
	"github.com/pak-app/gosuper/internal/config"
)

// ProcessState defines the current health of the service
type ProcessState string

const (
	StateStarting ProcessState = "STARTING"
	StateRunning  ProcessState = "RUNNING"
	StateStopped  ProcessState = "STOPPED"
	StateFatal    ProcessState = "FATAL" 
)

// Process holds the dynamic STATE of the service
type Process struct {
	Name         string       `json:"name"`
	PID          int          `json:"pid"`
	Status       ProcessState `json:"status"`
	StartTime    time.Time    `json:"start_time"`
	RestartCount int          `json:"restart_count"` // Current number of crashes, e.g., $C = 3$
	
    // We keep a reference to the config here so the Manager can read its rules,
    // but we use `json:"-"` so it does NOT get saved into the state.json file.
	Config       *config.ServiceConfig `json:"-"` 
	cmd          *exec.Cmd             `json:"-"`
}
