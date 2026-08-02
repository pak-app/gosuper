package server

import (
	"log"
	"net"
	"net/http"
	"os"
	"github.com/pak-app/gosuper/internal/logging"
	"github.com/pak-app/gosuper/internal/core"
)

var SocketPath string
var Server *http.Server
var daemonServer DaemonServerInterface
var serverLogger *log.Logger
var serverLogFile *os.File

func LoadSupervisors() {
	daemonServer =  &DaemonServer{
		Supervisors: make(map[string]core.SupervisorInterface),
		State: Stopped,
	}
}

func StartServer(socketPath string) {
	if socketPath == "" {
		socketPath = "tmp/gosuper.sock"
	}


	var err error
	serverLogger, serverLogFile, err = logging.NewFileLogger("log/daemon_server.log", "[SERVER] ", log.Ldate|log.Ltime|log.Lshortfile)
	if err != nil {
		panic(err)
	}

	serverLogger.Println("Daemon is running...")

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

	daemonServer.setState(Alive)
	daemonServer.setStartedAt()
	// 4. Serve HTTP over the Unix listener
	Server = &http.Server{Handler: mux}
	if err := Server.Serve(listener); err != nil {
		log.Panic(err)
	}

}
