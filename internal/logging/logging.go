package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	logFile        *os.File
	originalOutput io.Writer
	mu             sync.Mutex
	initialized    bool
)

func Init(logFilePath string) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil // Already initialized
	}

	originalOutput = log.Writer()

	if logFilePath == "" {
		logFilePath = "log/info.log"
	}

	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logFile = file
	log.SetOutput(file)
	initialized = true

	return nil
}

func Close() {
    mu.Lock()
    defer mu.Unlock()
    
    if !initialized {
        return
    }
    
    // Temporarily restore original output for shutdown message
    log.SetOutput(originalOutput)

    if logFile != nil {
        logFile.Close()
        logFile = nil
    }
    
    log.SetOutput(originalOutput)
    initialized = false
}
