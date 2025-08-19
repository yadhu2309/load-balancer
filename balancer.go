package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var current int
var mu sync.Mutex

func NextServer(servers []*Host) *Host {
	mu.Lock()
	defer mu.Unlock()
	start := current
	for {
		server := servers[current]
		current = (current + 1) % len(servers)
		server.mu.RLock()
		isHealthy := server.IsHealthy
		server.mu.RUnlock()
		if isHealthy {
			fmt.Println("next server", server.Address)
			return server
		}
		if start == current {
			return nil
		}

	}

}

func HealthCheck(servers []*Host) {
	// for {
	log.Println("=================Health Checking==================")
	for _, server := range servers {
		fmt.Println("address", server.Address)
		serverDial, err := net.DialTimeout("tcp", server.Address, 2*time.Second)
		// log.Println("goroutine handle", serverDial.RemoteAddr())
		server.mu.Lock()
		if err != nil {
			// panic(err)
			log.Println("❌", server.Address, "Unhealthy")
			server.IsHealthy = false
			// return

		} else {
			log.Println("✅", server.Address, "Is Healthy")
			server.IsHealthy = true
			serverDial.Close()
		}
		server.mu.Unlock()
		// defer serverDial.Close()
		// defer log.Println("connection closing after health check")

	}
	// }
}

func HandleConnection(client net.Conn, servers []*Host) {
	log.Println("in handle connection")
	defer client.Close()
	server := NextServer(servers).Address
	if server == "" {
		log.Println("No Healthy server")
		return
	}
	serverDial, err := net.Dial("tcp", server)
	log.Println("goroutine handle", serverDial.RemoteAddr())
	if err != nil {
		// panic(err)
		log.Println("error", err)
		return
	}
	defer serverDial.Close()
	go func() {

		_, err := io.Copy(serverDial, client)
		if err != nil {
			if strings.Contains(err.Error(), "broken pipe") ||
				strings.Contains(err.Error(), "reset by peer") {
				log.Println("⚠️ Connection closed:", err)
				return
			}
			log.Println("❌ Copy error:", err)
		}
	}()
	io.Copy(client, serverDial)
}
