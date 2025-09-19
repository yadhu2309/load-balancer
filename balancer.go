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
		if isHealthy != nil && *isHealthy {
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
		// fmt.Println("address", server.Address)
		serverDial, err := net.DialTimeout("tcp", server.Address, 2*time.Second)
		// log.Println("goroutine handle", serverDial.RemoteAddr())
		server.mu.Lock()
		if err != nil {
			// panic(err)
			log.Println("❌", server.Address, "Unhealthy")
			*server.IsHealthy = false
			// return

		} else {
			log.Println("✅", server.Address, "Is Healthy")
			*server.IsHealthy = true
			serverDial.Close()
		}
		server.mu.Unlock()
		// defer serverDial.Close()
		// defer log.Println("connection closing after health check")

	}
	// }
}

func HandleConnection(client net.Conn, servers []*Host) {
	log.Println("===================Handling Requests=================")
	defer func() {

		log.Println("✅ connection closed cleanly")
		client.Close()
	}()
	server := NextServer(servers)
	if server == nil {
		log.Println("No Healthy server")
		return
	}
	serverDial, err := net.Dial("tcp", server.Address)
	if serverDial == nil {
		fmt.Println("hell", err)
	}
	log.Println("goroutine handle", serverDial.RemoteAddr())
	if err != nil {
		// panic(err)
		log.Println("Error in Connection", err)
		return
	}
	defer serverDial.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	// client -> server
	go pipe(client, serverDial, "client->server", &wg)

	// server -> client
	go pipe(serverDial, client, "server->client", &wg)

	wg.Wait()

}

func MatchEndpoint(pools []*Pool, requestLine string) []*Host {
	var servers []*Host
	for _, pool := range pools {
		fmt.Println("poo;", pool.EndpointPrefix)
		for _, endpoint := range *pool.EndpointPrefix {
			if strings.HasPrefix(requestLine, "GET "+endpoint) ||
				strings.HasPrefix(requestLine, "POST "+endpoint) {
				servers = pool.Hosts
				break
			}
		}

	}
	return servers
}

func pipe(src, dst net.Conn, direction string, wg *sync.WaitGroup) {
	defer wg.Done()

	n, err := io.Copy(dst, src)
	log.Printf("%s copy done, bytes: %d, err: %v\n", direction, n, err)

	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") ||
			strings.Contains(err.Error(), "reset by peer") {
			log.Printf("⚠️ Connection closed (%s): %v\n", direction, err)
			return
		}
		log.Printf("❌ Copy error (%s): %v\n", direction, err)
	}

	// Half-close write to unblock the other direction
	if tcpConn, ok := dst.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
}
