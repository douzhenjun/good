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
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strconv"
)

type StorageController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	Service service.StorageService

	CommonService service.CommonService
}

func (sc *StorageController) addLog(level, content string) {
	sc.CommonService.AddLog(level, "system-storage", sc.Ctx.GetCookie("userName"), content)
}

func (sc *StorageController) PostUserAssign() mvc.Result {
	var err error
	id, err := sc.Ctx.PostValueInt("id")
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			sc.addLog("error", fmt.Sprintf("storage %v user assign error: %v", id, err))
		}
	}()
	if err != nil {
		return response.Error(response.ErrorParameter)
	}
	userId := sc.Ctx.PostValue("userId")
	err = sc.Service.UserAssign(id, userId)
	if err != nil {
		return response.Error(err)
	}
	sc.addLog("info", fmt.Sprintf("storage %v user assign successful", id))
	return response.Success(nil)
}

func (sc *StorageController) GetList() mvc.Result {
	userId, _ := sc.Ctx.URLParamInt("id")
	userTag := sc.Ctx.GetCookie("userTag")
	filter, _ := sc.Ctx.URLParamBool("filter")
	page := sc.Ctx.URLParamIntDefault("page", 0)
	pageSize := sc.Ctx.URLParamIntDefault("pagesize", 0)
	key := sc.Ctx.URLParam("key")
	list, count := sc.Service.List(page, pageSize, key, userId, userTag, filter)
	return response.Success(map[string]interface{}{
		"detail":   list,
		"all":      count,
		"page":     page,
		"pagesize": pageSize,
	})
}

func (sc *StorageController) PostAdd() mvc.Result {
	scName := sc.Ctx.PostValue("scName")
	var err error
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			sc.addLog("error", fmt.Sprintf("add sc '%v' error: %v", scName, err))
		}
	}()
	userIdString := sc.Ctx.PostValue("userId")
	remark := sc.Ctx.PostValue("remark")
	scType := sc.Ctx.PostValue("scType")
	nodeNum := sc.Ctx.PostValueIntDefault("nodeNum", 0)
	reclaimPolicy := sc.Ctx.PostValue("reclaimPolicy")
	orgTag := sc.Ctx.GetCookie("orgTag")
	userTag := sc.Ctx.GetCookie("userTag")
	_, err = sc.Service.Add(scName, reclaimPolicy, remark, orgTag, userTag, 0, scType, nodeNum, userIdString)
	if err != nil {
		return response.Error(err)
	}

	sc.addLog("info", fmt.Sprintf("add sc %v successful", scName))
	return response.Success(nil)
}

func (sc *StorageController) PostUpdate() mvc.Result {
	id, err := sc.Ctx.PostValueInt("id")
	if err != nil {
		return response.Error(response.ErrorParameter)
	}

	remark := sc.Ctx.PostValue("remark")
	nodeNum, err := sc.Ctx.PostValueInt("nodeNum")
	err = sc.Service.Update(id, remark, nodeNum)
	if err != nil {
		utils.LoggerError(err)
		sc.addLog("error", fmt.Sprintf("update sc %v error: %v", id, err))
		return response.Error(err)
	}
	sc.addLog("info", fmt.Sprintf("update sc %v successful", id))
	return response.Success(nil)
}

func (sc *StorageController) PostDelete() mvc.Result {
	id, err := sc.Ctx.PostValueInt("id")
	if err != nil {
		return response.Fail(response.ErrorParameter)
	}

	err = sc.Service.Delete(id)
	if err != nil {
		utils.LoggerError(err)
		sc.addLog("error", fmt.Sprintf("dalete sc %v error: %v", id, err))
		return response.Error(err)

	}
	sc.addLog("info", fmt.Sprintf("delete sc %v successful", id))
	return response.Success(nil)
}

