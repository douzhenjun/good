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
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type ClusterController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	Service service.ClusterService

	CommonService service.CommonService

	UserService service.UserService
}

func (cc *ClusterController) addLog(level, content string) {
	cc.CommonService.AddLog(level, "system-cluster", cc.Ctx.GetCookie("userName"), content)
}

func (cc *ClusterController) PostDisable() mvc.Result {
	clusterId, _ := cc.Ctx.PostValueInt("clusterId")
	err := cc.Service.ClusterDisable(clusterId)
	if err != nil {
		utils.LoggerError(err)
		cc.addLog("error", fmt.Sprintf("disable cluster %v error %v", clusterId, err))
		return response.Error(err)
	}
	cc.addLog("info", fmt.Sprintf("disable cluster %v successful", clusterId))
	return response.Success(nil)
}

func (cc *ClusterController) PostEnable() mvc.Result {
	clusterId, _ := cc.Ctx.PostValueInt("clusterId")
	err := cc.Service.ClusterEnable(clusterId)
	if err != nil {
		utils.LoggerError(err)
		cc.addLog("error", fmt.Sprintf("enable cluster %v error %v", clusterId, err))
		return response.Error(err)
	}
	cc.addLog("info", fmt.Sprintf("enable cluster %v successful", clusterId))
	return response.Success(nil)
}

func (cc *ClusterController) PostUpdate() mvc.Result {
	id, _ := cc.Ctx.PostValueInt("id")
	var err error
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			cc.addLog("error", fmt.Sprintf("update cluster %v error %v", id, err))
		}
	}()
	dataStr := cc.Ctx.PostValue("data")
	if len(dataStr) == 0 {
		err = errors.New("'data' field cannot be empty")
		return response.Error(err)
	}
	dataM := map[string]interface{}{}
	err = json.Unmarshal(utils.Str2bytes(dataStr), &dataM)
	if err != nil {
		err = fmt.Errorf("parse data to map error: %v", err)
		return response.Error(err)
	}
	err = cc.Service.Update(id, dataM)
	if err != nil {
		return response.Error(err)
	}
	cc.addLog("info", fmt.Sprintf("update cluster %v successful", id))
	return response.Success(nil)
}

func (cc *ClusterController) PostDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	id, err := cc.Ctx.PostValueInt("id")
	userName := cc.Ctx.GetCookie("userName")
	if err != nil {
		cc.CommonService.AddLog("error", "system-cluster", userName, fmt.Sprintf("delete cluster %v error %s ", id, response.ErrorParameter.En))
		return response.Fail(response.ErrorParameter)
	}
	keepPV, _ := cc.Ctx.PostValueBool("keepPv")
	err = cc.Service.Delete(id, keepPV)
	if err != nil {
		cc.CommonService.AddLog("error", "system-cluster", userName, fmt.Sprintf("delete cluster %v error %s ", id, err))
		return response.Error(err)
	}
	cc.CommonService.AddLog("info", "system-cluster", userName, fmt.Sprintf("delete cluster %v successful ", id))
	return response.Success(nil)
}

func (cc *ClusterController) PostAdd() mvc.Result {
	userName := cc.Ctx.GetCookie("userName")
	clusterName := cc.Ctx.PostValue("clusterName")
	var err error
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			cc.CommonService.AddLog("error", "system-cluster", userName, fmt.Sprintf("add cluster %v error %s ", clusterName, err))
		}
	}()
	orgTag := cc.Ctx.GetCookie("orgTag")
	password := cc.Ctx.PostValue("password")
	imageId, err := cc.Ctx.PostValueInt("imageId")
	if err != nil {
		err = fmt.Errorf("parse imageId to int error: %v", err)
		return response.Error(err)
	}
	userId, err := cc.Ctx.PostValueInt("userId")
	if err != nil {
		userTag := cc.Ctx.GetCookie("userTag")
		userId, err = cc.UserService.SelectIdByTag(userTag)
		if err != nil {
			return response.Error(err)
		}
	}
	storageMap := map[string]interface{}{}
	storage := cc.Ctx.PostValue("storage")
	err = json.Unmarshal([]byte(storage), &storageMap)
	if err != nil {
		err = fmt.Errorf("parse storage to json error: %v", err)
		return response.Error(err)
	}
	parameterMap := make([]map[string]interface{}, 0)
	parameter := cc.Ctx.PostValue("parameter")
	err = json.Unmarshal([]byte(parameter), &parameterMap)
	if err != nil {
		err = fmt.Errorf("parse parameter to json error: %v", err)
		return response.Error(err)
	}
	qos, err := ReadQos(cc.Ctx)
	if err != nil {
		return response.Error(err)
	}
	remark := cc.Ctx.PostValue("remark")
	comboId := cc.Ctx.PostValueIntDefault("comboId", 0)
	nodePort := cc.Ctx.PostValueIntDefault("nodeport", 0)
	_, err = cc.Service.Add(clusterName, password, storageMap, parameterMap, remark, userId, imageId, orgTag, "internal", qos, comboId, nodePort)
	if err != nil {
		return response.Error(err)
	}
	cc.CommonService.AddLog("info", "system-cluster", userName, fmt.Sprintf("add cluster %v successful ", clusterName))
	return response.Success(nil)
}

