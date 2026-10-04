package stack

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func Explore(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	ls := InitLifoStack(30i + 1)

	internalWg := &sync.WaitGroup{}
	internalWg.Add(1)
	go appendRandom(ls, 10_000, internalWg)
	internalWg.Add(1)
	go appendRandom(ls, 60_000, internalWg)
	internalWg.Add(1)
	go appendRandom(ls, 30_000, internalWg)
	internalWg.Add(1)
	go appendRandom(ls, 90_000, internalWg)
	internalWg.Add(1)
	go appendRandom(ls, 1_00_000, internalWg)
	internalWg.Add(1)
	go appendRandom(ls, 50_000, internalWg)

	internalWg.Wait()
	fmt.Println(ls.Peek().Data())
	ls.Lookup(
		func(n *StackNode[complex128]) bool { fmt.Println(n, n.Data()); return false },
	)
}

func appendRandom(ll *LifoStack[complex128], count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		ll.Push(complex(float64(i), rand.Float64()))
	}
}