// pv 新增
func (sc *StorageController) PostPvAdd() mvc.Result {
	storageId, err := sc.Ctx.PostValueInt("storageId")
	userName := sc.Ctx.GetCookie("userName")
	pvName := sc.Ctx.PostValue("pvName")
	if err != nil {
		utils.LoggerError(err)
		sc.CommonService.AddLog("error", "system-storage", userName, fmt.Sprintf("add pv %v error: %v", pvName, utils.ERROR_PARAMETER_EN))
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}

	mountPoint := sc.Ctx.PostValue("mountPoint")

	iqn := sc.Ctx.PostValue("iqn")
	lun, err := sc.Ctx.PostValueInt("lun")
	utils.LoggerError(err)
	size := sc.Ctx.PostValue("size")
	orgTag := sc.Ctx.GetCookie("orgTag")
	userTag := sc.Ctx.GetCookie("userTag")

	isSuccess, _, errMsg := sc.Service.PvAdd(storageId, pvName, mountPoint, iqn, lun, size, userTag, orgTag, "default")
	if isSuccess == false {
		sc.CommonService.AddLog("error", "system-storage", userName, fmt.Sprintf("add pv %v error: %v", pvName, errMsg))
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errMsg,
				"error_msg_zh": errMsg,
			},
		}
	}
	sc.CommonService.AddLog("info", "system-storage", userName, fmt.Sprintf("add pv %v successful", pvName))
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

// pv 删除
func (sc *StorageController) PostPvDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	id, err := sc.Ctx.PostValueInt("id")
	userName := sc.Ctx.GetCookie("userName")
	if err != nil {
		sc.CommonService.AddLog("error", "system-storage", userName, fmt.Sprintf("delete pv %v error: %v", id, response.ErrorParameter.En))
		return response.Fail(response.ErrorParameter)
	}

	err = sc.Service.PvDelete(id)
	if err != nil {
		sc.CommonService.AddLog("error", "system-storage", userName, fmt.Sprintf("add pv %v error: %s", id, err))
		return response.Error(err)
	}
	sc.CommonService.AddLog("info", "system-storage", userName, fmt.Sprintf("delete pv %v successful", id))
	return response.Success(nil)
}

// 初始化创建SC、PV
func (sc *StorageController) PostInitAdd() mvc.Result {
	userId := sc.Ctx.GetCookie("userId")
	operUsername := sc.Ctx.GetCookie("userName")
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}

	scName := sc.Ctx.PostValue("scName")
	reclaimPolicy := "Retain"
	remark := ""
	orgTag := sc.Ctx.GetCookie("orgTag")
	userTag := sc.Ctx.GetCookie("userTag")
	pvListStr := sc.Ctx.PostValue("pvList")
	scType := sc.Ctx.PostValue("scType")
	nodeNum, err := sc.Ctx.PostValueInt("nodeNum")
	pvList := make([]map[string]interface{}, 0)
	jsonerr := json.Unmarshal([]byte(pvListStr), &pvList)
	if jsonerr != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_DATA_TRANSFER_EN,
				"error_msg_zh": utils.ERROR_DATA_TRANSFER_ZH,
			},
		}
	}
	sameNameSc, _ := sc.Service.SelectOneScByName(scName)
	if sameNameSc.Name != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": fmt.Sprintf("Already has a SC with the same name: %s", scName),
				"error_msg_zh": fmt.Sprintf("已拥有相同名称: %s 的SC", scName),
			},
		}
	}
	if len(pvList) > 0 {
		for _, pvInfo := range pvList {
			sameNamePv, _ := sc.Service.SelectOnePvByName(pvInfo["pvName"].(string))
			if sameNamePv.Name != "" {
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": fmt.Sprintf("Already has a PV with the same name: %s", pvInfo["pvName"].(string)),
						"error_msg_zh": fmt.Sprintf("已拥有相同名称: %s 的PV", pvInfo["pvName"].(string)),
					},
				}
			}
		}
	}
	addSc, err := sc.Service.Add(scName, reclaimPolicy, remark, orgTag, userTag, int(userIdInt), scType, nodeNum, "")
	if err != nil {
		sc.CommonService.AddLog("error", "system-initadd", operUsername, fmt.Sprintf("delete sc %s failed", addSc.Name))
		return response.Error(err)
	}
	if len(pvList) > 0 {
		for _, pvInfo := range pvList {
			pvName := pvInfo["pvName"].(string)
			mountPoint := pvInfo["mountPoint"].(string)
			iqn := pvInfo["iqn"].(string)
			lun := int(pvInfo["lun"].(float64))
			size := strconv.FormatFloat(pvInfo["size"].(float64), 'G', -1, 64)
			isSuccess, _, errMsg := sc.Service.PvAdd(addSc.Id, pvName, mountPoint, iqn, lun, size, userTag, orgTag, "default")
			if isSuccess == false {
				//  添加操作日志
				sc.CommonService.AddLog("error", "system-initadd", operUsername, fmt.Sprintf("add pv %s failed", pvName))
				//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed ", host.Hostname), userName, "host")
				//删除所有已建pv
				for _, deletePvInfo := range pvList {
					deletePv, _ := sc.Service.SelectOnePvByName(deletePvInfo["pvName"].(string))
					err := sc.Service.PvDelete(deletePv.Id)
					if err != nil {
						sc.CommonService.AddLog("error", "system-initadd", operUsername, fmt.Sprintf("delete pv %s failed", deletePvInfo["pvName"].(string)))
						utils.LoggerError(err)
					}
				}
				//删除已建sc
				err = sc.Service.Delete(addSc.Id)
				if err != nil {
					utils.LoggerError(err)
					sc.CommonService.AddLog("error", "system-initadd", operUsername, fmt.Sprintf("delete sc %s failed", addSc.Name))
				}
				return response.Error(errors.New(errMsg))
			}
		}
	}

	sc.CommonService.AddLog("info", "system-initadd", operUsername, fmt.Sprintf("add sc %s successfully", addSc.Name))
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

