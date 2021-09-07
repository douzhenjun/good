package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

/**
	文件复制, 将path1处的文件复制到Path2处
 */
func CopyFile(path1, path2 string){
	//1.读取文件内容, 首先获得文件管理对象f1
	f1, err1 := os.Open(path1)
	if err1 != nil{
		fmt.Println("err1=", err1)
		return
	}
	//2.关闭文件f1
	defer f1.Close()

	//3.新建一个文件f2
	f2, err2 := os.Create(path2)
	if err2 != nil{
		fmt.Println("err2=", err2)
		return
	}
	//4.关闭文件f2
	defer f2.Close()

	//5.建立f1的文件读取缓存区, f2的文件写入
	r := bufio.NewReader(f1)
	w := bufio.NewWriterSize(f2, 1024)
	for{
		//6.buf返回一个字节数组, 每读取一行, 得到一个buf数组
		buf, err := r.ReadBytes('\n')

		//fmt.Printf("%v\n", buf)
		if err != nil{
			//读取到文件结尾则退出死循环
			if err == io.EOF{
				break
			}
			fmt.Println("err=", err)
		}

		//7.w向f2文件写入buf中的内容, 再刷新
		if _, err = w.Write(buf); err != nil{
			fmt.Println("err=", err)
		}
		if err = w.Flush(); err != nil{
			fmt.Println("err=", err)
		}
	}
}

/**
	文件移动
 */

func MoveFile(path1, path2 string){
	CopyFile(path1, path2)
	os.Remove(path1)
}
func main(){
	var path1 string = `D:\go_work\aaa.txt`
	var path2 string = `D:\go_work\bbb.txt`
	var path3 string = `D:\go_work\ccc.txt`

	CopyFile(path1, path2)
	MoveFile(path1, path3)
}
