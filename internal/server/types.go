package server

import (
	"github.com/pak-app/gosuper/internal/core"
	"sync"
)

type DaemonServerInterface interface {
	GetAllStatus() map[string]core.SupervisorStatus
	GetSupervisor(string) (core.SupervisorInterface, bool)
	StoreSupervisor(string, core.SupervisorInterface)
	SupervisorCount() int
	RemoveSupervisor(string)
}

type DaemonServer struct {
	mu          sync.RWMutex
	Supervisors map[string]core.SupervisorInterface
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
