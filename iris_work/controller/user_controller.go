package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"iris_work/service"
)

type UserController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//host功能实体
	Service	service.UserService
}

func (uc *UserController) PostAdd() mvc.Result {
	username := uc.Ctx.PostValue("name")
	mobile := uc.Ctx.PostValue("mobile")
	password := uc.Ctx.PostValue("password")
	cityName := uc.Ctx.PostValue("city_name")

	re := uc.Service.Add(username, mobile, password, cityName)

	if re.Code == 1 {
		return mvc.Response{
			Object: map[string]interface{} {
				"no": re.Code,
				"msg": "插入成功",
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{} {
			"errorno": re.Code,
			"error_msg_en": re.Error_msg_en.Error(),
			"error_msg_zh": re.Error_msg_zh.Error(),
		},
	}

}