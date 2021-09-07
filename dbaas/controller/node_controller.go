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
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type NodeController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	Service service.NodeService

	CommonService service.CommonService
}

func (nc *NodeController) PostOperatorImage() mvc.Result {
	err := nc.Service.OperatorImage()
	opUser := nc.Ctx.GetCookie("userName")
	if err != nil {
		utils.LoggerError(err)
		nc.CommonService.AddLog("error", "system-node", opUser, fmt.Sprintf("apply operator image error %v", err))
		return response.Error(err)
	}
	nc.CommonService.AddLog("info", "system-node", opUser, "apply operator image successful")
	return response.Success(nil)
}

func (nc *NodeController) GetList() mvc.Result {
	page, err := nc.Ctx.URLParamInt("page")
	if err != nil {
		page = 0
	}
	pageSize, err := nc.Ctx.URLParamInt("pagesize")
	if err != nil {
		pageSize = 0
	}
	key := nc.Ctx.URLParam("key")
	list, count := nc.Service.List(page, pageSize, key)
	mgmtTotal, computeTotal := nc.Service.ComputedSum()
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data": map[string]interface{}{
				"detail": list,
				"all":    count,

				"page":     page,
				"pagesize": pageSize,
			},
			"computeTotal": computeTotal,
			"mgmtTotal":    mgmtTotal,
			"error_msg_en": "",
			"error_msg_zh": "",
		},
	}
}

//  添加 删除 label
func (nc *NodeController) PostTag() mvc.Result {
	tagList := nc.Ctx.PostValue("tagList")
	tagListMap := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(tagList), &tagListMap)
	userName := nc.Ctx.GetCookie("userName")
	if err != nil {
		utils.LoggerError(err)
		nc.CommonService.AddLog("error", "system-node", userName, "node add or delete label error: "+utils.ERROR_PARAMETER_EN)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}

	for _, m := range tagListMap {
		success := true
		errMsg := ""
		if mgmtTag, ok := m["mgmtTag"]; ok {
			if mgmtTag.(bool) {
				success, errMsg = nc.Service.AddLabel(int(m["id"].(float64)), "iwhalecloud.dbassoperator", "mysqlha")
			} else {
				success, errMsg = nc.Service.DeleteLabel(int(m["id"].(float64)), "iwhalecloud.dbassoperator")
			}
		}

		if !success {
			nc.CommonService.AddLog("error", "system-node", userName, "node add or delete label error: "+errMsg)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}

		if mgmtTag, ok := m["computeTag"]; ok {
			if mgmtTag.(bool) {
				success, errMsg = nc.Service.AddLabel(int(m["id"].(float64)), "iwhalecloud.dbassnode", "mysql")
			} else {
				success, errMsg = nc.Service.DeleteLabel(int(m["id"].(float64)), "iwhalecloud.dbassnode")
			}
		}

		if !success {
			nc.CommonService.AddLog("error", "system-node", userName, "node add or delete label error: "+errMsg)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
	}
	nc.Service.AsyncDbLabel()
	nc.CommonService.AddLog("info", "system-node", userName, "node add or delete label successful")
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

// 部署operator
func (nc *NodeController) PostOperator() mvc.Result {
	mode := nc.Ctx.PostValue("mode")
	scName := nc.Ctx.PostValue("scName")
	userName := nc.Ctx.GetCookie("userName")
	err := nc.Service.Operator(mode, scName)
	if err != nil {
		nc.CommonService.AddLog("error", "system-node", userName, "operator error: "+err.Error())
		return response.Error(err)
	}
	nc.CommonService.AddLog("info", "system-node", userName, "operator successful")
	return response.Success(nil)
}

// operator
func (nc *NodeController) GetOperator() mvc.Result {
	page, err := nc.Ctx.URLParamInt("page")
	if err != nil {
		page = 0
	}
	pageSize, err := nc.Ctx.URLParamInt("pagesize")
	if err != nil {
		pageSize = 0
	}
	key := nc.Ctx.URLParam("key")
	list, errMsg, count := nc.Service.OperatorPodList(page, pageSize, key)
	if errMsg != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errMsg,
				"error_msg_zh": errMsg,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data": map[string]interface{}{
				"detail":   list,
				"all":      count,
				"page":     page,
				"pagesize": pageSize,
			},
			"error_msg_en": "",
			"error_msg_zh": "",
		},
	}
}

func (nc *NodeController) GetOperatorStatus() mvc.Result {
	data := nc.Service.GetOperatorStatus()
	return response.Success(data)
}

// operator 部署日志
func (nc *NodeController) GetLog() mvc.Result {
	key := nc.Ctx.URLParam("key")
	list, count := nc.Service.OperatorLogList(key)
	mode := nc.Service.GetOperatorMode()
	return response.Success(map[string]interface{}{
		"detail": list,
		"mode":   mode,
		"all":    count,
	})
}
