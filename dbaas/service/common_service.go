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
	"DBaas/datasource"
	"DBaas/models"
	"DBaas/utils"
	"DBaas/x/constant"
	"bufio"
	k8sContext "context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	influxClient "github.com/influxdata/influxdb/client/v2"
	"io"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	core1 "k8s.io/api/core/v1"
	storage1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CommonService interface {
	//k8s  api
	// 自定义扩展
	GetClientSet() *kubernetes.Clientset
	GetConfig() *rest.Config
	GetCtx() *k8sContext.Context
	// 静态资源crud
	GetResources(sourceType string, sourceName string, nameSpace string, opts interface{}) (error, *interface{})
	PatchOption(sourceType string, sourceName string, nameSpace string, playLoadBytes []byte, opts meta1.PatchOptions, patchType types.PatchType) error
	CreateOption(sourceType string, nameSpace string, sourceInterface interface{}, opts meta1.CreateOptions) error
	DeleteOption(sourceType string, sourceName string, nameSpace string, opts meta1.DeleteOptions) error
	// 动态资源crud
	CreateDynamicResource(deploymentYaml string) (*unstructured.Unstructured, error)
	DeleteDynamicResource(deploymentYaml string) error
	PatchDynamicResource(deploymentYaml string, controllerContent string) error
	UpdateDynamicResource(deploymentYaml string, updateData *unstructured.Unstructured) error
	GetDynamicResource(deploymentYaml string, sourceName string) (*unstructured.Unstructured, *unstructured.UnstructuredList, error)
	GetPodLogs(podName, container string) ([]string, error)
	GetLogsByLoki(pod, container string) (map[string]interface{}, error)
	GetNameSpace() string
	SetNameSpace()
	SetK8sConfig(config *rest.Config, clientSet *kubernetes.Clientset, ctx *k8sContext.Context, err error)
	ScanClusterPod(clusterId int, sourceName string, replicas int, isDeploy bool)
	// 同步PV状态
	AsyncPVStatus()
	PollingPVStatus(pvList []models.PersistentVolume, target core1.PersistentVolumePhase)
	FindPV(pod *core1.Pod)
	ChangePVPolicy(pvList []models.PersistentVolume) error
	AsyncNodeInfo()
	AsyncMetricsPods()
	AsyncOperatorLog()
	AsyncClusterInfo(id int)
	AsyncCommonInfo()
	AsyncImageStatus()
	ClientOperatorButton()
	ClearEvent(name string)
	// 清除无用的event
	ClearUselessEvents()
	TaskAdd(task interface{}, numTime int, key string, setMethodName string, getMethodName string)
	// 添加日志
	AddLog(level string, logSource string, people string, content string)

	// 公共模块
	SwitchCluster(id int, isAuto bool) (bool, string, string, models.ClusterInstance)
	FilterEvent(filter string) ([]models.PodLog, int64)
	GetEvent(name string) []models.PodLog
	CreatStatusTimeout(clusterId int)
	AsyncBackupJob()
	SetQosConfig(pvName string, clusterId int)
}

func NewCommonService(engine *xorm.Engine, config *rest.Config, clientSet *kubernetes.Clientset, ctx *k8sContext.Context, err error, logService LogService, influxdbClient influxClient.Client) CommonService {
	taskList := make([]map[string]interface{}, 0)
	clusterList := make([]models.ClusterInstance, 0)
	_ = engine.Where("operator != ?", "").Find(&clusterList)
	for _, instance := range clusterList {
		instance.Operator = ""
		_, _ = engine.Id(instance.Id).Cols("operator").Update(&instance)
	}
	common := commonService{
		Engine: engine,
		k8sClient: &k8sClient{
			Config:    config,
			ClientSet: clientSet,
			Ctx:       ctx,
			Err:       err,
		},
		CommonNameSpace: "default",
		LogService:      logService,
		InfluxdbClient:  influxdbClient,
		taskList:        taskList,
	}
	go common.ClientOperatorButton()
	go common.setClientTR()
	return &common
}

func (cs *commonService) setClientTR() {
	if cs.Config == nil || cs.clientTR != nil {
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(cs.Config.CAData)
	cliCrt, err := tls.X509KeyPair(cs.Config.CertData, cs.Config.KeyData)
	if err != nil {
		return
	}
	dialer := &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: time.Minute,
	}
	tr := &http.Transport{
		DialContext:     dialer.DialContext,
		IdleConnTimeout: time.Minute,
		MaxIdleConns:    60,
		MaxConnsPerHost: 20,
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	cs.clientTR = &http.Client{Transport: tr, Timeout: 10 * time.Second}
}

type commonService struct {
	Engine *xorm.Engine
	*k8sClient
	CommonNameSpace string
	LogService      LogService
	InfluxdbClient  influxClient.Client
	taskList        []map[string]interface{}
	clientTR        *http.Client
	stopListen      chan struct{}
}

type k8sClient struct {
	Config    *rest.Config
	ClientSet *kubernetes.Clientset
	Ctx       *k8sContext.Context
	Err       error
}

func (cs *commonService) SetQosConfig(pvName string, clusterId int) {
	var err error
	defer utils.LoggerErrorP(&err)
	qos := models.Qos{ClusterId: clusterId}
	exist, err := cs.Engine.Get(&qos)
	if !exist {
		return
	}
	setAction := func() bool {
		svc, conn, err := NewDiseService(cs.Engine)
		if err != nil {
			return false
		}
		defer CloseGrpc(conn)
		err = svc.SetQoSOfVolume(pvName, qos.QosLite)
		return err == nil
	}
	// 先执行一次设置动作，如果成功则退出，不成功则每10s执行1次，总执行20次
	if setAction() {
		return
	}
	utils.Polling(setAction, 10*time.Second, 20)
}

func (cs *commonService) GetLogsByLoki(pod, container string) (map[string]interface{}, error) {
	host := models.Sysparameter{ParamKey: "kubernetes_master_address"}
	exist, err := cs.Engine.Cols("param_value").Get(&host)
	if !exist {
		return nil, fmt.Errorf("not found kubernetes_master_address, error: %v", err)
	}
	end := time.Now()
	// 查询现在到7天前的日志
	start := end.AddDate(0, 0, -7)
	url := fmt.Sprintf(`http://%s:31000/loki/api/v1/query_range?query={pod="%s",container="%s"}&start=%v&end=%v`, host.ParamValue, pod, container, start.UnixNano(), end.UnixNano())
	data, err := utils.SimpleGet(url)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	if status, ok := m["status"]; !ok || status != "success" {
		return nil, errors.New("response status is unsuccessful")
	}
	ret, ok := m["data"].(map[string]interface{})
	if !ok {
		return nil, errors.New("there is no 'data' field in the response message")
	}
	return ret, nil
}

/*
AsyncBackupJob 同步备份作业
*/
func (cs *commonService) AsyncBackupJob() {
	var err error
	defer utils.LoggerErrorP(&err)
	dbTask := make([]map[string]interface{}, 0)
	err = cs.Engine.SQL("select bt.id, bt.name, ci.k8s_name cname from backup_task bt inner join cluster_instance ci on bt.cluster_id = ci.id").Find(&dbTask)
	if err != nil || len(dbTask) == 0 {
		return
	}

	_, backupList, err := cs.GetDynamicResource(constant.MysqlBackupYaml(cs.GetNameSpace(), ""), "")
	if err != nil || backupList == nil || len(backupList.Items) == 0 {
		return
	}

	dbTaskM := map[string]int64{}
	for i := range dbTask {
		taskId := dbTask[i]["id"].(int64)
		if n, ok := dbTask[i]["name"].(string); ok && n != "" {
			// name有值说明是手动备份
			dbTaskM[n] = taskId
			continue
		}
		dbTaskM[dbTask[i]["cname"].(string)] = taskId
	}

	dbJob := make([]models.BackupJob, 0)
	_ = cs.Engine.Find(&dbJob)
	dbJobM := map[string]int{}
	for i := range dbJob {
		dbJobM[dbJob[i].JobName] = i
	}

	err, podListSource := cs.GetResources("pod", "", cs.GetNameSpace(), meta1.ListOptions{LabelSelector: "job-name"})
	podList, ok := (*podListSource).(*core1.PodList)
	if err != nil || !ok || len(podList.Items) == 0 {
		return
	}
	for _, obj := range backupList.Items {
		var jobName = obj.GetName()
		var job models.BackupJob
		backupURL, ok := obj.Object["spec"].(map[string]interface{})["backupURL"].(string)
		if ok && backupURL != "" {
			job.BackupSet = backupURL[strings.LastIndex(backupURL, "/")+1:]
			backupPod := cs.getBackupPod(podList, jobName)
			if backupPod != nil {
				job.PodName = backupPod.Name
				for _, c := range backupPod.Status.ContainerStatuses {
					if c.Name != "backup" {
						continue
					}
					switch {
					case c.State.Waiting != nil:
						job.Status = c.State.Waiting.Reason
					case c.State.Running != nil:
						job.Status = "Running"
					case c.State.Terminated != nil:
						var diff = c.State.Terminated.FinishedAt.Sub(c.State.Terminated.StartedAt.Time)
						job.Duration = int(diff.Seconds())
						job.Status = c.State.Terminated.Reason
					}
					break
				}
			}
		}

		if i, ok := dbJobM[jobName]; ok {
			dbJobM[jobName] = -1 // 标记Job已插入，保持数据库和集群中的Job同步
			oldBackupSet := dbJob[i].BackupSet
			oldStatus := dbJob[i].Status
			oldDuration := dbJob[i].Duration
			oldPodName := dbJob[i].PodName
			if oldBackupSet != job.BackupSet || oldStatus != job.Status || oldDuration != job.Duration || oldPodName != job.PodName {
				_, err = cs.Engine.ID(dbJob[i].Id).Update(&job)
			}
		} else {
			job.JobName = jobName
			job.CreateTime = obj.GetCreationTimestamp().Time
			labels := obj.GetLabels()
			cycle := labels["recurrent"] == "true"
			var findTask = jobName
			if cycle {
				findTask = labels["cluster"]
			}
			if taskId, ok := dbTaskM[findTask]; ok {
				job.BackupTaskId = int(taskId)
			} else {
				utils.LoggerInfo("task id was not found when synchronizing the backup job")
				continue
			}
			_, err = cs.Engine.Insert(&job)
		}
	}

	for _, v := range dbJobM {
		if v != -1 {
			_, _ = cs.Engine.ID(dbJob[v].Id).Delete(new(models.BackupJob))
		}
	}
}

func (cs *commonService) getBackupPod(pods *core1.PodList, jobName string) *core1.Pod {
	exist, i := utils.SliceExist(len(pods.Items), func(i int) bool {
		label, ok := pods.Items[i].Labels["job-name"]
		return ok && label == fmt.Sprintf("%v-backup", jobName)
	})
	if !exist {
		return nil
	}
	return &pods.Items[i]
}

/*
设置集群Creating状态超时
*/
func (cs *commonService) CreatStatusTimeout(clusterId int) {
	<-time.After(time.Minute * 10)
	cluster := models.ClusterInstance{}
	find, _ := cs.Engine.ID(clusterId).Cols("name", "status").Get(&cluster)
	if !find || cluster.Status != models.ClusterStatusCreating {
		return
	}
	cluster.Status = models.ClusterStatusFalse
	_, _ = cs.Engine.ID(clusterId).Cols("status").Update(&cluster)
	statistics.ClusterTimeout(cluster.Name, clusterId)
}

func (cs *commonService) GetEvent(name string) []models.PodLog {
	events := make([]models.PodLog, 0)
	_ = cs.Engine.Where("name like ?", name+".%").Desc("oper_date").Find(&events)
	for i := range events {
		events[i].DateFormat()
	}
	return events
}

func (cs *commonService) FilterEvent(filter string) (events []models.PodLog, count int64) {
	var err error
	defer utils.LoggerErrorP(&err)
	events = make([]models.PodLog, 0)
	err = cs.Engine.Where("name like ?", filter+"%").Desc("oper_date").Find(&events)
	if err != nil {
		return
	}
	count, _ = cs.Engine.Where("name like ?", filter+"%").Count(new(models.PodLog))
	for i := range events {
		events[i].DateFormat()
	}
	return
}

func (cs *commonService) resetListen() {
	if cs.stopListen != nil {
		close(cs.stopListen)
	}
	cs.listenEvent()
}

func (cs *commonService) listenEvent() {
	if cs.Err != nil {
		return
	}
	watchlist := cache.NewListWatchFromClient(cs.ClientSet.CoreV1().RESTClient(), "events", cs.GetNameSpace(), fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&core1.Event{},
		time.Hour * 12,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if event, ok := obj.(*core1.Event); ok {
					if _, ok := eventSelector[event.InvolvedObject.Kind]; !ok {
						return
					}
					dbEvent := models.PodLog{Name: event.Name}
					exist, err := cs.Engine.Exist(&dbEvent)
					if !exist && err == nil {
						dbEvent.Type = event.Type
						dbEvent.Reason = event.Reason
						dbEvent.OperDate = event.FirstTimestamp.Time
						dbEvent.Message = event.Message
						dbEvent.From = fmt.Sprintf("%v, %v", event.Source.Component, event.Source.Host)
						_, _ = cs.Engine.Insert(&dbEvent)
					}
				}
			},
			//DeleteFunc: func(obj interface{}) {
			//	if event, ok := obj.(*core1.Event); ok {
			//		dbEvent := models.PodLog{Name: event.Name}
			//		_, _ = cs.Engine.Delete(&dbEvent)
			//	}
			//},
		},
	)
	stop := make(chan struct{})
	cs.stopListen = stop
	go controller.Run(stop)
}

