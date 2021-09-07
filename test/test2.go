package main

import "fmt"

func fibonacci() func(){
	var a,b,c int
	return func(){
		if a==0 && b==0 && c==0{
			fmt.Println(a)
			b = 1
		}else if a==0 && b==1 && c==0{
			fmt.Println(b)
			c = a+b
		}else{
			fmt.Println(c)
			a = b
			b = c
			c = a+b
		}
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		f()
	}
}
