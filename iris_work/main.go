package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"iris_work/configs"
	"iris_work/controller"
	"iris_work/datasource"
	"iris_work/service"
)

func main() {

	app := newApp()

	//应用App设置
	configation(app)



	config := configs.GetConfig()
	//路由设置
	mvcHandle(app,config)
	addr := ":" + config.Port
	app.Run(
		iris.Addr(addr),                               //在端口8080进行监听
		iris.WithoutServerError(iris.ErrServerClosed), //无服务错误提示
		iris.WithOptimizations,                        //对json数据序列化更快的配置
	)
}

//构建App
func newApp() *iris.Application {
	app := iris.New()

	//设置日志级别  开发阶段为debug
	app.Logger().SetLevel("debug")

	////注册静态资源
	//app.HandleDir("/static", "./static")
	//app.HandleDir("/manage/static", "./static")
	//app.HandleDir("/img", "./static/img")
	//
	////注册视图文件
	//app.RegisterView(iris.HTML("./static", ".html"))
	//app.Get("/", func(context context.Context) {
	//	context.View("index.html")
	//})

	return app
}


/**
* MVC 架构模式处理
*/
func mvcHandle(app *iris.Application, c *configs.AppConfig) {
	//初始化engine
	engine := datasource.NewPgEngine(c)

	userService := service.NewUserService(engine)

	// user模块
	user := mvc.New(app.Party("/api/user"))
	user.Register(userService)
	user.Handle(new(controller.UserController))

}

/**
 * 项目设置
 */
func configation(app *iris.Application) {

	//配置 字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))

	//错误配置
	//未发现错误
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    " not found ",
			"data":   iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    " interal error ",
			"data":   iris.Map{},
		})
	})
}