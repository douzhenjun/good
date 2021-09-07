package main

import (
	"github.com/kataras/iris/v12"
)

func main(){
	app := iris.New()

	//userParty := app.Party("/users", func(context context.Context) {
		// 处理下一级请求
	//	context.Next()
	//})

	app.Run(iris.Addr(":8004"), iris.WithoutServerError(iris.ErrServerClosed))
}
