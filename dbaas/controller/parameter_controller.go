/**
* @Description:
* @version:
* @Company: iwhalecloud
* @Author:  zhangwei
* @Date: 2021/02/24 10:30
* @LastEditors: zhangwei
* @LastEditTime: 2021/02/26 13:30
**/

package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	k8sContext "context"
	"fmt"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
	//"fmt"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type ParameterController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//host功能实体
	Service       service.ParameterService
	CommonService service.CommonService
}

//  获取系统参数列表
func (pc *ParameterController) GetList() mvc.Result {
	utils.LoggerInfo(" 获取系统参数列表 ")
	paramError := true
	pageSize, err := pc.Ctx.URLParamInt("pagesize")
	if err != nil {
		paramError = false
	}
	page, err := pc.Ctx.URLParamInt("page")
	if err != nil {
		paramError = false
	}
	key := pc.Ctx.URLParam("key")
	parameterList := make([]models.Sysparameter, 0)
	if paramError {
		parameterList, err = pc.Service.ListParameter(pageSize, (page-1)*pageSize, key)
	} else {
		parameterList, err = pc.Service.ListParameterAll(key)
	}
	if err != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_MSG_EN,
				"error_msg_zh": utils.ERROR_MSG_ZH,
			},
		}
	}
	parameterCount, err := pc.Service.GetParameterCount(key)
	returnData := make(map[string]interface{})
	returnData["page"] = page
	returnData["pagesize"] = pageSize
	returnData["all"] = parameterCount
	parameterDetailList := make([]interface{}, 0)
	if len(parameterList) > 0 {
		for _, parameter := range parameterList {
			parameterDetailMap := make(map[string]interface{})
			parameterDetailMap["id"] = parameter.Id
			parameterDetailMap["key"] = parameter.ParamKey
			parameterDetailMap["value"] = parameter.ParamValue
			parameterDetailMap["default"] = parameter.DefaultValue
			parameterDetailMap["isModifiable"] = parameter.IsModifiable
			parameterDetailList = append(parameterDetailList, parameterDetailMap)
		}
	}
	returnData["detail"] = parameterDetailList
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    returnData,
		},
	}
}

