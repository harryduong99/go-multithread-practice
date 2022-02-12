package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	money          = 100
	lock           = sync.Mutex{}
	moneyDeposited = sync.NewCond(&lock) // condition here !!!
)

func stingy() {
	for i := 0; i <= 1000; i++ {
		lock.Lock()
		fmt.Println("Stingy see balance of ", money)
		moneyDeposited.Signal()
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Stingy OK")
}

func spendy() {
	for i := 0; i <= 1000; i++ {
		lock.Lock()
		for money-20 < 0 {
			moneyDeposited.Wait() // Keep wating until we adding enough money
			// Locked until got signal from line 19
		}
		money -= 20
		fmt.Println("Spendy see balance of ", money)
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}
	println("Spendy Done")
}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3000 * time.Millisecond)
	print(money)
}
