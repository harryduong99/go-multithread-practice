package main

// Run this file and you can see that the time it took is so long

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

var (
	windRegex     = regexp.MustCompile(`\d* METAR.*EGLL \d*Z [A-Z ]*(\d{5}KT|VRB\d{2}KT).*=`)
	tafValidation = regexp.MustCompile(`.*TAF.*`)
	comment       = regexp.MustCompile(`\w*#.*`)
	metarClose    = regexp.MustCompile(`.*=`)
	variableWind  = regexp.MustCompile(`.*VRB\d{2}KT`)
	validWind     = regexp.MustCompile(`\d{5}KT`)
	windDirOnly   = regexp.MustCompile(`(\d{3})\d{2}KT`)
	windDist      [8]int
)

func parseToArray(text string) []string {
	lines := strings.Split(text, "\n")
	metarSlice := make([]string, 0, len(lines))
	metarStr := ""
	for _, line := range lines {
		if tafValidation.MatchString(line) {
			break
		}
		if !comment.MatchString(line) {
			metarStr += strings.Trim(line, " ")
		}
		if metarClose.MatchString(line) {
			metarSlice = append(metarSlice, metarStr)
			metarStr = ""
		}
	}
	return metarSlice
}

func mineWindDistribution(winds []string) {
	for _, wind := range winds {
		if variableWind.MatchString(wind) {
			for i := 0; i < 8; i++ {
				windDist[i]++
			}
		} else if validWind.MatchString(wind) {
			windStr := windDirOnly.FindAllStringSubmatch(wind, -1)[0][1]
			if d, err := strconv.ParseFloat(windStr, 64); err == nil {
				dirIndex := int(math.Round(d/45.0)) % 8
				windDist[dirIndex]++
			}
		}
	}
}

func extractWindDirection(metars []string) []string {
	winds := make([]string, 0, len(metars))
	for _, metar := range metars {
		if windRegex.MatchString(metar) {
			winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
		}
	}
	return winds
}

func main() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	abspath, _ := filepath.Abs(currentPath + "/../../metarfiles")
	files, _ := ioutil.ReadDir(abspath)

	start := time.Now()
	for _, file := range files {
		dat, err := ioutil.ReadFile(filepath.Join(abspath, file.Name()))
		if err != nil {
			panic(err)
		}
		text := string(dat)
		//1. Change to array, each meta report is a item in array
		metarsReports := parseToArray(text)

		// 2. Extract wind direction, 200904300150 METAR COR EGLL 300150Z 16005KT 9999 SCT023 SCT034 09/06 Q1014
		windDirections := extractWindDirection(metarsReports)

		// 3. Assign to N, NE, E, SE, S, SW, NW, ... -> E + 1
		mineWindDistribution(windDirections)

	}

	elapsed := time.Since(start)
	fmt.Printf("%v\n", windDist)
	fmt.Printf("Proccessing took", elapsed)
}
