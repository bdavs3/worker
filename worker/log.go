package worker

import (
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
}

// NewLog returns a log containing an initialized job entry map.
func NewLog() *Log {
	return &Log{
		entries: make(map[string]*logEntry),
	}
}

func (log *Log) addEntry(id string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id] = &logEntry{status: "active", output: ""}
}

func (log *Log) setStatus(id, status string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id].status = status
}

func (log *Log) getStatus(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	if entry, ok := log.entries[id]; ok {
		return entry.status, nil
	}

	return "", errors.New("job not found")
}

func (log *Log) setOutput(id, output string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id].output = output
}

func (log *Log) getOutput(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	if entry, ok := log.entries[id]; ok {
		return entry.output, nil
	}

	return "", errors.New("job not found")
}
