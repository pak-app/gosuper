package api

import (
	"context"
	"encoding/json"
	"github.com/pak-app/gosuper/internal/config"
	"log"
	"net/http"
	"os"
	"time"
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
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "running", "services": 2}`))
}

// /service/start route
func serviceStartController(w http.ResponseWriter, r *http.Request) {
	// 1. Create an empty struct to hold the incoming data
	var appConfig config.Config

	// 2. Decode the JSON body from the request into the struct
	err := json.NewDecoder(r.Body).Decode(&appConfig)
	if err != nil {
		// If the JSON is malformed or doesn't match the struct, return a 400 Bad Request
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Now you can use appConfig!
	log.Println("Starting service:", appConfig)

	// 3. Send a success response back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "services run successfully"}`))

}

// // /stop route
// func stopController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"stop": "running", "services": 2}`))
// }

// // /log route
// func logController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"log": "running", "services": 2}`))
// }
