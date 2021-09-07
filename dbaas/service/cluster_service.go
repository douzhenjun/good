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
	"DBaas/x/constant"
	"DBaas/x/response"
	"encoding/json"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"math/rand"
	"strconv"
	"strings"
	"time"

	mathEngine "github.com/dengsgo/math-engine/engine"
	"github.com/go-xorm/xorm"
	core1 "k8s.io/api/core/v1"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterService interface {
	PodDetail(clusterId int) (models.ClusterInstance, error)
	List(page int, pageSize int, key string, userId int, userTag string, isDeploy bool) ([]models.ClusterInstance, int64, error)
	Delete(id int, keepPV bool) error
	Patch(replicas int, id int) error
	Add(clusterName string, password string, storageMap map[string]interface{}, parameterMap []map[string]interface{}, remark string, userId int, imageId int, orgTag string, from string, qos *models.Qos, comboId, nodePort int) (int, error)
	Update(id int, dataMap map[string]interface{}) error
	SelectOne(id int) (models.ClusterInstance, string)
	CycleInfo(clusterId int) (*models.CycleInfo, error)
	ParameterEdit(clusterId int, list []map[string]interface{}) error
	ParamList(clusterId, page, pageSize int) ([]models.Clusterparameters, int, error)

	ClusterDisable(clusterId int) error
	ClusterEnable(clusterId int) error
	ApplyConfig(clusterId int) error
	NodePort(port int) (int, error)
}

type clusterService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func (cs *clusterService) NodePort(port int) (int, error) {
	min, max := 30000, 32767
	if port > 0 && (port < min || port > max) {
		return port, fmt.Errorf("node port must be between %v and %v", min, max)
	}
	err, svcAddr := cs.cs.GetResources("service", "", cs.cs.GetNameSpace(), meta1.ListOptions{})
	if err != nil {
		return port, err
	}
	svcList, ok := (*svcAddr).(*core1.ServiceList)
	if !ok {
		return port, errors.New("parse service list failed")
	}

	if port <= 0 {
		for {
			port = rand.Intn(max-min) + min
			exist := nodePortCheck(svcList, int32(port))
			if len(exist) == 0 {
				return port, nil
			}
		}
	}

	exist := nodePortCheck(svcList, int32(port))
	if len(exist) != 0 {
		return port, response.NewMsg(fmt.Sprintf("%v Port is occupied by '%s'", port, exist), fmt.Sprintf("%v端口被'%s'占用", port, exist))
	}
	return port, nil
}

func nodePortCheck(svcList *core1.ServiceList, port int32) string {
	for _, svc := range svcList.Items {
		for _, p := range svc.Spec.Ports {
			if p.NodePort != 0 && p.NodePort == port {
				return svc.Name
			}
		}
	}
	return ""
}

// 处理Mysql参数值
func handleMysqlConf(v string) interface{} {
	// 去除两边的双引号和单引号
	if v[0] == '"' && v[len(v)-1] == '"' || (v[0] == '\'' && v[len(v)-1] == '\'') {
		v = v[1 : len(v)-1]
	}
	// 如果可以转换为Int64, 则转为Int64
	if vi, ok := strconv.ParseInt(v, 10, 64); ok == nil {
		return vi
	}
	return v
}

func (cs *clusterService) ApplyConfig(clusterId int) error {
	if clusterId <= 0 {
		return errors.New("cluster id must > 0")
	}
	dbCluster := models.ClusterInstance{Id: clusterId}
	exist, err := cs.Engine.Cols("limit_mem", "k8s_name", "image_id").Get(&dbCluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", clusterId, err)
	}
	image := models.Images{Id: dbCluster.ImageId}
	exist, err = cs.Engine.Cols("image_type_id").Get(&image)
	if !exist {
		return fmt.Errorf("not found image %v, error: %v", image.Id, err)
	}
	paramList := make([]models.Defaultparameters, 0)
	err = cs.Engine.Where("image_type_id = ?", image.ImageTypeId).Find(&paramList)
	if err != nil {
		return err
	}
	dy, err := dyParam(int64(dbCluster.LimitMem), cs.Engine)
	if err != nil {
		return err
	}
	paramList = append(paramList, dy...)
	paramMap := make(map[string]interface{}, len(paramList))
	for i := range paramList {
		k, v := paramList[i].ParameterName, paramList[i].ParameterValue
		paramMap[k] = handleMysqlConf(v)
	}
	y := constant.MysqlClusterYaml(cs.cs.GetNameSpace(), "")
	cluster, _, err := cs.cs.GetDynamicResource(y, dbCluster.K8sName)
	if err != nil {
		return err
	}
	// test field
	//conf, _, _ := unstructured.NestedFieldNoCopy(cluster.Object, "spec", "mysqlConf")
	//confM := conf.(map[string]interface{})
	//for k, v := range confM {
	//	if paramMap[k] != v {
	//		fmt.Println(k, ":", v, "=>", paramMap[k], reflect.TypeOf(v), "=>", reflect.TypeOf(paramMap[k]))
	//	}
	//}
	//return nil
	err = unstructured.SetNestedField(cluster.Object, paramMap, "spec", "mysqlConf")
	if err != nil {
		return err
	}
	err = cs.cs.UpdateDynamicResource(y, cluster)
	if err != nil {
		return err
	}
	insertParam := make([]models.Clusterparameters, len(paramMap))
	index := 0
	for k, v := range paramMap {
		insertParam[index] = models.Clusterparameters{ParameterName: k, ParameterValue: fmt.Sprintf("%v", v), ClusterId: clusterId}
		index++
	}
	_, _ = cs.Engine.Where("cluster_id = ?", clusterId).Delete(new(models.Clusterparameters))
	_, err = cs.Engine.Insert(&insertParam)
	return err
}

