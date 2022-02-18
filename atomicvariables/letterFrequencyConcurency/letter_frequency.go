package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const allLeters = "abcdefghijklmnopqrstuvwxyz"

var lock = sync.Mutex{}

func countLetters(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		// lock.Lock() // using mutex to make sure the result of frequency (line25) will not affected
		// But we also can use atomic variable for this
		index := strings.Index(allLeters, c)
		if index >= 0 {
			atomic.AddInt32(&frequency[index], 1)
		}
		// lock.Unlock()
	}

	wg.Done()
}

func main() {
	var frequency [26]int32
	wg := sync.WaitGroup{} // remember that just using wait group here will give us an incorect result since thread could step on each other. Let solve it by mutex
	start := time.Now()
	for i := 1000; i <= 1200; i++ {
		wg.Add(1)
		go countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("Done")
	fmt.Printf("Processing took %s\n", elapsed)
	for i, f := range frequency {
		fmt.Println("%s -> %d\n", string(allLeters[i]), f)
	}
}
