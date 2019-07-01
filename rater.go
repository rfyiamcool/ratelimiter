package rateLimiter

import (
	"context"
	"sync/atomic"
	"time"
)

type RateLimiter interface {
	Start()

	TryAcquire() bool

	Acquire() bool

	Stop()
}

type StableRateLimiter struct {
	threshold        int64
	currentThreshold int64
	refillPeriod     time.Duration
	broadcastChannel chan bool
	quitChannel      chan bool
}

func NewRateLimiter(threshold int64, refillPeriod time.Duration) (rateLimiter *StableRateLimiter) {
	rateLimiter = &StableRateLimiter{
		threshold:        threshold,
		currentThreshold: threshold,
		refillPeriod:     refillPeriod,
		broadcastChannel: make(chan bool),
	}
	return rateLimiter
}

// Start to refill the bucket periodically.
func (limiter *StableRateLimiter) Start() {
	limiter.quitChannel = make(chan bool)
	quitChannel := limiter.quitChannel
	go func() {
		for {
			select {
			case <-quitChannel:
				return
			default:
				atomic.StoreInt64(&limiter.currentThreshold, limiter.threshold)
				time.Sleep(limiter.refillPeriod)
				close(limiter.broadcastChannel)
				limiter.broadcastChannel = make(chan bool)
			}
		}
	}()
}

func (limiter *StableRateLimiter) AcquireContext(ctx context.Context) (allow bool) {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
		}

		permit := atomic.AddInt64(&limiter.currentThreshold, -1)
		if permit < 0 {
			allow = false
			// block until the bucket is refilled
			select {
			case <-ctx.Done():
				return false
			case <-limiter.broadcastChannel:
				continue
			}
		}

		return true
	}
}

func (limiter *StableRateLimiter) Acquire() (allow bool) {
	for {
		permit := atomic.AddInt64(&limiter.currentThreshold, -1)
		if permit < 0 {
			allow = false
			select {
			case <-limiter.broadcastChannel:
				continue
			}
		}

		return true
	}
}

func (limiter *StableRateLimiter) TryAcquire() (allow bool) {
	permit := atomic.AddInt64(&limiter.currentThreshold, -1)
	if permit < 0 {
		allow = false
	} else {
		allow = true
	}
	return allow
}

// Stop the rate limiter.
func (limiter *StableRateLimiter) Stop() {
	close(limiter.quitChannel)
}
