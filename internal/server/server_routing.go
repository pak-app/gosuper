package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pak-app/gosuper/internal/config"
	"github.com/pak-app/gosuper/internal/core"
)

func gracefulShutdown() {

	// Give it a moment to ensure the response is sent.
	// A small delay can help, though flushing is the primary mechanism.
	time.Sleep(100 * time.Millisecond)

	// Create a context for the shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}

// /daemon/stop route
func daemonStopController(w http.ResponseWriter, r *http.Request) {

	if err := os.Remove(SocketPath); err != nil && !os.IsNotExist(err) {
		// If it's an error other than "file doesn't exist", something is wrong.
		log.Printf("Warning: could not remove existing socket: %v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "daemon shutdowned"}`))
	w.WriteHeader(http.StatusOK)

	go gracefulShutdown()
}

func daemonStatusController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "alive", "up_time": 1888, "start_date": "2025-12-12 12:00:00"}`))
}

// /service/status route
func serviceStatusController(w http.ResponseWriter, r *http.Request) {

	supervisorName := r.URL.Query().Get("supervisor_name")
	
	var supervisors map[string]core.SupervisorStatus

	// Send back all supervisors status
	if supervisorName == "" {
		supervisors = make(map[string]core.SupervisorStatus, len(daemonServer.Supervisors))
		supervisors = daemonServer.GetAllStatus()
	} else { // return target supervisor
		supervisors = make(map[string]core.SupervisorStatus, 1)
		sup := daemonServer.Supervisors[supervisorName]
		supervisors[supervisorName] = sup.Status()
	}

    jsonBytes, err := json.Marshal(supervisors)
    if err != nil {
        log.Println("Error:", err)
        return
    }

    jsonString := string(jsonBytes)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}

// /service/start route
func serviceStartController(w http.ResponseWriter, r *http.Request) {
	// 1. Create an empty struct to hold the incoming data
	var appConfig config.Config

	message := "services run successfully"

	// 2. Decode the JSON body from the request into the struct
	err := json.NewDecoder(r.Body).Decode(&appConfig)
	if err != nil {
		// If the JSON is malformed or doesn't match the struct, return a 400 Bad Request
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if appConfig.Supervisor.Name == "" {
		message = "supervisor name doesn't set in config file"
	} else {
		setupSupervisor(&appConfig)
	}

	// 3. Send a success response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, message)))

}

// /stop route
func serviceStopController(w http.ResponseWriter, r *http.Request) {

	responseBody := `{"message": "%s"}`
	message := "Services stopped"

	supervisorName := r.URL.Query().Get("supervisor_name")

	if supervisorName == "" {
		message = "supervisor name doesn't define in query body"
	} else if _, ok := daemonServer.Supervisors[supervisorName]; !ok {
		message = fmt.Sprintf("supervisor with name %s doesn't exist", supervisorName)
	} else {
		supervisor := daemonServer.Supervisors[supervisorName]
		supervisor.StopAllServices()
		delete(daemonServer.Supervisors, supervisorName)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(responseBody, message)))
}
