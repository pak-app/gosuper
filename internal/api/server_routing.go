package api

import (
	"net/http"
	"os"
	"log"
	"time"
	"context"
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

// /status route
func statusController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "running", "services": 2}`))
}

// // /start route
// func startController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"start": "running", "services": 2}`))
// }

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
