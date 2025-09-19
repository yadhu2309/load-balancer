package main

import (
	// "fmt"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
)

func TCPLoadBalancer(pools []*Pool, rateLimitConfig *RateLimit) {

	var wg sync.WaitGroup

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Println("Error Occured while connecting", err)
	}
	// fmt.Println(listener.Addr().String())
	defer listener.Close()
	// go AutoTune()

	for {
		conn, err := listener.Accept()
		// log.Println("connection>>>>>>>>", conn)
		if err != nil {
			log.Println("error occured ", err)
			continue
		}

		if rateLimitConfig != nil && rateLimitConfig.Capacity > 0 &&
			rateLimitConfig.RefillRatePerSec > 0 {
			ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
			// fmt.Println("ip", ip)
			ratelimiter := GetClientBucket(ip, rateLimitConfig.Capacity,
				rateLimitConfig.RefillRatePerSec)
			if ratelimiter == nil {
				// fmt.Println("nil")
				continue
			}
			if !ratelimiter.Allow() {
				log.Println("Rate limit exceeded for", ip)
				SendTCPError(conn, 429, "HTTP/1.1 429 Too Many Requests\r\n\r\n")
				// conn.Close()
				continue
			}
		}
		log.Println("Trying to acquire semaphore...")
		Sema.Aquire()
		log.Println("Acquired semaphore")

		// fmt.Println(Sema.SemaChannel, "sema")
		wg.Add(1)
		go func(conn net.Conn) {
			defer func() {
				log.Println("Releasing semaphore")
				Sema.Release()
				wg.Done()
			}()
			// log.Println("++++++++++++++++++++++++++++++++++++++")

			servers := pools[0].Hosts

			HandleConnection(conn, servers)

		}(conn)

	}
	// wg.Wait()

}

func SendTCPError(conn net.Conn, status int, message string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		status, http.StatusText(status), len(message), message,
	)
	conn.Write([]byte(response))
	conn.Close()
}
