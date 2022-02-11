package main

// thread could step on each other, this is the example for it,
// run it without lock Mutex then you will see the money at the result may will be diffirent from 100

import (
	"sync"
	"time"
)

var (
	money = 100
	lock  = sync.Mutex{}
)

func stingy() {
	for i := 0; i <= 1000; i++ {
		// lock.Lock()
		money += 10
		// lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Stingy OK")
}

func spendy() {
	for i := 0; i <= 1000; i++ {
		// lock.Lock()
		money -= 10
		// lock.Unlock()
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
