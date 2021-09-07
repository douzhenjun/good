package service

import (
	"DBaas/models"
	"DBaas/utils"
	"DBaas/x/response"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"math"
	"strconv"
	"strings"
	"time"
)

var path2market = map[string]string{"/external/cluster/add": "baas"}

type ExternalService interface {
	// VerifyVersion 验证api版本
	VerifyVersion(ctx iris.Context) error
	// VerifyStamp 验证时间戳, 返回工号和error
	VerifyStamp(userId string) (string, error)
	// OpenCluster 开通实例
	OpenCluster(username string, clusterName string, imageVersion string, clusterType string, storageType string, password string, storageMap map[string]interface{}) (int, error)
	// SelectCluster 查询实例
	SelectCluster(clusterId int) (models.ClusterInstance, error)
	// DeleteCluster 删除实例
	DeleteCluster(clusterId int) error
	// Login 免登录跳转
	Login(username string, clusterId int) (models.User, error)
	// CheckQuota 检测接口配额
	CheckQuota(storageMap map[string]interface{}, path string) error
	DisableCluster(clusterId int) error
	EnableCluster(clusterId int) error
}

func getApiUsage(path string, engine *xorm.Engine) (*models.ApiQuota, error) {
	market, ok := path2market[path]
	if !ok {
		return nil, errors.New("not found market")
	}
	clusterList := make([]models.ClusterInstance, 0)
	err := engine.SQL("select c.limit_mem, c.limit_cpu, c.storage from cluster_instance c inner join \"user\" u on c.user_id = u.id where u.auto_create = true and u.user_name like ? and c.deleted_at is null", market+"-%").Find(&clusterList)
	if err != nil {
		return nil, err
	}
	var cpuTotal, memTotal, storageTotal int
	for i := range clusterList {
		cpuTotal += clusterList[i].LimitCpu
		memTotal += clusterList[i].LimitMem
		storageTotal += clusterList[i].Storage
	}
	return &models.ApiQuota{Path: path, Cpu: cpuTotal, Memory: memTotal, Storage: storageTotal}, nil
}

func (es *externalService) CheckQuota(storageMap map[string]interface{}, path string) error {
	usage, err := getApiUsage(path, es.Engine)
	if err != nil {
		return err
	}
	usage.Cpu += int(storageMap["cpu"].(float64))
	usage.Memory += int(storageMap["mem"].(float64))
	usage.Storage += int(storageMap["size"].(float64))

	quota := models.ApiQuota{Path: path}
	exist, err := es.Engine.Get(&quota)
	if !exist {
		return fmt.Errorf("not found api quota: %v, error: %v", path, err)
	}
	if usage.Cpu > quota.Cpu {
		return fmt.Errorf("the total cpu of this api exceeds the quota limit, now: %v, max: %v", usage.Cpu, quota.Cpu)
	}
	if usage.Memory > quota.Memory {
		return fmt.Errorf("the total memory of this api exceeds the quota limit, now: %v, max: %v", usage.Memory, quota.Memory)
	}
	if usage.Storage > quota.Storage {
		return fmt.Errorf("the total storage of this api exceeds the quota limit, now: %v, max: %v", usage.Storage, quota.Storage)
	}
	return nil
}

type externalService struct {
	Engine  *xorm.Engine
	Cluster ClusterService
	Common  CommonService
}

func (es *externalService) DisableCluster(clusterId int) error {
	return es.Cluster.ClusterDisable(clusterId)
}

func (es *externalService) EnableCluster(clusterId int) error {
	return es.Cluster.ClusterEnable(clusterId)
}

func NewExternalService(db *xorm.Engine, cluster ClusterService, common CommonService) ExternalService {
	return &externalService{
		Engine:  db,
		Cluster: cluster,
		Common:  common,
	}
}

const externalApiVersion = "1.0"

func (es *externalService) VerifyVersion(ctx iris.Context) error {
	apiVer := ctx.URLParam("apiVersion")
	if externalApiVersion == apiVer {
		return nil
	}
	return fmt.Errorf("unsupported version %v", apiVer)
}

func (es *externalService) VerifyStamp(userId string) (string, error) {
	userId = strings.TrimSpace(userId)
	if len(userId) == 0 {
		return "", errors.New("userId is empty")
	}
	userId = strings.ReplaceAll(userId, " ", "+")
	decrypt := utils.AesDecryptApiDef(userId)
	if len(decrypt) <= 10 {
		return "", errors.New("user id length cannot be less than 10")
	}
	stamp, err := strconv.ParseFloat(decrypt[len(decrypt)-10:], 64)
	if err != nil {
		return "", fmt.Errorf("parse user id is fail: %s", err)
	}
	nowStamp := time.Now().Unix()
	diffStamp := float64(nowStamp) - stamp
	if math.Abs(diffStamp) > 300 {
		return "", errors.New("verify user id is fail")
	}
	return decrypt[:len(decrypt)-11], nil
}

