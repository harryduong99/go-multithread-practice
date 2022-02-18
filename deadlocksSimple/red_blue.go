package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.Mutex{}
	lock2 = sync.Mutex{}
)

func blueRobot() {
	for {
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock() // blue is holding lock1
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock() // blue is going to accquire lock2, but lock2 is used by red, so it need to wait till red release lock 2
		// but red also can't release Lock because it can't get lock1, which is accquired by lock 1
		fmt.Println("Blue: Both locks accquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Blue lock released")
	}
}

func redRobot() {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock2.Lock() // now red already got lock 2
		fmt.Println("Red: Acquiring lock2")
		lock1.Lock() // can't get lock 1 be cause lock 1 is holding it, waiign for lock 1 to release lock 1
		// Now red and blue are waiting for each other
		fmt.Println("Red: Both locks accquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: Red lock released")
	}
}

func main() {
	go blueRobot()
	go redRobot()
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
