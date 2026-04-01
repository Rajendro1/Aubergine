package logger

import (
	"fmt"
	"sync"
	"time"
)

// LogEntry represents a single log message structure
type LogEntry struct {
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Timestamp string            `json:"timestamp"`
	Fields    map[string]string `json:"fields,omitempty"`
}

// Logger handles asynchronous logging using a worker pool
type Logger struct {
	logChan     chan LogEntry
	workerCount int
	wg          sync.WaitGroup
	kafkaWriter *KafkaWriter
	fileWriter  *FileWriter
	stopChan    chan struct{}
}

// Config provides configuration for the Logger
type Config struct {
	WorkerCount int
	Brokers     []string
	Topic       string
	LogDir      string
	LogChanSize int
}

// New creates and starts a new Logger instance
func New(cfg Config) (*Logger, error) {
	if cfg.WorkerCount <= 0 {
		cfg.WorkerCount = 5
	}
	if cfg.LogChanSize <= 0 {
		cfg.LogChanSize = 1000
	}
	if cfg.LogDir == "" {
		cfg.LogDir = "logs"
	}

	fw, err := NewFileWriter(cfg.LogDir)
	if err != nil {
		return nil, fmt.Errorf("failed to init file writer: %w", err)
	}

	var kw *KafkaWriter
	if len(cfg.Brokers) > 0 && cfg.Topic != "" {
		kw = NewKafkaWriter(cfg.Brokers, cfg.Topic)
	}

	l := &Logger{
		logChan:     make(chan LogEntry, cfg.LogChanSize),
		workerCount: cfg.WorkerCount,
		kafkaWriter: kw,
		fileWriter:  fw,
		stopChan:    make(chan struct{}),
	}

	l.startWorkers()
	return l, nil
}

// Close gracefully waits for pending logs to be written and shuts down the logger
func (l *Logger) Close() {
	close(l.stopChan)
	close(l.logChan)
	l.wg.Wait()
	if l.kafkaWriter != nil {
		l.kafkaWriter.Close()
	}
}

func (l *Logger) push(level, msg string, fields map[string]string) {
	entry := LogEntry{
		Level:     level,
		Message:   msg,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Fields:    fields,
	}

	select {
	case l.logChan <- entry:
	default:
		fmt.Println("Warning: log channel is full, dropping log entry")
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string, fields map[string]string) {
	l.push("INFO", msg, fields)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields map[string]string) {
	l.push("ERROR", msg, fields)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields map[string]string) {
	l.push("DEBUG", msg, fields)
}
