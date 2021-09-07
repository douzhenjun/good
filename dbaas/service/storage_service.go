/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: ddh
 * @LastEditTime: 2021-02-07 16:32:07
 */

package service

import (
	"DBaas/models"
	"DBaas/utils"
	"DBaas/x/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	appsv1 "k8s.io/api/apps/v1"
	core1 "k8s.io/api/core/v1"
	storage1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
	"strings"
	"time"
)

type StorageService interface {
	UserAssign(id int, userIdStr string) error
	List(page int, pageSize int, key string, userId int, userTag string, isFilter bool) ([]models.ReturnSc, int64)
	Update(id int, remake string, nodeNum int) error
	Delete(id int) error
	Add(scName string, reclaimPolicy string, remark string, orgTag string, userTag string, userId int, scType string, nodeNum int, userIdStr string) (models.Sc, error)
	PvAdd(storageId int, pvName string, mountPoint string, iqn string, lun int, size string, userTag string, orgTag string, namespace string) (bool, models.PersistentVolume, string)
	PvDelete(id int) (err error)
	PVList(page int, pageSize int, key string, userTag string) ([]models.PersistentVolume, int64, error)
	CreateMysqlByPV(pvId int, storageMap map[string]interface{}, remark string, userId int, mysqlName string, qos *models.Qos) error
	SelectOneScByName(name string) (models.Sc, bool)
	SelectOnePvByName(name string) (models.PersistentVolume, bool)
	UserRegister(userId int, scList []map[string]interface{}) error
	DeletescUserbyUser(userId int) (bool, string)
	AddscUserbyuser(userId int, scList []map[string]interface{}) (bool, string)
}

type storageService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func NewStorageService(engine *xorm.Engine, cs CommonService) StorageService {
	return &storageService{
		Engine: engine,
		cs:     cs,
	}
}

func (ss *storageService) PVList(page int, pageSize int, key string, userTag string) (list []models.PersistentVolume, count int64, err error) {
	list = make([]models.PersistentVolume, 0)
	where := "name like ?"
	args := []interface{}{"%" + key + "%"}
	// AAAA为root用户可查看所有PV
	if userTag != "AAAA" {
		where += " AND user_tag = ?"
		args = append(args, userTag)
	}
	err = ss.Engine.Where(where, args...).Limit(pageSize, pageSize*(page-1)).Desc("id").Find(&list)
	if err != nil {
		return
	}
	count, _ = ss.Engine.Where(where, args...).Count(&models.PersistentVolume{})
	userCache := map[string]models.User{}
	for i := range list {
		tag := list[i].UserTag
		user, ok := userCache[tag]
		if !ok {
			user = models.User{UserTag: tag}
			_, _ = ss.Engine.Get(&user)
			userCache[tag] = user
		}
		list[i].UserId = user.Id
		list[i].Tenant = user.UserName
		list[i].CpuTotal = user.CpuAll
		list[i].MemTotal = int(user.MemAll)
		pod := models.Instance{Id: list[i].PodId}
		_, _ = ss.Engine.Cols("name").Get(&pod)
		list[i].PodName = pod.Name
		sc := models.Sc{Id: list[i].ScId}
		_, _ = ss.Engine.Cols("name").Get(&sc)
		list[i].SCName = sc.Name
	}
	return
}

