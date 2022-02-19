package main

import (
	"math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	totalAccounts  = 50000
	maxAmountMoved = 10
	initialMoney   = 100
	threads        = 4
)

func performMovements(ledger *[totalAccounts]int32, locks *[totalAccounts]sync.Locker, totalTrans *int64) {
	for {
		accountA := rand.Intn(totalAccounts)
		accountB := rand.Intn(totalAccounts)
		for accountA == accountB {
			accountB = rand.Intn(totalAccounts) // a must not equal to B
		}
		amountToMove := rand.Int31n(maxAmountMoved)
		toLock := []int{accountA, accountB}
		sort.Ints(toLock)
		// if we don't have this block, the total amount will not persisted
		locks[toLock[0]].Lock() // always locking the lower value first
		locks[toLock[1]].Lock()

		atomic.AddInt32(&ledger[accountA], -amountToMove)
		atomic.AddInt32(&ledger[accountB], amountToMove)
		atomic.AddInt64(totalTrans, 1)

		locks[toLock[1]].Unlock()
		locks[toLock[0]].Unlock()
	}
}

func main() {
	println("Total accounts:", totalAccounts, " total threads:", threads, "using SpinningLocks")
	var ledger [totalAccounts]int32
	var locks [totalAccounts]sync.Locker
	var totalTrans int64
	for i := 0; i < totalAccounts; i++ {
		ledger[i] = initialMoney
		locks[i] = NewSpinLock() // spinning lock will be faster than mutex, but not every case
	}

	for i := 0; i < threads; i++ {
		go performMovements(&ledger, &locks, &totalTrans) // create 4 threads
	}
	for {
		time.Sleep(2000 * time.Millisecond)
		var sum int32
		for i := 0; i < totalAccounts; i++ {
			locks[i].Lock()
		}
		for i := 0; i < totalAccounts; i++ {
			sum += ledger[i]
		}
		for i := 0; i < totalAccounts; i++ {
			locks[i].Unlock()
		}
		println(totalTrans, sum)
	}
}
