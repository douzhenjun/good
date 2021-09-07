package controller

import (
	"DBaas/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type TestController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//good
	Service        service.TestService
}

func (tc *TestController) GetHello() mvc.Result {
	return mvc.Response{
		Object: map[string]interface{}{
			"shuai": "go",
		},
	}
}
