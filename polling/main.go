package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(new(sync.Mutex))
	wg := new(sync.WaitGroup)

	// mock poll
	go func() {
		for id := range 2 {
			wg.Add(1)
			go func(id int) {
				// receiver
				s := time.Duration(rand.Uint64N(4))
				ctx, cancel := context.WithTimeout(context.Background(), s*time.Second)
				defer cancel()

				fmt.Printf("new request to %d, wait %d\n", id, s)
				msg := Poll(ctx, c, wg)

				fmt.Printf("received message from (%d): %s\n", id, msg)
			}(id)
		}
	}()

	// sender
	time.Sleep(2 * time.Second)
	// prevent race condition
	c.L.Lock()
	messages = append(messages, "Hi!")
	c.L.Unlock()
	// broadcast to waiter after send message
	c.Broadcast()

	wg.Wait()
}

func Poll(ctx context.Context, c *sync.Cond, wg *sync.WaitGroup) []string {
	defer wg.Done()
	c.L.Lock()
	defer c.L.Unlock()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		// prevent for-loop still run after context cancel
		case <-ctx.Done():
			c.Broadcast()
		case <-done:
			// do nothing because no loop wait anymore
		}
	}()

	// check condition to prevent missing wake up
	for len(messages) == 0 {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		c.Wait()
	}

	// do whatever
	return messages
}

var messages []string
