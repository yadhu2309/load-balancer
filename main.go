package main

import (
	"fmt"
	// "io"
	"log"
	"net"
	"time"
	// "strings"
)

//	type Host struct {
//		Address   string
//		IsHealthy bool
//		mu        sync.RWMutex
//	}
// var servers = ConfigLoader()

// var servers = []*Host{
// 	{Address: "localhost:8001", IsHealthy: true},
// 	{Address: "localhost:8002", IsHealthy: true},
// }

// var current int
// var mu sync.Mutex

func HandleRecover() {
	if r := recover(); r != nil {
		log.Println("error", r)
	}
}

// func NextServer(servers []*Host) *Host {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	start := current
// 	for {
// 		server := servers[current]
// 		current = (current + 1) % len(servers)
// 		server.mu.RLock()
// 		isHealthy := server.IsHealthy
// 		server.mu.RUnlock()
// 		if isHealthy {
// 			return server
// 		}
// 		if start == current {
// 			return nil
// 		}

// 	}

// }

// func HandleConnection(client net.Conn) {
// 	log.Println("in handle connection")
// 	defer client.Close()
// 	server := NextServer().Address
// 	if server == "" {
// 		log.Println("No Healthy server")
// 		return
// 	}
// 	serverDial, err := net.Dial("tcp", server)
// 	log.Println("goroutine handle", serverDial.RemoteAddr())
// 	if err != nil {
// 		// panic(err)
// 		log.Println("error", err)
// 		return
// 	}
// 	defer serverDial.Close()
// 	go func() {

// 		_, err := io.Copy(serverDial, client)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "broken pipe") ||
// 				strings.Contains(err.Error(), "reset by peer") {
// 				log.Println("⚠️ Connection closed:", err)
// 				return
// 			}
// 			log.Println("❌ Copy error:", err)
// 		}
// 	}()
// 	io.Copy(client, serverDial)
// }

// func HealthCheck() {
// 	// for {
// 	log.Println("=================Health Checking==================")
// 	for _, server := range servers {
// 		fmt.Println("address", server.Address)
// 		serverDial, err := net.DialTimeout("tcp", server.Address, 2*time.Second)
// 		// log.Println("goroutine handle", serverDial.RemoteAddr())
// 		server.mu.Lock()
// 		if err != nil {
// 			// panic(err)
// 			log.Println("❌", server.Address, "Unhealthy")
// 			server.IsHealthy = false
// 			// return

// 		} else {
// 			log.Println("✅", server.Address, "Is Healthy")
// 			server.IsHealthy = true
// 			serverDial.Close()
// 		}
// 		server.mu.Unlock()
// 		// defer serverDial.Close()
// 		// defer log.Println("connection closing after health check")

// 	}
// 	// }
// }

func main() {
	fmt.Println("Load balancer!!!")
	// file, err := AppLog()
	// if err == nil {
	// 	log.SetOutput(file)
	// }
	config := ConfigLoader()
	servers := config.Hosts
	interval := config.HealthCheckInterval

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	go func() {
		for range ticker.C {
			// fmt.Println("running", ticker)
			HealthCheck(servers)
		}
	}()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Println("Error Occured while connecting", err)
	}
	fmt.Println(listener.Addr().String())
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		log.Println("entered", conn.LocalAddr().Network())
		if err != nil {
			log.Println("error occured ", err)
			continue
		}
		go HandleConnection(conn, servers)
	}

}
