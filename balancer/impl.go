package balancer

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

const MaxProviderCount = 10
const RequiredHealthChecksToPass = 2
const MaxRequestsPerProvider = 3

type balancer struct {
	strategy Strategy
	pos       int

	providers   []Provider
	excludedIDs map[string]bool
	unhealthyIDs map[string]int
	inFlightRequests int
	mu sync.Mutex // synchronizes above structures

	hc HealthChecker
}

func New(s Strategy) LoadBalancer {
	b := &balancer{
		strategy:    s,
		pos:         -1,
		excludedIDs: map[string]bool{},
		unhealthyIDs: map[string]int{},
	}
	b.hc = HealthChecker{B: b}
	b.hc.Start()
	return b
}

func (b *balancer) Register(p Provider) error {
	fmt.Printf("Register provider %s \n", p.ID())

	b.mu.Lock()
	if len(b.providers) >= MaxProviderCount {
		return errors.New("provider's list capacity exceeded")
	}
	b.providers = append(b.providers, p)
	b.mu.Unlock()

	checkResult := p.Check()
	b.markProvider(p, checkResult)
	return nil
}

func (b *balancer) activeCount() int {
	return len(b.providers) - len(b.excludedIDs) - len(b.unhealthyIDs)
}

func (b *balancer) selectProvider() Provider {
	active := b.activeCount()
	if active < 1 {
		panic("no active provider")
	}
	skipFirstN := 0
	switch b.strategy {
	case StrategyRandom:
		skipFirstN = rand.Intn(active)
	case StrategyRoundRobin:
		b.pos = (b.pos + 1) % active
		skipFirstN = b.pos
	default:
		panic("unknown strategy")
	}

	for i, p := range b.providers {
		if b.excludedIDs[p.ID()] || b.unhealthyIDs[p.ID()] > 0 {
			continue
		}
		if skipFirstN == 0 {
			return b.providers[i]
		}
		skipFirstN--
	}
	return nil
}

func (b *balancer) Get() (string, error) {
	b.mu.Lock()
	if b.inFlightRequests >= MaxRequestsPerProvider * b.activeCount() {
		b.mu.Unlock()
		return "", errors.New("too many requests in flight")
	}
	b.inFlightRequests++
	p := b.selectProvider()
	b.mu.Unlock()

	defer func() {
		b.mu.Lock()
		b.inFlightRequests--
		b.mu.Unlock()
	}()

	return p.Get(), nil
}

func (b *balancer) Include(id string) {
	fmt.Printf("Include provider %s \n", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.excludedIDs, id)
}

func (b *balancer) Exclude(id string) {
	fmt.Printf("Exclude provider %s \n", id)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.excludedIDs[id] = true
	delete(b.unhealthyIDs, id)
}

func (b *balancer) HealthCheck() {
	for _, p := range b.providers {
		b.mu.Lock()
		excluded := b.excludedIDs[p.ID()]
		b.mu.Unlock()
		if excluded {
			continue
		}
		checkResult := p.Check()
		b.markProvider(p, checkResult)
	}
}

func (b *balancer) markProvider(p Provider, checkResult bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	id := p.ID()
	if b.excludedIDs[id] {
		return
	}
	if checkResult {
		v := b.unhealthyIDs[id]
		if v > 0 {
			v--
			if v > 0 {
				b.unhealthyIDs[id] = v
			} else {
				// healthy again
				delete(b.unhealthyIDs, id)
			}
		}
	} else {
		b.unhealthyIDs[id] = RequiredHealthChecksToPass
	}
}