func (sc *StorageController) GetPvList() mvc.Result {
	userTag := sc.Ctx.GetCookie("userTag")
	if len(userTag) != 4 {
		return response.Error(fmt.Errorf("UserTag is invalid: %s ", userTag))
	}
	page := sc.Ctx.URLParamIntDefault("page", 0)
	pageSize := sc.Ctx.URLParamIntDefault("pagesize", 0)
	key := sc.Ctx.URLParam("key")
	list, count, err := sc.Service.PVList(page, pageSize, key, userTag)
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

func (sc *StorageController) PostRestorePv() mvc.Result {
	var err error
	var pvId int
	userName := sc.Ctx.GetCookie("userName")
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			sc.CommonService.AddLog("error", "system-storage", userName, fmt.Sprintf("restore pv %v error: %s ", pvId, err))
		}
	}()

	pvId, err = sc.Ctx.PostValueInt("pvId")
	if err != nil {
		err = fmt.Errorf("parse 'pvId' error: %v", err)
		return response.Error(err)
	}
	userId, err := sc.Ctx.PostValueInt("userId")
	if err != nil {
		err = fmt.Errorf("parse 'userId' error: %v", err)
		return response.Error(err)
	}
	storageMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(sc.Ctx.PostValue("storage")), &storageMap)
	if err != nil {
		err = fmt.Errorf("parse 'storage' to json error: %v", err)
		return response.Error(err)
	}
	remark := sc.Ctx.PostValue("remark")
	mysqlName := sc.Ctx.PostValue("name")
	if len(mysqlName) == 0 {
		err = errors.New("mysql name cannot be empty")
		return response.Error(err)
	}
	qos, err := ReadQos(sc.Ctx)
	if err != nil {
		return response.Error(err)
	}
	err = sc.Service.CreateMysqlByPV(pvId, storageMap, remark, userId, mysqlName, qos)
	if err != nil {
		return response.Error(err)
	}
	sc.CommonService.AddLog("info", "system-storage", userName, fmt.Sprintf("restore pv %v successful", pvId))
	return response.Success(nil)
}
