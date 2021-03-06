package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	matrixSize = 100
)

var (
	matrixA = [matrixSize][matrixSize]int{}
	matrixB = [matrixSize][matrixSize]int{}
	result  = [matrixSize][matrixSize]int{}
)

func generateRandomMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			result[row][col] += rand.Intn(10) - 5
		}
	}
}

func workOutRow(row int) {
	for col := 0; col < matrixSize; col++ {
		for i := 0; i < matrixSize; i++ {
			result[row][col] += matrixA[row][i] * matrixB[i][col]
		}
	}
}

func main() {
	fmt.Println("Processing...")
	start := time.Now()
	for i := 0; i < 100; i++ {

		generateRandomMatrix(&matrixA)
		generateRandomMatrix(&matrixB)
		for row := 0; row < matrixSize; row++ {
			workOutRow(row) // each worker will handle this
			fmt.Printf("%v\n", result[row])
		}
	}
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Println("Processing took: ", elapsed)

}
