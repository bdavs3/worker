package worker

import (
	"bytes"
	"sync"
)

// log maps Linux processes to their status and output.
type log struct {
	entries map[string]*logEntry
	mu      sync.RWMutex
}

// logEntry represent's a Linux process's status and output.
type logEntry struct {
	status string
	output *bytes.Buffer
}

// newLog returns a log containing an initialized entry map.
func newLog() *log {
	return &log{
		entries: make(map[string]*logEntry),
	}
}

func (log *log) addEntry(id string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id] = &logEntry{status: statusActive, output: bytes.NewBuffer([]byte(nil))}
}

func (log *log) getEntryLocked(id string) (*logEntry, error) {
	entry, ok := log.entries[id]
	if !ok {
		return nil, &ErrJobNotFound{"job not found"}
	}

	return entry, nil
}

func (log *log) setStatus(id, status string) error {
	log.mu.Lock()
	defer log.mu.Unlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return err
	}

	entry.status = status
	return nil
}

func (log *log) getStatus(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return "", err
	}

	return entry.status, nil
}

func (log *log) getOutputBuffer(id string) (*bytes.Buffer, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return nil, err
	}

	return entry.output, nil
}
