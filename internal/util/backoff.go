package util

import (
	"math/rand"
	"time"
)

func RetryDelay(retries int) time.Duration {
	base := 50 * time.Millisecond
	max := 2 * time.Second

	delay := base * time.Duration(1<<retries)

	if delay > max {
		delay = max
	}

	jitter := time.Duration(rand.Intn(100)) * time.Millisecond

	return delay + jitter
}