func (cs *commonService) TaskAdd(task interface{}, numTime int, key string, setMethodName string, getMethodName string) {
	cs.taskList = append(cs.taskList, map[string]interface{}{
		"task":          task,
		"numTime":       numTime,
		"key":           key,
		"setMethodName": setMethodName,
		"getMethodName": getMethodName,
	})
}

func (cs *commonService) ClientOperatorButton() {
	for i, m := range cs.taskList {
		taskModel := reflect.ValueOf(m["task"])
		if m["numTime"].(int) <= 0 {
			// 数据库更新
			taskMethod := taskModel.MethodByName(m["setMethodName"].(string))
			params := make([]reflect.Value, 2)
			params[0] = reflect.ValueOf("")
			params[1] = reflect.ValueOf(cs.Engine)
			taskMethod.Call(params)
			// 释放
			if len(cs.taskList) <= 1 {
				cs.taskList = append(cs.taskList[0:0])
			} else {
				cs.taskList = append(cs.taskList[0:i], cs.taskList[i+1:]...)
			}
			go cs.AsyncClusterInfo(-1)
			continue
		}
		paramsGet := make([]reflect.Value, 1)
		paramsGet[0] = reflect.ValueOf(cs.Engine)
		currentModel := taskModel.MethodByName(m["getMethodName"].(string)).Call(paramsGet)
		m["numTime"] = m["numTime"].(int) - 1
		if fmt.Sprintf("%s", currentModel[0].Elem().FieldByName(m["key"].(string))) != fmt.Sprintf("%s", taskModel.Elem().FieldByName(m["key"].(string))) {
			m["numTime"] = 0
		}
	}
	<-time.After(time.Second)
	cs.ClientOperatorButton()
}

func (cs *commonService) ClearUselessEvents() {
	// 将需要保留的event加入map中
	usefulList := make(map[string]struct{})
	clusterList := make([]models.ClusterInstance, 0)
	_ = cs.Engine.Cols("k8s_name").Find(&clusterList)
	for i := range clusterList {
		usefulList[clusterList[i].K8sName] = struct{}{}
	}
	podList := make([]models.Instance, 0)
	_ = cs.Engine.Cols("name").Find(&podList)
	for i := range podList {
		usefulList[podList[i].Name] = struct{}{}
	}
	operatorList := make([]models.MysqlOperator, 0)
	_ = cs.Engine.Cols("name").Find(&operatorList)
	for i := range operatorList {
		usefulList[operatorList[i].Name] = struct{}{}
	}
	backupJobList := make([]models.BackupJob, 0)
	_ = cs.Engine.Cols("pod_name").Find(&backupJobList)
	for i := range backupJobList {
		usefulList[backupJobList[i].PodName] = struct{}{}
	}

	// 遍历event删除无用事件
	eventList := make([]models.PodLog, 0)
	_ = cs.Engine.Cols("id", "name").Find(&eventList)
	for i := range eventList {
		eventName := eventList[i].Name
		key := eventName[:strings.LastIndex(eventName, ".")]
		if _, has := usefulList[key]; has {
			continue
		}
		_, _ = cs.Engine.ID(eventList[i].Id).Delete(&eventList[i])
	}
}

func (cs *commonService) ClearEvent(name string) {
	events := make([]models.PodLog, 0)
	_ = cs.Engine.Where("name like ?", name+".%").Cols("id").Find(&events)
	for i := range events {
		_, _ = cs.Engine.ID(events[i].Id).Delete(&events[i])
	}
}

func (cs *commonService) SwitchCluster(id int, isAuto bool) (bool, string, string, models.ClusterInstance) {
	cluster := models.ClusterInstance{Id: id}
	_, err := cs.Engine.Omit("yaml_text").Get(&cluster)
	if err != nil {
		return false, utils.ERROR_SWITCH_EN, utils.ERROR_SWITCH_ZH, cluster
	}
	if !isAuto {
		if cluster.Operator != "" {
			return false, "This Cluster " + cluster.Operator, "This Cluster " + cluster.Operator, cluster
		} else {
			cluster.Operator = "Switching"
			cs.TaskAdd(&cluster, 10, "master", "SetOperator", "GetSelf")
		}
	}
	operatorList := make([]models.MysqlOperator, 0)
	err = cs.Engine.Find(&operatorList)
	if err != nil || len(operatorList) == 0 {
		return false, "No operator information was found", "未查询到operator信息", cluster
	}
	url := fmt.Sprintf(`http://%s:80/api/graceful-master-takeover-auto/%s.default`, operatorList[0].ServiceIP, cluster.K8sName)
	returnData, err := utils.SimpleGet(url)
	if err != nil {
		return false, err.Error(), err.Error(), cluster
	}
	returnMap := make(map[string]interface{})
	err = json.Unmarshal(returnData, &returnMap)
	if err != nil {
		utils.LoggerError(err)
		return false, utils.ERROR_DATA_TRANSFER_EN, utils.ERROR_DATA_TRANSFER_ZH, cluster
	}
	if _, ok := returnMap["Code"]; ok {
		if returnMap["Code"] == "ERROR" {
			if _, ok := returnMap["Message"]; ok {
				return false, fmt.Sprintf("%v", returnMap["Message"]), fmt.Sprintf("%v", returnMap["Message"]), cluster
			}
		}
	}
	_, err = cs.Engine.ID(cluster.Id).Update(&cluster)
	return true, "", "", cluster
}

func (cs *commonService) SetNameSpace() {
	k8sNamespace := models.Sysparameter{ParamKey: "kubernetes_namespace"}
	success, err := cs.Engine.Get(&k8sNamespace)
	namespace := &core1.Namespace{
		ObjectMeta: meta1.ObjectMeta{
			Name: k8sNamespace.ParamValue,
		},
	}
	if cs.Err == nil {
		_, _ = cs.ClientSet.CoreV1().Namespaces().Create(*cs.Ctx, namespace, meta1.CreateOptions{})
	}
	if success && err == nil {
		cs.CommonNameSpace = k8sNamespace.ParamValue
		cs.resetListen()
	}
}

