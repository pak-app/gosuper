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
	multiWriter    io.Writer
)

func Init(logFilePath string) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
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

	// Create a MultiWriter to write to both the terminal (os.Stdout or os.Stderr) and the file
	multiWriter = io.MultiWriter(os.Stdout, file)

	// Set the standard logger to use the MultiWriter
	log.SetOutput(multiWriter)

	initialized = true

	return nil
}

func Writer() io.Writer {
	mu.Lock()
	defer mu.Unlock()

	if !initialized {
		return os.Stdout
	}
	return multiWriter
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

// NewFileLogger creates a brand new *log.Logger writing to the specified file and stdout.
// It also returns the *os.File handle so it can be closed by the caller, and any initialization error.
func NewFileLogger(logFilePath string, prefix string, flag int) (*log.Logger, *os.File, error) {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return nil, nil, err
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	writer := io.MultiWriter(os.Stdout, file)
	logger := log.New(writer, prefix, flag)
	return logger, file, nil
}
