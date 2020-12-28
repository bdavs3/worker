package auth

import "sync"

type ownershipTracker struct {
	ownerships map[string]map[string]struct{}
	mu         sync.Mutex
}

func newOwnershipTracker() *ownershipTracker {
	return &ownershipTracker{
		ownerships: make(map[string]map[string]struct{}),
	}
}

func (ot *ownershipTracker) setOwner() {

}

func (ot *ownershipTracker) isOwner() bool {

}
