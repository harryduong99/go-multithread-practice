package main

// This implementation is much faster than the nomal search that using on thread in fileseach !!!

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func fileSearchConcurency(root string, filename string) {
	fmt.Println("Searching in", root)
	files, _ := ioutil.ReadDir(root)
	for _, file := range files {
		if strings.Contains(file.Name(), filename) {
			lock.Lock() // ensure multiple thread won't write on each other, then the matches won't be wrong
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}

		if file.IsDir() {
			waitgroup.Add(1) // update the wait group
			go fileSearchConcurency(filepath.Join(root, file.Name()), filename)
		}
	}
	waitgroup.Done()
}

func main() {
	waitgroup.Add(1)
	go fileSearchConcurency("/home/duongnam/work/src/github.com/duongnam99", "README.md")
	waitgroup.Wait()
	for _, file := range matches {
		fmt.Println("Matched", file)
	}
}
