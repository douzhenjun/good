/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date:  2021-02-22 09:44:07
 * @LastEditors: ddh
 * @LastEditTime: 2021-02-22 09:44:07
 */

package controller

import (
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type LogController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	Service service.LogService
	Common  service.CommonService
}

func (lc *LogController) GetList() mvc.Result {
	page := lc.Ctx.URLParamIntDefault("page", 0)
	pageSize := lc.Ctx.URLParamIntDefault("pagesize", 0)
	key := lc.Ctx.URLParam("key")
	typeString := lc.Ctx.URLParam("type")
	operatorName := lc.Ctx.URLParam("operatorName")
	container := lc.Ctx.URLParam("container")
	if typeString == "mysql-operator" {
		if container == "" || operatorName == "" {
			return response.Fail(response.ErrorParameter)
		}
		log, err := lc.Common.GetLogsByLoki(operatorName, container)
		if err != nil {
			utils.LoggerError(err)
			return response.Error(err)
		}
		return response.Success(log)
	} else {
		list, count, err := lc.Service.List(page, pageSize, key, typeString)
		if err != nil {
			utils.LoggerError(err)
			return response.Error(err)
		}
		return response.Success(map[string]interface{}{
			"detail":   list,
			"all":      count,
			"page":     page,
			"pagesize": pageSize,
		})
	}
}