func (cs *commonService) SetK8sConfig(config *rest.Config, clientSet *kubernetes.Clientset, ctx *k8sContext.Context, err error) {
	cs.Err = err
	cs.ClientSet = clientSet
	cs.Config = config
	cs.Ctx = ctx
	cs.setClientTR()
	cs.resetListen()
}

func (cs *commonService) operatorLog(flag int) {
	operatorStatus := true
	var notFoundOperator []string
	for i := 0; i < flag; i++ {
		// 部署成功之后写入数据库
		podName := fmt.Sprintf("mysql-operator-%v", i)
		err, podsAddr := cs.GetResources("pod", podName, cs.GetNameSpace(), meta1.GetOptions{})
		if err != nil {
			if i == 0 {
				operatorStatus = false
			}
			utils.LoggerError(err)
			operator := models.MysqlOperator{Name: podName}
			exist, _ := cs.Engine.Get(&operator)
			if exist {
				operator.Status = "NotFound"
				operator.Ready = "False"
				operator.ContainerStatus = "0/0"
				_, _ = cs.Engine.ID(operator.Id).Update(&operator)
				notFoundOperator = append(notFoundOperator, podName)
			}
			continue
		}
		pod := (*podsAddr).(*core1.Pod)
		if pod.Name == "" {
			operatorStatus = false
			continue
		}
		var serviceIP string
		err, svcAddr := cs.GetResources("service", "mysql-operator", cs.GetNameSpace(), meta1.GetOptions{})
		utils.LoggerError(err)
		if err == nil {
			svc := (*svcAddr).(*core1.Service)
			serviceIP = svc.Spec.ClusterIP
		}

		ready := "False"
		sumContainer := 0
		for _, condition := range pod.Status.ContainerStatuses {
			if condition.Ready {
				sumContainer += 1
			}
		}
		containerStatus := fmt.Sprintf("%v/%v", sumContainer, len(pod.Status.ContainerStatuses))
		if sumContainer == len(pod.Status.ContainerStatuses) && string(pod.Status.Phase) == "Running" {
			ready = "True"
		} else {
			operatorStatus = false
		}

		operator := models.MysqlOperator{Name: pod.Name}
		exist, _ := cs.Engine.Get(&operator)
		if !exist {
			operator.DeploymentStatus = "Deploying"
		}
		if string(pod.Status.Phase) == "Running" {
			operator.DeploymentStatus = "Deployment complete"
		}
		operator.ContainerStatus = containerStatus
		operator.Status = string(pod.Status.Phase)
		operator.Ready = ready
		operator.NodeName = pod.Spec.NodeName
		operator.Ip = pod.Status.PodIP
		operator.Replicas = strconv.Itoa(flag)
		operator.ServiceIP = serviceIP
		if exist {
			_, err = cs.Engine.ID(operator.Id).Update(&operator)
			utils.LoggerError(err)
		} else {
			_, err = cs.Engine.Insert(&operator)
			utils.LoggerError(err)
		}
	}

	if operatorStatus {
		creating, _ := models.GetConfigBool("operator@creating", cs.Engine)
		if creating {
			_ = models.SetConfigBool("operator@creating", false, cs.Engine)
		}
	}
	if notFoundOperator != nil {
		_ = models.SetConfig("operator@reason", fmt.Sprint("Not found operator:", strings.Join(notFoundOperator, ", ")), cs.Engine)
	}
	dbStatus, _ := models.GetConfigBool("operator@status", cs.Engine)
	if dbStatus != operatorStatus {
		_ = models.SetConfigBool("operator@status", operatorStatus, cs.Engine)
	}
}

/*
扫描集群中pod并入库
*/
func (cs *commonService) ScanClusterPod(clusterId int, sourceName string, replicas int, isDeploy bool) {
	select {
	case <-time.After(2 * time.Second):
		cluster := models.ClusterInstance{Id: clusterId}
		hasCluster, _ := cs.Engine.Exist(&cluster)
		// 集群有可能被删除
		if !hasCluster {
			return
		}
		syncCount := 0
		for i := 0; i < replicas; i++ {
			mysqlPod := models.Instance{}
			if isDeploy {
				mysqlPod.Name = sourceName
			} else {
				mysqlPod.Name = fmt.Sprintf("%v-mysql-%v", sourceName, i)
			}
			err, podsAddr := cs.GetResources("pod", mysqlPod.Name, cs.GetNameSpace(), meta1.GetOptions{})
			if err != nil {
				utils.LoggerError(err)
				continue
			} else {
				exist, err := cs.Engine.Exist(&mysqlPod)
				if err != nil {
					utils.LoggerError(err)
					continue
				}
				if exist {
					syncCount++
					continue
				}
				pods := (*podsAddr).(*core1.Pod)
				status := string(pods.Status.Phase)
				if status == "false" {
					status = "Error"
				}
				mysqlPod.ClusterId = clusterId
				mysqlPod.Name = pods.Name
				mysqlPod.Version = pods.ResourceVersion
				mysqlPod.Status = status
				_, err = cs.Engine.Insert(&mysqlPod)
				if err != nil {
					utils.LoggerError(err)
					continue
				}
				syncCount++
				if !isDeploy {
					go cs.FindPV(pods)
				} else {
					cluster := models.ClusterInstance{Id: clusterId}
					_, err = cs.Engine.Cols("pv_id").Get(&cluster)
					if err == nil {
						pv := models.PersistentVolume{Id: cluster.PvId, PodId: mysqlPod.Id}
						_, _ = cs.Engine.Id(pv.Id).Cols("pod_id").Update(&pv)
					}
				}
			}
		}
		if syncCount == replicas {
			go cs.AsyncClusterInfo(clusterId)
		}
	}
}

func (cs *commonService) AsyncMetricsPods() {
	if cs.Err != nil {
		return
	}
	influxTags := map[string]string{"metrics": "pod"}
	bp, err := datasource.NewBatchPoints()
	if err != nil {
		utils.LoggerError(err)
		return
	}
	mysqlInstance := make([]models.Instance, 0)
	_ = cs.Engine.Find(&mysqlInstance)
	wg := sync.WaitGroup{}
	wg.Add(len(mysqlInstance))
	for _, instance := range mysqlInstance {
		go func(instance models.Instance) {
			defer wg.Done()
			//生成要访问的url
			url := cs.Config.Host + "/apis/metrics.k8s.io/v1beta1/namespaces/" + cs.CommonNameSpace + "/pods/" + instance.Name
			//提交请求
			request, err := http.NewRequest("GET", url, nil)
			utils.LoggerError(err)
			//处理返回结果
			response, err := cs.clientTR.Do(request)
			if err != nil {
				utils.LoggerError(err)
				return
			}
			defer response.Body.Close()
			//返回的状态码
			status := response.StatusCode
			if status != 200 {
				return
			}
			result, err := ioutil.ReadAll(response.Body)
			resultMap := map[string]interface{}{}
			err = json.Unmarshal(result, &resultMap)
			if err != nil {
				return
			}
			metricsPod := map[string]float64{}
			metricsPod["cpu"] = float64(0)
			metricsPod["mem"] = float64(0)
			for _, container := range resultMap["containers"].([]interface{}) {
				usage := container.(map[string]interface{})["usage"].(map[string]interface{})
				cpuStr := usage["cpu"].(string)
				// 去除结尾的n
				cpuUse, _ := strconv.ParseInt(cpuStr[:len(cpuStr)-1], 10, 64)
				if cpuUse > 0 {
					metricsPod["cpu"] += float64(cpuUse)
				}
				memStr := usage["memory"].(string)
				var memUse int64
				if len(memStr) > 1 {
					// 去除结尾的Ki
					memUse, _ = strconv.ParseInt(memStr[:len(memStr)-2], 10, 64)
				}
				if memUse > 0 {
					metricsPod["mem"] += float64(memUse)
				}
			}
			metricsPod["mem"] /= 1048576 // 1024 * 1024
			metricsPod["cpu"] /= 1000000000
			for k, v := range metricsPod {
				f := map[string]interface{}{
					k:     v,
					"pod": instance.Name,
				}
				pt, err := influxClient.NewPoint("metrics_pod_"+k, influxTags, f, time.Now())
				if err != nil {
					utils.LoggerError(err)
					continue
				}
				bp.AddPoint(pt)
			}
		}(instance)
	}
	wg.Wait()
	err = cs.InfluxdbClient.Write(bp)
	utils.LoggerError(err)
}

//  operator 信息入库
func (cs *commonService) AsyncOperatorLog() {
	replicas, err := models.GetConfigInt("operator@replicas", cs.Engine)
	if err != nil {
		return
	}
	cs.operatorLog(replicas)
}

/*
定时同步PV状态
*/
func (cs *commonService) AsyncPVStatus() {
	dbPVList := make([]models.PersistentVolume, 0)
	err := cs.Engine.Find(&dbPVList)
	if err != nil {
		return
	}
	for i := range dbPVList {
		err, pvSource := cs.GetResources("pv", dbPVList[i].Name, cs.GetNameSpace(), meta1.GetOptions{})
		if err != nil {
			if utils.ErrorContains(err, "not found") {
				_, _ = cs.Engine.ID(dbPVList[i].Id).Delete(&dbPVList[i])
			}
			continue
		}
		pv := (*pvSource).(*core1.PersistentVolume)
		oldStatus, oldPolicy := dbPVList[i].Status, dbPVList[i].ReclaimPolicy
		newStatus, newPolicy := string(pv.Status.Phase), string(pv.Spec.PersistentVolumeReclaimPolicy)
		if oldStatus == newStatus && oldPolicy == newPolicy {
			continue
		}
		dbPVList[i].Status = newStatus
		dbPVList[i].ReclaimPolicy = newPolicy
		_, _ = cs.Engine.Id(dbPVList[i].Id).Cols("status", "reclaim_policy").Update(&dbPVList[i])
	}
}

