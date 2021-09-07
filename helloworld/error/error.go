package error

import (
	"fmt"
	"os"
)

func Error(){
	_, err := os.Open("filename.txt")
	if err != nil {
		fmt.Println(err)
	}
}

