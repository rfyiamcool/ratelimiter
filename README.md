# rateLimiter

a simple and not lock golang rate limiter

## Usage:

no block mode

```
rlimiter := rateLimiter.NewRateLimiter(5, 1 * time.Second)
for {
	 allow := rlimiter.TryAcquire()
	 if allow {
		 task.Fn()
	 }
}
```

block mode with context

```
limiter := rateLimiter.NewRateLimiter(10, 1 * time.Second)
for {
	 allow := limiter.AcquireContext(ctx)
	 if allow{
		 task.Fn()
	 }
}
```

block mode

```
limiter := rateLimiter.NewRateLimiter(10, 1 * time.Second)
for {
	 allow := limiter.Acquire()
	 if allow{
		 task.Fn()
	 }
}
```