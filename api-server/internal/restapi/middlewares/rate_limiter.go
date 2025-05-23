package middleware

import (
	"net/http"
	"time"
)

type Limiter struct {
	id       int
	interval time.Duration
	numReq   int
	cr       map[int][]time.Time
}

func NewRateLimiter(id int, interval time.Duration, numReq int) *Limiter {
	return &Limiter{
		id:       id,
		interval: interval,
		numReq:   numReq,
		cr:       make(map[int][]time.Time),
	}
}

func (rl *Limiter) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oldReq, exists := rl.cr[rl.id]
		if !exists {
			
		}
		next.ServeHTTP(w, r)
	})
}