/*
短时间间隔同步PV状态为目标值
*/
func (cs *commonService) PollingPVStatus(pvList []models.PersistentVolume, target core1.PersistentVolumePhase) {
	utils.Polling(func() bool {
		var completeCount int
		for i := range pvList {
			err, pvSource := cs.GetResources("pv", pvList[i].Name, cs.GetNameSpace(), meta1.GetOptions{})
			if err != nil {
				continue
			}
			pv := (*pvSource).(*core1.PersistentVolume)
			if pv.Status.Phase != target {
				continue
			}
			pvList[i].Status = string(pv.Status.Phase)
			pvList[i].ReclaimPolicy = string(pv.Spec.PersistentVolumeReclaimPolicy)
			_, err = cs.Engine.ID(pvList[i].Id).Cols("status", "reclaim_policy").Update(&pvList[i])
			if err == nil {
				completeCount++
			}
		}
		return completeCount == len(pvList)
	}, 5*time.Second, 60)
}

/*
集群创建的时候发现PV并入库
*/
func (cs *commonService) FindPV(pod *core1.Pod) {
	var pvcName string
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == "data" {
			pvcName = volume.PersistentVolumeClaim.ClaimName
		}
	}
	utils.Polling(func() bool {
		err, pvcSource := cs.GetResources("pvc", pvcName, cs.GetNameSpace(), meta1.GetOptions{})
		if err != nil {
			return false
		}
		pvc := (*pvcSource).(*core1.PersistentVolumeClaim)
		pvName := pvc.Spec.VolumeName
		if pvName == "" {
			return false
		}
		exist, _ := cs.Engine.Exist(&models.PersistentVolume{Name: pvName})
		if exist {
			return true
		}
		err, pvSource := cs.GetResources("pv", pvName, cs.GetNameSpace(), meta1.GetOptions{})
		if err != nil {
			return false
		}
		pv := (*pvSource).(*core1.PersistentVolume)
		dbPod := models.Instance{Name: pod.Name}
		hasPod, err := cs.Engine.Cols("id", "cluster_id").Get(&dbPod)
		if !hasPod {
			return true
		}
		dbCluster := models.ClusterInstance{Id: dbPod.ClusterId}
		hasCluster, err := cs.Engine.Cols("user_tag", "org_tag").Get(&dbCluster)
		if !hasCluster {
			return true
		}
		dbSC := models.Sc{Name: pv.Spec.StorageClassName}
		hasSC, err := cs.Engine.Get(&dbSC)
		if !hasSC {
			return true
		}
		var dbPV models.PersistentVolume
		dbPV.Name = pv.Name
		dbPV.PvcName = pvcName
		dbPV.PodId = dbPod.Id
		dbPV.ScId = dbSC.Id
		dbPV.UserTag = dbCluster.UserTag
		dbPV.OrgTag = dbCluster.OrgTag
		capacity := pv.Spec.Capacity["storage"]
		dbPV.Capacity = capacity.String()
		dbPV.Status = string(pv.Status.Phase)
		dbPV.ReclaimPolicy = string(pv.Spec.PersistentVolumeReclaimPolicy)
		_, err = cs.Engine.Insert(&dbPV)
		if err != nil {
			return false
		}
		go cs.SetQosConfig(pv.Name, dbCluster.Id)
		return true
	}, 5*time.Second, 60)
}

/*
修改PV策略, 将PV的Delete策略改为Retain
*/
func (cs *commonService) ChangePVPolicy(pvList []models.PersistentVolume) error {
	for i := range pvList {
		err, pvSource := cs.GetResources("pv", pvList[i].Name, cs.GetNameSpace(), meta1.GetOptions{})
		if err != nil {
			return err
		}
		pv := (*pvSource).(*core1.PersistentVolume)
		if pv.Spec.PersistentVolumeReclaimPolicy == core1.PersistentVolumeReclaimRetain {
			continue
		}
		pv.Spec.PersistentVolumeReclaimPolicy = core1.PersistentVolumeReclaimRetain
		_, err = cs.ClientSet.CoreV1().PersistentVolumes().Update(*cs.GetCtx(), pv, meta1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func getSvc(svcList *core1.ServiceList, svcName string) *core1.Service {
	if svcList == nil || len(svcList.Items) == 0 {
		return nil
	}
	exist, i := utils.SliceExist(len(svcList.Items), func(i int) bool {
		return svcList.Items[i].Name == svcName
	})
	if !exist {
		return nil
	}
	return &svcList.Items[i]
}

func getCluster(clusterList *unstructured.UnstructuredList, clusterName string) *unstructured.Unstructured {
	if clusterList == nil || len(clusterList.Items) == 0 {
		return nil
	}
	exist, i := utils.SliceExist(len(clusterList.Items), func(i int) bool {
		return clusterList.Items[i].GetName() == clusterName
	})
	if !exist {
		return nil
	}
	return &clusterList.Items[i]
}

func getPod(podList *core1.PodList, podName string) *core1.Pod {
	if podList == nil || len(podList.Items) == 0 {
		return nil
	}
	exist, i := utils.SliceExist(len(podList.Items), func(i int) bool {
		return podList.Items[i].Name == podName
	})
	if !exist {
		return nil
	}
	return &podList.Items[i]
}

func (cs *commonService) AsyncCommonInfo() {
	cs.AsyncClusterInfo(-1)
}
func (cs *commonService) AsyncClusterInfo(id int) {
	if cs.Err != nil {
		return
	}
	dbClusterList := make([]models.ClusterInstance, 0)
	if id != -1 {
		err := cs.Engine.Where("id = ?", id).Omit("yaml_text").Find(&dbClusterList)
		utils.LoggerError(err)
	} else {
		err := cs.Engine.Find(&dbClusterList)
		utils.LoggerError(err)
	}

	_, svcListSource := cs.GetResources("service", "", cs.GetNameSpace(), meta1.ListOptions{LabelSelector: "mysql.presslabs.org/cluster"})
	svcList, _ := (*svcListSource).(*core1.ServiceList)
	_, clusterList, _ := cs.GetDynamicResource(constant.MysqlClusterYaml(cs.GetNameSpace(), ""), "")
	_, podListSource := cs.GetResources("pod", "", cs.GetNameSpace(), meta1.ListOptions{LabelSelector: "mysql.presslabs.org/cluster"})
	podList, _ := (*podListSource).(*core1.PodList)
	for _, instance := range dbClusterList {
		if instance.Status == models.ClusterStatusDisable {
			continue
		}
		var updateFlag bool
		var svc *core1.Service
		if instance.IsDeploy {
			err, svcAddr := cs.GetResources("service", fmt.Sprintf("%v-svc", instance.Name), cs.GetNameSpace(), meta1.GetOptions{})
			utils.LoggerError(err)
			if svcAddr != nil && err == nil {
				svc = (*svcAddr).(*core1.Service)
			}
		} else {
			svc = getSvc(svcList, fmt.Sprintf("%v-mysql-master", instance.K8sName))
		}
		if svc != nil {
			conn := fmt.Sprintf("%v", svc.Spec.Ports[0].NodePort)
			inConn := fmt.Sprintf("%v:%v", svc.Spec.ClusterIP, svc.Spec.Ports[0].Port)
			if instance.ConnectString != conn || instance.InnerConnectString != inConn {
				updateFlag = true
				instance.ConnectString = conn
				instance.InnerConnectString = inConn
				for _, port := range svc.Spec.Ports {
					if port.Name == "sidecar-ttyd" {
						instance.ConsolePort = strconv.Itoa(int(port.NodePort))
						break
					}
				}
			}
		}

		oldStatus := instance.Status
		clusterRoleMap := map[string]interface{}{}
		if !instance.IsDeploy {
			instanceAddr := getCluster(clusterList, instance.K8sName)
			if instanceAddr == nil {
				instance.Status = models.ClusterStatusNotFound
				_, _ = cs.Engine.ID(instance.Id).Omit("yaml_text").Update(&instance)
				continue
			}
			if insStatus, _ := instanceAddr.Object["status"]; insStatus != nil {
				// 查找master节点名称
				if nodes, _ := insStatus.(map[string]interface{})["nodes"]; nodes != nil {
					for _, v := range nodes.([]interface{}) {
						value := v.(map[string]interface{})
						conditions, ok := value["conditions"].([]interface{})
						if !ok {
							continue
						}
						var find bool
						for _, i := range conditions {
							iv := i.(map[string]interface{})
							if iv["type"] == "Master" && iv["status"] == "True" {
								clusterRoleMap[value["name"].(string)] = 1
								find = true
								break
							}
						}
						if find {
							break
						}
					}
				}
				// 查找实例状态
				if conditions := insStatus.(map[string]interface{})["conditions"]; conditions != nil {
					for _, v := range conditions.([]interface{}) {
						value := v.(map[string]interface{})
						if value["type"] == "Ready" {
							status := value["status"].(string)
							if status == models.ClusterStatusTrue && instance.Status != models.ClusterStatusTrue {
								go statistics.ClusterComplete(instance.Name, instance.Id)
							}
							if instance.Status != models.ClusterStatusCreating || status == models.ClusterStatusTrue {
								instance.Status = status
							}
							break
						}
					}
				}
			}
		} else {
			deploySource, err := cs.GetClientSet().AppsV1().Deployments(cs.GetNameSpace()).Get(*cs.Ctx, instance.Name, meta1.GetOptions{})
			if err != nil {
				utils.LoggerError(err)
				instance.Status = models.ClusterStatusNotFound
				_, _ = cs.Engine.ID(instance.Id).Omit("yaml_text").Update(&instance)
				continue
			}
			for _, status := range deploySource.Status.Conditions {
				if status.Type == appsv1.DeploymentAvailable {
					if instance.Status != models.ClusterStatusCreating || string(status.Status) == models.ClusterStatusTrue {
						instance.Status = string(status.Status)
					}
					break
				}
			}
		}

		mysqlInstanceList := make([]models.Instance, 0)
		err := cs.Engine.Where("cluster_id = ?", instance.Id).OrderBy("id").Find(&mysqlInstanceList)
		utils.LoggerError(err)
		replicas, _ := strconv.Atoi(instance.Replicas)
		if len(mysqlInstanceList) < replicas {
			go cs.ScanClusterPod(instance.Id, instance.K8sName, replicas, instance.IsDeploy)
		}
		oldMaster := instance.Master
		oldActualReplicas := instance.ActualReplicas
		instance.ActualReplicas = strconv.Itoa(len(mysqlInstanceList))
		podStatusSlice := make([]map[string]interface{}, len(mysqlInstanceList))
		for i, m := range mysqlInstanceList {
			var mysql *core1.Pod
			if instance.IsDeploy {
				err, mysqlAddr := cs.GetResources("pod", m.Name, cs.GetNameSpace(), meta1.GetOptions{})
				if err == nil && mysqlAddr != nil {
					mysql = (*mysqlAddr).(*core1.Pod)
				}
			} else {
				mysql = getPod(podList, m.Name)
			}
			if mysql == nil {
				_, _ = cs.Engine.ID(m.Id).Delete(&m)
				_, _ = cs.Engine.Where("pod_id = ?", m.Id).Delete(new(models.PersistentVolume))
			} else {
				containerCount := 0
				var podStatus string
				var initCount int
				for _, s := range mysql.Status.InitContainerStatuses {
					if s.Ready {
						initCount++
					}
				}
				var maxInit = len(mysql.Status.InitContainerStatuses)
				if initCount != maxInit {
					podStatus = fmt.Sprintf("Init:%v/%v", initCount, maxInit)
				}
				for _, s := range mysql.Status.ContainerStatuses {
					if s.Ready {
						containerCount++
					}
					if len(podStatus) != 0 {
						continue
					}
					if s.State.Waiting != nil {
						podStatus = s.State.Waiting.Reason
					} else if s.State.Terminated != nil {
						podStatus = s.State.Terminated.Reason
					}
				}
				if len(podStatus) == 0 {
					podStatus = string(mysql.Status.Phase)
				}
				podStatusSlice[i] = map[string]interface{}{"podName": mysql.Name, "containerStatus": fmt.Sprintf("%v/%v", containerCount, len(mysql.Status.ContainerStatuses)), "status": podStatus}
				baseInfo := map[string]interface{}{}
				baseInfo["name"] = mysql.Name
				baseInfo["namespace"] = mysql.Namespace
				baseInfo["status"] = mysql.Status.Phase
				baseInfo["node"] = mysql.Spec.NodeName
				baseInfo["ip"] = mysql.Status.PodIP
				baseInfo["labels"] = mysql.Labels
				baseInfo["annotations"] = mysql.Annotations
				initContainer, err := json.Marshal(mysql.Spec.InitContainers)
				utils.LoggerError(err)
				baseInfoMarshal, err := json.Marshal(baseInfo)
				utils.LoggerError(err)
				volume, err := json.Marshal(mysql.Spec.Volumes)
				utils.LoggerError(err)
				container, err := json.Marshal(mysql.Spec.Containers)
				utils.LoggerError(err)
				m.Status = string(mysql.Status.Phase)
				m.Version = mysql.ResourceVersion
				m.InitContainer = string(initContainer)
				m.BaseInfo = string(baseInfoMarshal)
				m.Volume = string(volume)
				m.ContainerInfo = string(container)
				m.DomainName = mysql.Status.PodIP
				if !instance.IsDeploy {
					if clusterRoleMap[fmt.Sprintf("%v.mysql.default", m.Name)] == 1 {
						m.Role = "Master"
						instance.Master = m.Name
					} else {
						m.Role = "Replicate"
					}
				}
				_, err = cs.Engine.ID(m.Id).Update(&m)
				utils.LoggerError(err)
			}
		}
		podStatus, err := json.Marshal(podStatusSlice)
		utils.LoggerError(err)
		if updateFlag || oldStatus != instance.Status || instance.ActualReplicas != oldActualReplicas || instance.Master != oldMaster || instance.PodStatus != string(podStatus) {
			instance.PodStatus = string(podStatus)
			_, err = cs.Engine.ID(instance.Id).Omit("yaml_text").Update(&instance)
			utils.LoggerError(err)
		}
	}
}

func (cs *commonService) AsyncImageStatus() {
	var err error
	defer utils.LoggerErrorP(&err)
	imageList := make([]models.Images, 0)
	err = cs.Engine.Omit("description").Find(&imageList)
	if err != nil {
		return
	}
	session := cs.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return
	}
	address := getImageAddress(cs.Engine)
	for _, image := range imageList {
		oldVersion := image.Status
		image.SetStatus(address)
		if oldVersion != image.Status {
			// 状态发生变化时才进行更新
			_, err := session.ID(image.Id).Cols("status").Update(image)
			utils.LoggerError(err)
		}
	}
	err = session.Commit()
}

func (cs *commonService) AddLog(level string, logSource string, people string, content string) {
	err := cs.LogService.Add(level, logSource, people, content)
	utils.LoggerError(err)
}

func (cs *commonService) GetNameSpace() string {
	return cs.CommonNameSpace
}

func (cs *commonService) AsyncNodeInfo() {
	var err error
	defer utils.LoggerErrorP(&err)
	err, nodesAddr := cs.GetResources("node", "", cs.GetNameSpace(), meta1.ListOptions{})
	if err != nil {
		return
	}
	nodes, ok := (*nodesAddr).(*core1.NodeList)
	if !ok {
		return
	}
	dbNodes := make([]models.Node, 0)
	err = cs.Engine.Find(&dbNodes)
	if err != nil {
		return
	}
	dbNodeM := make(map[string]int)
	for i := range dbNodes {
		dbNodeM[dbNodes[i].NodeName] = i
	}
	for _, item := range nodes.Items {
		var node models.Node
		labels, _ := json.Marshal(item.Labels)
		var nodeStatus = "NotReady"
		for _, condition := range item.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				nodeStatus = "Ready"
				break
			}
		}
		nodeAge := time.Now().Sub(item.CreationTimestamp.Time).String()
		if i := strings.Index(nodeAge, "."); i > 0 {
			nodeAge = nodeAge[:i]
		}

		node.Age = nodeAge
		node.Label = string(labels)
		node.Status = nodeStatus
		if i, ok := dbNodeM[item.Name]; ok {
			dbNodeM[item.Name] = -1
			_, err = cs.Engine.ID(dbNodes[i].Id).Update(&node)
		} else {
			node.NodeName = item.Name
			node.OrgTag = "AA"
			node.UserTag = "AA"
			_, err = cs.Engine.Insert(&node)
		}
		utils.LoggerError(err)
	}
	for _, i := range dbNodeM {
		if i != -1 {
			_, _ = cs.Engine.ID(dbNodes[i]).Delete(new(models.Node))
		}
	}
}

