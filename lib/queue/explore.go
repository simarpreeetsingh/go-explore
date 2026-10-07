package queue

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func Explore(wg *sync.WaitGroup) {
	start := time.Now()
	fmt.Println("Exploring queues")
	if wg != nil {
		defer wg.Done()
	}
	q := InitQueue(30i + 1)

	internalWg := &sync.WaitGroup{}
	internalWg.Add(1)
	go appendRandom(q, 10_000, internalWg)
	internalWg.Add(1)
	go appendRandom(q, 60_000, internalWg)
	internalWg.Add(1)
	go appendRandom(q, 30_000, internalWg)
	internalWg.Add(1)
	go appendRandom(q, 90_000, internalWg)
	internalWg.Add(1)
	go appendRandom(q, 1_00_000, internalWg)
	internalWg.Add(1)
	go appendRandom(q, 50_000, internalWg)

	internalWg.Wait()
	fmt.Println(q.Peek().Data())
	q.Lookup(
		func(n *QueueNode[complex128]) bool { fmt.Println(n, n.Data()); return false },
	)

	fmt.Println("Queue size:", q.Size(), "time taken (ms):", time.Since(start).Milliseconds())
	fmt.Println("Finished exploring queues")
}

func appendRandom(ll *Queue[complex128], count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		ll.Enqueue(complex(float64(i), rand.Float64()))
	}
}
