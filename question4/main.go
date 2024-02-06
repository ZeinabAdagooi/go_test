package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"
)

func wirter(ctx context.Context, ch chan string, e chan bool, nw int) {

	for i := 1; i <= nw; i++ {

		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Println("Error acquiring semaphore:", err)
			return
		}
		if err := qem.Acquire(ctx, 1); err != nil {
			fmt.Println("Error acquiring semaphore:", err)
			return
		}

		message := fmt.Sprintf("Message %d from the producer.", i)
		ch <- message

		e <- false

		sem.Release(1)
		qem.Release(1)

		//time.Sleep(1 * time.Second)
	}
	e <- true
}

func reader(ctx context.Context, ch chan string, done chan struct{}, e chan bool, mr int) {
	defer close(done)

	for {

		if err := qem.Acquire(ctx, 1); err != nil {
			fmt.Println("Error acquiring semaphore:", err)
			return
		}
		if err := sem.Acquire(ctx, 1); err != nil {
			fmt.Println("Error acquiring semaphore:", err)
			return
		}
		ed := <-e
		if ed {

			return
		}
		message := <-ch

		fmt.Println("Received:", message)

		qem.Release(1)

		sem.Release(1)

		time.Sleep(1 * time.Second)
	}

}

var (
	qem *semaphore.Weighted
	sem *semaphore.Weighted
	eem *semaphore.Weighted
)

func main() {
	// Create a buffered channel with a capacity of 5

	sem = semaphore.NewWeighted(int64(5))
	qem = semaphore.NewWeighted(int64(5))
	eem = semaphore.NewWeighted(int64(5))

	// Create a context to cancel the goroutines when the main program finishes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	n := 2
	m := 8

	messageChannel := make(chan string, 5)

	done := make(chan struct{})
	e := make(chan bool, 1)

	// Start goroutines
	go wirter(ctx, messageChannel, e, n)
	go reader(ctx, messageChannel, done, e, m)

	select {
	case <-done:
		fmt.Printf(" received done \n")
		return
	}
}
