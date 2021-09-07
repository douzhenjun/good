package main

import (
	"net/http"
	"sync"
)

var wg sync.WaitGroup

var urls = []string{
	"http://www.golang.org/",
	"http://www.google.com/",
	"http://www.qq.com/",
}

func main() {
	for _, url := range urls {
		//每一个URL启动一个goroutine,同时给wg加1
		wg.Add(1)
		go func(url string) {
			//当前goroutine结束后给wg减1,wg.Done()等价于wg.Add(-1)
			defer wg.Done()
			//发送http get请求并打印http返回码
			resp, err := http.Get(url)
			if err == nil {
				println(resp.Status)
			}
		}(url)
	}
	//等待所有请求结束
	wg.Wait()
}