// 根据副本数匹配可用的SC, userId可选
func matchSc(copy, userId int, engine *xorm.Engine) (string, error) {
	if copy <= 0 {
		return "", errors.New("copy must > 0")
	}
	session := engine.Cols("sc_type", "node_num", "name", "sc.id")
	if userId > 0 {
		session.Join("LEFT OUTER", "sc_user", "sc.id = sc_user.sc_id").
			Where("user_id = ?", userId)
	}
	scList := make([]models.Sc, 0)
	err := session.Find(&scList)

	if err != nil {
		return "", err
	}
	for i := range scList {
		scList[i].CheckNodeNum(engine)
		if scList[i].NodeNum >= copy {
			return scList[i].Name, nil
		}
	}
	return "", response.NewMsg("Did not find a matching SC", "未找到符合要求的SC")
}

// CheckNetworkPolicy 检查 NetworkPolicy 策略是否存在，不存在则新建
func (cs *clusterService) checkNetworkPolicy() error {
	y := fmt.Sprintf(`apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: block-all
  namespace: %v
spec:
  podSelector: 
    matchLabels:
      block-network: all
  policyTypes:
  - Ingress
  - Egress`, cs.cs.GetNameSpace())
	_, _, err := cs.cs.GetDynamicResource(y, "block-all")
	if err == nil || !utils.ErrorContains(err, "not found") {
		return err
	}
	_, err = cs.cs.CreateDynamicResource(y)
	return err
}

func (cs *clusterService) ClusterDisable(clusterId int) error {
	if clusterId <= 0 {
		return errors.New("cluster id must > 0")
	}
	if err := cs.checkNetworkPolicy(); err != nil {
		return err
	}
	pods := make([]models.Instance, 0)
	err := cs.Engine.Where("cluster_id = ?", clusterId).Cols("name").Find(&pods)
	if err != nil || len(pods) == 0 {
		return fmt.Errorf("not found pod in the cluster %v， error: %v", clusterId, err)
	}
	for i := range pods {
		patchData := `{"metadata":{"labels":{"block-network":"all"}}}`
		err = cs.cs.PatchOption("pod", pods[i].Name, cs.cs.GetNameSpace(), utils.Str2bytes(patchData), meta1.PatchOptions{}, types.StrategicMergePatchType)
		utils.LoggerError(err)
	}
	_, err = cs.Engine.ID(clusterId).Cols("status").Update(&models.ClusterInstance{Status: models.ClusterStatusDisable})
	return err
}

func (cs *clusterService) ClusterEnable(clusterId int) error {
	if clusterId <= 0 {
		return errors.New("cluster id must > 0")
	}
	pods := make([]models.Instance, 0)
	err := cs.Engine.Where("cluster_id = ?", clusterId).Cols("name").Find(&pods)
	if err != nil || len(pods) == 0 {
		return fmt.Errorf("not found pod in the cluster %v， error: %v", clusterId, err)
	}
	for i := range pods {
		patchData := `{"metadata":{"labels":{"block-network":null}}}`
		err = cs.cs.PatchOption("pod", pods[i].Name, cs.cs.GetNameSpace(), utils.Str2bytes(patchData), meta1.PatchOptions{}, types.StrategicMergePatchType)
		utils.LoggerError(err)
	}
	_, err = cs.Engine.ID(clusterId).Cols("status").Update(&models.ClusterInstance{Status: models.ClusterStatusTrue})
	if err != nil {
		cs.cs.AsyncClusterInfo(clusterId)
	}
	return err
}

