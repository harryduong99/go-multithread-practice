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
	"time"
)

type Point2D struct {
	x int
	y int
}

var (
	r = regexp.MustCompile(`\((\d*),(\d*)\)`)
)

func findArea(pointsStr string) {
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
func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	absPath, _ := filepath.Abs(currentPath + "/../")
	dat, _ := ioutil.ReadFile(filepath.Join(absPath, "polygons.txt"))
	text := string(dat)

	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		// line := "(4,10),(12,8),(10,3),(2,2),(7,5)"
		findArea(line)
	}

	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)

}
