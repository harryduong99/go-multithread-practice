package main

import "time"

func waitOnBarrier(name string, timeToSleep int, barrier *Barrier) {
	for {
		println(name, "Running")
		time.Sleep(time.Duration(timeToSleep) * time.Second)
		println(name, "is waiting on barrier")
		barrier.Wait()
	}
}
func main() {
	barrier := NewBarrier(2)
	go waitOnBarrier("red", 4, barrier)    // "Red is wating on barrier" after 4 second, then
	go waitOnBarrier("green", 10, barrier) //"Green is wating on barrier" after 10 second, then, both of them running again, then thing happen agian, with the same behavior
	time.Sleep(time.Duration(100) * time.Second)
}
