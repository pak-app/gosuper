package server

import (
	"github.com/pak-app/gosuper/internal/core"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

type DaemonServer struct {
	mu          sync.RWMutex
	Supervisors map[string]*core.Supervisor
}

var SocketPath string
var Server *http.Server
var daemonServer *DaemonServer

func LoadSupervisors() {
	daemonServer =  &DaemonServer{
		Supervisors: make(map[string]*core.Supervisor),
	}
}

func StartServer(socketPath string) {
	if socketPath == "" {
		socketPath = "tmp/gosuper.sock"
	}

	LoadSupervisors()

	SocketPath = socketPath

	// 1. Remove the socket file if it already exists
	if err := os.RemoveAll(socketPath); err != nil {
		panic(err)
	}

	// 2. Create the Unix socket listener
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// 3. Define your HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("POST /daemon/stop", daemonStopController)
	mux.HandleFunc("GET /daemon/status", daemonStatusController)
	mux.HandleFunc("POST /service/start", serviceStartController)
	mux.HandleFunc("GET /service/status", serviceStatusController)
	mux.HandleFunc("POST /service/stop", serviceStopController)
	// mux.HandleFunc("/log", logController)

	// 4. Serve HTTP over the Unix listener
	Server = &http.Server{Handler: mux}
	if err := Server.Serve(listener); err != nil {
		log.Panic(err)
	}
}
