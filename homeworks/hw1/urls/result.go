package urls

import "sync"

type defaultAggregationResult struct {
	mu     sync.RWMutex
	counts map[int]int
	done   bool
}

func NewDefaultAggregationResult() *defaultAggregationResult {
	return &defaultAggregationResult{
		counts: make(map[int]int),
	}
}

func (r *defaultAggregationResult) Done() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.done
}

func (r *defaultAggregationResult) GetResponsesCount(code int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counts[code]
}

func (r *defaultAggregationResult) GetResult() map[int]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[int]int, len(r.counts))
	for k, v := range r.counts {
		out[k] = v
	}
	return out
}

func (r *defaultAggregationResult) add(code int) {
	r.mu.Lock()
	r.counts[code]++
	r.mu.Unlock()
}

func (r *defaultAggregationResult) markDone() {
	r.mu.Lock()
	r.done = true
	r.mu.Unlock()
}
