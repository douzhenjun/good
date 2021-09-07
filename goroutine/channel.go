package main
import "fmt"
func Count(ch chan int, i int) {
	ch <- i
	//fmt.Println("Counting")
}
func main() {
	//chs := make([]chan int, 10)
	//for i := 0; i < 10; i++ {
	//	chs[i] = make(chan int)
	//	go Count(chs[i], i)
	//}
	//for _, ch := range(chs) {
	//	data := <-ch
	//	fmt.Println(data)
	//}

	ch := make(chan int, 1)
	for {
		select {
		case ch <- 0:
		case ch <- 1:
		}
		i := <-ch
		fmt.Println("Value received:", i)
	}
}
