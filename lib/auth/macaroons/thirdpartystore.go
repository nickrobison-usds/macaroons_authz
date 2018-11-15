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
	l.m[loc] = info
}

func (l MemThirdPartyStore) ThirdPartyInfo(ctx context.Context, loc string) (bakery.ThirdPartyInfo, error) {
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
