/**
* @Description:
* @version:
* @Company: iwhalecloud
* @Author:  zhangwei
* @Date: 2020/11/16 10:30
* @LastEditors: zhangwei
* @LastEditTime: 2020/11/16 13:30
**/
package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"errors"
	"fmt"
	"strings"

	//"fmt"
	"encoding/json"
	"github.com/kataras/iris/v12"

	"github.com/kataras/iris/v12/mvc"
	"strconv"
	"time"
)

type UserController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//host功能实体
	Service        service.UserService
	StorageService service.StorageService
	CommonService  service.CommonService
}

//  获取用户列表
func (uc *UserController) GetList() mvc.Result {
	utils.LoggerInfo(" 获取用户列表 ")
	paramError := true
	pageSize, err := uc.Ctx.URLParamInt("pagesize")
	if err != nil {
		paramError = false
	}
	page, err := uc.Ctx.URLParamInt("page")
	if err != nil {
		paramError = false
	}
	key := uc.Ctx.URLParam("key")
	userList := make([]models.User, 0)
	if paramError {
		userList, err = uc.Service.ListUser(pageSize, (page-1)*pageSize, key)
	} else {
		userList, err = uc.Service.ListUserAll(key)
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
	userCount, err := uc.Service.GetUserCount(key)
	storageNum, err := uc.Service.GetAllStorageCount()
	clusterNum, err := uc.Service.GetClusterInstanceCount(0)
	userDetailMap := make(map[string]interface{})
	userMap := make(map[int]interface{})
	returnDataList := make([]interface{}, 0)
	userDetailData := make(map[string]interface{})
	var userDetailList []interface{}
	userDetailData["page"] = page
	userDetailData["pagesize"] = pageSize
	userDetailData["all"] = userCount
	userDetailData["clusterNum"] = clusterNum
	userDetailData["storageNum"] = storageNum

	if len(userList) > 0 {
		signaluserch := make(chan map[string]interface{}, len(userList)+1)
		for _, user := range userList {
			cpuUsed, memUsed, storageUsed, err := uc.Service.GetClusterCpuMemStorage(user.Id)
			utils.LoggerError(err)
			userClusterNum, err := uc.Service.GetClusterInstanceCount(user.Id)
			utils.LoggerError(err)
			userstorageNum, err := uc.Service.GetStorageCount(user.Id)
			utils.LoggerError(err)
			useBackup, err := uc.Service.GetUseBackup(user.Id)
			utils.LoggerError(err)
			go SignalUser(user, cpuUsed, memUsed, storageUsed, userClusterNum, userstorageNum, useBackup, signaluserch, uc)
		}
		for {
			time.Sleep(10 * time.Millisecond)
			userDetailMap = <-signaluserch
			userMap[userDetailMap["id"].(int)] = userDetailMap
			userDetailList = append(userDetailList, userDetailMap)
			if len(userDetailList) == len(userList) {
				break
			}
		}
		close(signaluserch)
		for _, signaluser := range userList {
			returnDataList = append(returnDataList, userMap[signaluser.Id])
		}
	}
	userDetailData["detail"] = returnDataList
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    userDetailData,
		},
	}
}

func SignalUser(user models.User, cpuUsed int64, memUsed int64, storageUsed int64, userClusterNum int64, userstorageNum int64, useBackup int, signaluserch chan map[string]interface{}, uc *UserController) {
	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	userBaseInformation := make(map[string]interface{})
	userBaseInformation["id"] = user.Id
	userBaseInformation["cpuTotal"] = user.CpuAll
	userBaseInformation["cpuUsed"] = cpuUsed
	userBaseInformation["memTotal"] = user.MemAll
	userBaseInformation["memUsed"] = memUsed
	userBaseInformation["storTotal"] = user.StorageAll
	userBaseInformation["storUsed"] = storageUsed
	userBaseInformation["clusterNum"] = userClusterNum
	userBaseInformation["storageNum"] = userstorageNum
	userBaseInformation["remark"] = user.Remarks
	userBaseInformation["copyUsed"] = useBackup
	userBaseInformation["copyTotal"] = user.BackupMax
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
				err := uc.Service.ModifyUser(updateUser, user.Id, false)
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
	signaluserch <- userBaseInformation
}

