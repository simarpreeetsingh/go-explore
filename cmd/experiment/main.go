package main

import (
	"fmt"
	"math/rand/v2"
	"sync"

	"github.com/simarpreeetsingh/go-explore/lib/linkedlist"
)

func main() {
	sllHead := linkedlist.InitSinglyLinkedList(1798345)
	sllHead.Print()

	sllHead.Append(231134).Append(45367).Append(43509).Append(7590)

	sllHead.PrintList()

	sllHead.Append(64523123)

	fmt.Println("Deleting next node...")
	sllHead.DeleteNext()

	sllHead.PrintList()

	sllHead.Print()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go appendRandom(sllHead, 10_000, wg)
	wg.Add(1)
	go appendRandom(sllHead, 60_000, wg)
	wg.Add(1)
	go appendRandom(sllHead, 30_000, wg)
	wg.Add(1)
	go appendRandom(sllHead, 90_000, wg)
	wg.Add(1)
	go appendRandom(sllHead, 1_00_000, wg)
	wg.Add(1)
	go appendRandom(sllHead, 50_000, wg)

	wg.Wait()

	sllHead.Print()
	sllHead.PrintList()
}

func appendRandom(ll *linkedlist.SinglyLinkedList[int], count int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < count; i++ {
		ll.Append(rand.IntN(600))
	}
}
