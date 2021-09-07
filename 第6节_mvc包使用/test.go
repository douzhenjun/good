package main

import "fmt"

func main(){
	c := add(3, 4)
	fmt.Println(c)
}


func add(a, b int) (sum int){
	//sum = a + b
	//return
	anonymous := func(x, y int) int {
		return x+y
	}
	return anonymous(a,b)
	//sum := a+b
	//return sum
}