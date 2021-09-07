package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"pra-iris/configs"
	"pra-iris/controller"
	"pra-iris/datasource"
	"pra-iris/service"
)

func main() {
	app := iris.New()
	c := configs.GetConfig()
	mvcHandle(app, c)
	addr := ":" + c.Port
	// 启动
	_ = app.Run(
		iris.Addr(addr),
		)
}

func mvcHandle(app *iris.Application, c *configs.AppConfig) {
	//初始化engine
	engine := datasource.NewPgEngine(c)

	userService := service.NewUserService(engine)

	// user模块
	user := mvc.New(app.Party("/api/user"))
	user.Register(userService)
	user.Handle(new(controller.UserController))

}