//  获取权限信息列表
func (hc *UserController) GetRoleList() mvc.Result {
	utils.LoggerInfo(" 获取权限信息列表 ")
	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	//user := hc.Ctx.GetCookie("userId")
	var roleList []map[string]interface{}
	roleInfos, ErrorMsgEn, ErrorMsgZh := capricornService.GetRoleResources("", "")
	if len(roleInfos) > 0 {
		roleList = roleInfos
	}
	if ErrorMsgEn != "" && ErrorMsgZh != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": ErrorMsgEn,
				"error_msg_zh": ErrorMsgZh,
			},
		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    roleList,
		},
	}
}

//  用户新增
func (uc *UserController) PostAdd() mvc.Result {
	utils.LoggerInfo(" 用户新增")
	username := uc.Ctx.PostValue("username")
	password := uc.Ctx.PostValue("password")
	cpu, err := uc.Ctx.PostValueInt("cpu")
	storage, err := uc.Ctx.PostValueInt("stor")
	//scType := uc.Ctx.PostValue("scType")
	//nodeNum, err := uc.Ctx.PostValueInt("nodeNum")
	utils.LoggerError(err)
	mem, err := uc.Ctx.PostValueInt("mem")
	utils.LoggerError(err)
	roleId := uc.Ctx.PostValue("roleId")
	remark := uc.Ctx.PostValue("remark")
	scListStr := uc.Ctx.PostValue("scList")
	backMax := uc.Ctx.PostValueIntDefault("copyTotal", 0)
	scList := make([]map[string]interface{}, 0)
	jsonerr := json.Unmarshal([]byte(scListStr), &scList)
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
	operUsername := uc.Ctx.GetCookie("userName")
	orgTag := uc.Ctx.GetCookie("organizationTag")
	userTag := uc.Ctx.GetCookie("userTag")
	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)

	sameNameUser, _ := uc.Service.SelectOneByName(username)
	if sameNameUser.UserName != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": "Already has a user with the same name",
				"error_msg_zh": "已拥有相同名称的用户",
			},
		}
	}
	sameUserInst, errorMsgEn, errorMsgZh := capricornService.GetUserResources("", username, "")
	if len(sameUserInst) > 0 {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": "Already has a user with the same name in capricorn",
				"error_msg_zh": "用户模块已拥有相同名称的用户",
			},
		}
	} else if errorMsgEn != "" && errorMsgZh != "" {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add user %s failed: %s", username, errorMsgEn+errorMsgZh))
		//uc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed: %s", name, errorMsgEn+errorMsgZh), userName, "host")
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errorMsgEn,
				"error_msg_zh": errorMsgZh,
			},
		}
	}

	if len(scList) > 0 {
		for _, scInfo := range scList {
			if scInfo["type"] == "custom" {
				sameNameSc, _ := uc.StorageService.SelectOneScByName(scInfo["scName"].(string))
				if sameNameSc.Name != "" {
					return mvc.Response{
						Object: map[string]interface{}{
							"errorno":      utils.RECODE_FAIL,
							"error_msg_en": fmt.Sprintf("Already has a SC with the same name: %s", scInfo["scName"].(string)),
							"error_msg_zh": fmt.Sprintf("已拥有相同名称: %s 的SC", scInfo["scName"].(string)),
						},
					}
				}
			}
			if _, ok := scInfo["pvList"]; ok {
				pvList := scInfo["pvList"].([]interface{})
				if len(pvList) > 0 {
					for _, pvInfo := range pvList {
						sameNamePv, _ := uc.StorageService.SelectOnePvByName(pvInfo.(map[string]interface{})["pvName"].(string))
						if sameNamePv.Name != "" {
							return mvc.Response{
								Object: map[string]interface{}{
									"errorno":      utils.RECODE_FAIL,
									"error_msg_en": fmt.Sprintf("Already has a PV with the same name: %s", pvInfo.(map[string]interface{})["pvName"].(string)),
									"error_msg_zh": fmt.Sprintf("已拥有相同名称: %s 的PV", pvInfo.(map[string]interface{})["pvName"].(string)),
								},
							}
						}
					}
				}
			}
		}
	}

	userInst, errorMsgEn, errorMsgZh := capricornService.AddUserResources(roleId, operUsername, username, password, "", "")
	if _, ok := userInst["user_id"]; ok {
		user := models.User{ZdcpId: int(userInst["user_id"].(float64)), MemAll: int64(mem), CpuAll: cpu, Remarks: remark, UserName: username, StorageAll: storage, UserTag: userInst["user_tag"].(string), BackupMax: backMax}
		isSuccess := uc.Service.SaveUser(&user)
		//hostInstIdstring = strconv.Itoa(int(hostInst["inst_id"].(float64)))
		if !isSuccess {
			//删除用户
			isSuccess, errorMsgEn, errorMsgZh := uc.Service.DeleteUserAndUserinst(user.Id, operUsername)
			if !isSuccess {
				uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete user %s failed", username))
				//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
				iris.New().Logger().Error(errorMsgEn + errorMsgZh)
			}
			//  添加操作日志
			uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add user %s failed: %s", username, errorMsgEn+errorMsgZh))
			//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed ", host.Hostname), userName, "host")
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_ADD_EN,
					"error_msg_zh": utils.ERROR_ADD_ZH,
				},
			}
		}
	} else if errorMsgEn != "" && errorMsgZh != "" {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add user %s failed: %s", username, errorMsgEn+errorMsgZh))
		//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed: %s", name, errorMsgEn+errorMsgZh), userName, "host")
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errorMsgEn,
				"error_msg_zh": errorMsgZh,
			},
		}
	} else {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add user %s failed: %s", username, errorMsgEn+errorMsgZh))
		//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed: %s", name, errorMsgEn+errorMsgZh), userName, "host")
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_ADD_EN,
				"error_msg_zh": utils.ERROR_ADD_ZH,
			},
		}
	}
	addUser, _ := uc.Service.SelectOneByName(username)
	if len(scList) > 0 {
		firstScInfo := scList[0]
		if firstScInfo["type"] == "ready" {
			issuccess, errMsg := uc.StorageService.AddscUserbyuser(addUser.Id, scList)
			if !issuccess {
				//  添加操作日志
				uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("user assign sc failed"))
				//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed ", host.Hostname), userName, "host")
				iris.New().Logger().Error(errMsg)
				//删除scuser
				deleRe, errM := uc.StorageService.DeletescUserbyUser(addUser.Id)
				if !deleRe {
					uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete scUser %s failed: %s", addUser.UserName, errM))
					//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
					iris.New().Logger().Error(errM)
				}
				//删除用户
				isSuccess, errorMsgEn, errorMsgZh := uc.Service.DeleteUserAndUserinst(addUser.Id, operUsername)
				if !isSuccess {
					uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete user %s failed", addUser.UserName))
					//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
					iris.New().Logger().Error(errorMsgEn + errorMsgZh)
				}
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_ADD_EN,
						"error_msg_zh": utils.ERROR_ADD_ZH,
					},
				}
			}
		} else if firstScInfo["type"] == "custom" {
			for _, scInfo := range scList {
				if scInfo["type"] == "custom" {
					remark := ""
					if _, ok := scInfo["remark"]; ok {
						remark = scInfo["remark"].(string)
					}
					nodeNum := 0
					if _, ok := scInfo["nodeNum"]; ok {
						nodeNum = int(scInfo["nodeNum"].(float64))
					}
					scName := scInfo["scName"].(string)
					reclaimPolicy := scInfo["reclaimPolicy"].(string)
					scType := scInfo["scType"].(string)
					_, err = uc.StorageService.Add(scName, reclaimPolicy, remark, orgTag, userTag, addUser.Id, scType, nodeNum, "")
					if err != nil {
						//  添加操作日志
						uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add sc %s failed", scInfo["scName"].(string)))
						//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed ", host.Hostname), userName, "host")
						//删除所有已建sc
						for _, deleteScInfo := range scList {
							deleteSc, _ := uc.StorageService.SelectOneScByName(deleteScInfo["scName"].(string))
							err = uc.StorageService.Delete(deleteSc.Id)
							if err != nil {
								utils.LoggerError(err)
								uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete sc %s failed", deleteScInfo["scName"].(string)))
							}
						}
						//删除scuser
						deleRe, errM := uc.StorageService.DeletescUserbyUser(addUser.Id)
						if !deleRe {
							uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete scUser %s failed: %s", addUser.UserName, errM))
							//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
							utils.LoggerError(errors.New(errM))
						}
						//删除用户
						isSuccess, errorMsgEn, _ := uc.Service.DeleteUserAndUserinst(addUser.Id, operUsername)
						if !isSuccess {
							uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete user %s failed", addUser.UserName))
							//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
							utils.LoggerError(errors.New(errorMsgEn))
						}
						return response.Error(err)
					}
					pvList := scInfo["pvList"].([]interface{})
					if len(pvList) > 0 {
						addSc, _ := uc.StorageService.SelectOneScByName(scName)
						for _, pvInfo := range pvList {
							pvName := pvInfo.(map[string]interface{})["pvName"].(string)
							mountPoint := pvInfo.(map[string]interface{})["mountPoint"].(string)
							iqn := pvInfo.(map[string]interface{})["iqn"].(string)
							lun := int(pvInfo.(map[string]interface{})["lun"].(float64))
							size := strconv.FormatFloat(pvInfo.(map[string]interface{})["size"].(float64), 'G', -1, 64)
							isSuccess, _, errMsg := uc.StorageService.PvAdd(addSc.Id, pvName, mountPoint, iqn, lun, size, userTag, orgTag, "default")
							if isSuccess == false {
								//  添加操作日志
								uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("add pv %s failed", pvInfo.(map[string]interface{})["pvName"].(string)))
								//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s failed ", host.Hostname), userName, "host")
								//删除所有已建pv
								for _, deletePvInfo := range pvList {
									deletePv, _ := uc.StorageService.SelectOnePvByName(deletePvInfo.(map[string]interface{})["pvName"].(string))
									err := uc.StorageService.PvDelete(deletePv.Id)
									if err != nil {
										uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete pv %s failed", deletePvInfo.(map[string]interface{})["pvName"].(string)))
										utils.LoggerError(err)
									}
								}
								//删除所有已建sc
								for _, deleteScInfo := range scList {
									deleteSc, _ := uc.StorageService.SelectOneScByName(deleteScInfo["scName"].(string))
									err = uc.StorageService.Delete(deleteSc.Id)
									if err != nil {
										utils.LoggerError(err)
										uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete sc %s failed", deleteScInfo["scName"].(string)))
									}
								}
								//删除用户
								isSuccess, errorMsgEn, errorMsgZh := uc.Service.DeleteUserAndUserinst(addUser.Id, operUsername)
								if !isSuccess {
									uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("delete user %s failed", addUser.UserName))
									//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
									iris.New().Logger().Error(errorMsgEn + errorMsgZh)
								}
								return mvc.Response{
									Object: map[string]interface{}{
										"errorno":      utils.RECODE_FAIL,
										"error_msg_en": errMsg,
										"error_msg_zh": errMsg,
									},
								}
							}
						}
					}
				}

			}
		}
	}
	uc.CommonService.AddLog("info", "system-user", operUsername, fmt.Sprintf("add user %s successfully", username))
	//hc.OperLogService.SaveOptLog(fmt.Sprintf("add host %s successfully", name), userName, "host")
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

