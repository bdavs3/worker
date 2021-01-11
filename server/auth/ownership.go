package auth

import (
	"sync"
)

type ownershipTracker struct {
	// The empty struct allows the inner map to be treated like a set.
	ownerships map[string]map[string]struct{}
	mu         sync.RWMutex
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
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	user, ok := ot.ownerships[username]
	if !ok { // User has not been registered as owner of any resource.
		return false
	}

	_, ok = user[id]

	return ok
}
