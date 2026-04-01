package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// FileWriter defines the structure that handles synchronous disk writes for logs
type FileWriter struct {
	dir string
	mu  sync.Mutex
}

// NewFileWriter makes sure the directory exists or fails
func NewFileWriter(dir string) (*FileWriter, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed creating log directory: %w", err)
	}
	return &FileWriter{
		dir: dir,
	}, nil
}

// Write handles creating appending logs to daily files using mutexes to prevent goroutine conflicts
func (fw *FileWriter) Write(msg []byte) error {
	fw.mu.Lock()
	defer fw.mu.Unlock()

	fileName := filepath.Join(fw.dir, time.Now().UTC().Format("2006-01-02")+".log")
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open log file %s: %w", fileName, err)
	}
	defer f.Close()

	// Ensure the message has a newline
	if len(msg) == 0 || msg[len(msg)-1] != '\n' {
		msg = append(msg, '\n')
	}

	_, err = f.Write(msg)
	return err
}
