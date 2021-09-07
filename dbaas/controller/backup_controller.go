package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"fmt"
	"github.com/iris-contrib/schema"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type BackupController struct {
	Ctx     iris.Context
	Service service.BackupService
	Common  service.CommonService
}

func (bc *BackupController) addLog(level, content string) {
	bc.Common.AddLog(level, "system-backup", bc.Ctx.GetCookie("userName"), content)
}

func (bc *BackupController) GetStorageType() mvc.Result {
	list, err := bc.Service.StorageType()
	if err != nil {
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"detail": list})
}

func (bc *BackupController) GetStorageList() mvc.Result {
	page := bc.Ctx.URLParamIntDefault("page", 0)
	pageSize := bc.Ctx.URLParamIntDefault("pagesize", 0)
	userId, _ := bc.Ctx.URLParamInt("userId")
	search := bc.Ctx.URLParam("search")
	list, count, err := bc.Service.StorageList(userId, page, pageSize, search)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	var last int
	clusterId, _ := bc.Ctx.URLParamInt("clusterId")
	if clusterId > 0 {
		last, err = bc.Service.StorageLast(clusterId)
		if err != nil {
			utils.LoggerError(err)
			return response.Error(err)
		}
	}
	return response.Success(map[string]interface{}{
		"all":          count,
		"detail":       list,
		"page":         page,
		"pagesize":     pageSize,
		"lastSelected": last,
	})
}

func (bc *BackupController) PostStorageCreate() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	userIdStr := bc.Ctx.PostValue("userId")
	var storage = new(models.BackupStorage)
	values := bc.Ctx.FormValues()
	delete(values, "userId")
	err = schema.DecodeForm(values, storage)
	if err != nil {
		return response.Error(err)
	}
	err = bc.Service.StorageCreate(storage, userIdStr)
	if err != nil {
		bc.addLog("error", fmt.Sprintf("add backup storage %v error: %v", storage.Name, err))
		return response.Error(err)
	}
	bc.addLog("info", fmt.Sprintf("add backup storage %v successful", storage.Name))
	return response.Success(nil)
}

func (bc *BackupController) PostStorageDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	storageId, _ := bc.Ctx.PostValueInt("storageId")
	err = bc.Service.StorageDelete(storageId)
	if err != nil {
		bc.addLog("error", fmt.Sprintf("delete backup storage %v error: %v", storageId, err))
		return response.Error(err)
	}
	bc.addLog("info", fmt.Sprintf("delete backup storage %v successful", storageId))
	return response.Success(nil)
}

func (bc *BackupController) GetStorageReconnect() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	storageId, err := bc.Ctx.URLParamInt("storageId")
	if storageId <= 0 {
		return response.Fail(response.ErrorParameter)
	}
	err = bc.Service.StorageReconnect(storageId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (bc *BackupController) GetList() mvc.Result {
	page := bc.Ctx.URLParamIntDefault("page", 0)
	pageSize := bc.Ctx.URLParamIntDefault("pagesize", 0)
	storageId, _ := bc.Ctx.URLParamInt("storageId")
	clusterId, _ := bc.Ctx.URLParamInt("clusterId")
	startTime, _ := bc.Ctx.URLParamInt64("startTime")
	endTime, _ := bc.Ctx.URLParamInt64("endTime")
	search := bc.Ctx.URLParam("key")
	status := bc.Ctx.URLParam("status")
	t := bc.Ctx.URLParam("type")
	userTag := bc.Ctx.GetCookie("userTag")
	list, count, err := bc.Service.List(page, pageSize, storageId, clusterId, startTime, endTime, search, status, t, userTag)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{
		"all":      count,
		"detail":   list,
		"page":     page,
		"pagesize": pageSize,
	})
}