func (cs *clusterService) getQos(clusterId int) *models.QosLite {
	qos := models.Qos{ClusterId: clusterId}
	_, _ = cs.Engine.Get(&qos)
	return qos.QosLite
}

func (cs *clusterService) CycleInfo(clusterId int) (*models.CycleInfo, error) {
	if clusterId <= 0 {
		return nil, errors.New("cluster id must > 0")
	}
	var info = new(models.CycleInfo)
	_, err := cs.Engine.SQL("select bt.crontab, bt.keep_copy, bs.name storage_name, bt.set_type, bt.set_date, bt.set_time from backup_task bt inner join backup_storage bs on bt.storage_id = bs.id where bt.close = false and bt.cluster_id = ? and bt.type = ?", clusterId, "cycle").Get(info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (cs *clusterService) ParamList(clusterId, page, pageSize int) ([]models.Clusterparameters, int, error) {
	if clusterId <= 0 {
		return nil, 0, errors.New("cluster id must > 0")
	}
	list := make([]models.Clusterparameters, 0)
	session := cs.Engine.Where("cluster_id = ?", clusterId).Asc("id")
	count, err := pageFind(page, pageSize, &list, session, new(models.Clusterparameters))
	if err != nil {
		return nil, 0, err
	}
	return list, int(count), nil
}

func (cs *clusterService) ParameterEdit(clusterId int, list []map[string]interface{}) error {
	dbCluster := models.ClusterInstance{Id: clusterId}
	exist, err := cs.Engine.Cols("k8s_name").Get(&dbCluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", clusterId, err)
	}
	cluster, _, err := cs.cs.GetDynamicResource(constant.MysqlClusterYaml(cs.cs.GetNameSpace(), ""), dbCluster.K8sName)
	if err != nil {
		return err
	}
	conf, ok := cluster.Object["spec"].(map[string]interface{})["mysqlConf"].(map[string]interface{})
	if !ok {
		return errors.New("parse spec -> mysqlConf is fail")
	}
	session := cs.Engine.NewSession()
	defer session.Close()
	err = session.Begin()
	if err != nil {
		return err
	}
	for i := range list {
		key, value := list[i]["paramKey"].(string), list[i]["paramValue"].(string)
		if _, ok := conf[key]; !ok {
			_ = session.Rollback()
			return fmt.Errorf("not found key %v in mysqlConf", key)
		}
		conf[key] = handleMysqlConf(value)
		_, _ = session.Where("cluster_id = ?", clusterId).And("parameter_name = ?", key).Cols("parameter_value").
			Update(&models.Clusterparameters{ParameterValue: fmt.Sprintf("%v", value)})
	}
	err = cs.cs.UpdateDynamicResource(constant.MysqlClusterYaml(cs.cs.GetNameSpace(), ""), cluster)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}

func NewClusterService(engine *xorm.Engine, cs CommonService) ClusterService {
	return &clusterService{
		Engine: engine,
		cs:     cs,
	}
}

func (cs *clusterService) Patch(replicas int, id int) (err error) {
	cluster := models.ClusterInstance{Id: id}
	_, err = cs.Engine.Omit("yaml_text").Get(&cluster)
	if err != nil {
		return
	}
	if cluster.Status != models.ClusterStatusTrue {
		return errors.New("This Cluster Ready Is Not True ")
	}
	if cluster.Operator != "" {
		return errors.New("This Cluster " + cluster.Operator)
	}
	if cluster.Replicas == strconv.Itoa(replicas) {
		return
	}

	mysqlList := make([]models.Instance, 0)
	err = cs.Engine.Where("cluster_id = ?", cluster.Id).Cols("role", "name").Find(&mysqlList)
	if err != nil {
		return
	}
	masterNum := ""
	for i := range mysqlList {
		if mysqlList[i].Role == "Master" {
			masterNum = mysqlList[i].Name[strings.LastIndex(mysqlList[i].Name, "-")+1:]
			break
		}
	}
	masterNumInt, err := strconv.Atoi(masterNum)
	if err != nil {
		return
	}
	y := constant.MysqlClusterYaml(cs.cs.GetNameSpace(), "")
	updateData, _, err := cs.cs.GetDynamicResource(y, cluster.K8sName)
	if err != nil {
		return err
	}
	cluster.Operator = "Scaling"
	cs.cs.TaskAdd(&cluster, 10, "ActualReplicas", "SetOperator", "GetSelf")
	_, _ = cs.Engine.ID(cluster.Id).Cols("operator").Update(&cluster)
	if replicas < masterNumInt+1 {
		updateData.Object["spec"].(map[string]interface{})["replicas"] = masterNumInt + 1
		err = cs.cs.UpdateDynamicResource(y, updateData)
		utils.LoggerError(err)
		cluster.Replicas = strconv.Itoa(masterNumInt + 1)
		cluster.ComboId = 0
		_, _ = cs.Engine.ID(cluster.Id).Cols("replicas", "combo_id").Update(&cluster)
		cs.cs.AsyncClusterInfo(id)
		go func() {
			success, msg, _, _ := cs.cs.SwitchCluster(id, true)
			<-time.After(time.Second * 40)
			if success {
				_ = cs.Patch(replicas, id)
			} else {
				cs.cs.AddLog("error", "system-cluster", "system", "scale cluster "+cluster.Name+" error: "+msg)
			}
		}()
		return nil
	}
	updateData.Object["spec"].(map[string]interface{})["replicas"] = replicas
	err = cs.cs.UpdateDynamicResource(y, updateData)
	if err != nil {
		return err
	}
	cluster.Replicas = strconv.Itoa(replicas)
	cluster.ComboId = 0
	_, _ = cs.Engine.ID(cluster.Id).Cols("replicas", "combo_id").Update(&cluster)
	cs.cs.AsyncClusterInfo(id)
	return nil
}

/*
获取用户的当前可用的cpu, memory, storage数量 (剩余)
*/
func getUserSurplus(userId int, engine *xorm.Engine) (int, int, int) {
	u := models.User{Id: userId}
	_, err := engine.Cols("cpu_all", "mem_all", "storage_all").Get(&u)
	utils.LoggerError(err)
	cluster := make([]models.ClusterInstance, 0)
	err = engine.Cols("limit_cpu", "limit_mem", "storage").Where("user_id = ?", userId).Find(&cluster)
	utils.LoggerError(err)
	var cpu, mem, storage int
	for i := range cluster {
		cpu += cluster[i].LimitCpu
		mem += cluster[i].LimitMem
		storage += cluster[i].Storage
	}
	return u.CpuAll - cpu, int(u.MemAll) - mem, u.StorageAll - storage
}

func getUserResource(id int, needCpu int, needMem int, needStorage int, engine *xorm.Engine) (bool, string, models.User) {
	user := models.User{Id: id}
	_, err := engine.Get(&user)
	if err != nil {
		return false, err.Error(), user
	}

	cluster := make([]models.ClusterInstance, 0)
	err = engine.Where("user_id = ?", user.Id).Cols("limit_cpu", "limit_mem", "storage").Find(&cluster)
	if err != nil {
		return false, err.Error(), user
	}
	for i := range cluster {
		needCpu += cluster[i].LimitCpu
		needMem += cluster[i].LimitMem
		needStorage += cluster[i].Storage
	}
	if needCpu > user.CpuAll {
		return false, "This User Not Enough CPU", user
	}
	if needMem > int(user.MemAll) {
		return false, "This User Not Enough Memory", user
	}
	if needStorage > user.StorageAll {
		return false, "This User Not Enough Storage", user
	}

	return true, "", user
}

// 动态设置mysql的参数
func dyParam(mem int64, engine *xorm.Engine) ([]models.Defaultparameters, error) {
	mem *= 1024 * 1024 * 1024
	sys := make([]models.Sysparameter, 0)
	err := engine.Cols("param_key", "param_value").Where("sysparameter.param_key like ?", "mp_%").Find(&sys)
	if err != nil {
		return nil, err
	}
	params := make([]models.Defaultparameters, len(sys))
	for i := range sys {
		params[i].ParameterName = sys[i].ParamKey[3:]
		f := strings.ReplaceAll(sys[i].ParamValue, "{m}", strconv.FormatInt(mem, 10))
		res, err := mathEngine.ParseAndExec(f)
		if err != nil {
			return nil, err
		}
		var val string
		// 特殊参数处理
		switch params[i].ParameterName {
		case "innodb_buffer_pool_size":
			// 将单位换算成M
			val = fmt.Sprintf("%vm", int64(res)/1024/1024)
		default:
			val = strconv.FormatInt(int64(res), 10)
		}
		params[i].ParameterValue = val
	}
	return params, nil
}

func (cs *clusterService) Add(clusterName string, password string, storageMap map[string]interface{}, parameterMap []map[string]interface{}, remark string, userId int, imageId int, orgTag string, from string, qos *models.Qos, comboId, nodePort int) (clusterId int, err error) {
	limitMem, limitCpu, storage, replicas := int(storageMap["mem"].(float64)), int(storageMap["cpu"].(float64)), int(storageMap["size"].(float64)), int(storageMap["copy"].(float64))
	enough, msg, u := getUserResource(userId, limitCpu, limitMem, storage, cs.Engine)
	if !enough {
		return 0, errors.New(msg)
	}
	scName, ok := storageMap["scName"].(string)
	if !ok || len(scName) == 0 {
		scName, err = matchSc(replicas, u.Id, cs.Engine)
		if err != nil {
			return
		}
	}
	mysqlImage := models.Images{Id: imageId}
	_, err = cs.Engine.Get(&mysqlImage)
	if err != nil {
		return
	}
	mysqlConf, err := getImageParam(imageId, cs.Engine)
	if err != nil {
		return
	}
	if parameterMap != nil {
		for _, m := range parameterMap {
			for i := range mysqlConf {
				if m["key"] == mysqlConf[i].ParameterName {
					mysqlConf[i].ParameterValue = fmt.Sprintf("%v", m["value"])
					break
				}
			}
		}
	}

	dp, err := dyParam(int64(limitMem), cs.Engine)
	if err != nil {
		return 0, err
	}
	mysqlConf = append(mysqlConf, dp...)

	var mysqlConfString string
	for i := range mysqlConf {
		k, v := mysqlConf[i].ParameterName, mysqlConf[i].ParameterValue
		if i == 0 {
			mysqlConfString = fmt.Sprintf(`%v: %v`, k, v)
		} else {
			mysqlConfString = fmt.Sprintf(`%v
    %v: %v`, mysqlConfString, k, v)
		}
	}

	secretName := "secret-" + strconv.Itoa(userId) + clusterName
	mysqlSecret := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
 namespace: %v
 name: %v #密码对象名称
type: Opaque
data:
 # root password is required to be specified
 ROOT_PASSWORD: %v   #base64编码
 # a name for database that will be created, not required
 # DATABASE: dXNlcmRi #base64编码`, cs.cs.GetNameSpace(), secretName, password)
	_, err = cs.cs.CreateDynamicResource(mysqlSecret)
	utils.LoggerError(err)

	mysqlRepositories := ""
	if mysqlImage.Status == "Invalid" {
		mysqlRepositories = fmt.Sprintf("%v:%v", mysqlImage.ImageName, mysqlImage.Version)
	} else {
		mysqlRepositories = fmt.Sprintf("%v/%v:%v", getImageAddress(cs.Engine), mysqlImage.ImageName, mysqlImage.Version)
	}

	k8sName := clusterName + strconv.Itoa(userId)
	clusterYaml := fmt.Sprintf(`
apiVersion: mysql.presslabs.org/v1alpha1
kind: MysqlCluster
metadata:
  # 集群名字
  namespace: %v
  name: %v
spec:
  # 集群副本数目
  # 注意：0副本表示关闭集群
  replicas: %v
  # 集群密码配置
  secretName: %v
  # 集群mysql镜像配置
  image: %v
  # operator 忽略readonly/readwrite设置，创建/删除/启动/停止 设置为false，启动完成后设置成true
  ignoreReadOnly: false
  # 集群nodeport配置
  masterServiceSpec:
    serviceType: NodePort
    nodePort: %v
  # mysql配置
  mysqlConf: 
    %v
  # pod配置  
  podSpec:
    # pod节点标签选择
    nodeSelector:
      iwhalecloud.dbassnode: mysql
    # pod资源配置
    resources:
      # 所需最小配置
      requests:
        cpu: 1000m
        memory: 1024Mi
      # 最大可使用配置
      limits:
        cpu: %v
        memory: %vGi
  # 存储配置
  volumeSpec:
    persistentVolumeClaim:
      # 存储使用的storageClass名字
      storageClassName: %v
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          # 存储容量
          storage: %vGi 
`, cs.cs.GetNameSpace(), k8sName, replicas, secretName, mysqlRepositories, nodePort, mysqlConfString, limitCpu, limitMem, scName, storage)
	_, err = cs.cs.CreateDynamicResource(clusterYaml)
	if err != nil {
		return
	}

	clusterSession := cs.Engine.NewSession()
	defer clusterSession.Close()
	_ = clusterSession.Begin()
	secret := fmt.Sprintf(`{"ROOT_PASSWORD":"%v"}`, password)
	clusterInstance := models.ClusterInstance{
		Name:          clusterName,
		K8sName:       k8sName,
		ConnectString: "",
		Status:        models.ClusterStatusCreating,
		Storage:       storage,
		SecretName:    secretName,
		UserId:        userId,
		ImageId:       imageId,
		ScName:        scName,
		Replicas:      strconv.Itoa(replicas),
		LimitCpu:      limitCpu,
		LimitMem:      limitMem,
		Remark:        remark,
		YamlText:      clusterYaml,
		OrgTag:        orgTag,
		UserTag:       u.UserTag,
		Secret:        secret,
		ComboId:       comboId,
	}
	_, err = clusterSession.Insert(&clusterInstance)
	if err != nil {
		_ = clusterSession.Rollback()
		return
	}
	if qos != nil {
		qos.ClusterId = clusterInstance.Id
		_, _ = clusterSession.Insert(qos)
	}
	// 把集群的参数入库
	clusterParams := make([]*models.Clusterparameters, len(mysqlConf))
	for i := range mysqlConf {
		clusterParams[i] = &models.Clusterparameters{
			ParameterName:  mysqlConf[i].ParameterName,
			ParameterValue: mysqlConf[i].ParameterValue,
			ClusterId:      clusterInstance.Id,
		}
	}
	_, _ = clusterSession.Insert(&clusterParams)
	err = clusterSession.Commit()
	if err != nil {
		_ = clusterSession.Rollback()
		return
	}
	go statistics.ClusterDeploy(clusterInstance.Name, clusterInstance.Id, replicas, from)
	if replicas > 0 {
		go cs.cs.ScanClusterPod(clusterInstance.Id, k8sName, replicas, false)
	}
	go cs.cs.CreatStatusTimeout(clusterInstance.Id)
	return clusterInstance.Id, nil
}

func (cs *clusterService) Update(id int, dataMap map[string]interface{}) error {
	if id <= 0 {
		return errors.New("cluster id must > 0")
	}
	if len(dataMap) == 0 {
		return nil
	}
	dbCluster := models.ClusterInstance{Id: id}
	exist, err := cs.Engine.Cols("k8s_name", "replicas").Get(&dbCluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", id, err)
	}
	cluster, _, err := cs.cs.GetDynamicResource(constant.MysqlClusterYaml(cs.cs.GetNameSpace(), ""), dbCluster.K8sName)
	if err != nil {
		return err
	}
	if cpu, ok := dataMap["cpu"].(float64); ok && cpu > 0 {
		err = unstructured.SetNestedField(cluster.Object, cpu, "spec", "podSpec", "resources", "limits", "cpu")
		if err != nil {
			return err
		}
		dbCluster.LimitCpu = int(cpu)
	}
	if mem, ok := dataMap["mem"].(float64); ok && mem > 0 {
		err = unstructured.SetNestedField(cluster.Object, fmt.Sprintf("%vGi", mem), "spec", "podSpec", "resources", "limits", "memory")
		if err != nil {
			return err
		}
		dbCluster.LimitMem = int(mem)
	}
	// 副本数
	if c, ok := dataMap["copy"].(float64); ok && c > 0 && strconv.Itoa(int(c)) != dbCluster.Replicas {
		err = unstructured.SetNestedField(cluster.Object, c, "spec", "replicas")
		if err != nil {
			return err
		}
		dbCluster.Replicas = strconv.Itoa(int(c))
	}
	// node port
	if port, ok := dataMap["nodeport"].(float64); ok && port != 0 {
		err = unstructured.SetNestedField(cluster.Object, port, "spec", "masterServiceSpec", "nodePort")
		if err != nil {
			return err
		}
		dbCluster.ConnectString = fmt.Sprintf("%v", port)
	}

	err = cs.cs.UpdateDynamicResource(constant.MysqlClusterYaml(cs.cs.GetNameSpace(), dbCluster.K8sName), cluster)
	if err != nil {
		return err
	}
	dbCluster.Remark, _ = dataMap["remark"].(string)
	comboId, _ := dataMap["comboId"].(float64)
	dbCluster.ComboId = int(comboId)
	_, _ = cs.Engine.ID(id).MustCols("combo_id").Update(&dbCluster)
	if qosM, ok := dataMap["qos"].(map[string]interface{}); ok && len(qosM) != 0 {
		// map转为struct: map->string->struct
		qosStr, _ := json.Marshal(qosM)
		lite := new(models.QosLite)
		err = json.Unmarshal(qosStr, lite)
		if err != nil {
			return err
		}
		qos := models.Qos{ClusterId: id, QosLite: lite}
		exist, err = cs.Engine.Where("cluster_id = ?", id).Exist(new(models.Qos))
		if err != nil {
			return err
		}
		if exist {
			_, err = cs.Engine.Where("cluster_id = ?", id).AllCols().Update(&qos)
		} else {
			_, err = cs.Engine.Insert(&qos)
		}
		if err != nil {
			return err
		}
		pvList := make([]models.PersistentVolume, 0)
		err = cs.Engine.SQL("select pv.name from persistent_volume pv inner join instance ins on pv.pod_id = ins.id where cluster_id = ?", id).Find(&pvList)
		if err != nil {
			return err
		}
		for i := range pvList {
			go cs.cs.SetQosConfig(pvList[i].Name, id)
		}
	}

	return nil
}

func (cs *clusterService) Delete(id int, keepPV bool) (err error) {
	cluster := models.ClusterInstance{Id: id}
	_, err = cs.Engine.Get(&cluster)
	if err != nil {
		return
	}
	if !cluster.IsDeploy {
		mysqlSecret := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
 namespace: %v
 name: %v #密码对象名称
type: Opaque`, cs.cs.GetNameSpace(), cluster.SecretName)
		err = cs.cs.DeleteDynamicResource(mysqlSecret)
		utils.LoggerError(err)
	}

	pvList := make([]models.PersistentVolume, 0)
	// 更改pv策略
	if keepPV {
		err = cs.Engine.
			Cols("persistent_volume.id", "persistent_volume.name", "persistent_volume.pod_id").
			Join("LEFT OUTER", "instance", "persistent_volume.pod_id = instance.id").
			Where("instance.cluster_id = ?", id).Find(&pvList)
		if err != nil {
			return errors.New(fmt.Sprintf("found pv by cluster_id:%v is error: %v", id, err))
		}
		err = cs.cs.ChangePVPolicy(pvList)
		if err != nil {
			return
		}
	}

	if !cluster.IsDeploy {
		err = cs.cs.DeleteDynamicResource(cluster.YamlText)
		if err != nil && cluster.Status != models.ClusterStatusNotFound {
			return
		}
	} else {
		// not found 错误说明资源不存在, 相当于删除成功
		err = cs.cs.GetClientSet().AppsV1().Deployments(cs.cs.GetNameSpace()).Delete(*cs.cs.GetCtx(), cluster.Name, meta1.DeleteOptions{})
		if err != nil && !utils.ErrorContains(err, "not found") {
			return
		}
		err = cs.cs.DeleteOption("service", cluster.Name+"-svc", cs.cs.GetNameSpace(), meta1.DeleteOptions{})
		if err != nil && !utils.ErrorContains(err, "not found") {
			return
		}
		err = cs.cs.DeleteOption("pvc", cluster.Name+"-pvc", cs.cs.GetNameSpace(), meta1.DeleteOptions{})
		if err != nil && !utils.ErrorContains(err, "not found") {
			return
		}
		if !keepPV {
			dbPV := models.PersistentVolume{Id: cluster.PvId}
			_, _ = cs.Engine.Cols("name").Get(&dbPV)
			_ = cs.cs.DeleteOption("pv", dbPV.Name, cs.cs.GetNameSpace(), meta1.DeleteOptions{})
		}
	}

	_, err = cs.Engine.ID(cluster.Id).Delete(&cluster)
	if err != nil {
		return
	}
	_, _ = cs.Engine.Where("cluster_id = ?", cluster.Id).Delete(new(models.Qos))
	_, _ = cs.Engine.Where("cluster_id = ?", cluster.Id).And("type = ?", models.BackupTypeCycle).Cols("close").Update(&models.BackupTask{Close: true})
	go statistics.ClusterDelete(cluster.Name, cluster.Id)
	if keepPV {
		go cs.cs.PollingPVStatus(pvList, core1.VolumeReleased)
	}
	pods := make([]models.Instance, 0)
	_ = cs.Engine.Cols("id", "name").Where("cluster_id = ?", cluster.Id).Find(&pods)
	go cs.cs.ClearEvent(cluster.K8sName)
	for _, instance := range pods {
		_, _ = cs.Engine.Id(instance.Id).Delete(&instance)
		if !keepPV {
			_, _ = cs.Engine.Where("pod_id = ?", instance.Id).Delete(&models.PersistentVolume{})
		}
		go cs.cs.ClearEvent(instance.Name)
	}
	return nil
}

func (cs *clusterService) List(page int, pageSize int, key string, userId int, userTag string, isDeploy bool) (clusterList []models.ClusterInstance, count int64, err error) {
	clusterList = make([]models.ClusterInstance, 0)
	like := "%" + key + "%"
	where := "(k8s_name like ? OR remark like ?) AND is_deploy = ? "
	args := []interface{}{like, like, isDeploy}
	if userId > 0 {
		where += "AND user_id = ?"
		args = append(args, userId)
	} else if userTag != "AAAA" {
		where += "And user_tag = ?"
		args = append(args, userTag)
	}
	err = cs.Engine.Where(where, args...).Omit("yaml_text").Limit(pageSize, pageSize*(page-1)).Desc("id").Find(&clusterList)
	if err != nil {
		return
	}
	count, _ = cs.Engine.Where(where, args...).Count(&models.ClusterInstance{})

	for i := range clusterList {
		userModel := models.User{Id: clusterList[i].UserId}
		_, _ = cs.Engine.Get(&userModel)
		clusterList[i].UserName = userModel.UserName
		imageModel := models.Images{Id: clusterList[i].ImageId}
		_, _ = cs.Engine.Get(&imageModel)
		clusterList[i].ImageName = imageModel.ImageName + ":" + imageModel.Version
		clusterList[i].PodStatusMap = utils.RawJson(clusterList[i].PodStatus)
		events := cs.cs.GetEvent(clusterList[i].K8sName)
		clusterList[i].Events = events
		clusterList[i].Qos = cs.getQos(clusterList[i].Id)
		if clusterList[i].IsDeploy {
			pv := models.PersistentVolume{Id: clusterList[i].PvId}
			_, _ = cs.Engine.Cols("name", "pod_id", "sc_id").Get(&pv)
			clusterList[i].PvName = pv.Name
			pod := models.Instance{Id: pv.PodId}
			_, _ = cs.Engine.Cols("name").Get(&pod)
			clusterList[i].PodName = pod.Name
		} else {
			sc := models.Sc{Name: clusterList[i].ScName}
			_, _ = cs.Engine.Get(&sc)
			scNodes := sc.NodeNum
			if sc.ScType == "unique-storage" {
				pvCount, _ := cs.Engine.Where("sc_id = ?", sc.Id).Count(&models.PersistentVolume{})
				scNodes = int(pvCount)
			}
			clusterList[i].ScNodes = scNodes
			clusterList[i].ScType = sc.ScType
			clusterList[i].CycleInfo, err = cs.CycleInfo(clusterList[i].Id)
			utils.LoggerError(err)
		}
	}
	return
}

func (cs *clusterService) PodDetail(clusterId int) (models.ClusterInstance, error) {
	if clusterId <= 0 {
		return models.ClusterInstance{}, errors.New("cluster id must > 0")
	}
	cluster := models.ClusterInstance{Id: clusterId}
	_, err := cs.Engine.Omit("yaml_text").Get(&cluster)
	if err != nil {
		return cluster, err
	}
	u := models.User{Id: cluster.UserId}
	_, _ = cs.Engine.Cols("user_name").Get(&u)
	cluster.UserName = u.UserName
	i := models.Images{Id: cluster.ImageId}
	_, _ = cs.Engine.Cols("image_name").Get(&i)
	cluster.ImageName = i.ImageName
	mysql := make([]models.Instance, 0)
	err = cs.Engine.Where(" cluster_id = ?", clusterId).Desc("id").Find(&mysql)
	if err != nil {
		return cluster, err
	}
	//  数据库密码
	cluster.SecretMap = utils.RawJson(cluster.Secret)
	for i := range mysql {
		mysql[i].VolumeObject = utils.RawJson(mysql[i].Volume)
		mysql[i].InitContainerObject = utils.RawJson(mysql[i].InitContainer)
		mysql[i].ContainerInfoObject = utils.RawJson(mysql[i].ContainerInfo)
		mysql[i].BaseInfoObject = utils.RawJson(mysql[i].BaseInfo)
		mysql[i].Events = cs.cs.GetEvent(mysql[i].Name)
	}
	sc := models.Sc{Name: cluster.ScName}
	_, err = cs.Engine.Get(&sc)
	utils.LoggerError(err)
	if sc.ScType == "unique-storage" {
		count, _ := cs.Engine.Where("sc_id = ?", sc.Id).Count(new(models.PersistentVolume))
		cluster.ScNodes = int(count)
	} else {
		cluster.ScNodes = sc.NodeNum
	}
	cluster.Instance = mysql
	cluster.Qos = cs.getQos(clusterId)
	if !cluster.IsDeploy {
		cluster.CycleInfo, err = cs.CycleInfo(clusterId)
		utils.LoggerError(err)
	}
	cluster.Events = cs.cs.GetEvent(cluster.K8sName)
	return cluster, nil
}

func (cs *clusterService) SelectOne(id int) (models.ClusterInstance, string) {
	var cluster models.ClusterInstance
	_, err := cs.Engine.Where(" id = ? ", id).Get(&cluster)
	if err != nil {
		utils.LoggerError(err)
		return cluster, err.Error()
	}
	return cluster, ""
}