//  修改系统参数信息
func (pc *ParameterController) PostUpdate() mvc.Result {
	utils.LoggerInfo(" 修改系统参数信息")
	params := pc.Ctx.PostValue("params")
	if params == "" {
		return response.Fail(response.ErrorParameter)
	}
	opUser := pc.Ctx.GetCookie("userName")
	paramsList := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(params), &paramsList)
	utils.LoggerError(err)
	k8sConfigModify := false
	harborConfigModify := false
	namespaceModify := false
	k8sConfigValue := ""
	harborAddress := ""
	k8sMasterAddressValue := ""
	if len(paramsList) > 0 {
		for _, param := range paramsList {
			key := param["key"].(string)
			value := param["value"].(string)
			p := pc.Service.SelectOneByKey(key)
			if p.ParamValue != value {
				if key == "kubernetes_config" {
					k8sConfigModify = true
					k8sConfigValue = value
				} else if key == "harbor_address" {
					harborConfigModify = true
					harborAddress = value
				} else if key == "kubernetes_namespace" {
					namespaceModify = true
				} else if key == "kubernetes_master_address" {
					k8sConfigModify = true
					k8sMasterAddressValue = value
				}
			}
		}
	}
	harborAddressFailTag := 0
	//验证镜像地址是否可连接
	if harborConfigModify {
		if harborAddress == "" {
			return response.FailMsg("Failed to get harbor_address", "获取harbor地址失败")
		}
		url := "http://" + fmt.Sprintf(`%s/api/health`, harborAddress)
		returnData, err := utils.SimpleGet(url)
		if err != nil {
			harborAddressFailTag += 1
			pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "Failed to use harbor_address parameter value failed to connect to harbor. Please confirm whether the harbor configuration parameters filled in are correct"))
		} else {
			returnMap := make(map[string]interface{})
			err = json.Unmarshal(returnData, &returnMap)
			utils.LoggerError(err)
			if _, ok := returnMap["status"]; ok {
				if returnMap["status"] != "healthy" {
					harborAddressFailTag += 1
					pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "Failed to use harbor_address parameter value failed to connect to harbor. Please confirm whether the harbor configuration parameters filled in are correct"))
				}
			} else {
				harborAddressFailTag += 1
				pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "Failed to use harbor_address parameter value failed to connect to harbor. Please confirm whether the harbor configuration parameters filled in are correct"))
			}
		}
	}
	k8sConfigValueConnectTag := 0
	//验证k8s参数是否可连接
	if k8sConfigModify {
		if k8sConfigValue == "" {
			kubernetesConfigParam := pc.Service.SelectOneByKey("kubernetes_config")
			k8sConfigValue = kubernetesConfigParam.ParamValue
		}
		if k8sMasterAddressValue == "" {
			kubernetesMasterAddressParam := pc.Service.SelectOneByKey("kubernetes_master_address")
			k8sMasterAddressValue = kubernetesMasterAddressParam.ParamValue
		}
		if k8sConfigValue == "" && k8sMasterAddressValue == "" {
			return response.FailMsg("Failed to get kubernetes address and configuration parameters", "获取kubernetes地址、配置参数失败")
		}
		if k8sConfigValue != "" {
			k8sConfig, clientSet, ctx, err := service.InitK8s(k8sConfigValue)
			if err != nil {
				k8sConfigValueConnectTag += 1
				pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, err.Error()))
			} else {
				k8sNodeCh := make(chan interface{})
				k8sErrCh := make(chan interface{})
				go TestGetNode(clientSet, ctx, k8sNodeCh, k8sErrCh)
				select {
				case <-time.After(time.Second * 3):
					k8sConfigValueConnectTag += 1
					pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "Using k8s configuration parameters to connect k8s timeout, please confirm whether the k8s configuration parameters are correct"))
				case node := <-k8sNodeCh:
					k8sErr := <-k8sErrCh
					if k8sErr != nil {
						k8sConfigValueConnectTag += 1
						pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, k8sErr.(error).Error()))
					}
					if node == nil {
						k8sConfigValueConnectTag += 1
						pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "Failed to use k8s configuration parameters to connect k8s and query node information. Please confirm whether the k8s configuration parameters are correct"))
					}
				}
			}
			pc.CommonService.SetK8sConfig(k8sConfig, clientSet, ctx, err)
		}

		if k8sConfigValue != "" && k8sMasterAddressValue != "" {
			if strings.Contains(k8sConfigValue, k8sMasterAddressValue) == false {
				k8sConfigValueConnectTag += 1
				pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed: %s", params, "kubernetes_master_address and kubernetes_config does not match, please confirm whether the k8s configuration parameter is correct"))
			}
		}
		pc.CommonService.AsyncNodeInfo()
	}

	if harborAddressFailTag > 0 {
		if k8sConfigValueConnectTag > 0 {
			return response.FailMsg(
				"Using the value of harbor_address failed to connect to harbor, and the value of kubernetes configuration  failed to connect to kubernetes. Please confirm whether the parameters of harbor and kubernetes are correct",
				"使用harbor_address 参数值连接harbor失败，使用kubernetes配置参数连接kubernetes失败，请确认填写的harbor、kubernetes配置参数是否正确")
		} else {
			return response.FailMsg(
				"Using the value of harbor_address failed to connect to harbor. Please confirm whether the harbor configuration parameters filled in are correct",
				"使用失败harbor_address 参数值连接harbor失败，请确认填写的harbor配置参数是否正确")
		}
	}
	if k8sConfigValueConnectTag > 0 {
		return response.FailMsg(
			"Failed to connect kubernetes with kubernetes address and configuration parameters. Please confirm whether the filled kubernetes address and configuration parameters are correct",
			"使用kubernetes地址、配置参数连接kubernetes失败，请确认填写的kubernetes地址、配置参数是否正确")
	}
	isSuccess := pc.Service.ModifyParameter(paramsList)
	if !isSuccess {
		//  添加操作日志
		pc.CommonService.AddLog("error", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s failed", params))
		return response.Fail(response.ErrorUpdate)
	}
	pc.CommonService.AddLog("info", "system-systemparameter", opUser, fmt.Sprintf("update systemparameter %s successfully", params))
	if namespaceModify {
		pc.CommonService.SetNameSpace()
	}
	return response.Success(nil)
}

func TestGetNode(clientSet *kubernetes.Clientset, ctx *k8sContext.Context, k8sNodeCh chan interface{}, k8sErrCh chan interface{}) {
	k8sNode, k8sErr := clientSet.CoreV1().Nodes().List(*ctx, meta1.ListOptions{})
	k8sNodeCh <- k8sNode
	k8sErrCh <- k8sErr
}

//  还原默认值
func (pc *ParameterController) PostReset() mvc.Result {
	utils.LoggerInfo(" 还原默认值 ")
	strId := pc.Ctx.PostValue("id")
	operUsername := pc.Ctx.GetCookie("userName")
	idInt, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}
	isSuccess := pc.Service.ResetParameter(int(idInt))
	oldParam := pc.Service.SelectOne(int(idInt))
	if !isSuccess {
		//  添加操作日志
		pc.CommonService.AddLog("error", "system-systemparameter", operUsername, fmt.Sprintf("reset systemparameter %s failed", oldParam.ParamKey))
		//hc.OperLogService.SaveOptLog(fmt.Sprintf("update host %s failed ", host.Name), userName, "host")
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_UPDATE_EN,
				"error_msg_zh": utils.ERROR_UPDATE_ZH,
			},
		}
	}

	//  添加操作日志
	pc.CommonService.AddLog("info", "system-systemparameter", operUsername, fmt.Sprintf("reset systemparameter %s successfully", oldParam.ParamKey))
	//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno":      utils.RECODE_OK,
			"error_msg_en": "",
			"error_msg_zh": "",
		},
	}
}
