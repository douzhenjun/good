package main

import "runtime"

func main(){
	c := make(chan struct{})
	go func(i chan struct{}){
		sum := 0
		for i := 0; i < 100; i++{
			sum += i
		}
		println(sum)
		c <- struct{}{}
	}(c)

	println("NumGoroutine=", runtime.NumGoroutine())
	<- c
}
