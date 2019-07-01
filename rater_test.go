package rateLimiter

import (
	"testing"
	"time"
)

func TestStableRateLimiter(t *testing.T) {
	rateLimiter := NewRateLimiter(1, 10*time.Millisecond)
	rateLimiter.Start()
	defer rateLimiter.Stop()

	allow := rateLimiter.TryAcquire()
	if !allow {
		t.Error("Unexpected blocked by rate limiter")
	}
	allow = rateLimiter.TryAcquire()
	if allow {
		t.Error("Should be blocked")
	}
}
