package macaroons

import (
	"context"
	"regexp"
	"sync"

	"gopkg.in/macaroon-bakery.v2/bakery"
)

type MemThirdPartyStore struct {
	reg *regexp.Regexp
	mu  sync.RWMutex
	m   map[string]bakery.ThirdPartyInfo
}

func (l *MemThirdPartyStore) AddInfo(loc string, info bakery.ThirdPartyInfo) {
	l.mu.Lock()
	defer l.mu.Unlock()

	filtered := l.filterLocation(loc)
	log.Debug("Adding key: ", filtered)

	// Check if it already exists, if so, don't overwrite it
	_, ok := l.m[filtered]
	if ok {
		log.Debugf("Already have key for %s", filtered)
		return
	}
	l.m[filtered] = info
}

func (l MemThirdPartyStore) ThirdPartyInfo(ctx context.Context, loc string) (bakery.ThirdPartyInfo, error) {
	log.Debugf("Getting key for location: %s", loc)
	log.Debugf("Filtered loc: %s", l.filterLocation(loc))
	l.mu.RLock()
	defer l.mu.RUnlock()

	info, ok := l.m[l.filterLocation(loc)]
	if !ok {
		return info, bakery.ErrNotFound
	}
	return info, nil
}

func (l MemThirdPartyStore) filterLocation(loc string) string {
	return l.reg.ReplaceAllString(loc, "")
}

func NewMemThirdPartyStore() MemThirdPartyStore {
	return MemThirdPartyStore{
		m:   make(map[string]bakery.ThirdPartyInfo),
		mu:  sync.RWMutex{},
		reg: regexp.MustCompile("([0-9a-zA-z]+-)+[0-9a-zA-Z]+/"),
	}
}
