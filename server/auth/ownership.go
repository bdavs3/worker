package auth

import "sync"

type ownershipTracker struct {
	// The inner map employs the empty struct so it may be treated like a set.
	ownerships map[string]map[string]struct{}
	mu         sync.Mutex
}

func newOwnershipTracker() *ownershipTracker {
	return &ownershipTracker{
		ownerships: make(map[string]map[string]struct{}),
	}
}

func (ot *ownershipTracker) setOwner(username, id string) {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if _, ok := ot.ownerships[username]; !ok {
		ot.ownerships[username] = make(map[string]struct{})
	}
	ot.ownerships[username][id] = struct{}{}
}

func (ot *ownershipTracker) isOwner(username, id string) bool {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	_, ok := ot.ownerships[username][id]

	return ok
}
