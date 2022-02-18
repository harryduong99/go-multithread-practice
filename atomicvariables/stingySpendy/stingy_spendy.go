package main

// Using atomic var instead of mutex

import (
	"sync/atomic"
	"time"
)

var (
	money int32 = 100
)

func stingy() {
	for i := 0; i <= 1000; i++ {
		atomic.AddInt32(&money, 10)
		time.Sleep(1 * time.Millisecond)
	}
	println("Stingy OK")
}

func spendy() {
	for i := 0; i <= 1000; i++ {
		atomic.AddInt32(&money, -10)
		time.Sleep(1 * time.Millisecond)
	}
	println("Spendy OK")
}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3000 * time.Millisecond)
	print(money)
}