func (cs *commonService) GetClientSet() *kubernetes.Clientset {
	return cs.ClientSet
}

func (cs *commonService) GetConfig() *rest.Config {
	return cs.Config
}

func (cs *commonService) GetCtx() *k8sContext.Context {
	return cs.Ctx
}

/**
 * 新增资源 支持类型: node, pod, sc, pv, pvc, service
 */
func (cs *commonService) CreateOption(sourceType string, nameSpace string, sourceInterface interface{}, opts meta1.CreateOptions) error {
	if cs.Err != nil {
		return cs.Err
	}
	var err error
	switch sourceType {
	case "node":
		if source, ok := sourceInterface.(*core1.Node); ok {
			_, err = cs.ClientSet.CoreV1().Nodes().Create(*cs.Ctx, source, opts)
			if err != nil {
				utils.LoggerError(err)
			}
		}
	case "pod":
		if source, ok := sourceInterface.(*core1.Pod); ok {
			_, err = cs.ClientSet.CoreV1().Pods(nameSpace).Create(*cs.Ctx, source, opts)
			if err != nil {
				utils.LoggerError(err)
			}
		}
	case "sc":
		if source, ok := sourceInterface.(*storage1.StorageClass); ok {
			_, err = cs.ClientSet.StorageV1().StorageClasses().Create(*cs.Ctx, source, opts)
			if err != nil {
				utils.LoggerError(err)
			}
		}
	case "pv":
		if source, ok := sourceInterface.(*core1.PersistentVolume); ok {
			_, err = cs.ClientSet.CoreV1().PersistentVolumes().Create(*cs.Ctx, source, opts)
			if err != nil {
				utils.LoggerError(err)
			}
		}
	case "pvc":
		if source, ok := sourceInterface.(*core1.PersistentVolumeClaim); ok {
			_, err = cs.ClientSet.CoreV1().PersistentVolumeClaims(nameSpace).Create(*cs.Ctx, source, opts)
			utils.LoggerError(err)
		}
	case "service":
		if source, ok := sourceInterface.(*core1.Service); ok {
			_, err = cs.ClientSet.CoreV1().Services(nameSpace).Create(*cs.Ctx, source, opts)
			if err != nil {
				utils.LoggerError(err)
			}
		}
	}
	return err
}

/**
 * 资源删除  支持类型: node, pod, sc, pv, pvc, service
 */
