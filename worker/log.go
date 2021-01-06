package worker

import (
	"bytes"
	"sync"
)

type syncBuffer struct {
	mu sync.RWMutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.b.String()
}

// A log contains information about Linux processes being executed by the worker
// library. Use newLog to create a new instance.
type log struct {
	entries map[string]*logEntry
	mu      sync.RWMutex
}

// A logEntry contains data relevant to a single Linux process.
type logEntry struct {
	status       string
	outputBuffer *syncBuffer
	killC        chan bool
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

	log.entries[id] = &logEntry{status: statusActive, outputBuffer: &syncBuffer{}}
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

func (log *log) getOutputBuffer(id string) (*syncBuffer, error) {
	log.mu.RLock()
	defer log.mu.RUnlock()

	entry, err := log.getEntryLocked(id)
	if err != nil {
		return nil, err
	}

	return entry.outputBuffer, nil
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