// PostUpdate 修改用户信息
func (uc *UserController) PostUpdate() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	const notSet = "not_set"
	id, _ := uc.Ctx.PostValueInt("id")
	if id <= 0 {
		return response.Error(response.ErrorParameter)
	}

	scListStr := uc.Ctx.PostValueDefault("scList", notSet)
	var scList []map[string]interface{}
	if scListStr != notSet {
		scList = make([]map[string]interface{}, 0)
		err = json.Unmarshal(utils.Str2bytes(scListStr), &scList)
		if err != nil {
			err = fmt.Errorf("parse 'scList' to json arrary error: %v", err)
			return response.Error(err)
		}
	}

	operUsername := uc.Ctx.GetCookie("userName")
	oldUser, errM := uc.Service.SelectOne(id)
	if errM != "" {
		err = errors.New(errM)
		return response.Error(err)
	}
	oldUserName := oldUser.UserName
	if scList != nil {
		err = uc.StorageService.UserRegister(id, scList)
		if err != nil {
			uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("update user %s assign sc failed error %v", oldUserName, err))
			return response.Error(err)
		}
	}

	// 更新用户角色
	roleId := uc.Ctx.PostValueDefault("roleId", notSet)
	if roleId != notSet {
		userInstIdString := strconv.Itoa(oldUser.ZdcpId)
		capricornService, conn := service.NewCapricornService()
		defer service.CloseGrpc(conn)
		_, errorMsgEn, errorMsgZh := capricornService.UpdateUserResources(userInstIdString, roleId, operUsername, "", "", "", "")
		if errorMsgEn != "" || errorMsgZh != "" {
			uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("update user %s failed: %s", oldUserName, errorMsgEn))
			err = errors.New(errorMsgEn)
			return response.Error(err)
		}
	}

	cpu, _ := uc.Ctx.PostValueInt("cpu")
	mem, _ := uc.Ctx.PostValueInt("mem")
	storage, _ := uc.Ctx.PostValueInt("stor")
	backupMax, _ := uc.Ctx.PostValueInt("copyTotal")
	remark := uc.Ctx.PostValueDefault("remark", notSet)

	var user models.User
	var cpuUsed, memUsed, storageUsed int64
	if cpu != -1 || mem != -1 || storage != -1 {
		cpuUsed, memUsed, storageUsed, err = uc.Service.GetClusterCpuMemStorage(id)
		if err != nil {
			return response.Error(err)
		}
	}
	switch {
	case cpu != -1:
		if cpu < int(cpuUsed) {
			return response.FailMsg("The modified CPU quota is less than the used CPU quota", "修改的CPU配额小于已使用的CPU配额")
		}
		user.CpuAll = cpu
		fallthrough
	case mem != -1:
		if mem < int(memUsed) {
			return response.FailMsg("The modified MEM quota is less than the used MEM quota", "修改的MEM配额小于已使用的MEM配额")
		}
		user.MemAll = int64(mem)
		fallthrough
	case storage != -1:
		if storage < int(storageUsed) {
			return response.FailMsg("The modified Storage quota is less than the used Storage quota", "修改的Storage配额小于已使用的Storage配额")
		}
		user.StorageAll = storage
		fallthrough
	case remark != notSet:
		user.Remarks = remark
		fallthrough
	case backupMax != -1:
		user.BackupMax = backupMax
	}

	err = uc.Service.ModifyUser(user, id, user.Remarks != notSet)
	if err != nil {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("update user %s failed: %v", oldUserName, err))
		return response.Error(err)
	}
	uc.CommonService.AddLog("info", "system-user", operUsername, fmt.Sprintf("update user %s successfully", oldUserName))
	return response.Success(nil)
}

