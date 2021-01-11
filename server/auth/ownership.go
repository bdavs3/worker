package auth

import (
	"sync"
)

// OwnershipRecorder assists in request authorization by tracking resource ownership.
type OwnershipRecorder interface {
	SetOwner(username, id string)
	IsOwner(username, id string) bool
}

// Owners is the OwnershipRecorder used by the auth layer. Use NewOwners to create a
// new instance.
type Owners struct {
	// The empty struct allows the inner map to be treated like a set.
	ownerships map[string]map[string]struct{}
	mu         sync.RWMutex
}

// NewOwners creates a new instance of the owner log.
func NewOwners() *Owners {
	return &Owners{
		ownerships: make(map[string]map[string]struct{}),
	}
}

// DummyOwners is an OwnershipRecorder intended only for testing dependent functions.
type DummyOwners struct{}

func (do *DummyOwners) SetOwner(username, id string)     {}
func (do *DummyOwners) IsOwner(username, id string) bool { return false }

// SetOwner registers the given user as the owner of the resource with the given id.
func (ot *Owners) SetOwner(username, id string) {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if _, ok := ot.ownerships[username]; !ok {
		ot.ownerships[username] = make(map[string]struct{})
	}
	ot.ownerships[username][id] = struct{}{}
}

// IsOwner returns true only if the given user is the registered owner of the resource
// represented by the given id.
func (ot *Owners) IsOwner(username, id string) bool {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	user, ok := ot.ownerships[username]
	if !ok { // User has not been registered as owner of any resource.
		return false
	}

	_, ok = user[id]

	return ok
}
