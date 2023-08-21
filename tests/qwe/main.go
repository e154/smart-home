package main

import (
	"fmt"
	"sync"
)

const countWorkers = 1000
const countTasks = 10000000

type Result struct {
	value    float64
	hasError bool
}

func MakeTasks(count int) <-chan Result {
	//var ch chan Result
	ch := make(chan Result, 100)
	go func() {
		for i := 0; i < count; i++ {
			ch <- Result{
				value:    float64(i) * 2.42,
				hasError: (i % 10) == 0,
			}

		}
		close(ch)
	}()
	return ch
}

func ProcessUsingMutex(ch <-chan Result, countWorkers int) (float64, int64) {
	var wg sync.WaitGroup
	var errMu sync.Mutex
	var mu sync.Mutex
	var countErrors int64
	var result float64

	for i := 0; i < countWorkers; i++ {
		wg.Add(1)
		go func() {
			for item := range ch {
				if item.hasError {
					errMu.Lock()
					countErrors++
					errMu.Unlock()
				} else {
					mu.Lock()
					result += item.value
					mu.Unlock()
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return result, countErrors
}

func ProcessUsingChannels(ch <-chan Result, countWorkers int) (float64, int64) {

	var countErrors int64
	var result float64
	retport := make(chan interface{})

	var wg sync.WaitGroup
	wg.Add(1)
	go func(){
		for item := range retport {
			switch v := item.(type) {
			case float64:
				result += v
			case int64:
				countErrors += v
			}
		}
		wg.Done()
	}()

	var worker sync.WaitGroup
	for i := 0; i < countWorkers; i++ {
		worker.Add(1)
		go func() {
			for item := range ch {
				if item.hasError {
					retport <- 1
				} else {
					retport <- item.value
				}
			}
			worker.Done()
		}()
	}
	worker.Wait()

	close(retport)

	wg.Wait()

	return result, countErrors
}

func main() {
	ch := MakeTasks(countTasks)
	//fmt.Println(ProcessUsingMutex(ch, countWorkers))
	fmt.Println(ProcessUsingChannels(ch, countWorkers))
}

