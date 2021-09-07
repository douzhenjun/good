// closure.go
package main

import (
	"fmt"
)

func main() {
	var j int = 5

	a := func() func() {
		var i int = 10
		return func() {
			fmt.Printf("i, j: %d, %d\n", i, j)
		}
	}() //末尾的括号表明匿名函数被调用，并将返回的函数指针赋给变量a

	a()

	j *= 2

	a()
}