package worker

import (
	"sync"
)

// A log contains information about Linux processes. A zero value of this
// type is invalid - Use newLog to create a new instance.
type log struct {
	entries map[string]*logEntry
	mu      sync.RWMutex
}

// A logEntry contains data relevant to a singular Linux process.
type logEntry struct {
	status string
	output string
	killC  chan bool
}

// newLog creates a new instance of the process log.
func newLog() *log {
	return &log{
		entries: make(map[string]*logEntry),
	}
}

func (log *log) addEntry(id string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	log.entries[id] = &logEntry{status: statusActive}
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

func (log *log) appendOutput(id string, output []byte) error {
	log.mu.Lock()
	defer log.mu.Unlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return err
	}

	entry.output += string(output)

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

func (log *log) makeKillC(id string) chan bool {
	log.mu.Lock()
	defer log.mu.Unlock()

	entry, _ := log.getEntryLocked(id)

	entry.killC = make(chan bool)
	return entry.killC
}

func (log *log) nullifyKillC(id string) {
	log.mu.Lock()
	defer log.mu.Unlock()

	entry, _ := log.getEntryLocked(id)

	entry.killC = nil
}

func (log *log) getKillC(id string) (chan bool, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return nil, err
	}
	if entry.killC == nil {
		return nil, &ErrJobNotActive{"job not active"}
	}

	return entry.killC, nil
}
