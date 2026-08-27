package linkedlist

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

func Explore(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	sllHead := InitSinglyLinkedList(1798345)
	sllHead.Print()

	sllHead.Append(231134).Append(45367).Append(43509).Append(7590)

	sllHead.PrintList()

	sllHead.Append(64523123)

	fmt.Println("Deleting next node...")
	sllHead.DeleteNext()

	sllHead.PrintList()

	sllHead.Print()

	internalWg := &sync.WaitGroup{}

	internalWg.Add(1)
	go appendRandom(sllHead, 10_000, internalWg)
	internalWg.Add(1)
	go appendRandom(sllHead, 60_000, internalWg)
	internalWg.Add(1)
	go appendRandom(sllHead, 30_000, internalWg)
	internalWg.Add(1)
	go appendRandom(sllHead, 90_000, internalWg)
	internalWg.Add(1)
	go appendRandom(sllHead, 1_00_000, internalWg)
	internalWg.Add(1)
	go appendRandom(sllHead, 50_000, internalWg)

	internalWg.Wait()

	sllHead.Print()
	sllHead.PrintList()
}

func appendRandom(ll *SinglyLinkedList[int], count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		ll.Append(rand.IntN(600))
	}
}
