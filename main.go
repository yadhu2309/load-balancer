package main

import (

	// "io"

	"log"
	"runtime"

	// "strings"
	"time"
)

func HandleRecover() {
	if r := recover(); r != nil {
		log.Println("error", r)
	}
}

var Sema = InitSemaphore(runtime.NumCPU() * 4)

func main() {
	log.Println("Load balancer!!!!!!!!!!!!!!!!!")
	
	config := ConfigLoader()
	if config == nil {
		log.Println("Config.json is empty or not found")
		return
	}
	tcpConfig, exists := (*config)["tcp"]
	if !exists {
		log.Println("No TCP configuration found in config.json")
		return
	}
	pools := tcpConfig.Pool
	if len(*pools) == 0 {
		log.Println("No pool found in config.json")
		return
	}
	interval := tcpConfig.HealthCheckInterval

	// HealthCheck(servers)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	servers := (*pools)[0].Hosts
	go func() {
		for range ticker.C {
			//
			HealthCheck(servers)
		}
	}()
	loadingStrategy := tcpConfig.LoadingStrategy
	log.Println("Loading Strategy:", loadingStrategy)
	TCPLoadBalancer(*pools, tcpConfig.RateLimit, loadingStrategy)
	// ratelimiter := GetClientBucket("127.0.0.1")
	// if !ratelimiter.Allow() {
	// 	fmt.Println("ratelimiter.........")
	// }

}

// interval := config.HealthCheckInterval
// isMultiplePool := config.IsMultiplePool

// // size := runtime.NumCPU() // number of logical CPUs
// // fmt.Println("cpu size", size)
// // numG := runtime.NumGoroutine()
// // fmt.Print("goruntime", numG)

// ticker := time.NewTicker(time.Duration(interval) * time.Second)
// go AutoTune()
// go func() {
// 	for range ticker.C {
// 		if isMultiplePool {

// 			for _, pool := range pools {
// 				servers := pool.Hosts
// 				HealthCheck(servers)
// 			}
// 		} else {
// 			log.Println("IsmultiplePool", false)
// 			servers := pools[0].Hosts
// 			HealthCheck(servers)
// 		}
// 		// fmt.Println("running", ticker)
// 	}
// }()
