package main

import "fmt"

type Student struct {
	name string
	age int
	Class string
}
func main() {
	var stu1 Student
	stu1.age = 22
	stu1.name = "wd"
	stu1.Class = "class1"
	fmt.Println(stu1.name)  //wd

	var stu2 *Student = new(Student)
	stu2.name = "jack"
	stu2.age = 33
	fmt.Println(stu2.name,(*stu2).name)//jack jack

	var stu3 *Student = &Student{ name:"rose",age:18,Class:"class3"}
	fmt.Println(stu3.name,(*stu3).name) //rose  rose


}