//  用户删除
func (uc *UserController) PostDelete() mvc.Result {
	userId, err := uc.Ctx.PostValueInt("id")
	if err != nil {
		utils.LoggerError(err)
		return response.Fail(response.ErrorParameter)
	}
	opUser := uc.Ctx.GetCookie("userName")
	clusterList, err := uc.Service.GetClustersByUser(userId)
	utils.LoggerError(err)
	if len(clusterList) > 0 {
		cls := make([]string, 0)
		for _, cl := range clusterList {
			cls = append(cls, cl.Name)
		}
		clsStr := strings.Join(cls, ",")
		return response.FailMsg(
			fmt.Sprintf("The current user is associated with a database instance: %s. Please delete the corresponding database instance first", clsStr),
			fmt.Sprintf("当前用户下关联有数据库实例: %s，请先删除相应的数据库实例", clsStr))
	}
	oldUser, errM := uc.Service.SelectOne(userId)
	if errM != "" {
		return response.Error(errors.New(errM))
	}
	oldUsername := oldUser.UserName
	isSuccess, errorMsgEn, errorMsgZh := uc.Service.DeleteUserAndUserinst(userId, opUser)
	if isSuccess {
		//  添加操作日志
		uc.CommonService.AddLog("info", "system-user", opUser, fmt.Sprintf("delete user %s successfully", oldUsername))
		return response.Success(nil)
	} else {
		uc.CommonService.AddLog("error", "system-user", opUser, fmt.Sprintf("delete user %s failed", oldUsername))
		return response.FailMsg(errorMsgEn, errorMsgZh)
	}
}

