package main

import (
	"fmt"
	"sync"
	"time"
)

// Worker function that simulates a task
func worker(id int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done() // Notify the WaitGroup that this goroutine is done

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second) // Simulate some work
		message := fmt.Sprintf("Worker %d: task %d", id, i)
		ch <- message // Send message to the channel
	}
}

// Main function
func main() {
	// Create a channel to communicate between goroutines
	messageChannel := make(chan string)

	// Create three WaitGroups for different groups of workers
	var wgGroup1 sync.WaitGroup
	var wgGroup2 sync.WaitGroup
	var wgGroup3 sync.WaitGroup

	// Launch 3 worker goroutines for group 1
	for i := 1; i <= 3; i++ {
		wgGroup1.Add(1) // Increment the WaitGroup counter
		go worker(i, messageChannel, &wgGroup1)
	}

	// Launch 3 worker goroutines for group 2
	for i := 4; i <= 6; i++ {
		wgGroup2.Add(1) // Increment the WaitGroup counter
		go worker(i, messageChannel, &wgGroup2)
	}

	// Launch 3 worker goroutines for group 3
	for i := 7; i <= 9; i++ {
		wgGroup3.Add(1) // Increment the WaitGroup counter
		go worker(i, messageChannel, &wgGroup3)
	}

	// Launch a goroutine to print messages from the channel
	go func() {
		for msg := range messageChannel {
			fmt.Println(msg)
		}
	}()

	// Wait for all worker goroutines to complete for all groups
	go func() {
		wgGroup1.Wait()
		wgGroup2.Wait()
		wgGroup3.Wait()
		close(messageChannel) // Close the channel after all workers are done
	}()

	// Wait for a while to allow message printing
	time.Sleep(10 * time.Second) // Adjust this based on your needs

	fmt.Println("All workers completed.")
}
