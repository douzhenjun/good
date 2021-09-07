package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)
	ch3 := make(chan int, 3)


	wg.Add(3)
	go say11(wg, ch1, ch2)
	go say12(wg, ch2, ch3)
	go say13(wg, ch1, ch3)
	wg.Wait()
	time.Sleep(1 * time.Second)

}


func say11(wg *sync.WaitGroup, ch1 chan int, ch2 chan int) {
	defer wg.Done()
	for i := 1; i <= 7; i++ {
		ch2 <- 3*i-1
		fmt.Println(<-ch1)
	}
}

func say12(wg *sync.WaitGroup, ch2 chan int, ch3 chan int) {
	defer wg.Done()
	for i := 1; i <= 7; i++ {
		ch3 <- 3*i
		fmt.Println(<-ch2)
	}
}

func say13(wg *sync.WaitGroup, ch1 chan int, ch3 chan int){
	defer wg.Done()
	for i := 1; i <= 7; i++{
		ch1 <- 3*i-2
		if i < 7{
			fmt.Println(<-ch3)
		}
	}
}