func (ss *storageService) CreateMysqlByPV(pvId int, storageMap map[string]interface{}, remark string, userId int, mysqlName string, qos *models.Qos) (err error) {
	nameExist, _ := ss.Engine.Where("is_deploy = true").Exist(&models.ClusterInstance{Name: mysqlName})
	if nameExist {
		return errors.New("This name already exists ")
	}
	exist, _ := ss.Engine.Exist(&models.ClusterInstance{PvId: pvId})
	if exist {
		return errors.New("The MySQL created by the current pv already exists ")
	}
	dbPV := models.PersistentVolume{Id: pvId}
	_, err = ss.Engine.Get(&dbPV)
	if err != nil {
		return
	}
	// 检查pv是否存在
	err, pvSource := ss.cs.GetResources("pv", dbPV.Name, ss.cs.GetNameSpace(), meta1.GetOptions{})
	if err != nil {
		return fmt.Errorf("pv does not exist in k8s, err: %v", err)
	}

	// 根据pv查询已删除的pod
	pod := models.Instance{Id: dbPV.PodId}
	existPod, err := ss.Engine.Unscoped().Cols("cluster_id").Get(&pod)
	if !existPod {
		return errors.New(fmt.Sprintf("Found pod is error: %v", err))
	}
	cluster := models.ClusterInstance{Id: pod.ClusterId}
	existCluster, err := ss.Engine.Unscoped().Cols("user_id", "image_id", "storage", "user_tag", "org_tag", "secret").Get(&cluster)
	if !existCluster {
		return errors.New(fmt.Sprintf("Found cluster is error: %v", err))
	}
	limitMem, limitCpu := int(storageMap["mem"].(float64)), int(storageMap["cpu"].(float64))
	enough, msg, _ := getUserResource(userId, limitCpu, limitMem, cluster.Storage, ss.Engine)
	if !enough {
		return errors.New(msg)
	}

	pvcName := mysqlName + "-pvc"
	scName := mysqlName + "-sc"
	svcName := mysqlName + "-svc"

	mysqlImage := models.Images{Id: cluster.ImageId}
	hasImage, err := ss.Engine.Get(&mysqlImage)
	if !hasImage {
		return errors.New(fmt.Sprintf("Found mysql image is error: %v", err))
	}
	var imageURL string
	if mysqlImage.Status == "Invalid" {
		imageURL = fmt.Sprintf("%v:%v", mysqlImage.ImageName, mysqlImage.Version)
	} else {
		imageURL = fmt.Sprintf("%v/%v:%v", getImageAddress(ss.Engine), mysqlImage.ImageName, mysqlImage.Version)
	}

	// 设置pv的sc名称
	pv := (*pvSource).(*core1.PersistentVolume)
	pv.Spec.ClaimRef = nil
	pv.Spec.StorageClassName = scName
	_, err = ss.cs.GetClientSet().CoreV1().PersistentVolumes().Update(*ss.cs.GetCtx(), pv, meta1.UpdateOptions{})
	if err != nil {
		return
	}

	pvcConfig := core1.PersistentVolumeClaim{
		TypeMeta:   meta1.TypeMeta{Kind: "PersistentVolumeClaim", APIVersion: "v1"},
		ObjectMeta: meta1.ObjectMeta{Name: pvcName},
		Spec: core1.PersistentVolumeClaimSpec{
			StorageClassName: &scName,
			AccessModes:      []core1.PersistentVolumeAccessMode{core1.ReadWriteOnce},
			Resources: core1.ResourceRequirements{
				Requests: map[core1.ResourceName]resource.Quantity{
					"storage": resource.MustParse(dbPV.Capacity),
				},
			},
		},
	}
	err = ss.cs.CreateOption("pvc", ss.cs.GetNameSpace(), &pvcConfig, meta1.CreateOptions{})
	if err != nil {
		return
	}
	svcConfig := core1.Service{
		TypeMeta:   meta1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		ObjectMeta: meta1.ObjectMeta{Name: svcName},
		Spec: core1.ServiceSpec{
			Type: core1.ServiceTypeNodePort,
			Ports: []core1.ServicePort{
				{Name: "mysqlport", Port: 3306, TargetPort: intstr.FromInt(3306)},
				{Name: "sidecar-ttyd", Port: 7681, TargetPort: intstr.FromInt(7681)},
			},
			Selector: map[string]string{"app": "mysql"},
		},
	}
	err = ss.cs.CreateOption("service", ss.cs.GetNameSpace(), &svcConfig, meta1.CreateOptions{})
	if err != nil {
		return
	}
	secretMap := map[string]string{}
	err = json.Unmarshal([]byte(cluster.Secret), &secretMap)
	if err != nil {
		return
	}
	mysqlConfig := appsv1.Deployment{
		TypeMeta:   meta1.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		ObjectMeta: meta1.ObjectMeta{Name: mysqlName},
		Spec: appsv1.DeploymentSpec{
			Selector: &meta1.LabelSelector{MatchLabels: map[string]string{"app": "mysql"}},
			Strategy: appsv1.DeploymentStrategy{Type: appsv1.RecreateDeploymentStrategyType},
			Template: core1.PodTemplateSpec{
				ObjectMeta: meta1.ObjectMeta{Labels: map[string]string{"app": "mysql"}},
				Spec: core1.PodSpec{
					Containers: []core1.Container{
						{
							Name:  "mysql",
							Image: imageURL,
							Resources: core1.ResourceRequirements{
								Limits: core1.ResourceList{
									"memory": resource.MustParse(fmt.Sprintf("%vGi", limitMem)),
									"cpu":    resource.MustParse(strconv.Itoa(limitCpu)),
								},
							},
							Env:          []core1.EnvVar{{Name: "MYSQL_ROOT_PASSWORD", Value: secretMap["ROOT_PASSWORD"]}},
							Ports:        []core1.ContainerPort{{Name: "mysql", ContainerPort: 3306}},
							VolumeMounts: []core1.VolumeMount{{Name: "mysql-persistent-storage", MountPath: "/var/lib/mysql"}},
						},
						{
							Name:  "sidecar",
							Image: "10.45.10.107:8099/k8s/mysql-sidecar:ttyd",
							Ports: []core1.ContainerPort{
								{Name: "mysql", ContainerPort: 7681},
							},
						},
					},
					Volumes: []core1.Volume{
						{Name: "mysql-persistent-storage", VolumeSource: core1.VolumeSource{
							PersistentVolumeClaim: &core1.PersistentVolumeClaimVolumeSource{ClaimName: pvcName},
						}},
					},
				},
			},
		},
	}
	_, err = ss.cs.GetClientSet().AppsV1().Deployments(ss.cs.GetNameSpace()).Create(*ss.cs.GetCtx(), &mysqlConfig, meta1.CreateOptions{})
	if err != nil {
		return err
	}

	// 查找pod名称
	var podName string
	for i := 0; i < 10; i++ {
		err, podListSource := ss.cs.GetResources("pod", "", ss.cs.GetNameSpace(), meta1.ListOptions{LabelSelector: "app=mysql"})
		if err != nil {
			utils.LoggerError(err)
			<-time.After(time.Second)
			continue
		}
		podList := (*podListSource).(*core1.PodList)
		for _, pod := range podList.Items {
			for _, volume := range pod.Spec.Volumes {
				if volume.PersistentVolumeClaim == nil {
					continue
				}
				if volume.PersistentVolumeClaim.ClaimName == pvcName {
					podName = pod.Name
					goto FindPodNameEnd
				}
			}
		}
		<-time.After(time.Second)
	}
FindPodNameEnd:
	clusterIns := models.ClusterInstance{
		Name:     mysqlName,
		K8sName:  podName,
		Status:   models.ClusterStatusCreating,
		Storage:  cluster.Storage,
		UserId:   cluster.UserId,
		ImageId:  cluster.ImageId,
		ScName:   scName,
		Replicas: "1",
		LimitCpu: limitCpu,
		LimitMem: limitMem,
		Remark:   remark,
		UserTag:  cluster.UserTag,
		OrgTag:   cluster.OrgTag,
		IsDeploy: true,
		Secret:   cluster.Secret,
		PvId:     pvId,
	}
	_, err = ss.Engine.Insert(&clusterIns)
	if err == nil {
		if qos != nil {
			qos.ClusterId = clusterIns.Id
			_, _ = ss.Engine.Insert(qos)
			go ss.cs.SetQosConfig(dbPV.Name, clusterIns.Id)
		}
		go ss.cs.CreatStatusTimeout(clusterIns.Id)
		go ss.cs.ScanClusterPod(clusterIns.Id, clusterIns.K8sName, 1, true)
		go ss.cs.PollingPVStatus([]models.PersistentVolume{dbPV}, core1.VolumeBound)
		dbPV.PvcName = pvcName
		_, _ = ss.Engine.ID(dbPV.Id).Cols("pvc_name").Update(&dbPV)
	}
	return err
}

