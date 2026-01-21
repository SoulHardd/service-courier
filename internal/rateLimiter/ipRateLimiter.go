package rateLimiter

import "sync"

type IPRateLimiter struct {
	limiters map[string]*RateLimiter
	mu       sync.Mutex
}

func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

func (i *IPRateLimiter) Get(ip string, maxTokens float64, refillRate float64) *RateLimiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, ok := i.limiters[ip]
	if !ok {
		limiter = NewRateLimiter(maxTokens, refillRate)
		i.limiters[ip] = limiter
	}

	return limiter
}
