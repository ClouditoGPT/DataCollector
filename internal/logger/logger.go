package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	mu     sync.Mutex
	logger *log.Logger
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[INFO] "+msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[ERROR] "+msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logger.Printf("[DEBUG] "+msg, args...)
}

var defaultLogger = NewLogger("")

func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}