// UserAssign userIdStr为-1时表示ALL，为空时表示无租户，有租户时以逗号分割
func (ss *storageService) UserAssign(id int, userIdStr string) error {
	if id <= 0 {
		return errors.New("storage id must > 0")
	}

	sc := models.Sc{Id: id}
	exist, err := ss.Engine.Cols("name", "sc_type").Get(&sc)
	if !exist {
		return fmt.Errorf("not found sc %v, error: %v", id, err)
	}

	assignAll := userIdStr == "-1"
	if assignAll {
		if sc.ScType == models.ScTypeUnique {
			return errors.New("the unique storage cannot be set to ALL")
		}
		_, _ = ss.Engine.Where("sc_id = ?", id).Delete(new(models.ScUser))
	} else {
		dbUsers := make([]models.ScUser, 0)
		err = ss.Engine.Where("sc_id = ?", id).Find(&dbUsers)
		if err != nil {
			return err
		}
		userStrList := strings.Split(userIdStr, ",")
		if sc.ScType == models.ScTypeUnique {
			// 处理独占存储
			var clusterCount int64
			clusterCount, err = ss.Engine.Where("sc_name = ?", sc.Name).Count(new(models.ClusterInstance))
			if err != nil {
				return err
			}
			if clusterCount > 0 {
				return response.NewMsg("This storage is already in use and cannot be modified", "此存储已被使用，无法修改")
			}
			if len(userStrList) <= 0 {
				_, err = ss.Engine.Where("sc_id = ?", id).Delete(new(models.ScUser))
				return err
			}
			if len(userStrList) != 1 {
				return errors.New("unique-storage only assign one user")
			}
			var userId int
			userId, err = strconv.Atoi(userStrList[0])
			if err != nil {
				return err
			}
			scUser := models.ScUser{UserId: userId, ScId: id}
			if len(dbUsers) == 0 {
				_, err = ss.Engine.Insert(&scUser)
			} else {
				_, err = ss.Engine.Where("sc_id = ?", id).Update(&scUser)
			}
		} else {
			// <userId,index>
			dbUsersM := map[int]int{}
			for i := range dbUsers {
				dbUsersM[dbUsers[i].UserId] = i
			}

			insertList := make([]models.ScUser, 0)
			for i := range userStrList {
				if len(userStrList[i]) == 0 {
					continue
				}
				var userId, err = strconv.Atoi(userStrList[i])
				if err != nil {
					utils.LoggerError(err)
					continue
				}
				if _, ok := dbUsersM[userId]; ok {
					// 数据库已存在直接跳过，并标记为-1
					dbUsersM[userId] = -1
					continue
				}
				insertList = append(insertList, models.ScUser{ScId: id, UserId: userId})
			}
			for _, v := range dbUsersM {
				if v != -1 {
					clusterCount, err := ss.Engine.Where("sc_name = ?", sc.Name).And("user_id = ?", dbUsers[v].UserId).Count(new(models.ClusterInstance))
					if err != nil {
						return err
					}
					if clusterCount > 0 {
						u := models.User{}
						_, _ = ss.Engine.ID(dbUsers[v].UserId).Cols("user_name").Get(&u)
						return response.NewMsg(fmt.Sprintf("%v have already used this storage, this user cannot be deleted", u.UserName), fmt.Sprintf("%v已使用此存储，不能删除此用户", u.UserName))
					}
				}
			}
			if len(insertList) > 0 {
				_, err := ss.Engine.Insert(&insertList)
				if err != nil {
					return err
				}
			}
			for _, v := range dbUsersM {
				if v != -1 {
					_, _ = ss.Engine.ID(dbUsers[v].Id).Delete(new(models.ScUser))
				}
			}
		}
	}
	_, err = ss.Engine.ID(id).Cols("assign_all").Update(&models.Sc{AssignAll: assignAll})
	return err
}

