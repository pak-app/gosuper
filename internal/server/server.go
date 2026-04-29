package server

import (
	"log"
	"net"
	"net/http"
	"os"
)

var SocketPath string
var Server *http.Server

func StartServer(socketPath string) {
	if socketPath == "" {
		socketPath = "tmp/gosuper.sock"
	}

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

	mux.HandleFunc("/daemon/stop", daemonStopController)
	mux.HandleFunc("/daemon/status", daemonStatusController)
	mux.HandleFunc("/service/start", serviceStartController)
	mux.HandleFunc("/service/status", serviceStatusController)
	// mux.HandleFunc("/services/stop", )
	// mux.HandleFunc("/log", logController)

	// 4. Serve HTTP over the Unix listener
	Server = &http.Server{Handler: mux}
	if err := Server.Serve(listener); err != nil {
		log.Panic(err)
	}
}
