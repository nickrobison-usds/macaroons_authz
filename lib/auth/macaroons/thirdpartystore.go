package macaroons

import (
	"context"
	"sync"

	"gopkg.in/macaroon-bakery.v2/bakery"
)

type MemThirdPartyStore struct {
	mu sync.RWMutex
	m  map[string]bakery.ThirdPartyInfo
}

func (l *MemThirdPartyStore) AddInfo(loc string, info bakery.ThirdPartyInfo) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Check if it already exists, if so, don't overwrite it
	_, ok := l.m[loc]
	if ok {
		log.Debugf("Already have key for %s", loc)
		return
	}
	l.m[loc] = info
}

func (l MemThirdPartyStore) ThirdPartyInfo(ctx context.Context, loc string) (bakery.ThirdPartyInfo, error) {
	log.Debugf("Getting key for location: %s", loc)
	l.mu.RLock()
	defer l.mu.RUnlock()

	info, ok := l.m[loc]
	if !ok {
		return info, bakery.ErrNotFound
	}
	return info, nil
}

func NewMemThirdPartyStore() MemThirdPartyStore {
	return MemThirdPartyStore{
		m: make(map[string]bakery.ThirdPartyInfo),
	}
}
