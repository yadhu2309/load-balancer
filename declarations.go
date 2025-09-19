package main

import (
	"sync"
	"time"
)

var (
	buckets  = make(map[string]*TokenBucket)
	bucketmu sync.Mutex
)

type Host struct {
	Address   string `json:"address"`
	IsHealthy *bool  `json:"is_healthy"`
	mu        sync.RWMutex
}

type Pool struct {
	Hosts          []*Host   `json:"hosts"`
	EndpointPrefix *[]string `json:"endpoint_prefix"`
	Name           string    `json:"name"`
	Listen         string    `json:"listen"`
}

// type Config struct {
// 	TcpConfig  *PoolConfig `json:"tcp"`
// 	HttpConfig *PoolConfig `json:"http"`
// }

var Config = make(map[string]*PoolConfig)

type PoolConfig struct {
	Pool                *[]*Pool   `json:"pools"`
	HealthCheckInterval int        `json:"healthcheck_interval"`
	IsMultiplePool      bool       `json:"is_multiple_pool"`
	RateLimit           *RateLimit `json:"rate_limit"`
}

type RateLimit struct {
	Capacity         int `json:"capacity"`
	RefillRatePerSec int `json:"refill_rate_per_sec"`
}

type Semaphore struct {
	SemaChannel chan struct{}
}

type TokenBucket struct {
	Capacity         int
	RefillRatePerSec time.Duration
	Token            int
	mu               sync.Mutex
}
