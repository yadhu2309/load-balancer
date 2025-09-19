package main

import (
	// "fmt"
	"fmt"
	"log"
	"time"
)

func ClientLimiter(capacity int, refillRatePerSec int) *TokenBucket {
	tb := &TokenBucket{Capacity: capacity,
		Token: capacity,
		// RefillRatePerSec: func() *time.Duration {
		// 	d := time.Second / time.Duration(refillRatePerSec)
		// 	return &d
		// }(),
		RefillRatePerSec: time.Second / time.Duration(refillRatePerSec),
	}

	go tb.RefillToken()

	return tb

}

func (tb *TokenBucket) RefillToken() {
	ticker := time.NewTicker(tb.RefillRatePerSec)

	for range ticker.C {
		fmt.Print("yadhu")
		tb.mu.Lock()
		if tb.Token < tb.Capacity {
			log.Println("Refilling Token in every", tb.RefillRatePerSec)
			tb.Token++
		}
		tb.mu.Unlock()
	}

}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.Token > 0 {
		tb.Token--
		log.Println("Token left", tb)
		return true
	}
	return false
}

func GetClientBucket(ip string, capacity int, refillRatePerSec int) *TokenBucket {
	bucketmu.Lock()
	defer bucketmu.Unlock()
	if bucket, exists := buckets[ip]; exists {
		log.Println("buckets", bucket)
		return bucket
	}
	bucket := ClientLimiter(capacity, refillRatePerSec)
	buckets[ip] = bucket
	log.Println("buckets", bucket)

	return bucket
}
