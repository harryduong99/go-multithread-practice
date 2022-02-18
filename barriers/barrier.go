package main

import "sync"

type Barrier struct {
	total int
	count int
	mutex *sync.Mutex
	cond  *sync.Cond
}

func NewBarrier(size int) *Barrier {
	lockToUse := &sync.Mutex{}
	condToUse := sync.NewCond(lockToUse)
	return &Barrier{size, size, lockToUse, condToUse} // initial, count == size of thread
}

func (b *Barrier) Wait() {
	b.mutex.Lock()    // safely modify the count
	b.count -= 1      // 1 more thread is waiting
	if b.count == 0 { // if all if the threads have called wait function
		b.count = b.total
		b.cond.Broadcast()
	} else {
		b.cond.Wait() // keep wating
	}
	b.mutex.Unlock()
}
