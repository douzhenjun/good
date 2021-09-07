package main1

import "fmt"

func Var1(){
	arr := [5]int{1, 2, 3, 4, 5}
	for i := 0; i < len(arr); i++ {
		arr[i] += 100
	}
	fmt.Println(arr)  // [101 102 103 104 105]
}
