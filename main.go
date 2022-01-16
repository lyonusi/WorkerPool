package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const VALIDATE_BUFFER = 100
const DELETE_BUFFER = 100

var validateWg sync.WaitGroup
var deleteWg sync.WaitGroup
var ids []string = []string{"1", "2", "3", "4", "5"}
var deleteChan = make(chan string)
var validateSemaphore = make(chan struct{}, VALIDATE_BUFFER)
var deleteSemaphore = make(chan struct{}, DELETE_BUFFER)
var validateCounter int32
var validateDoneCounter int32
var deleteCounter int32
var deleteDoneCounter int32
var validToDeleteCounter int32
var validateMu sync.Mutex

func validate(id string) {
	// validateMu.Lock()
	// validateCounter++
	// fmt.Println("validated ", validateCounter)
	// validateMu.Unlock()

	atomic.AddInt32(&validateCounter, 1)
	// fmt.Println("validate counter - ", atomic.AddInt32(&validateCounter, 1))
	defer validateWg.Done()
	// fmt.Println("___validating ", id)

	randInt := (rand.Intn(100) + 100)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	// fmt.Println("___[validate] task ", id, " done")
	rand := rand.Intn(100)

	if rand < 50 {
		// fmt.Println("___deleteChan <- ", id)
		deleteWg.Add(1)
		deleteChan <- id
		atomic.AddInt32(&validToDeleteCounter, 1)
		fmt.Println("___deleteChan <- ", id, " done")
	}
	<-validateSemaphore
	// fmt.Println("validate done counter - ", atomic.AddInt32(&validateDoneCounter, 1))
	atomic.AddInt32(&validateDoneCounter, 1)
	return
}

func delete(id string) {
	// fmt.Println("delete counter", atomic.AddInt32(&deleteCounter, 1))
	atomic.AddInt32(&deleteCounter, 1)
	defer deleteWg.Done()
	// fmt.Println("... deleting ", id)
	randInt := (rand.Intn(100) + 100)
	time.Sleep(time.Millisecond * time.Duration(randInt))
	fmt.Println("... [delete] task ", id, " done")
	<-deleteSemaphore
	// fmt.Println("delete done counter", atomic.AddInt32(&deleteDoneCounter, 1))
	atomic.AddInt32(&deleteDoneCounter, 1)
	return
}

func main() {
	run()
}

func run() {
	go func() {
		for i := range deleteChan {
			deleteSemaphore <- struct{}{}
			go delete(i)
		}
	}()

	for i := 0; i < 10000; i++ {
		validateSemaphore <- struct{}{}
		validateWg.Add(1)
		go validate(fmt.Sprint(i))
	}
	validateWg.Wait()
	// close(deleteChan)
	// fmt.Println("wait here!!!!")
	deleteWg.Wait()
	fmt.Println("validateCounter = ", validateCounter)
	fmt.Println("validateDoneCounter = ", validateDoneCounter)
	fmt.Println("deleteCounter = ", deleteCounter)
	fmt.Println("deleteDoneCounter = ", deleteDoneCounter)
	fmt.Println("validToDeleteCounter = ", validToDeleteCounter)
}
