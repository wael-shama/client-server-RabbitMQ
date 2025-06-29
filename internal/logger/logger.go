package logger

import (
	"os"
	"sync"
)

type SafeLogger struct {
	filename string
	mu       sync.Mutex
}

func NewSafeLogger(filename string) *SafeLogger {
	return &SafeLogger{filename: filename}
}

func (l *SafeLogger) WriteString(s string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(s + "\n")
	return err
}
