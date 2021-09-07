package controller

import (
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ExternalController struct {
	Ctx           iris.Context
	Service       service.ExternalService
	CommonService service.CommonService
}

func (ec *ExternalController) PostClusterAdd() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	err = ec.Service.VerifyVersion(ec.Ctx)
	if err != nil {
		return response.Error(err)
	}
	param := make(map[string]interface{})
	_ = ec.Ctx.ReadJSON(&param)
	userId := utils.ReadMapString(param, "userId")
	userId, err = ec.Service.VerifyStamp(userId)
	if err != nil {
		return response.Error(err)
	}
	clusterName := utils.ReadMapString(param, "clusterName")
	imageVersion := utils.ReadMapString(param, "imageVersion")
	clusterType := utils.ReadMapString(param, "clusterType")
	storageType := utils.ReadMapString(param, "storageType")
	password := utils.ReadMapString(param, "password")
	storageMap, ok := param["storage"].(map[string]interface{})
	if !ok || !utils.MustMap(storageMap, "mem", "cpu", "size", "copy") {
		err = errors.New("storage map is incomplete")
		return response.Error(err)
	}
	err = ec.Service.CheckQuota(storageMap, "/external/cluster/add")
	if err != nil {
		return response.Error(err)
	}
	clusterId, err := ec.Service.OpenCluster(userId, clusterName, imageVersion, clusterType, storageType, password, storageMap)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"clusterId": clusterId})
}

func (ec *ExternalController) GetClusterSelect() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	err = ec.Service.VerifyVersion(ec.Ctx)
	if err != nil {
		return response.Error(err)
	}
	clusterId, _ := ec.Ctx.URLParamInt("clusterId")
	cluster, err := ec.Service.SelectCluster(clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(cluster)
}

func (ec *ExternalController) PostClusterDelete() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	err = ec.Service.VerifyVersion(ec.Ctx)
	if err != nil {
		return response.Error(err)
	}
	param := make(map[string]interface{})
	_ = ec.Ctx.ReadJSON(&param)
	userId := utils.ReadMapString(param, "userId")
	_, err = ec.Service.VerifyStamp(userId)
	if err != nil {
		return response.Error(err)
	}
	clusterId := utils.ReadMapInt(param, "clusterId")
	if clusterId == 0 {
		err = errors.New("cluster ID parameter error")
		return response.Error(err)
	}
	err = ec.Service.DeleteCluster(clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (ec *ExternalController) GetLogin() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	userId := ec.Ctx.URLParam("userId")
	userId, err = ec.Service.VerifyStamp(userId)
	if err != nil {
		return response.Error(err)
	}
	clusterId, _ := ec.Ctx.URLParamInt("clusterId")
	user, err := ec.Service.Login(userId, clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(map[string]string{"username": user.UserName, "password": user.Password})
}

func (ec *ExternalController) PostClusterDisable() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	err = ec.Service.VerifyVersion(ec.Ctx)
	if err != nil {
		return response.Error(err)
	}
	param := make(map[string]interface{})
	_ = ec.Ctx.ReadJSON(&param)
	userId := utils.ReadMapString(param, "userId")
	_, err = ec.Service.VerifyStamp(userId)
	if err != nil {
		return response.Error(err)
	}
	clusterId := utils.ReadMapInt(param, "clusterId")
	err = ec.Service.DisableCluster(clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (ec *ExternalController) PostClusterEnable() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	err = ec.Service.VerifyVersion(ec.Ctx)
	if err != nil {
		return response.Error(err)
	}
	param := make(map[string]interface{})
	_ = ec.Ctx.ReadJSON(&param)
	userId := utils.ReadMapString(param, "userId")
	_, err = ec.Service.VerifyStamp(userId)
	if err != nil {
		return response.Error(err)
	}
	clusterId := utils.ReadMapInt(param, "clusterId")
	err = ec.Service.EnableCluster(clusterId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}
