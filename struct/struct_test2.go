package main

import "fmt"

type Student1 struct {
	name string
	age int
	Class string
}

func Newstu(name1 string,age1 int,class1 string) Student1 {
	return Student1{name:name1,age:age1,Class:class1}
}
func main() {
	stu1 := Newstu("wd",22,"math")
	fmt.Println(stu1.name) // wd
}
