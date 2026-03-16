package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	c := sync.NewCond(new(sync.Mutex))
	wg := new(sync.WaitGroup)

	go func() {
		for i := range 5 {
			wg.Add(1)
			go worker(i, c, wg)
		}
	}()

	time.Sleep(5 * time.Second)
	c.Broadcast()

	wg.Wait()
}

func worker(i int, c *sync.Cond, wg *sync.WaitGroup) {
	defer wg.Done()

	c.L.Lock()
	fmt.Println(i, "Wait")
	c.Wait()
	c.L.Unlock()

	s := rand.Uint64N(3)
	time.Sleep(time.Duration(s) * time.Second)
	fmt.Printf("%d Woke up in %ds\n", i, s)
}
