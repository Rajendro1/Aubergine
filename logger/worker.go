package logger

import (
	"encoding/json"
	"fmt"
)

// startWorkers initializes the worker pool
func (l *Logger) startWorkers() {
	l.wg.Add(l.workerCount)
	for i := 0; i < l.workerCount; i++ {
		go l.worker()
	}
}

// worker is launched in a goroutine to read from the log channel
func (l *Logger) worker() {
	defer l.wg.Done()
	for entry := range l.logChan {
		l.processEntry(entry)
	}
}

// processEntry receives one log entry and decides its destination (Kafka or file)
func (l *Logger) processEntry(entry LogEntry) {
	data, err := json.Marshal(entry)
	if err != nil {
		fmt.Printf("failed to marshal log entry: %v\n", err)
		return
	}

	success := false
	if l.kafkaWriter != nil {
		err = l.kafkaWriter.Write(data)
		if err == nil {
			success = true
		} else {
			// If Kafka fails, we fallback entirely to file
			fmt.Printf("Kafka write failed, falling back to file: %v\n", err)
		}
	}

	if !success {
		err = l.fileWriter.Write(data)
		if err != nil {
			fmt.Printf("Fallback file write failed: %v\n", err)
		}
	}
}
