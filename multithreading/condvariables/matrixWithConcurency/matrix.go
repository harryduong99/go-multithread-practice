package main

// this example showing the way we multiplcation for a pair of matrix.
// We will have to calculate 100 pairs of them.
// instead of create 100 pairs then calculate, we will do it paralel, create and calculate at the same time

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	matrixSize = 100
)

var (
	matrixA   = [matrixSize][matrixSize]int{}
	matrixB   = [matrixSize][matrixSize]int{}
	result    = [matrixSize][matrixSize]int{}
	rwLock    = sync.RWMutex{}
	cond      = sync.NewCond(rwLock.RLocker())
	waitGroup = sync.WaitGroup{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			result[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	rwLock.RLock()
	for { // do it "forever"
		waitGroup.Done() // done from the previous calculation
		cond.Wait()      // waiting for the broadcast from the main thread
		for col := 0; col < matrixSize; col++ {
			for i := 0; i < matrixSize; i++ {
				result[row][col] += matrixA[row][i] * matrixB[i][col]
			}
		}
	}

}

func main() {
	fmt.Println("Processing...")
	waitGroup.Add(matrixSize)
	for row := 0; row < matrixSize; row++ {
		go workOutRow(row) // each worker will handle this, the number of thread is the size of matrix
	}

	start := time.Now()
	for i := 0; i < 100; i++ {
		waitGroup.Wait() // wating for the worker (calculation for previous pair) done (line 38)
		rwLock.Lock()    // Lock the writing
		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		waitGroup.Add(matrixSize) // since waitgroup was released (Done()) (line 38), it was assinged to zero, we need to reset it
		rwLock.Unlock()
		cond.Broadcast() // signal to line 39 that the creating matrix is done, it can be calculate now
	}
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Println("Processing took: ", elapsed)

}
