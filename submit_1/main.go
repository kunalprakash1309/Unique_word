package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	timeStart := time.Now()

	file, err := os.Open("input.txt")
	if err != nil {
		log.Println("Error in opening the file")
	}
	defer file.Close()



	result := wordFrequency(file)

	log.Println(result)

	uniqueWord(result)

	timeEnd := time.Now()
	log.Println("Total time taken is : ", timeEnd.Sub(timeStart))
}

func wordFrequency(file *os.File) (wordMap map[string]int) {

	chunkSize := 1024*1024
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	wordMap = make(map[string]int)

	doneChannel := make(chan bool)
	wordChannel := make(chan string)

	scanner := bufio.NewScanner(file)

	buf := make([]byte,chunkSize)
	scanner.Buffer(buf, chunkSize)
	// scanner.Buffer(buf, chunkSize)
	scanner.Split(bufio.ScanWords)

	go func() {
		for {
			select {
			case w := <- wordChannel:
				wordMap[w]++
			case <-doneChannel:
				return
			}
		}
	}()

	for scanner.Scan() {
		word := scanner.Text()
		wg.Add(1)
		
		go func(w string){
			defer wg.Done()
			w = reg.ReplaceAllString(w, "")
			if w != "" {
				wordChannel <- w
			}
		}(word)
	}

	
	wg.Wait()
	doneChannel <- true
	
	return 
}

func uniqueWord(Map map[string]int) {
	for key, value := range Map {
		if value == 1 {
			log.Println("First Unique word ", key)
			break
		}
	}
}

