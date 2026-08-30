package middleware

import (
	"testing"
	"time"
)

// The limiter's opportunistic sweep (triggered past 1024 keys) must evict
// keys whose attempts all fell out of the window, while keeping keys that are
// still live — otherwise a spray of unique IPs grows the map without bound.
func TestRateLimiterSweepEvictsExpiredKeys(t *testing.T) {
	rl := &LoginRateLimiter{
		attempts: make(map[string][]time.Time),
		limit:    5,
		window:   5 * time.Minute,
	}

	// 1100 unique keys, one attempt each → crosses the sweep threshold.
	for i := 0; i < 1100; i++ {
		if !rl.Allow("spray-ip-" + string(rune('a'+i%26)) + string(rune('0'+i%10)) + string(rune(i))) {
			t.Fatalf("first attempt for key %d must pass", i)
		}
	}
	// All attempts are within the window → nothing evicted yet, map stays live.
	if len(rl.attempts) == 0 {
		t.Fatal("expected live entries")
	}

	// Age every entry past the window, then trigger one more Allow → sweep runs.
	expired := time.Now().Add(-2 * rl.window)
	for k := range rl.attempts {
		rl.attempts[k] = []time.Time{expired}
	}
	if !rl.Allow("fresh-ip") {
		t.Fatal("fresh key must pass")
	}

	// Swept: expired keys are gone; only the fresh key remains.
	for k := range rl.attempts {
		if k != "fresh-ip" {
			t.Fatalf("expired key %q survived the sweep", k)
		}
	}
}