func (ss *storageService) List(page int, pageSize int, key string, userId int, userTag string, isFilter bool) ([]models.ReturnSc, int64) {
	scPv := make([]models.Sc, 0)
	var count int64
	if userId <= 0 && len(userTag) != 0 {
		u := models.User{}
		_, _ = ss.Engine.Where("user_tag = ?", userTag).Cols("id").Get(&u)
		userId = u.Id
	}
	if userId > 0 {
		scList := make([]models.ScUser, 0)
		err := ss.Engine.Where("user_id = ?", userId).Find(&scList)
		utils.LoggerError(err)

		for _, v := range scList {
			sc := models.Sc{Id: v.ScId}
			_, err = ss.Engine.Where("name like ? or describe like ? ", "%"+key+"%", "%"+key+"%").Get(&sc)
			utils.LoggerError(err)
			if len(sc.Name) > 0 {
				scPv = append(scPv, sc)
			}
		}

		assignAllSc := make([]models.Sc, 0)
		err = ss.Engine.Where("assign_all = true").Find(&assignAllSc)
		utils.LoggerError(err)
		if len(assignAllSc) > 0 {
			scPv = append(scPv, assignAllSc...)
		}

		count = int64(len(scPv))
		if utils.MustInt(page, pageSize) {
			min := pageSize * (page - 1)
			max := min + pageSize
			scPv = scPv[min:utils.Min(max, len(scPv))]
		}
	} else {
		err := ss.Engine.Where("name like ? or describe like ? ", "%"+key+"%", "%"+key+"%").Limit(pageSize, pageSize*(page-1)).Desc("id").Find(&scPv)
		utils.LoggerError(err)
		count, _ = ss.Engine.Where("name like ? or describe like ? ", "%"+key+"%", "%"+key+"%").Count(&models.Sc{})
	}

	scReturn := make([]models.ReturnSc, len(scPv))
	for i, sc := range scPv {
		pv := make([]models.PersistentVolume, 0)
		if sc.ScType == models.ScTypeUnique {
			err := ss.Engine.Where(" sc_id = ?", sc.Id).OrderBy("id").Find(&pv)
			utils.LoggerError(err)
			sc.NodeNum = len(pv)
		}

		cluster := make([]models.ClusterInstance, 0)
		err := ss.Engine.Where("sc_name = ?", sc.Name).Omit("yaml_text").Find(&cluster)
		utils.LoggerError(err)

		var scUserRaw json.RawMessage
		if sc.AssignAll {
			scUserRaw = []byte("-1")
		} else {
			scUser := make([]models.ScUser, 0)
			err = ss.Engine.Where("sc_id = ?", sc.Id).Find(&scUser)
			utils.LoggerError(err)
			for i, user := range scUser {
				u := models.User{Id: user.UserId}
				_, err = ss.Engine.Get(&u)
				utils.LoggerError(err)
				scUser[i].UserName = u.UserName
			}
			scUserRaw, _ = json.Marshal(scUser)
		}

		scReturn[i] = models.ReturnSc{Sc: sc, Children: pv, Cluster: cluster, ScUser: scUserRaw}
	}

	if isFilter {
		for i, sc := range scReturn {
			if sc.ScType == models.ScTypeUnique && len(sc.Cluster) > 0 {
				scReturn = append(scReturn[0:i], scReturn[i+1:]...)
			}
		}
	}
	return scReturn, count
}

