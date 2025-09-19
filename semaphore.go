package main

import (
	"fmt"
	"runtime"
	"time"
)

func InitSemaphore(size int) *Semaphore {
	return &Semaphore{SemaChannel: make(chan struct{}, size)}
}

func (Sema *Semaphore) Aquire() {
	Sema.SemaChannel <- struct{}{}
}

func (Sema *Semaphore) Release() {
	<-Sema.SemaChannel
}

func AutoTune() {
	fmt.Println("==================== Auto Tune ===========================")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		numCpu := runtime.NumCPU()
		fmt.Println("Number of CPU", numCpu)
		numGo := runtime.NumGoroutine()
		fmt.Println("Number of Goroutine", numGo)

	}
}
