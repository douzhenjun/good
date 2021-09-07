package service

import (
	"DBaas/models"
	"DBaas/utils"
	"github.com/go-xorm/xorm"
)

type HomeService interface {
	CommonUserBaseInfo(userName string) (map[string]interface{},error)
	Cluster3DInfo(userName string) ([]models.ClusterInstance,error)
	GetMasterPodByCluster(clusterId int) (models.Instance, error)
	Pod3DInfo(userName string) ([]models.Instance,error)
	SelectOneCluster(id int) (models.ClusterInstance,error)
	User3DInfo(userName string) ([]models.User,error)
	GetOperatorStatus() (string, error)
}
type homeService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func NewHomeService(db *xorm.Engine, cs CommonService) HomeService {
	return &homeService{
		Engine: db,
		cs:     cs,
	}
}


func (hs *homeService) CommonUserBaseInfo(userName string) (map[string]interface{},error) {
	CommonUserBaseInfoMap := make(map[string]interface{})
	user := models.User{UserName: userName}
	_, err := hs.Engine.Get(&user)
	clusterCount, err := hs.Engine.Where(" user_id = ? ", user.Id).Count(new(models.ClusterInstance))
	onlineNum, err := hs.Engine.Where(" user_id = ? ", user.Id).And(" status = ? ", "True").Count(new(models.ClusterInstance))
	errorNum, err := hs.Engine.Where(" user_id = ? ", user.Id).And(" status != ? ", "True").Count(new(models.ClusterInstance))
	cpuUsed,err := hs.Engine.Where(" user_id = ? ", user.Id).SumInt(new(models.ClusterInstance),"limit_cpu")
	memUsed,err := hs.Engine.Where(" user_id = ? ", user.Id).SumInt(new(models.ClusterInstance),"limit_mem")
	storUsed,err := hs.Engine.Where(" user_id = ? ", user.Id).SumInt(new(models.ClusterInstance),"storage")
	if err != nil {
		utils.LoggerError(err)
		return CommonUserBaseInfoMap, err
	}
	CommonUserBaseInfoMap["clusterNum"]=clusterCount
	CommonUserBaseInfoMap["onlineNum"]=onlineNum
	//CommonUserBaseInfoMap["offlineNum"]=offlineNum
	CommonUserBaseInfoMap["errorNum"]=errorNum
	CommonUserBaseInfoMap["cpuTotal"]=user.CpuAll
	CommonUserBaseInfoMap["cpuUsed"]=cpuUsed
	CommonUserBaseInfoMap["memTotal"]=user.MemAll
	CommonUserBaseInfoMap["memUsed"]=memUsed
	CommonUserBaseInfoMap["storTotal"]=user.StorageAll
	CommonUserBaseInfoMap["storUsed"]=storUsed
	return CommonUserBaseInfoMap, nil
}


func (hs *homeService) User3DInfo(userName string) ([]models.User,error) {
	userList := make([]models.User, 0)
	var err error
	if userName !=""{
		err = hs.Engine.Where(" user_name = ? ", userName).OrderBy("-id").Find(&userList)
	}else{
		err = hs.Engine.OrderBy("-id").Find(&userList)
	}
	if err != nil {
		utils.LoggerError(err)
		return userList, err
	}
	return userList, err
}

func (hs *homeService) Cluster3DInfo(userName string) ([]models.ClusterInstance,error) {
	ClusterList := make([]models.ClusterInstance, 0)
	var err error
	if userName !=""{
		user := &models.User{UserName:userName}
		_, err =  hs.Engine.Get(user)
		err = hs.Engine.Where(" user_id = ? ", user.Id ).OrderBy("-id").Find(&ClusterList)
	}else{
		err = hs.Engine.OrderBy("-id").Find(&ClusterList)
	}

	if err != nil {
		utils.LoggerError(err)
		return ClusterList, err
	}
	return ClusterList, err
}

func (hs *homeService) GetMasterPodByCluster(clusterId int) (models.Instance, error) {
	var masterPod models.Instance
	_, err := hs.Engine.Where(" cluster_id = ? ", clusterId).And("role = ? ","Master").Get(&masterPod)
	if err != nil {
		utils.LoggerError(err)
		return masterPod,err
	}
	return masterPod, err
}

func (hs *homeService) Pod3DInfo(userName string) ([]models.Instance,error) {
	PodList := make([]models.Instance, 0)
	var err error
	if userName !=""{
		ClusterList := make([]models.ClusterInstance, 0)
		user := &models.User{UserName:userName}
		_, err =  hs.Engine.Get(user)
		err = hs.Engine.Where(" user_id = ? ", user.Id ).OrderBy("-id").Find(&ClusterList)
		for _, cluster := range ClusterList {
			SignalPodList := make([]models.Instance, 0)
			err = hs.Engine.Where(" cluster_id = ? ", cluster.Id ).OrderBy("-id").Find(&SignalPodList)
			for _, signalPod := range SignalPodList {
				PodList = append(PodList,signalPod)
			}
		}
	}else{
		err = hs.Engine.OrderBy("-id").Find(&PodList)
	}
	if err != nil {
		utils.LoggerError(err)
		return PodList, err
	}
	return PodList, err
}

func (hs *homeService) SelectOneCluster(id int) (models.ClusterInstance,error) {
	var cluster models.ClusterInstance
	_, err := hs.Engine.Where(" id = ? ", id).Omit("yaml_text").Get(&cluster)
	if err != nil {
		utils.LoggerError(err)
		return cluster,err
	}
	return cluster,err
}

func (hs *homeService) GetOperatorStatus() (string, error) {
	return models.GetConfig("operator@status", hs.Engine)
}