package server

import (
	"log"
	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/core"
)

func setupSupervisor(cfg *config.Config) {

	var supervisor *core.Supervisor

	supervisor, ok := daemonServer.Supervisors[cfg.Supervisor.Name]

	if !ok {
		supervisor = core.NewSupervisor()
		supervisor.LoadServices(cfg)
	}

	err := supervisor.RunServices()

	if err != nil {
		log.Println("Error during running services: ", err)
	}

	daemonServer.Supervisors[cfg.Supervisor.Name] = supervisor
}