func (cc *ClusterController) PostOperate() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	userName := cc.Ctx.GetCookie("userName")
	replicas := cc.Ctx.PostValueIntDefault("replicas", 0)
	id, err := cc.Ctx.PostValueInt("id")
	if err != nil {
		cc.CommonService.AddLog("error", "system-cluster", userName, fmt.Sprintf("Operate cluster %v error %s ", id, response.ErrorParameter.En))
		return response.Fail(response.ErrorParameter)
	}
	err = cc.Service.Patch(replicas, id)
	if err != nil {
		cc.CommonService.AddLog("error", "system-cluster", userName, fmt.Sprintf("Operate cluster %v error %s ", id, err.Error()))
		return response.Error(err)
	}
	cc.CommonService.AddLog("info", "system-cluster", userName, fmt.Sprintf("Operate cluster %v successful", id))
	return response.Success(nil)
}

func (cc *ClusterController) GetList() mvc.Result {
	userTag := cc.Ctx.GetCookie("userTag")
	userId := cc.Ctx.URLParamIntDefault("id", 0)
	page := cc.Ctx.URLParamIntDefault("page", 0)
	pageSize := cc.Ctx.URLParamIntDefault("pagesize", 0)
	key := cc.Ctx.URLParam("key")
	isDeploy, _ := cc.Ctx.URLParamBool("isDeploy")
	list, count, err := cc.Service.List(page, pageSize, key, userId, userTag, isDeploy)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{
		"detail":   list,
		"all":      count,
		"page":     page,
		"pagesize": pageSize,
	})
}

func (cc *ClusterController) GetParamList() mvc.Result {
	clusterId, _ := cc.Ctx.URLParamInt("clusterId")
	page, pageSize := cc.Ctx.URLParamIntDefault("page", 0), cc.Ctx.URLParamIntDefault("pagesize", 0)
	list, count, err := cc.Service.ParamList(clusterId, page, pageSize)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{
		"all":      count,
		"page":     page,
		"pagesize": pageSize,
		"detail":   list,
	})
}

func (cc *ClusterController) PostParameterEdit() mvc.Result {
	opUser := cc.Ctx.GetCookie("userName")
	clusterId, _ := cc.Ctx.PostValueInt("clusterId")
	paramList := make([]map[string]interface{}, 0)
	paramStr := cc.Ctx.PostValue("paramList")
	err := json.Unmarshal(utils.Str2bytes(paramStr), &paramList)
	if clusterId <= 0 || err != nil {
		return response.Fail(response.ErrorParameter)
	}
	if len(paramList) == 0 {
		return response.Success(nil)
	}
	err = cc.Service.ParameterEdit(clusterId, paramList)
	if err != nil {
		utils.LoggerError(err)
		cc.CommonService.AddLog("error", "system-cluster", opUser, fmt.Sprintf("cluster %v edit parameter '%v' error: %v", clusterId, paramStr, err))
		return response.Error(err)
	}
	cc.CommonService.AddLog("info", "system-cluster", opUser, fmt.Sprintf("cluster %v edit parameter '%v' successful", clusterId, paramStr))
	return response.Success(nil)
}