func (es *externalService) OpenCluster(username string, clusterName string, imageVersion string, clusterType string, storageType string, password string, storageMap map[string]interface{}) (clusterId int, err error) {
	image := models.Images{Version: imageVersion}
	has, _ := es.Engine.Cols("id").Get(&image)
	if !has {
		return 0, fmt.Errorf("no image for %v found", imageVersion)
	}
	user := models.User{UserName: username}
	exist, err := es.Engine.Cols("id", "user_tag").Get(&user)
	if err != nil {
		return
	}
	var scName string
	if !exist {
		// 用户不存在, 自动创建
		scConfig := models.Sysparameter{ParamKey: "zcm_sc_" + storageType}
		_, _ = es.Engine.Get(&scConfig)
		scList := strings.Split(scConfig.ParamValue, ",")
		if len(scList) == 0 {
			return 0, errors.New("not found sc configs in system parameter")
		}
		scName = scList[0] // TODO: 暂时取第一个
		sc := models.Sc{Name: scName}
		has, _ = es.Engine.Cols("id", "sc_type").Get(&sc)
		if !has {
			return 0, fmt.Errorf("not found sc: %s", scName)
		}
		if sc.ScType != "shared-storage" {
			return 0, errors.New("the sc type in the configuration is not shared-storage")
		}

		capricorn, conn := NewCapricornService()
		defer CloseGrpc(conn)
		// 获取角色列表信息
		roleList, _, _ := capricorn.GetRoleResources("", "")
		var roleId string
		for i := range roleList {
			if roleList[i]["rolename"] == "dbaas" { // TODO: dbaas-zcm
				roleId = strconv.Itoa(int(roleList[i]["id"].(float64)))
			}
		}
		if len(roleId) == 0 {
			return 0, errors.New("dbaas-zcm role not found")
		}

		var userPass = utils.AesEncryptPassword(utils.RandCode(6))
		userInfo, errMsg, _ := capricorn.AddUserResources(roleId, "root", username, userPass, "", "")
		if errMsg != "" {
			return 0, errors.New(errMsg)
		}
		user.ZdcpId = int(userInfo["user_id"].(float64))
		// zcm租户配额不限制
		user.MemAll = 99999
		user.CpuAll = 99999
		user.StorageAll = 999999
		user.UserTag = userInfo["user_tag"].(string)
		user.AutoCreate = true
		user.Password = userPass
		_, err = es.Engine.Insert(&user)
		if err != nil {
			capricorn.DeleteUserResources(strconv.Itoa(user.ZdcpId), "root")
			return 0, fmt.Errorf("auto create user is fail: %s", err)
		}
		scUser := models.ScUser{ScId: sc.Id, UserId: user.Id}
		_, _ = es.Engine.Insert(&scUser)
	} else {
		_, _ = es.Engine.SQL("select sc.name from sc left outer join sc_user su on sc.id = su.sc_id where su.user_id = ?", user.Id).Get(&scName)
		if len(scName) == 0 {
			return 0, errors.New("not found sc in user")
		}
	}
	storageMap["scName"] = scName
	return es.Cluster.Add(clusterName, password, storageMap, nil, "", user.Id, image.Id, "", "external", nil, 0, 0)
}

func (es *externalService) SelectCluster(clusterId int) (models.ClusterInstance, error) {
	cluster := models.ClusterInstance{Id: clusterId}
	exist, _ := es.Engine.
		Omit("yaml_text", "console_port", "inner_connect_string", "sc_name", "remark", "master", "user_tag").
		Get(&cluster)
	if !exist {
		return cluster, fmt.Errorf("cannot find cluster with id %v", clusterId)
	}
	masterAddress := models.Sysparameter{ParamKey: "kubernetes_master_address"}
	has, _ := es.Engine.Get(&masterAddress)
	if has {
		cluster.ConnectString = fmt.Sprintf("%v:%v", masterAddress.ParamValue, cluster.ConnectString)
	}
	image := models.Images{Id: cluster.ImageId}
	has, _ = es.Engine.Get(&image)
	if has {
		cluster.ImageName = fmt.Sprintf("%v:%v", image.ImageName, image.Version)
	}
	cluster.PodStatusMap = utils.RawJson(cluster.PodStatus)
	if cluster.Status == models.ClusterStatusFalse {
		cluster.Events = es.Common.GetEvent(cluster.K8sName)
	}
	return cluster, nil
}

func (es *externalService) DeleteCluster(clusterId int) error {
	return es.Cluster.Delete(clusterId, false)
}

func (es *externalService) Login(username string, clusterId int) (models.User, error) {
	if clusterId > 0 {
		var cluster models.ClusterInstance
		_, _ = es.Engine.ID(clusterId).Cols("status").Get(&cluster)
		if cluster.Status == models.ClusterStatusDisable {
			return models.User{}, response.NewMsg("the instance has been disabled", "该实例已被禁用")
		}
	}
	user := models.User{UserName: username}
	exist, err := es.Engine.Cols("password").Get(&user)
	if !exist {
		return models.User{}, fmt.Errorf("not found user %v error: %v", username, err)
	}
	return user, nil
}
