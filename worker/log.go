package worker

import (
	"context"
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

// NewLog returns a log containing an initialized entry map.
func NewLog() *Log {
	return &Log{
		entries: make(map[string]*logEntry),
	}
}

// NotFoundErr occurs when a job cannot be found in the worker log.
type NotFoundErr struct{ msg string }

func (e *NotFoundErr) Error() string { return e.msg }

// NotActiveErr occurs when termination is attempted on a job that
// doesn't exist.
type NotActiveErr struct{ msg string }

func (e *NotActiveErr) Error() string { return e.msg }

func (log *Log) addEntry(id string, cancel context.CancelFunc) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id] = &logEntry{status: statusActive, output: "", cancel: cancel}
}

func (log *Log) setStatus(id, status string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id].status = status
}

func (log *Log) getStatus(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return "", err
	}

	return entry.status, nil
}

func (log *Log) appendOutput(id string, output []byte) error {
	log.mu.Lock()
	defer log.mu.Unlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return err
	}

	bytes := []byte(entry.output)
	entry.output = string(append(bytes, output...))

	return nil
}

func (log *Log) getOutput(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return "", err
	}

	return entry.output, nil
}

func (log *Log) getCancelFunc(id string) (func(), error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return func() {}, err
	}
	if entry.status != "active" {
		return func() {}, &NotActiveErr{"Job not active."}
	}

	return entry.cancel, nil
}

func (log *Log) getEntryLocked(id string) (*logEntry, error) {
	entry, ok := log.entries[id]
	if !ok {
		return &logEntry{}, &NotFoundErr{"Job not found."}
	}

	return entry, nil
}
