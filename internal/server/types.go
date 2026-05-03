package server

import (
	"github.com/pak-app/gosuper/internal/core"
	"sync"
)

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
