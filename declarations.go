package main

import "sync"

type Host struct {
	Address   string `json:"address"`
	IsHealthy bool   `json:"is_healthy"`
	mu        sync.RWMutex
}

type Pool struct {
	Hosts []*Host `json:"hosts"`
}

type Config struct {
	Pool                []*Pool `json:"pools"`
	HealthCheckInterval int     `json:"healthcheck_interval"`
	IsMultiplePool bool `json:"is_multiple_pool"`
}
