package server

import (
	"sync"
	"time"

	"github.com/pak-app/gosuper/internal/core"
)

type DaemonServerInterface interface {
	GetAllStatus() map[string]core.SupervisorStatus
	GetSupervisor(string) (core.SupervisorInterface, bool)
	StoreSupervisor(string, core.SupervisorInterface)
	SupervisorCount() int
	RemoveSupervisor(string)
	Status() DaemonServerStatus
	setState(DaemonServerState)
}

type DaemonServerState string

const (
	Stopped  DaemonServerState = "Stopped"
	Alive    DaemonServerState = "Alive"
	Starting DaemonServerState = "Starting"
	Stopping DaemonServerState = "Stopping"
)

type DaemonServer struct {
	mu          sync.RWMutex
	Supervisors map[string]core.SupervisorInterface
	StartedAt   time.Time
	State       DaemonServerState
}

type DaemonServerStatus struct {
	SupervisorsCounts int               `json:"supervisor_counts"`
	StartedAt         int64             `json:"started_at"`
	State             DaemonServerState `json:"state"`
}

func (ds *DaemonServer) GetAllStatus() map[string]core.SupervisorStatus {

	ds.mu.RLock()
	supervisors := make([]core.SupervisorInterface, len(ds.Supervisors))
	for _, sup := range ds.Supervisors {
		supervisors = append(supervisors, sup)
	}
	ds.mu.RUnlock()

	num := len(supervisors)

	results := make(chan core.SupervisorStatus, num)

	var wg sync.WaitGroup

	wg.Add(num)

	for _, sup := range supervisors {
		go func(s core.SupervisorInterface) {
			defer wg.Done()
			results <- s.Status()
		}(sup)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	supMap := make(map[string]core.SupervisorStatus, num)

	for status := range results {
		supMap[status.Name] = status
	}

	return supMap
}

func (ds *DaemonServer) GetSupervisor(name string) (core.SupervisorInterface, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	sup, ok := ds.Supervisors[name]
	return sup, ok
}

func (ds *DaemonServer) StoreSupervisor(name string, sup core.SupervisorInterface) {
	ds.mu.Lock()
	ds.Supervisors[name] = sup
	ds.mu.Unlock()
}

func (ds *DaemonServer) SupervisorCount() int {
	return len(ds.Supervisors)
}

func (ds *DaemonServer) RemoveSupervisor(name string) {
	ds.mu.Lock()
	delete(ds.Supervisors, name)
}

func (ds *DaemonServer) Status() DaemonServerStatus {
	return DaemonServerStatus{
		SupervisorsCounts: len(ds.Supervisors),
		StartedAt:         ds.StartedAt.UnixMilli(),
		State:             ds.State,
	}
}

func (ds *DaemonServer) setState(state DaemonServerState) {
	ds.State = state
}
