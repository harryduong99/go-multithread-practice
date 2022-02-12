package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const numberOfThreads int = 8 // how many core your processor have ?

type Point2D struct {
	x int
	y int
}

var (
	r         = regexp.MustCompile(`\((\d*),(\d*)\)`)
	waitGroup = sync.WaitGroup{}
)

func findArea(inputchannel chan string) {
	for pointsStr := range inputchannel { // consume the work that master feeded

		var points []Point2D
		for _, p := range r.FindAllStringSubmatch(pointsStr, -1) {
			x, _ := strconv.Atoi(p[1])
			y, _ := strconv.Atoi(p[2])
			points = append(points, Point2D{x, y})
		}

		area := 0.0
		for i := 0; i < len(points); i++ {
			a, b := points[i], points[(i+1)%len(points)] // last point will go back to zero (eg: 5%5 = 0) to not goout af slice
			area += float64(a.x*b.y) - float64(a.y*b.x)
		}
		fmt.Println(math.Abs(area) / 2.0)
	}
	waitGroup.Done()
}

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	absPath, _ := filepath.Abs(currentPath + "/../")
	dat, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	text := string(dat)

	inputChannel := make(chan string, 1000) // main
	for i := 0; i < numberOfThreads; i++ {
		go findArea(inputChannel)
	}

	waitGroup.Add(numberOfThreads) // add  theads

	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		inputChannel <- line // Master feeding request for the worker
	}
	close(inputChannel) // signal to workers that work is done

	waitGroup.Wait() // wait till 8 thead is done

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)

}