func (ss *storageService) Add(scName string, reclaimPolicy string, remark string, orgTag string, userTag string, userId int, scType string, nodeNum int, userIdStr string) (models.Sc, error) {
	sc := models.Sc{
		Name:          scName,
		ScType:        scType,
		NodeNum:       nodeNum,
		ReclaimPolicy: reclaimPolicy,
		Describe:      remark,
		OrgTag:        orgTag,
		UserTag:       userTag,
		AssignAll:     userIdStr == "-1",
	}
	namespace := ss.cs.GetNameSpace()
	// 独有存储，在k8s里面新建
	if scType == "unique-storage" {
		reclaimPolicyCore := core1.PersistentVolumeReclaimPolicy(reclaimPolicy)
		scConfig := storage1.StorageClass{
			TypeMeta: meta1.TypeMeta{
				Kind:       "StorageClass",
				APIVersion: "storage.k8s.io/v1",
			},
			ObjectMeta: meta1.ObjectMeta{
				Name: scName,
			},
			Provisioner:   "kubernetes.io/no-provisioner",
			ReclaimPolicy: &reclaimPolicyCore,
		}
		err := ss.cs.CreateOption("sc", namespace, &scConfig, meta1.CreateOptions{})
		if err != nil {
			return sc, err
		}
	} else {
		err, scAddr := ss.cs.GetResources("sc", scName, namespace, meta1.GetOptions{})
		if err != nil {
			return sc, err
		}
		if value, ok := (*scAddr).(*storage1.StorageClass); ok {
			sc.ReclaimPolicy = string(*value.ReclaimPolicy)
		}
	}
	_, err := ss.Engine.Insert(&sc)
	if err != nil || sc.AssignAll {
		return sc, err
	}

	userIds := strings.Split(userIdStr, ",")
	su := make([]models.ScUser, 0)
	for i := range userIds {
		if len(userIds[i]) == 0 {
			continue
		}
		id, err := strconv.Atoi(userIds[i])
		if err != nil {
			utils.LoggerError(err)
			continue
		}
		su = append(su, models.ScUser{UserId: id, ScId: sc.Id})
	}
	if userId > 0 {
		su = append(su, models.ScUser{UserId: userId, ScId: sc.Id})
	}

	if len(su) > 0 {
		_, _ = ss.Engine.Insert(&su)
	}
	return sc, nil
}

