/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date:  2021-03-01 09:44:07
 * @LastEditors: zhangwei
 * @LastEditTime: 2021-03-01 09:44:07
 */

package controller

import (
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type PodController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx              iris.Context
	CommonService    service.CommonService
	Service          service.PodService
}

func (pc *PodController) GetLog() mvc.Result {
	podId, _ := pc.Ctx.URLParamInt("podId")
	log, err := pc.Service.GetLog(podId)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(log)
}

func (pc *PodController) PostSwitch() mvc.Result {
	id, err := pc.Ctx.PostValueInt("id")
	if err != nil {
		return response.Fail(response.ErrorParameter)
	}
	opUser := pc.Ctx.GetCookie("userName")
	success, errMsgEn, errMsgZn, cluster := pc.CommonService.SwitchCluster(id, false)
	if success {
		pc.CommonService.AddLog("info", "system-cluster", opUser, fmt.Sprintf("switch cluster %v successfully ", cluster.Name))
		return response.Success(nil)
	} else {
		pc.CommonService.AddLog("error", "system-cluster", opUser, fmt.Sprintf("switch cluster %v error: %s", cluster.Name, errMsgEn))
		return response.FailMsg(errMsgEn, errMsgZn)
	}
}

func (pc *PodController) GetDetail() mvc.Result {
	selectType := pc.Ctx.URLParam("type")
	condition := pc.Ctx.URLParam("condition")
	conditionMap := map[string]interface{}{}
	if len(condition) != 0 {
		_ = json.Unmarshal([]byte(condition), &conditionMap)
	}
	podId, _ := pc.Ctx.URLParamInt("podId")
	attrId := pc.Ctx.URLParamIntDefault("attrId", 0)
	modelId := pc.Ctx.URLParamIntDefault("modelId", 0)
	time, _ := pc.Ctx.URLParamInt("time")
	if time <= 0 {
		time = 3000
	}
	detail, err := pc.Service.GetDetail(podId, attrId, modelId, time, selectType, conditionMap)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(detail)
}
