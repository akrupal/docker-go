package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Limiter struct {
	id       int
	interval time.Duration
	numReq   int
	cr       map[int][]time.Time
	sync.Mutex
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
		rl.Lock()
		defer rl.Unlock()
		now := time.Now()
		oldReq, exists := rl.cr[rl.id]
		if !exists {
			oldReq = []time.Time{}
		}
		newReq := []time.Time{}

		for _, e := range oldReq {
			if now.Sub(e) <= rl.interval {
				newReq = append(newReq, e)
			}
		}
		newReq = append(newReq, now)
		if len(newReq) >= rl.numReq {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode("Too many requests, try after some time!")
			return
		}
		rl.cr[rl.id] = newReq
		next.ServeHTTP(w, r)
	})
}