func (ss *storageService) Update(id int, remake string, nodeNum int) error {
	sc := models.Sc{
		Id:       id,
		Describe: remake,
		NodeNum:  nodeNum,
	}
	_, err := ss.Engine.ID(sc.Id).Update(&sc)
	return err
}

func (ss *storageService) Delete(id int) error {
	if id <= 0 {
		return errors.New("storage id must > 0")
	}
	sc := models.Sc{
		Id: id,
	}
	exist, err := ss.Engine.Get(&sc)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}

	clusterCount, err := ss.Engine.Where("sc_name = ?", sc.Name).Count(new(models.ClusterInstance))
	if err != nil {
		return err
	}

	if clusterCount > 0 {
		return response.NewMsg("This storage is occupied by the cluster and cannot be deleted", "此存储被集群占用，无法删除！")
	}

	if sc.ScType == models.ScTypeUnique {
		pvCount, err := ss.Engine.Where("sc_id = ?", id).Count(new(models.PersistentVolume))
		if err != nil {
			return err
		}
		if pvCount > 0 {
			return response.NewMsg("This store has pv and cannot be deleted", "请先删除此存储下的pv！")
		}
		err = ss.cs.DeleteOption("sc", sc.Name, ss.cs.GetNameSpace(), meta1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	_, err = ss.Engine.ID(sc.Id).Delete(&sc)
	_, _ = ss.Engine.Where("sc_id = ?", sc.Id).Delete(new(models.ScUser))
	return err
}

func (ss *storageService) PvAdd(storageId int, pvName string, mountPoint string, iqn string, lun int, size string, userTag string, orgTag string, namespace string) (bool, models.PersistentVolume, string) {
	ipAddr := ""
	port := ""

	if strings.Contains(mountPoint, ":") {
		ipAddr = strings.Split(mountPoint, ":")[0]
		port = strings.Split(mountPoint, ":")[1]
	} else {
		return false, models.PersistentVolume{}, "mountPoint format error"
	}

	pv := models.PersistentVolume{
		Name:     pvName,
		ScId:     storageId,
		Lun:      lun,
		Capacity: size,
		Iqn:      iqn,
		IpAddr:   ipAddr,
		Port:     port,
		UserTag:  userTag,
		OrgTag:   orgTag,
	}

	sc := models.Sc{Id: storageId}

	success, err := ss.Engine.Get(&sc)
	utils.LoggerError(err)

	if !success {
		if err != nil {
			return success, pv, err.Error()
		} else {
			return success, pv, ""
		}
	}

	scConfig := core1.PersistentVolume{
		TypeMeta: meta1.TypeMeta{
			Kind:       "PersistentVolume",
			APIVersion: "v1",
		},
		ObjectMeta: meta1.ObjectMeta{
			Name: pvName,
		},
		Spec: core1.PersistentVolumeSpec{
			Capacity: core1.ResourceList{
				core1.ResourceName("storage"): resource.MustParse(fmt.Sprintf("%sGi", size)),
			},
			AccessModes:      []core1.PersistentVolumeAccessMode{core1.PersistentVolumeAccessMode("ReadWriteOnce")},
			StorageClassName: sc.Name,
			PersistentVolumeSource: core1.PersistentVolumeSource{
				ISCSI: &core1.ISCSIPersistentVolumeSource{
					TargetPortal:   mountPoint,
					IQN:            iqn,
					ISCSIInterface: "iser",
					Lun:            int32(lun),
					FSType:         "xfs",
					ReadOnly:       false,
				},
			},
		},
	}

	err = ss.cs.CreateOption("pv", namespace, &scConfig, meta1.CreateOptions{})
	if err != nil {
		return false, pv, err.Error()
	}

	_, err = ss.Engine.Insert(&pv)
	if err != nil {
		return false, pv, err.Error()
	}

	return true, pv, ""
}

func (ss *storageService) PvDelete(id int) (err error) {
	pv := models.PersistentVolume{Id: id}
	hasPV, err := ss.Engine.Get(&pv)
	if !hasPV {
		return fmt.Errorf("Not found pv %v error: %s ", id, err)
	}
	if pv.Status == string(core1.VolumeBound) {
		return errors.New("Bound status cannot be deleted ")
	}
	sc := models.Sc{Id: pv.ScId}
	hasSc, err := ss.Engine.Get(&sc)
	if !hasSc {
		return fmt.Errorf("Not found sc %v error: %s ", pv.ScId, err)
	}
	if sc.ScType != "shared-storage" {
		// 判断没有没cluster集群占用
		existCluster, err := ss.Engine.Exist(&models.ClusterInstance{ScName: sc.Name})
		if err != nil {
			return err
		}
		if existCluster {
			return errors.New("There are clusters in this PV ")
		}
	}

	err = ss.cs.DeleteOption("pv", pv.Name, "", meta1.DeleteOptions{})
	if err != nil && !utils.ErrorContains(err, "not found") {
		return
	}
	_, err = ss.Engine.ID(pv.Id).Delete(&pv)
	if err != nil {
		return
	}
	return nil
}

func (ss *storageService) SelectOneScByName(name string) (models.Sc, bool) {
	var sc models.Sc
	_, err := ss.Engine.Where(" name = ? ", name).Get(&sc)
	utils.LoggerError(err)
	return sc, err == nil
}

func (ss *storageService) SelectOnePvByName(name string) (models.PersistentVolume, bool) {
	var pv models.PersistentVolume
	_, err := ss.Engine.Where(" name = ? ", name).Get(&pv)
	utils.LoggerError(err)
	return pv, err == nil
}

func (ss *storageService) UserRegister(userId int, scList []map[string]interface{}) error {
	session := ss.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}
	scUser := models.ScUser{UserId: userId}
	_, err := session.Delete(&scUser)
	if err != nil {
		return err
	}
	if len(scList) > 0 {
		for _, scInfo := range scList {
			if scInfo["type"] == "ready" {
				scId := int(scInfo["id"].(float64))
				sc := models.Sc{Id: scId}
				success, err := session.Cols("sc_type").Get(&sc)
				if !success {
					_ = session.Rollback()
					return fmt.Errorf("not found sc %v, error: %v", scId, err)
				}

				count, err := session.Where(" sc_id = ? ", scId).Count(new(models.ScUser))
				if err != nil {
					_ = session.Rollback()
					return err
				}

				if sc.ScType == "unique-storage" && count > 1 {
					return errors.New("unique-storage only assign one user")
				}
				scUser := models.ScUser{UserId: userId, ScId: sc.Id}
				_, err = session.Insert(&scUser)
				if err != nil {
					_ = session.Rollback()
					return err
				}
			}
		}
	}
	err = session.Commit()
	return err
}

