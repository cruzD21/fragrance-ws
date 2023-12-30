package main

import (
	"math/rand"
	"sync"
	"time"
)

type RateLimiter struct {
	delay       time.Duration
	randomDelay time.Duration
	counter     int
	mutex       sync.Mutex
}

var maxRequest = 9

func (rl *RateLimiter) Increment() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.counter++
}

func (rl *RateLimiter) ShouldDelay() time.Duration {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Logic to increase delay based on the counter
	if rl.counter >= maxRequest {
		rl.delay = 5 * time.Second        // New base delay
		rl.randomDelay = 10 * time.Second // New random delay
		// Optionally reset the counter if you want the delay to be temporary
		rl.counter = 0
	}

	// Return the delay with some randomization
	return rl.delay + time.Duration(rand.Int63n(int64(rl.randomDelay)))

}
