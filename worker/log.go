package worker

import (
	"sync"
)

// log contains a map that pairs jobs with their status and output.
type log struct {
	entries map[string]*logEntry
	mu      sync.RWMutex
}

// logEntry is a struct containing the status and output of a process.
type logEntry struct {
	status string
	output string
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

	log.entries[id] = &logEntry{status: statusActive, output: ""}
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

func (log *log) appendOutput(id string, output []byte) error {
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

func (log *log) getOutput(id string) (string, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return "", err
	}

	return entry.output, nil
}

func (log *log) getEntryLocked(id string) (*logEntry, error) {
	entry, ok := log.entries[id]
	if !ok {
		return nil, &ErrJobNotFound{"job not found"}
	}

	return entry, nil
}