func (ss *storageService) DeletescUserbyUser(userId int) (bool, string) {
	scUser := models.ScUser{UserId: userId}
	_, err := ss.Engine.Delete(&scUser)
	if err != nil {
		utils.LoggerError(err)
		return true, err.Error()
	}
	return true, ""
}

func (ss *storageService) AddscUserbyuser(userId int, scList []map[string]interface{}) (bool, string) {
	session := ss.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		iris.New().Logger().Info(err.Error())
	}
	if len(scList) > 0 {
		for _, scInfo := range scList {
			if scInfo["type"] == "ready" {
				scId := int(scInfo["id"].(float64))
				sc := models.Sc{Id: scId}
				success, err := session.Get(&sc)
				if !success {
					if err != nil {
						session.Rollback()
						return false, err.Error()
					} else {
						session.Rollback()
						return false, ""
					}
				}

				scUserList := make([]models.ScUser, 0)
				err = session.Where(" sc_id = ? ", scId).Find(&scUserList)
				if err != nil {
					session.Rollback()
					return false, err.Error()
				}

				if sc.ScType == "unique-storage" {
					if len(scUserList) > 1 {
						return false, "unique-storage only assign one user"
					}
				}
				scUser := models.ScUser{UserId: userId, ScId: sc.Id}
				_, err = session.Insert(&scUser)
				if err != nil {
					session.Rollback()
					utils.LoggerError(err)
					return false, err.Error()
				}
			}
		}
	}
	session.Commit()
	return true, ""
}
