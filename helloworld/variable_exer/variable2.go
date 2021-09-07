package main1

import "fmt"

func Var2(){

	slice := []int{1,2,3,4,5}
	// 创建新的切片，其长度为 2 个元素，容量为 4 个元素
	mySlice := slice[1:3]
	// 使用原有的容量来分配一个新元素，将新元素赋值为 40
	mySlice1 := append(mySlice, 40)

	printSlice(mySlice1) //len=3 cap=4 0xc00000a4b8 [2 3 40]

	printSlice(slice) //len=5 cap=5 0xc00000a4b0 [1 2 3 40 5]

	mySlice2 := slice[1:5]

	mySlice2 = append(mySlice2, 40)

	printSlice(mySlice2) //len=5 cap=8 0xc00000c240 [2 3 40 5 40]

	printSlice(slice) //len=5 cap=5 0xc00000a4b0 [1 2 3 40 5]
}

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %p %v\n", len(s), cap(s), s, s)
}