// 用户启用禁用
func (uc *UserController) PatchOperate() mvc.Result {
	iris.New().Logger().Info(" 用户启用禁用 ")
	strId := uc.Ctx.URLParam("id")
	operUsername := uc.Ctx.GetCookie("userName")
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

	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	oldUser, errM := uc.Service.SelectOne(int(idInt))
	if errM != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errM,
				"error_msg_zh": errM,
			},
		}
	}
	userInstIdString := strconv.Itoa(oldUser.ZdcpId)
	_, errorMsgEn, errorMsgZh := capricornService.OperateUserResources(userInstIdString, operUsername)
	if errorMsgEn != "" && errorMsgZh != "" {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("operate user %s failed: %s", oldUser.UserName, errorMsgEn))
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errorMsgEn,
				"error_msg_zh": errorMsgZh,
			},
		}
	}
	//  添加操作日志
	uc.CommonService.AddLog("info", "system-user", operUsername, fmt.Sprintf("operate user %s successfully", oldUser.UserName))
	//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

//  重置密码
func (uc *UserController) PostPwdkeyReset() mvc.Result {
	utils.LoggerInfo(" 重置密码 ")
	strId := uc.Ctx.PostValue("id")
	operUsername := uc.Ctx.GetCookie("userName")
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

	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	oldUser, errM := uc.Service.SelectOne(int(idInt))
	if errM != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}
	userInstIdString := strconv.Itoa(oldUser.ZdcpId)
	returnDataMap := make(map[string]interface{})
	returnData, errorMsgEn, errorMsgZh := capricornService.GetRandomPassword(userInstIdString, operUsername)
	if _, ok := returnData["enPassword"]; ok {
		returnDataMap["password"] = returnData["enPassword"]
	}
	if errorMsgEn != "" && errorMsgZh != "" {
		uc.CommonService.AddLog("error", "system-user", operUsername, fmt.Sprintf("reset user`s: %s password failed: %s", oldUser.UserName, errorMsgEn))
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errorMsgEn,
				"error_msg_zh": errorMsgZh,
			},
		}
	}
	//  添加操作日志
	uc.CommonService.AddLog("info", "system-user", operUsername, fmt.Sprintf("reset user`s: %s password successfully", oldUser.UserName))
	//hc.OperLogService.SaveOptLog(fmt.Sprintf("delete host  %v  successful ", hostname), userName, "host")
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno":      utils.RECODE_OK,
			"error_msg_en": "",
			"error_msg_zh": "",
			"data":         returnDataMap,
		},
	}
}
