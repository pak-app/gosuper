package server

import (
	"fmt"
	"log"

	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/core"
)

var newSupervisor = func () core.SupervisorInterface {
	return core.NewSupervisor()
}

func setupSupervisor(cfg *config.Config) {

	var supervisor core.SupervisorInterface

	supervisor, ok := daemonServer.GetSupervisor(cfg.Supervisor.Name)


	if !ok {
		supervisor = newSupervisor()
		supervisor.LoadServices(cfg)
	}

	err := supervisor.RunServices()

	if err != nil {
		log.Println("Error during running services: ", err)
	}

	daemonServer.StoreSupervisor(cfg.Supervisor.Name, supervisor)
}

func simpleMessageResponse(message string) string {
	return fmt.Sprintf(`{"message":"%s"}`, message)
}