func (cc *ClusterController) GetPodList() mvc.Result {
	clusterId, _ := cc.Ctx.URLParamInt("id")
	cluster, err := cc.Service.PodDetail(clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(cluster)
}

func (cc *ClusterController) GetUserinfo() mvc.Result {
	userName := cc.Ctx.GetCookie("userName")
	clusterId := cc.Ctx.URLParam("clusterId")
	userBaseInformation := make(map[string]interface{})
	returnData := make(map[string]interface{})
	returnData["detail"] = userBaseInformation
	if userName == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_OK,
				"data":         returnData,
				"error_msg_en": "",
				"error_msg_zh": "",
			},
		}
	}
	user, result := cc.UserService.SelectOneByName(userName)
	if result == false {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"data":         returnData,
				"error_msg_en": "Failed to get user information",
				"error_msg_zh": "获取用户信息失败",
			},
		}
	}
	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	cpuUsed, memUsed, storageUsed, err := cc.UserService.GetClusterCpuMemStorage(user.Id)
	utils.LoggerError(err)
	if clusterId != "" {
		clusterIdInt, err := strconv.ParseInt(clusterId, 10, 64)
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
		if clusterIdInt > 0 {
			cluster, err := cc.Service.SelectOne(int(clusterIdInt))
			if err != "" {
				utils.LoggerInfo(err)
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"data":         returnData,
						"error_msg_en": err,
						"error_msg_zh": err,
					},
				}
			}
			cpuUsed = int64(int(cpuUsed) - cluster.LimitCpu)
			memUsed = int64(int(memUsed) - cluster.LimitMem)
			storageUsed = int64(int(storageUsed) - cluster.Storage)
		}
	}

	userBaseInformation["id"] = user.Id
	userBaseInformation["cpuTotal"] = user.CpuAll
	userBaseInformation["cpuUsed"] = cpuUsed
	userBaseInformation["memTotal"] = user.MemAll
	userBaseInformation["memUsed"] = memUsed
	userBaseInformation["storTotal"] = user.StorageAll
	userBaseInformation["storUsed"] = storageUsed
	userBaseInformation["remark"] = user.Remarks
	userIdString := strconv.Itoa(user.ZdcpId)

	hostInfo := make(map[string]interface{}, 0)
	userList, ErrorMsgEn, ErrorMsgZh := capricornService.GetUserResources(userIdString, "", "")

	if ErrorMsgEn != "" && ErrorMsgZh != "" {
		iris.New().Logger().Error(ErrorMsgEn)
	}
	if len(userList) > 0 {
		hostInfo = userList[0]
		if _, ok := hostInfo["username"]; ok {
			userBaseInformation["username"] = hostInfo["username"]
			if user.UserName != hostInfo["username"] {
				updateUser := models.User{UserName: hostInfo["username"].(string)}
				err = cc.UserService.ModifyUser(updateUser, user.Id, false)
				if err != nil {
					utils.LoggerInfo("同步capricorn模块用户信息")
				}
			}

		}
		if _, ok := hostInfo["status"]; ok {
			if hostInfo["status"] == "inactive" {
				userBaseInformation["status"] = "disable"
			} else {
				userBaseInformation["status"] = "enable"
			}
		}
		if _, ok := hostInfo["roleList"]; ok {
			userBaseInformation["roleList"] = hostInfo["roleList"].([]interface{})
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"errorno":      utils.RECODE_OK,
			"data":         returnData,
			"error_msg_en": "",
			"error_msg_zh": "",
		},
	}
}

func (cc *ClusterController) PostConfigApply() mvc.Result {
	clusterId, _ := cc.Ctx.PostValueInt("clusterId")
	err := cc.Service.ApplyConfig(clusterId)
	if err != nil {
		utils.LoggerError(err)
		cc.addLog("error", fmt.Sprintf("apply cluster (id: %v) configs error: %v", clusterId, err))
		return response.Error(err)
	}
	cc.addLog("info", fmt.Sprintf("apply cluster (id: %v) configs successful", clusterId))
	return response.Success(nil)
}

func (cc *ClusterController) PostNodeportCheck() mvc.Result {
	port, _ := cc.Ctx.PostValueInt("nodeport")
	port, err := cc.Service.NodePort(port)
	result := map[string]interface{}{}
	if err != nil {
		msg := response.IsMsg(err)
		if msg == nil {
			return response.Error(err)
		}
		result["error_msg_en"] = msg.En
		result["error_msg_zh"] = msg.Zh
	}
	result["nodeport"] = port
	return response.Success(result)
}
