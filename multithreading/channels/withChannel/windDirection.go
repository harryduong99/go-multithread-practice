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

func parseToArray(textChannel chan string, metarChannel chan []string) {
	for text := range textChannel { // read from channel
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
		metarChannel <- metarSlice
	}
	close(metarChannel)
}

func mineWindDistribution(windsChannel chan []string, distChannel chan [8]int) {
	for winds := range windsChannel {
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
	distChannel <- windDist
	close(distChannel)
}

func extractWindDirection(metarChannel chan []string, windsChannel chan []string) {
	for metars := range metarChannel {
		winds := make([]string, 0, len(metars))
		for _, metar := range metars {
			if windRegex.MatchString(metar) {
				winds = append(winds, windRegex.FindAllStringSubmatch(metar, -1)[0][1])
			}
		}
		windsChannel <- winds
	}
	close(windsChannel)
}

func main() {
	textChannel := make(chan string)
	metarChannel := make(chan []string)
	windsChannel := make(chan []string)
	resultsChannel := make(chan [8]int)
	//1. Change to array, each meta report is a item in array
	go parseToArray(textChannel, metarChannel)

	// 2. Extract wind direction, 200904300150 METAR COR EGLL 300150Z 16005KT 9999 SCT023 SCT034 09/06 Q1014
	go extractWindDirection(metarChannel, windsChannel)

	// 3. Assign to N, NE, E, SE, S, SW, NW, ... -> E + 1
	go mineWindDistribution(windsChannel, resultsChannel)

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
		textChannel <- text // the main thread will be blocked here, waiting for this text in textChannel is consumed

	}
	close(textChannel) // signal to other threards that the processing is done
	results := <-resultsChannel

	elapsed := time.Since(start)
	fmt.Printf("%v\n", results)
	fmt.Printf("Proccessing took", elapsed)
}
