package main

import "sync"

type Host struct {
	Address   string `json:"address"`
	IsHealthy bool   `json:"is_healthy"`
	mu        sync.RWMutex
}

type Config struct {
	Hosts []*Host `json:"hosts"`
	HealthCheckInterval int `json:"healthcheck_interval"`
}
