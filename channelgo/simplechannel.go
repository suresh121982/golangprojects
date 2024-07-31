package main

import (
	"fmt"
	"sync"
	"time"
)

// Function to simulate work and send a message to the channel
func worker(id int, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second) // Simulate work
		msg := fmt.Sprintf("Worker %d: Message %d", id, i)
		ch <- msg // Send message to the channel
	}
	close(ch) // Close the channel when done
}

func main() {

	var wg sync.WaitGroup

	// Create a channel for string messages
	ch := make(chan string)

	// Start 3 worker goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, ch, &wg)
	}
	go func() {
		wg.Wait() // Wait for all workers to finish
		close(ch) // Close the channel when done
	}()
	// Receive and print messages from the channel
	for msg := range ch {
		fmt.Println(msg)
	}

	fmt.Println("All messages received.")
}