func (cs *commonService) DeleteOption(sourceType string, sourceName string, nameSpace string, opts meta1.DeleteOptions) error {
	if cs.Err != nil {
		return cs.Err
	}
	var err error
	switch sourceType {
	case "node":
		err = cs.ClientSet.CoreV1().Nodes().Delete(*cs.Ctx, sourceName, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "pod":
		err = cs.ClientSet.CoreV1().Pods(nameSpace).Delete(*cs.Ctx, sourceName, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "sc":
		err = cs.ClientSet.StorageV1().StorageClasses().Delete(*cs.Ctx, sourceName, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "pv":
		err = cs.ClientSet.CoreV1().PersistentVolumes().Delete(*cs.Ctx, sourceName, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "pvc":
		err = cs.ClientSet.CoreV1().PersistentVolumeClaims(cs.GetNameSpace()).Delete(*cs.Ctx, sourceName, opts)
		utils.LoggerError(err)
	case "service":
		err = cs.ClientSet.CoreV1().Services(nameSpace).Delete(*cs.Ctx, sourceName, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	}
	return err
}

/**
 *  关于修改的api  支持的类型:node pod sc pv service
 *  修改和新增为一个接口
 *  patch 一般为types.StrategicMergePatchType  但是有些删除操作需要是types.JSONPatchType
 */
func (cs *commonService) PatchOption(sourceType string, sourceName string, nameSpace string, playLoadBytes []byte, opts meta1.PatchOptions, patchType types.PatchType) error {
	if cs.Err != nil {
		return cs.Err
	}
	var err error
	switch sourceType {
	case "node":
		_, err = cs.ClientSet.CoreV1().Nodes().Patch(*cs.Ctx, sourceName, patchType, playLoadBytes, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "pod":
		_, err = cs.ClientSet.CoreV1().Pods(nameSpace).Patch(*cs.Ctx, sourceName, patchType, playLoadBytes, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "sc":
		_, err = cs.ClientSet.StorageV1().StorageClasses().Patch(*cs.Ctx, sourceName, patchType, playLoadBytes, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "pv":
		_, err = cs.ClientSet.CoreV1().PersistentVolumes().Patch(*cs.Ctx, sourceName, patchType, playLoadBytes, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	case "service":
		_, err = cs.ClientSet.CoreV1().Services(nameSpace).Patch(*cs.Ctx, sourceName, patchType, playLoadBytes, opts)
		if err != nil {
			utils.LoggerError(err)
		}
	}
	return err
}

/**
 * 获取k8s静态类型的资源的api
 * 支持类型: node, pod, sc, pv, pvc, service， svc
 */
func (cs *commonService) GetResources(sourceType string, sourceName string, nameSpace string, opts interface{}) (error, *interface{}) {
	var source interface{}
	var err error
	if cs.Err != nil {
		return cs.Err, &source
	}
	var timeout int64 = 5
	switch sourceType {
	case "node":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().Nodes().List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().Nodes().Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "pod":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().Pods(nameSpace).List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().Pods(nameSpace).Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "sc":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.StorageV1().StorageClasses().List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.StorageV1().StorageClasses().Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "pv":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().PersistentVolumes().List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().PersistentVolumes().Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "pvc":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().PersistentVolumeClaims(nameSpace).List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().PersistentVolumeClaims(nameSpace).Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "service":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().Services(nameSpace).List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().Services(nameSpace).Get(*cs.Ctx, sourceName, opts)
			}
		}
	case "secret":
		if sourceName == "" {
			if opts, ok := opts.(meta1.ListOptions); ok {
				opts.TimeoutSeconds = &timeout
				source, err = cs.ClientSet.CoreV1().Secrets(nameSpace).List(*cs.Ctx, opts)
			}
		} else {
			if opts, ok := opts.(meta1.GetOptions); ok {
				source, err = cs.ClientSet.CoreV1().Secrets(nameSpace).Get(*cs.Ctx, sourceName, opts)
			}
		}
	}
	return err, &source
}

/**
 * 动态资源类型客户端  DynamicResource
 */
func (cs *commonService) CreateDynamicResource(deploymentYaml string) (*unstructured.Unstructured, error) {
	// 1. Prepare a RESTMapper to find GVR
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dc, err := discovery.NewDiscoveryClientForConfig(cs.Config)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cs.Config)
	if err != nil {
		return nil, err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(deploymentYaml), nil, obj)
	if err != nil {
		return nil, err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	source, err := dr.Create(*cs.Ctx, obj, meta1.CreateOptions{})
	return source, err
}

/**
 * 动态资源类型客户端  删除DynamicResource
 */
func (cs *commonService) DeleteDynamicResource(deploymentYaml string) error {
	// 1. Prepare a RESTMapper to find GVR
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dc, err := discovery.NewDiscoveryClientForConfig(cs.Config)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cs.Config)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(deploymentYaml), nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}
	err = dr.Delete(*cs.Ctx, obj.GetName(), meta1.DeleteOptions{})
	return err
}

/**
 * 动态资源类型客户端  修改DynamicResource,  修改和集群启停用一个接口
 */
func (cs *commonService) PatchDynamicResource(deploymentYaml string, controllerContent string) error {
	// 1. Prepare a RESTMapper to find GVR
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dc, err := discovery.NewDiscoveryClientForConfig(cs.Config)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cs.Config)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(deploymentYaml), nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}
	// 6. Marshal object into JSON
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Patch(*cs.Ctx, obj.GetName(), types.ApplyPatchType, data, meta1.PatchOptions{
		FieldManager: controllerContent,
	})
	return err
}

func (cs *commonService) UpdateDynamicResource(deploymentYaml string, updateData *unstructured.Unstructured) error {
	// 1. Prepare a RESTMapper to find GVR
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dc, err := discovery.NewDiscoveryClientForConfig(cs.Config)
	if err != nil {
		return err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cs.Config)
	if err != nil {
		return err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(deploymentYaml), nil, obj)
	if err != nil {
		return err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	// 7. Create or Update the object with SSA
	//     types.ApplyPatchType indicates SSA.
	//     FieldManager specifies the field owner ID.
	_, err = dr.Update(*cs.Ctx, updateData, meta1.UpdateOptions{})
	return err
}

/**
 * 动态资源类型客户端  查询DynamicResource,  传集群名字查单个,  传空查所有
 */
func (cs *commonService) GetDynamicResource(deploymentYaml string, sourceName string) (*unstructured.Unstructured, *unstructured.UnstructuredList, error) {
	// 1. Prepare a RESTMapper to find GVR
	var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dc, err := discovery.NewDiscoveryClientForConfig(cs.Config)
	if err != nil {
		return nil, nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(cs.Config)
	if err != nil {
		return nil, nil, err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode([]byte(deploymentYaml), nil, obj)
	if err != nil {
		return nil, nil, err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, nil, err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}
	if sourceName == "" {
		var t int64 = 5
		opts := meta1.ListOptions{TimeoutSeconds: &t}
		source, err := dr.List(*cs.Ctx, opts)
		return nil, source, err
	} else {
		sourceList, err := dr.Get(*cs.Ctx, sourceName, meta1.GetOptions{})
		return sourceList, nil, err
	}

}

func (cs *commonService) GetPodLogs(podName, container string) ([]string, error) {
	var req = cs.ClientSet.CoreV1().Pods(cs.GetNameSpace()).GetLogs(podName, &core1.PodLogOptions{Container: container})
	s, err := req.Stream(k8sContext.TODO())
	if err != nil {
		return nil, err
	}
	var ret = make([]string, 0)
	br := bufio.NewReader(s)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		ret = append(ret, string(a))
	}
	_ = s.Close()
	return ret, err
}

// 获取采集的性能数据
func GetPerformanceData(modelId int32, selectType string, podName string, attrId int32, attrName string, timeInt int64, condition map[string]interface{}, cluster models.ClusterInstance) map[string]interface{} {
	cmdService, cmdConn := NewCmdService()
	collectService, conn := NewCollectService()
	defer CloseGrpc(cmdConn)
	defer CloseGrpc(conn)
	hostDetailInformation := make(map[string]interface{})
	hostDetailInformation["modelId"] = modelId
	hostDetailInformation["podName"] = podName
	//  筛选条件
	conditionString := ""
	for s, i := range condition {
		conditionString += fmt.Sprintf(" AND %s = '%s' ", s, i)
	}
	detail := make([]map[string]interface{}, 0)
	cmdbModelField, _, _ := cmdService.GetCmdbModelField(modelId, selectType, 0, attrId, attrName, "", "K8sMySQLPod")
	oneBodyChan := make(chan map[string]interface{}, 20)
	if len(cmdbModelField) > 0 {
		var count int
		for _, m := range cmdbModelField {
			go getSignalDetail(m, collectService, timeInt, podName, oneBodyChan, conditionString, cluster)
		}
		for {
			oneBody := <-oneBodyChan
			if oneBody != nil {
				detail = append(detail, oneBody)
			}
			count++
			if count == len(cmdbModelField) {
				break
			}
		}
	}
	//数组排序
	sort.Slice(detail, func(i, j int) bool {
		return detail[i]["attrId"].(float64) < detail[j]["attrId"].(float64)
	})
	hostDetailInformation["detail"] = detail
	return hostDetailInformation
}

//获取步长
func getStepInterval(timeSecond float64, intervalTime float64) string {
	interTime := timeSecond / (60 * 60)
	interval := "60s"
	if 0 <= interTime && interTime < 0.5 {
		if intervalTime > 5 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "5s"
		}
	} else if 0.5 <= interTime && interTime < 1 {
		if intervalTime > 10 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "10s"
		}
	} else if 1 <= interTime && interTime < 3 {
		if intervalTime > 20 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "20s"
		}
	} else if 3 <= interTime && interTime < 6 {
		if intervalTime > 60 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "60s"
		}
	} else if 6 <= interTime && interTime < 12 {
		if intervalTime > 120 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "120s"
		}
	} else if 12 <= interTime && interTime < 48 {
		if intervalTime > 300 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = "300s"
		}
	} else if 48 <= interTime && interTime < 72 {
		if intervalTime > 10*60 {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = strconv.FormatInt(10*60, 10) + "s"
		}
	} else {
		index := timeSecond - (60 * 60 * 72)
		if intervalTime > index {
			interval = fmt.Sprintf("%vs", intervalTime)
		} else {
			interval = strconv.FormatInt(int64(600+(10*60)*index), 10) + "s"
		}
	}
	return interval
}

