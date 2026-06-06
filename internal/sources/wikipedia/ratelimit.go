package wikipedia

import "time"

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(delay time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(delay),
	}
}

func (r *RateLimiter) Wait() {
	<-r.ticker.C
}