func (bc *BackupController) PostCreate() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	var task = new(models.BackupTask)
	err = bc.Ctx.ReadForm(task)
	if err != nil {
		return response.Error(err)
	}
	err = bc.Service.Create(task)
	if err != nil {
		bc.addLog("error", fmt.Sprintf("cluster %v create %v backup error: %v", task.ClusterId, task.Type, err))
		return response.Error(err)
	}
	bc.addLog("info", fmt.Sprintf("cluster %v create %v backup successful", task.ClusterId, task.Type))
	return response.Success(nil)
}

func (bc *BackupController) PostDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	jobId, err := bc.Ctx.PostValueInt("jobId")
	if jobId <= 0 {
		return response.Fail(response.ErrorParameter)
	}
	err = bc.Service.Delete(jobId)
	if err != nil {
		bc.addLog("error", fmt.Sprintf("delete backup job %v error: %v", jobId, err))
		return response.Error(err)
	}
	bc.addLog("info", fmt.Sprintf("delete backup job %v successful", jobId))
	return response.Success(nil)
}

func (bc *BackupController) PostCycleDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	clusterId, err := bc.Ctx.PostValueInt("clusterId")
	if err != nil {
		return response.Error(err)
	}
	err = bc.Service.DeleteCycle(clusterId)
	if err != nil {
		bc.addLog("error", fmt.Sprintf("delete cycle backup task of cluster %v error %v", clusterId, err))
		return response.Error(err)
	}
	bc.addLog("info", fmt.Sprintf("delete cycle backup task of cluster %v successful", clusterId))
	return response.Success(nil)
}

func (bc *BackupController) PostRecovery() mvc.Result {
	operatorName := bc.Ctx.GetCookie("userName")
	var err error
	jobId, _ := bc.Ctx.PostValueInt("jobId")
	defer func() {
		if err != nil {
			utils.LoggerError(err)
			bc.Common.AddLog("error", "system-backup", operatorName, fmt.Sprintf("recovery cluster from backup job %v error %s ", jobId, err))
		}
	}()
	clusterName := bc.Ctx.PostValue("clusterName")
	password := bc.Ctx.PostValue("password")
	storageStr := bc.Ctx.PostValue("storage")
	storageMap, err := utils.S2JMap(storageStr)
	if err != nil {
		return response.Error(err)
	}
	paramStr := bc.Ctx.PostValue("parameter")
	paramMap, err := utils.S2JMap2(paramStr)
	if err != nil {
		return response.Error(err)
	}
	remark := bc.Ctx.PostValue("remark")
	qos, err := ReadQos(bc.Ctx)
	if err != nil {
		return response.Error(err)
	}
	comboId := bc.Ctx.PostValueIntDefault("comboId", 0)
	nodePort := bc.Ctx.PostValueIntDefault("nodeport", 0)
	err = bc.Service.Recovery(clusterName, password, remark, storageMap, paramMap, jobId, qos, comboId, nodePort)
	if err != nil {
		return response.Error(err)
	}
	bc.Common.AddLog("info", "system-backup", operatorName, fmt.Sprintf("recovery cluster from backup job %v successful", jobId))
	return response.Success(nil)
}

func (bc *BackupController) GetEvent() mvc.Result {
	jobId, _ := bc.Ctx.URLParamInt("taskId")
	list, err := bc.Service.Event(jobId)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"detail": list})
}

func (bc *BackupController) GetLog() mvc.Result {
	jobId, _ := bc.Ctx.URLParamInt("taskId")
	data, err := bc.Service.Logs(jobId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"detail": data})
}

func (bc *BackupController) PostStorageUser() mvc.Result {
	userIdStr := bc.Ctx.PostValue("userIds")
	storageId, _ := bc.Ctx.PostValueInt("storageId")
	err := bc.Service.StorageUser(userIdStr, storageId)
	if err != nil {
		utils.LoggerError(err)
		bc.addLog("error", fmt.Sprintf("assign backup storage users (%v) error: %v", userIdStr, err))
		return response.Error(err)
	}
	bc.addLog("error", fmt.Sprintf("assign backup storage users (%v) successful", userIdStr))
	return response.Success(nil)
}
