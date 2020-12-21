package worker

import (
	"context"
	"errors"
	"sync"
)

// Log contains a map that pairs jobs with their status and output.
type Log struct {
	entries map[string]*logEntry
	mu      sync.RWMutex
}

// logEntry is a struct containing the status and output of a process.
type logEntry struct {
	status string
	output string
	cancel context.CancelFunc
}

// NewLog returns a log containing an initialized job entry map.
func NewLog() *Log {
	return &Log{
		entries: make(map[string]*logEntry),
	}
}

func (log *Log) addEntry(id string, cancel context.CancelFunc) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id] = &logEntry{status: "active", output: "", cancel: cancel}
}

func (log *Log) setStatus(id, status string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id].status = status
}

func (log *Log) getStatus(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntry(id)
	if err != nil {
		return "", err
	}

	return entry.status, nil
}

func (log *Log) setOutput(id, output string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id].output = output
}

func (log *Log) getOutput(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntry(id)
	if err != nil {
		return "", err
	}

	return entry.output, nil
}

func (log *Log) getCancelFunc(id string) (func(), error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntry(id)
	if err != nil {
		return func() {}, err
	}
	if entry.status != "active" {
		return func() {}, errors.New("job not active")
	}

	return entry.cancel, nil
}

func (log *Log) getEntry(id string) (*logEntry, error) {
	entry, ok := log.entries[id]
	if !ok {
		return &logEntry{}, errors.New("job not found")
	}

	return entry, nil
}
