package utilities

import "sync"

type ABool struct {
	value bool
	mu    *sync.RWMutex
}

func (b *ABool) Set(val bool) {
	b.mu.Lock()
	b.value = val
	b.mu.Unlock()
}

func (b *ABool) Get() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.value
}

func NewABool() *ABool {
	return &ABool{
		mu: &sync.RWMutex{},
	}
}