func getSignalDetail(mm interface{}, collectService CollectService, timeInt int64, podName string, oneBodyChan chan map[string]interface{}, conditionString string, cluster models.ClusterInstance) {
	var sql string
	m := mm.(map[string]interface{})
	oneBody := make(map[string]interface{})
	collect := m["collect"].(map[string]interface{})

	intervalIntTime := "5s"
	collectInt := int64(0)
	if timeInt == 0 && collect["interval"] != nil {
		timeInt = 2 * int64(collect["interval"].(float64))
	}
	if collect["interval"] != nil && timeInt < int64(collect["interval"].(float64)) {
		timeInt = int64(collect["interval"].(float64))
		collectInt = int64(collect["interval"].(float64))
	}
	if collect["interval"] != nil {
		intervalIntTime = getStepInterval(float64(timeInt), collect["interval"].(float64))
	}
	if collectInt == 0 {
		collectInt = timeInt
	}
	fieldType := m["field_type"].(map[string]interface{})
	chart := m["chart"].([]interface{})
	if collectionIndex, ok := m["collection_index"]; ok {
		fieldString := ""
		menFieldString := ""
		displayFields, ok := collect["display_fields"].([]interface{})
		if !ok {
			displayFields = make([]interface{}, 0)
		}
		var tableName = "tableName"
		_, tableNameOk := m["collect"].(map[string]interface{})["measurement"].(string)
		if tableNameOk {
			tableName = strings.Replace(m["collect"].(map[string]interface{})["measurement"].(string), "_k8s_salve", "", 1)
		}
		typeString := fieldType["name"].(string)
		attrId := m["id"]
		unit := ""
		if collect["result_type"] != nil {
			unit = collect["result_type"].(string)
		} else {
			if v, ok := m["unit"].(string); ok {
				unit = v
			}
		}

		if methodsCollect, ok := collect["method"]; ok {
			if methodsCollect == "PostgreSQL" {
				unit = ""
			}
		}
		nameZh := m["name_zh"]
		name := m["name"]
		for i, i2 := range displayFields {
			fieldString += fmt.Sprintf(" \"%v\"", i2.(string))
			menFieldString += fmt.Sprintf("MEAN( \"%v\" ) as \"%v\"", i2.(string), i2.(string))
			if i != len(displayFields)-1 {
				fieldString += ", "
				menFieldString += ", "
			}
		}
		isSummationResult := isSummation(tableName)
		if collectionIndex == "collect" {
			//  不管是什么图形，  公用的部分
			typeChartString := ""
			if len(chart) > 0 {
				typeChartString = chart[0].(map[string]interface{})["chart_type"].(string)
			}
			oneBody["type"] = typeString
			oneBody["chartType"] = typeChartString
			oneBody["attrId"] = attrId
			oneBody["code"] = name
			oneBody["unit"] = unit
			oneBody["name"] = nameZh
			//  文本类型
			if typeChartString == "Text" {
				if isSummationResult {
					sql = fmt.Sprintf(`SELECT  DIFFERENCE( %v ) as %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=60s %v`, fieldString, fieldString, tableName, podName, conditionString)
				} else {
					sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs  %v  order by time desc limit 1;`, fieldString, tableName, podName, timeInt, conditionString)
				}
				result := collectService.GetInfluxDbData(sql, "1")
				if len(result) > 0 {
					result := result[len(result)-1 : len(result)]
					if typeString == "Json" {
						sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time >= %v AND time <= %v  %v`, fieldString, tableName, podName, fmt.Sprintf("%f", result[0]["time"].(float64)-500000), fmt.Sprintf("%f", result[0]["time"].(float64)+500000), conditionString)
						resultAll := collectService.GetInfluxDbData(sql, "")
						resultHandle := make([]map[string]interface{}, 0)
						for i, m2 := range resultAll {
							resultOneDetail := make(map[string]interface{}, 0)
							resultOne := make([]map[string]interface{}, 0)
							resultStatus := make(map[string]interface{}, 0)
							string := m2[displayFields[0].(string)].(string)
							if strings.Index(string, "@@@") != -1 {
								status := strings.Split(string, "@@@")[1]
								string = strings.Split(string, "@@@")[0]
								string = strings.ReplaceAll(string, "@@", "\"")
								status = strings.ReplaceAll(status, "@@", "\"")
								err := json.Unmarshal([]byte(status), &resultStatus)
								utils.LoggerError(err)
								err = json.Unmarshal([]byte(string), &resultOne)
								utils.LoggerError(err)
							}
							resultOneDetail["name"] = strconv.FormatInt(int64(i+1), 10)
							resultOneDetail["attr"] = resultOne
							resultOneDetail["status"] = resultStatus
							resultHandle = append(resultHandle, resultOneDetail)
						}
						oneBody["data"] = resultHandle

					} else {
						resultOne := make([]map[string]interface{}, 0)
						for _, field := range displayFields {
							if valueString, ok := result[0][field.(string)].(string); ok {
								resultOneMiddle := make([]map[string]interface{}, 0)
								if strings.Index(valueString, "@@@") != -1 {
									if strings.Index(valueString, "@@@@") != -1 {
										valueString = strings.ReplaceAll(valueString, "@@@@", "\"\"")
									}
									valueString = strings.Split(valueString, "@@@")[0]
									valueString = strings.ReplaceAll(valueString, "@@", "\"")
									err := json.Unmarshal([]byte(valueString), &resultOneMiddle)
									utils.LoggerError(err)
									result[0][field.(string)] = resultOneMiddle
								} else if strings.Index(valueString, "@@") != -1 {
									valueString = strings.ReplaceAll(valueString, "@@", "\"")
									err := json.Unmarshal([]byte(valueString), &resultOneMiddle)
									utils.LoggerError(err)
									result[0][field.(string)] = resultOneMiddle
								} else {
									sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time >= %v AND time <= %v  %v `, fieldString, tableName, podName, fmt.Sprintf("%f", result[0]["time"].(float64)-500000), fmt.Sprintf("%f", result[0]["time"].(float64)+500000), conditionString)
									resultAll := collectService.GetInfluxDbData(sql, "")
									for _, m2 := range resultAll {
										resultOne = append(resultOne, m2)
									}
									break
								}
							}
						}
						if len(displayFields) == 1 {
							if value, ok := result[0][displayFields[0].(string)].([]map[string]interface{}); ok {
								resultOne = value
							} else {
								resultOne = append(make([]map[string]interface{}, 0), map[string]interface{}{"value": result[0][displayFields[0].(string)]})
							}
						} else {
							if len(resultOne) <= 0 {
								resultOne = append(resultOne, result[0])
							}
						}
						oneBody["data"] = resultOne
					}
				} else {
					oneBody["data"] = make([]map[string]interface{}, 0)
				}
			} else if typeChartString == "Table" {
				// 屏蔽表格形式的数据， 主要是为了屏蔽Largest_Tables_by_Size和Largest_Tables_by_Row_Count
				oneBody = nil

				//if tableName == "mysql_info_schema_table_rows" || tableName == "mysql_info_schema_table_size" {
				//	conditionString = `AND "schema" != 'sys' AND "schema" != 'sys_operator' `
				//}
				////  对表格的处理
				//sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs  %v order by time desc limit 1;`, fieldString, tableName, podName, timeInt, conditionString)
				//result := collectService.GetInfluxDbData(sql, "1")
				//oneBodyData := make(map[string]interface{}, 0)
				//columns := make([]map[string]interface{}, 0)
				//if len(result) > 0 {
				//	sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time = %v  %v;`, fieldString, tableName, podName, fmt.Sprintf("%f", result[0]["time"].(float64)), conditionString)
				//	resultAll := collectService.GetInfluxDbData(sql, "")
				//	if len(result) > 0 {
				//		sort.Slice(resultAll, func(i, j int) bool {
				//			return resultAll[i]["value"].(float64) > resultAll[j]["value"].(float64)
				//		})
				//		if len(resultAll) > 10 {
				//			oneBodyData["detail"] = resultAll[:10]
				//		} else {
				//			oneBodyData["detail"] = resultAll
				//		}
				//	}
				//	if chartVal, ok := m["chart_fields"]; ok && len(chartVal.([]interface{})) > 0 {
				//		chartVal := chartVal.([]interface{})[0].(map[string]interface{})
				//		tableField := strings.Split(chartVal["table_field"].(string), ",")
				//		alias := strings.Split(chartVal["alias"].(string), ",")
				//		for i, s := range tableField {
				//			columns = append(columns, map[string]interface{}{"name": s, "alias": alias[i]})
				//		}
				//		oneBodyData["columns"] = columns
				//	}
				//} else {
				//	oneBodyData["detail"] = result
				//	oneBodyData["columns"] = columns
				//}
				//oneBody["data"] = oneBodyData
			} else if typeChartString == "Bar" {
				chartDataOne := make([]map[string]interface{}, 0)
				barLabel := make([]interface{}, 0)
				if chartFields, ok := m["chart_fields"].([]interface{}); ok && len(chartFields) > 0 && tableNameOk {
					chartFields := chartFields[0].(map[string]interface{})
					yField := strings.Split(chartFields["y_field"].(string), ",")
					var groupValues []map[string]interface{}
					xField, ok := chartFields["x_field"].(string)
					// 分组
					groupField, ok := chartFields["group_field"].(string)
					if xField != "time" {
						groupField = xField
					}
					groupFields := make([]map[string]interface{}, 0)
					if len(groupField) > 0 && ok {
						sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' %v limit 30;`, groupField, tableName, podName, conditionString)
						groupFields = collectService.GetInfluxDbData(sql, "")
						for _, v := range groupFields {
							repeat := false
							for _, values := range groupValues {
								if v[groupField] == values[groupField] {
									repeat = true
									break
								}
							}
							if !repeat {
								groupValues = append(groupValues, map[string]interface{}{groupField: v[groupField]})
							}
						}
					}
					if len(groupValues) > 0 {
						for _, field := range yField {
							barResult := make(map[string]interface{}, 0)
							barResult["name"] = field
							for _, groupValue := range groupValues {
								middleResult := make(map[string]interface{}, 0)
								sql = fmt.Sprintf(`SELECT  %v  FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND %v = '%v'  AND time>=now()-%vs %v fill(0) order by time desc limit 1;`, field, tableName, podName, groupField, groupValue[groupField], timeInt, conditionString)
								results := collectService.GetInfluxDbData(sql, "")
							resultsLoop:
								for i, result := range results {
									if result[field].(float64) == 0 && (i == 0 || i == len(results)-1) {
										continue resultsLoop
									}
									if xField == "time" {
										middleResult["xValue"] = result["time"]
										middleResult["label"] = groupValue[groupField]
									} else {
										middleResult["xValue"] = groupValue[groupField]
									}
									_, err := result[field].(string)
									if err {
										middleResult["yValue"] = fmt.Sprintf("%.2f", result[field])
									} else {
										middleResult["yValue"] = result[field]
									}
									barLabel = append(barLabel, middleResult["xValue"].(string))
									barResult[middleResult["xValue"].(string)] = middleResult["yValue"]
								}
							}
							chartDataOne = append(chartDataOne, barResult)
						}
					} else {
						barResult := make(map[string]interface{}, 0)
						barResult["name"] = name
						for _, field := range yField {
							middleResult := make(map[string]interface{}, 0)
							sql = fmt.Sprintf(`SELECT %v  FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs  %v  fill(0) order by time desc limit 1;`, field, tableName, podName, timeInt, conditionString)
							results := collectService.GetInfluxDbData(sql, "")
						resultsLoop2:
							for i, result := range results {
								if result[field].(float64) == 0 && (i == 0 || i == len(results)-1) {
									continue resultsLoop2
								}
								middleResult["xValue"] = result["time"].(string)
								_, err := result[field].(string)
								if err {
									middleResult["yValue"] = fmt.Sprintf("%.2f", result[field])
								} else {
									middleResult["yValue"] = result[field]
								}
								barLabel = append(barLabel, middleResult["xValue"].(string))
								barResult[middleResult["xValue"].(string)] = middleResult["yValue"]
							}
							chartDataOne = append(chartDataOne, barResult)
						}
					}
				}
				oneBody["data"] = chartDataOne
				oneBody["label"] = getSlice(barLabel)
			} else if typeChartString == "Line" {
				chartComputed := make([]map[string]interface{}, 0)
				chartDataOne := make([]map[string]interface{}, 0)
				if chartFields, ok := m["chart_fields"].([]interface{}); ok && len(chartFields) > 0 && tableNameOk {
					chartFields := chartFields[0].(map[string]interface{})
					yField := strings.Split(chartFields["y_field"].(string), ",")
					var groupValues []map[string]interface{}
					xField, ok := chartFields["x_field"].(string)
					// 分组
					groupField, ok := chartFields["group_field"].(string)
					if xField != "time" {
						groupField = xField
					}
					groupFields := make([]map[string]interface{}, 0)
					if len(groupField) > 0 && ok {
						sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' %v limit 30;`, groupField, tableName, podName, conditionString)
						groupFields = collectService.GetInfluxDbData(sql, "")
						for _, v := range groupFields {
							repeat := false
							for _, values := range groupValues {
								if v[groupField] == values[groupField] {
									repeat = true
									break
								}
							}
							if !repeat {
								groupValues = append(groupValues, map[string]interface{}{groupField: v[groupField]})
							}
						}
					}
					if len(groupValues) > 0 {
						for _, field := range yField {
							for _, groupValue := range groupValues {
								middleResult := make(map[string]interface{}, 0)
								sql = fmt.Sprintf(`SELECT  MEAN( %v ) as %v  FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND %v = '%v'  AND time>=now()-%vs %v group by time(%v)  fill(0) order by time desc;`, field, field, tableName, podName, groupField, groupValue[groupField], timeInt, conditionString, intervalIntTime)
								results := collectService.GetInfluxDbData(sql, "")
								labelComputed := make(map[string]interface{}, 0)
								labelComputed["label"] = groupValue[groupField]
								labelComputed["max"] = float64(0)
								labelComputed["min"] = float64(0)
							resultsLineLoop:
								for i, result := range results {
									if result[field].(float64) == 0 && (i == 0 || i == len(results)-1) {
										continue resultsLineLoop
									}
									middleResult = make(map[string]interface{}, 0)
									if xField == "time" {
										middleResult["xValue"] = result["time"]
										middleResult["label"] = groupValue[groupField]
									} else {
										middleResult["xValue"] = groupValue[groupField]
									}
									compare, err := result[field].(float64)
									if err {
										if compare > labelComputed["max"].(float64) {
											labelComputed["max"] = compare
										}
										if compare < labelComputed["min"].(float64) {
											labelComputed["min"] = compare
										}
									}
									_, err = result[field].(string)
									if err {
										middleResult["yValue"] = fmt.Sprintf("%.2f", result[field])
									} else {
										middleResult["yValue"] = result[field]
									}
									chartDataOne = append(chartDataOne, middleResult)
								}
								chartComputed = append(chartComputed, labelComputed)
							}
						}
					} else {
						for _, field := range yField {
							middleResult := make(map[string]interface{}, 0)
							if isSummationResult {
								sql = fmt.Sprintf(`SELECT  DIFFERENCE( %v ) as %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs %v`, field, field, tableName, podName, timeInt, conditionString)

							} else {
								if tableName == "metrics_pod_cpu" || tableName == "metrics_pod_mem" {
									sql = fmt.Sprintf(`SELECT  MEAN( %v ) as %v  FROM zdcp.autogen.%v WHERE "pod"='%v' AND time>=now()-%vs %v group by time(%v)  fill(0) order by time desc;`, field, field, tableName, podName, timeInt, conditionString, intervalIntTime)
								} else {
									sql = fmt.Sprintf(`SELECT  MEAN( %v ) as %v  FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs %v group by time(%v)  fill(0) order by time desc;`, field, field, tableName, podName, timeInt, conditionString, intervalIntTime)
								}
							}
							results := collectService.GetInfluxDbData(sql, "")
							if tableName == "metrics_pod_cpu" {
								for _, result := range results {
									cpuUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(result[field].(float64))/float64(cluster.LimitCpu))*100), 64)
									if err != nil {
										fmt.Println(err)
									}
									result[field] = cpuUsage
								}
							}
							if tableName == "metrics_pod_mem" {
								for _, result := range results {
									memUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(result[field].(float64))/float64(cluster.LimitMem))*100), 64)
									if err != nil {
										fmt.Println(err)
									}
									result[field] = memUsage
								}
							}

							labelComputed := make(map[string]interface{}, 0)
							labelComputed["label"] = field
							labelComputed["max"] = float64(0)
							labelComputed["min"] = float64(0)
						resultsLineLoop2:
							for i, result := range results {
								if result[field].(float64) == 0 && (i == 0 || i == len(results)-1) {
									continue resultsLineLoop2
								}
								middleResult = make(map[string]interface{}, 0)
								middleResult["xValue"] = result["time"].(string)
								compare, err := result[field].(float64)
								if err {
									if compare > labelComputed["max"].(float64) {
										labelComputed["max"] = compare
									}
									if compare < labelComputed["min"].(float64) {
										labelComputed["min"] = compare
									}
								}
								_, err = result[field].(string)
								if err {
									middleResult["yValue"] = fmt.Sprintf("%.2f", result[field])
								} else {
									middleResult["yValue"] = result[field]
								}
								middleResult["label"] = field
								chartDataOne = append(chartDataOne, middleResult)
							}
							chartComputed = append(chartComputed, labelComputed)
						}
					}
				}
				oneBody["computed"] = chartComputed
				oneBody["data"] = chartDataOne
			} else if typeChartString == "Ring" || typeChartString == "Pie" {
				chartDataOne := make([]map[string]interface{}, 0)
				for _, field := range displayFields {
					sql = fmt.Sprintf(`SELECT %v  FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs %v fill(0) order by time desc limit 1;`, field, tableName, podName, timeInt, conditionString)
					results := collectService.GetInfluxDbData(sql, "")
					if len(results) > 0 {
						middleResult := make(map[string]interface{}, 0)
						middleResult["count"] = results[0][field.(string)]
						middleResult["item"] = field
						chartDataOne = append(chartDataOne, middleResult)
					}
				}
				oneBody["data"] = chartDataOne
			} else {
				sql = fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs %v group by time(%v) fill(0);`, menFieldString, tableName, podName, timeInt, conditionString, intervalIntTime)
				results := collectService.GetInfluxDbData(sql, "")
				for _, result := range results {
					middleResult := make(map[string]interface{}, 0)
					middleResult["xValue"] = result["time"].(string)
					middleResult["label"] = displayFields[0].(string)
					middleResult["yValue"] = result[displayFields[0].(string)]
					oneBody["data"] = middleResult
				}
			}
		}
	}
	oneBodyChan <- oneBody
}

//去重
func getSlice(slice []interface{}) []interface{} {
	result := make([]interface{}, 0)
	tempMap := map[string]byte{}
	for _, e := range slice {
		l := len(tempMap)
		tempMap[e.(string)] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

//采集项是否是累加值
func isSummation(tableName string) bool {
	result := false
	var summationTable = []string{"mysql_global_status_questions", "mysql_global_status_slow_queries", "mysql_global_status_table_locks_waited", "mysql_global_status_bytes_received", "mysql_global_status_bytes_sent"}
	for _, sumTable := range summationTable {
		if strings.Contains(tableName, sumTable) {
			result = true
		}
	}
	return result
}

/**
 * k8s初始化
 */
func InitK8s(configString string) (*rest.Config, *kubernetes.Clientset, *k8sContext.Context, error) {
	filename := "./k8sconfig"
	var f *os.File
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在
		f, err = os.Create(filename)
	} else {
		err = os.Remove(filename)
		f, err = os.Create(filename)
	}
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(f, configString)
	if err != nil {
		panic(err)
	}
	ctx := k8sContext.TODO()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", filename)
	if err != nil {
		return config, nil, &ctx, err
	}
	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return config, clientSet, &ctx, err
	}
	return config, clientSet, &ctx, err
}
