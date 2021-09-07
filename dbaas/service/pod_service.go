package service

import (
	"DBaas/models"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
)

type PodService interface {
	SelectPodsForAlarm() (map[string]interface{}, error)
	GetLog(podId int) (map[string]interface{}, error)
	GetDetail(podId, attrId, modelId, time int, selectType string, con map[string]interface{}) (map[string]interface{}, error)
}

func NewPodService(db *xorm.Engine, cs CommonService) PodService {
	return &podService{
		Engine: db,
		cs:     cs,
	}
}

type podService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func (ps *podService) GetDetail(podId, attrId, modelId, time int, selectType string, con map[string]interface{}) (map[string]interface{}, error) {
	if podId <= 0 {
		return nil, errors.New("podId must > 0")
	}
	pod := models.Instance{Id: podId}
	exist, err := ps.Engine.Cols("name", "cluster_id").Get(&pod)
	if !exist {
		return nil, fmt.Errorf("not found pod %v, error: %v", podId, err)
	}
	cluster := models.ClusterInstance{Id: pod.ClusterId}
	exist, err = ps.Engine.Omit("yaml_text").Get(&cluster)
	if !exist {
		return nil, fmt.Errorf("not found cluster %v, error: %v", pod.ClusterId, err)
	}
	return GetPerformanceData(int32(modelId), selectType, pod.Name, int32(attrId), "", int64(time), con, cluster), nil
}

func (ps *podService) GetLog(podId int) (map[string]interface{}, error) {
	if podId <= 0 {
		return nil, errors.New("pod id must > 0")
	}
	pod := models.Instance{Id: podId}
	exist, err := ps.Engine.Cols("name").Get(&pod)
	if !exist {
		return nil, fmt.Errorf("not found pod %v, error: %v", podId, err)
	}
	return ps.cs.GetLogsByLoki(pod.Name, "mysql")
}

func (ps *podService) SelectPodsForAlarm() (ret map[string]interface{}, err error) {
	ret = make(map[string]interface{})
	podList := make([]models.Instance, 0)
	err = ps.Engine.OrderBy("-id").Find(&podList)
	if err != nil {
		return
	}
	if len(podList) > 0 {
		for _, instance := range podList {
			instanceInfo := make(map[string]interface{})
			instanceInfo["model_object"] = "K8sMySQLPod"
			instanceInfo["module_name"] = "DBaaS"
			instanceInfo["inst_id"] = instance.Id
			instanceInfo["inst_name"] = instance.Name
			instanceInfo["host_ip"] = instance.DomainName
			cluster := models.ClusterInstance{Id: instance.ClusterId}
			_, err = ps.Engine.Cols("user_tag", "limit_cpu", "limit_mem").Get(&cluster)
			if err != nil {
				return
			}
			instanceInfo["user_tag"] = cluster.UserTag
			instanceInfo["cpuTotal"] = cluster.LimitCpu
			instanceInfo["memTotal"] = cluster.LimitMem
			ret[instance.Name] = instanceInfo
		}
	}
	return ret, nil
}
