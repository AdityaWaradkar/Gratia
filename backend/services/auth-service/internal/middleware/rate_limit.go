package middleware

import (
	"net/http"
	"sync"
	"time"
)

// Simple in-memory rate limiter by IP address
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int           // max requests
	burst    int           // burst capacity
	ttl      time.Duration // how long to keep visitor data
}

type visitor struct {
	remaining int
	lastSeen  time.Time
}

func NewRateLimiter(rate, burst int, ttl time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		ttl:      ttl,
	}

	go rl.cleanupVisitors()

	return rl
}

func (rl *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.ttl {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			remaining: rl.burst,
			lastSeen:  time.Now(),
		}
		rl.visitors[ip] = v
	}
	return v
}

func (rl *rateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr // or use X-Forwarded-For header if behind proxy

		v := rl.getVisitor(ip)

		rl.mu.Lock()
		defer rl.mu.Unlock()

		if v.remaining <= 0 {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		v.remaining--
		v.lastSeen = time.Now()

		next.ServeHTTP(w, r)
	})
}
