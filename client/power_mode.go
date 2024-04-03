package main

import "sync"

type PowerMode struct {
	factor int
	mu     sync.RWMutex
}

func NewPowerMode(factor int) *PowerMode {
	return &PowerMode{factor: factor}
}

func (p *PowerMode) Change(factor int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.factor = factor
}

func (p *PowerMode) GetConsumption() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.factor
}
