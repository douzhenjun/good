package service

import (
	"DBaas/models"
	"DBaas/utils"
	"DBaas/x/constant"
	"DBaas/x/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"strconv"
	"strings"
)

type BackupService interface {
	StorageType() ([]string, error)
	StorageList(userId, page, pageSize int, search string) ([]models.BackupStorage, int64, error)
	StorageCreate(backupStorage *models.BackupStorage, userIdStr string) error
	StorageDelete(storageId int) error
	StorageReconnect(storageId int) error
	StorageLast(clusterId int) (int, error)
	StorageUser(userIdStr string, storageId int) error

	List(page, pageSize, storageId, clusterId int, startTime, endTime int64, search, status, t, userTag string) ([]models.BackupJobView, int, error)
	Create(backupTask *models.BackupTask) error
	Delete(jobId int) error
	DeleteCycle(clusterId int) error
	Recovery(clusterName, password, remark string, storageMap map[string]interface{}, param []map[string]interface{}, jobId int, qos *models.Qos, comboId, nodePort int) error

	Event(jobId int) ([]models.PodLog, error)
	Logs(jobId int) ([]string, error)
}

type backupService struct {
	engine *xorm.Engine
	cs     CommonService
}

func NewBackupService(db *xorm.Engine, cs CommonService) BackupService {
	return &backupService{db, cs}
}

func (bs *backupService) StorageUser(userIdStr string, storageId int) error {
	if storageId <= 0 {
		return errors.New("storage id must > 0")
	}
	// userId为-1时表示ALL，为空时表示无租户，有租户时以逗号分割
	if userIdStr == "-1" {
		_, _ = bs.engine.Where("storage_id = ?", storageId).Delete(new(models.BackupStorageUser))
	} else {
		dbUser := make([]models.BackupStorageUser, 0)
		err := bs.engine.Where("storage_id = ?", storageId).Find(&dbUser)
		if err != nil {
			return err
		}
		dbUserM := map[int]int{}
		for i := range dbUser {
			dbUserM[dbUser[i].UserId] = i
		}
		userStrList := strings.Split(userIdStr, ",")
		insertList := make([]models.BackupStorageUser, 0)
		for i := range userStrList {
			if len(userStrList[i]) == 0 {
				continue
			}
			var userId, err = strconv.Atoi(userStrList[i])
			if err != nil {
				utils.LoggerError(err)
				continue
			}
			if _, ok := dbUserM[userId]; ok {
				dbUserM[userId] = -1
				continue
			}
			insertList = append(insertList, models.BackupStorageUser{StorageId: storageId, UserId: userId})
		}
		if len(insertList) > 0 {
			_, err := bs.engine.Insert(&insertList)
			if err != nil {
				return err
			}
		}
		for _, v := range dbUserM {
			if v != -1 {
				_, _ = bs.engine.ID(dbUser[v].Id).Delete(new(models.BackupStorageUser))
			}
		}
	}
	_, err := bs.engine.ID(storageId).Cols("assign_all").Update(&models.BackupStorage{AssignAll: userIdStr == "-1"})
	return err
}

func GetUseBackup(userId int, engine *xorm.Engine) (int, error) {
	findTask := make([]models.BackupTask, 0)
	err := engine.Where("user_id = ?", userId).Cols("type", "keep_copy", "close", "cluster_id").Find(&findTask)
	if err != nil {
		return 0, fmt.Errorf("query task is error: %v", err)
	}
	var useBackup int
	for i := range findTask {
		if findTask[i].Type == models.BackupTypeOnce {
			useBackup++
		} else {
			if findTask[i].Close {
				// 如果已经关闭则查询剩下的备份数量
				var count int64
				_, _ = engine.SQL("select count(*) from backup_job bj inner join backup_task bt on bt.id = bj.backup_task_id where bt.cluster_id = ? and bt.type = ?", findTask[i].ClusterId, models.BackupTypeCycle).Get(&count)
				useBackup += int(count)
			} else {
				useBackup += findTask[i].KeepCopy
			}
		}
	}
	return useBackup, nil
}

func (bs *backupService) Logs(jobId int) ([]string, error) {
	if jobId <= 0 {
		return nil, errors.New("backup job id must > 0")
	}
	job := models.BackupJob{Id: jobId}
	exist, err := bs.engine.Cols("pod_name").Get(&job)
	if !exist {
		return nil, fmt.Errorf("not found backup job %v, error: %v", jobId, err)
	}
	return bs.cs.GetPodLogs(job.PodName, "backup")
}

func (bs *backupService) Event(jobId int) ([]models.PodLog, error) {
	if jobId <= 0 {
		return nil, errors.New("backup job id must > 0")
	}
	job := models.BackupJob{Id: jobId}
	exist, err := bs.engine.Cols("pod_name").Get(&job)
	if !exist {
		return nil, fmt.Errorf("not found backup job %v, error: %v", jobId, err)
	}
	return bs.cs.GetEvent(job.PodName), nil
}

func (bs *backupService) StorageLast(clusterId int) (int, error) {
	// 首先选择周期备份使用的存储
	task := models.BackupTask{ClusterId: clusterId, Type: models.BackupTypeCycle}
	exist, _ := bs.engine.Cols("storage_id").Get(&task)
	if exist {
		return task.StorageId, nil
	}
	task = models.BackupTask{}
	_, err := bs.engine.Where("cluster_id = ?", clusterId).Cols("storage_id").Desc("id").Get(&task)
	return task.StorageId, err
}

func (bs *backupService) StorageType() ([]string, error) {
	var ret = make([]string, 0)
	err := bs.engine.SQL("select type from backup_storage_type").Find(&ret)
	return ret, err
}

func (bs *backupService) Recovery(clusterName, password, remark string, storageMap map[string]interface{}, param []map[string]interface{}, jobId int, qos *models.Qos, comboId, nodePort int) error {
	if jobId == 0 {
		return errors.New("backup job id cannot be 0")
	}
	if !utils.StringLength(clusterName, 1, 35) {
		return errors.New("cluster name length must in 1-35")
	}
	if !utils.StringLength(password, 1, 30) {
		return errors.New("mysql password length must in 1-30")
	}
	if !utils.MustMap(storageMap, "cpu", "mem", "size", "copy") {
		return errors.New("storage map is incomplete")
	}

	job := models.BackupJob{Id: jobId}
	exist, err := bs.engine.Get(&job)
	if !exist {
		return fmt.Errorf("not found backup job %v, error: %v", jobId, err)
	}
	if job.Status != models.BackupStatusCompleted || job.BackupSet == "" {
		return errors.New("backup job is not completed")
	}
	task := models.BackupTask{Id: job.BackupTaskId}
	exist, err = bs.engine.Get(&task)
	if !exist {
		return fmt.Errorf("not found backup task %v, error: %v", task.Id, err)
	}
	backupStorage := models.BackupStorage{Id: task.StorageId}
	exist, err = bs.engine.Get(&backupStorage)
	if !exist {
		return fmt.Errorf("not found backup storage %v, error: %v", task.StorageId, err)
	}

	limitMem, limitCpu, storage, replicas := int(storageMap["mem"].(float64)), int(storageMap["cpu"].(float64)), int(storageMap["size"].(float64)), int(storageMap["copy"].(float64))
	enough, msg, u := getUserResource(task.UserId, limitCpu, limitMem, storage, bs.engine)
	if !enough {
		return errors.New(msg)
	}
	scName, ok := storageMap["scName"].(string)
	if !ok || len(scName) == 0 {
		scName, err = matchSc(replicas, u.Id, bs.engine)
		if err != nil {
			return err
		}
	}
	cluster := models.ClusterInstance{Id: task.ClusterId}
	exist, err = bs.engine.Unscoped().Cols("image_id").Get(&cluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", cluster.Id, err)
	}

	mysqlConf, err := getImageParam(cluster.ImageId, bs.engine)
	if err != nil {
		return fmt.Errorf("query mysql configs is error: %v", err)
	}
	if param != nil {
		for _, m := range param {
			for i := range mysqlConf {
				if m["key"] == mysqlConf[i].ParameterName {
					mysqlConf[i].ParameterValue = fmt.Sprintf("%v", m["value"])
					break
				}
			}
		}
	}
	dp, err := dyParam(int64(limitMem), bs.engine)
	if err != nil {
		return fmt.Errorf("setting dynamic parameters is error: %v", err)
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

	secretName := fmt.Sprintf("secret-%v%v", task.UserId, clusterName)
	mysqlSecret := fmt.Sprintf(`
apiVersion: v1
kind: Secret
metadata:
 namespace: %v
 name: %v
type: Opaque
data:
 ROOT_PASSWORD: %v`, bs.cs.GetNameSpace(), secretName, password)
	_, err = bs.cs.CreateDynamicResource(mysqlSecret)
	if err != nil {
		return fmt.Errorf("create secret is error: %v", err)
	}
	deleteSecret := func() { _ = bs.cs.DeleteDynamicResource(mysqlSecret) }

	k8sName := fmt.Sprintf("%v-b-%v", clusterName, task.UserId)
	var initBucketURL string
	if backupStorage.Bucket[len(backupStorage.Bucket)-1] == '/' {
		initBucketURL = backupStorage.Bucket + job.BackupSet
	} else {
		initBucketURL = fmt.Sprintf("%v/%v", backupStorage.Bucket, job.BackupSet)
	}
	clusterYaml := fmt.Sprintf(`apiVersion: mysql.presslabs.org/v1alpha1
kind: MysqlCluster
metadata:
  namespace: %v
  name: %v
spec:
  replicas: %v
  secretName: %v
  # 连接备份集
  initBucketURL: %v
  initBucketSecretName: %v
  ignoreReadOnly: false
  masterServiceSpec:
    serviceType: NodePort
    nodePort: %v
  mysqlConf: 
    %v
  podSpec:
    nodeSelector:
      iwhalecloud.dbassnode: mysql
    resources:
      requests:
        cpu: 1000m
        memory: 1024Mi
      limits:
        cpu: %v
        memory: %vGi
  volumeSpec:
    persistentVolumeClaim:
      storageClassName: %v
      accessModes: [ "ReadWriteOnce" ]
      resources:
        requests:
          storage: %vGi
`, bs.cs.GetNameSpace(), k8sName, replicas, secretName, initBucketURL, backupStorage.Name, nodePort, mysqlConfString, limitCpu, limitMem, scName, storage)
	_, err = bs.cs.CreateDynamicResource(clusterYaml)
	if err != nil {
		deleteSecret()
		return fmt.Errorf("create cluster is error: %v", err)
	}
	secret := fmt.Sprintf(`{"ROOT_PASSWORD":"%v"}`, password)
	dbCluster := models.ClusterInstance{
		Name:       clusterName,
		K8sName:    k8sName,
		Status:     models.ClusterStatusCreating,
		Storage:    storage,
		SecretName: secretName,
		UserId:     task.UserId,
		ImageId:    cluster.ImageId,
		ScName:     scName,
		Replicas:   strconv.Itoa(replicas),
		LimitCpu:   limitCpu,
		LimitMem:   limitMem,
		YamlText:   clusterYaml,
		UserTag:    u.UserTag,
		Secret:     secret,
		Remark:     remark,
		ComboId:    comboId,
	}
	_, err = bs.engine.Insert(&dbCluster)
	if err != nil {
		_ = bs.cs.DeleteDynamicResource(clusterYaml)
		deleteSecret()
		return fmt.Errorf("insert error: %v", err)
	}
	if qos != nil {
		qos.ClusterId = dbCluster.Id
		_, _ = bs.engine.Insert(qos)
	}
	go statistics.ClusterDeploy(clusterName, dbCluster.Id, replicas, "internal")
	// 集群参数入库
	clusterParams := make([]*models.Clusterparameters, len(mysqlConf))
	for i := range mysqlConf {
		clusterParams[i] = &models.Clusterparameters{
			ParameterName:  mysqlConf[i].ParameterName,
			ParameterValue: mysqlConf[i].ParameterValue,
			ClusterId:      dbCluster.Id,
		}
	}
	_, _ = bs.engine.Insert(&clusterParams)
	if replicas > 0 {
		go bs.cs.ScanClusterPod(dbCluster.Id, k8sName, replicas, false)
	}
	go bs.cs.CreatStatusTimeout(dbCluster.Id)
	return nil
}

func (bs *backupService) StorageList(userId, page, pageSize int, search string) ([]models.BackupStorage, int64, error) {
	var session = bs.engine.Where("name like ?", "%"+search+"%")
	if userId > 0 {
		session.
			Join("LEFT OUTER", "backup_storage_user", "backup_storage_user.storage_id = backup_storage.id").
			And("backup_storage.assign_all = true").Or("backup_storage_user.user_id = ?", userId)
	}
	var list = make([]models.BackupStorage, 0)
	count, err := pageFind(page, pageSize, &list, session, new(models.BackupStorage))
	if err != nil {
		return nil, 0, err
	}
	serial := pageSize * (page - 1)
	for i := range list {
		serial++
		list[i].Serial = serial
		if list[i].AssignAll {
			list[i].UserIds = utils.RawJson(`"-1"`)
		} else {
			var userIds = make([]models.User, 0)
			err = bs.engine.SQL("select u.id, u.user_name from backup_storage_user bsu inner join \"user\" u on u.id = bsu.user_id where bsu.storage_id = ?", list[i].Id).Find(&userIds)
			utils.LoggerError(err)
			s, _ := json.Marshal(userIds)
			list[i].UserIds = utils.RawJson(utils.Bytes2str(s))
		}
	}
	return list, count, nil
}

func (bs *backupService) StorageCreate(backupStorage *models.BackupStorage, userIdStr string) error {
	backupStorage.SetStatus()
	y := fmt.Sprintf(`apiVersion: v1
kind: Secret
metadata:
  namespace: %v
  name: %v
type: Opaque
data:
    AWS_ACCESS_KEY_ID: %v
    AWS_SECRET_KEY: %v
    S3_PROVIDER: %v
    S3_ENDPOINT: %v`, bs.cs.GetNameSpace(), backupStorage.Name, utils.Base64En(backupStorage.AccessKey), utils.Base64En(backupStorage.SecretKey), utils.Base64En(backupStorage.Type), utils.Base64En(backupStorage.EndPoint))
	_, err := bs.cs.CreateDynamicResource(y)
	if err != nil {
		return err
	}
	backupStorage.AssignAll = userIdStr == "-1"
	_, err = bs.engine.Insert(backupStorage)
	if err != nil || backupStorage.AssignAll {
		return err
	}
	userIds := strings.Split(userIdStr, ",")
	bsu := make([]models.BackupStorageUser, len(userIds))
	for i := range userIds {
		userId, err := strconv.Atoi(userIds[i])
		if err != nil {
			utils.LoggerError(err)
			continue
		}
		bsu[i] = models.BackupStorageUser{StorageId: backupStorage.Id, UserId: userId}
	}
	if len(bsu) > 0 {
		_, _ = bs.engine.Insert(&bsu)
	}
	return nil
}

func (bs *backupService) StorageDelete(storageId int) error {
	if storageId <= 0 {
		return errors.New("storage id must > 0")
	}
	s := models.BackupStorage{Id: storageId}
	exist, err := bs.engine.Cols("name").Get(&s)
	if !exist {
		return fmt.Errorf("not found backup storage %v, error: %v", s.Id, err)
	}
	exist, err = bs.engine.Where("storage_id = ?", storageId).And("close = false").Exist(new(models.BackupTask))
	if err != nil {
		return err
	}
	if exist {
		return errors.New("the backup storage is in use")
	}
	err = bs.cs.DeleteDynamicResource(constant.SecretYaml(bs.cs.GetNameSpace(), s.Name))
	if err != nil && !utils.ErrorContains(err, "not found") {
		return err
	}
	_, err = bs.engine.ID(storageId).Delete(new(models.BackupStorage))
	if err != nil {
		return err
	}
	_, _ = bs.engine.Where("storage_id = ?", storageId).Delete(new(models.BackupStorageUser))
	return nil
}

func (bs *backupService) StorageReconnect(storageId int) error {
	var storage = models.BackupStorage{Id: storageId}
	var exist, err = bs.engine.Get(&storage)
	if !exist || err != nil {
		return fmt.Errorf("not found storage %v, error: %v", storageId, err)
	}
	storage.SetStatus()
	return nil
}

func (bs *backupService) List(page, pageSize, storageId, clusterId int, startTime, endTime int64, search, status, t, userTag string) ([]models.BackupJobView, int, error) {
	var sql = NewSql("select bj.*, bt.type, bt.keep_copy, bt.crontab, ci.name cluster_name, bs.name storage_name, bt.user_id, ci.image_id, ci.storage old_storage, ci.connect_string from backup_job bj inner join backup_task bt on bj.backup_task_id = bt.id inner join cluster_instance ci on bt.cluster_id = ci.id inner join backup_storage bs on bt.storage_id = bs.id where 1=1")
	if len(userTag) != 0 && userTag != "AAAA" {
		u := models.User{UserTag: userTag}
		exist, err := bs.engine.Cols("id").Get(&u)
		if !exist {
			return nil, 0, fmt.Errorf("not found user tag %v, error: %v", userTag, err)
		}
		sql.And("bt.user_id = ?", u.Id)
	}
	if len(search) != 0 {
		sql.And("bj.job_name like ?", "%"+search+"%")
	}
	if storageId > 0 {
		sql.And("bt.storage_id = ?", storageId)
	}
	if clusterId > 0 {
		sql.And("bt.cluster_id = ?", clusterId)
	}
	if startTime > 0 && endTime > 0 {
		sql.And("bj.create_time between TO_TIMESTAMP(?) and TO_TIMESTAMP(?)", startTime, endTime)
	}
	if len(status) != 0 {
		sql.And("bj.status = ?", status)
	}
	if len(t) != 0 {
		sql.And("bt.type = ?", t)
	}
	sql.Raw("order by bj.id desc")
	var list = make([]models.BackupJobView, 0)
	err := sql.Session(bs.engine).Find(&list)
	if err != nil {
		return nil, 0, err
	}
	count := len(list)
	if utils.MustInt(page, pageSize) {
		min := pageSize * (page - 1)
		max := min + pageSize
		list = list[min:utils.Min(max, len(list))]
	}
	for i := range list {
		list[i].FormatDate()
	}
	return list, count, nil
}

func (bs *backupService) Create(task *models.BackupTask) (err error) {
	if !utils.MustInt(task.StorageId, task.ClusterId) {
		return errors.New("some parameters cannot be 0, [clusterId, storageId]")
	}
	var dbCluster = models.ClusterInstance{Id: task.ClusterId}
	exist, err := bs.engine.Cols("k8s_name", "user_id").Get(&dbCluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", dbCluster.Id, err)
	}

	var u = models.User{Id: dbCluster.UserId}
	exist, err = bs.engine.Cols("backup_max").Get(&u)
	if !exist {
		return fmt.Errorf("not found user %v, error: %v", u.Id, err)
	}
	useBackup, err := GetUseBackup(u.Id, bs.engine)
	if err != nil {
		return err
	}

	// 先判断是否已经创建周期备份，如果存在周期备份，手工备份的备份存储需要和周期备份的保持一致
	cycleTask := models.BackupTask{ClusterId: task.ClusterId, Type: models.BackupTypeCycle}
	existCycle, err := bs.engine.Cols("id", "storage_id", "close", "keep_copy").Get(&cycleTask)
	if err != nil {
		return
	}
	if existCycle && !cycleTask.Close && task.Type == models.BackupTypeOnce && cycleTask.StorageId != task.StorageId {
		return errors.New("cycle backup has been created, backup storage is different")
	}
	switch task.Type {
	case models.BackupTypeOnce:
		useBackup++
	case models.BackupTypeCycle:
		useBackup += task.KeepCopy
		if existCycle {
			useBackup -= cycleTask.KeepCopy
		}
	}
	if useBackup > u.BackupMax {
		return response.ErrorBackupMax
	}

	// 当不存在周期备份或者备份存储不一致时需要更新存储
	var updateStorage = !existCycle || cycleTask.StorageId != task.StorageId
	backupStorage := models.BackupStorage{Id: task.StorageId}
	if updateStorage {
		exist, err = bs.engine.Cols("name", "bucket").Get(&backupStorage)
		if !exist {
			return fmt.Errorf("not found backup storage %v, error: %v", backupStorage.Id, err)
		}
	}

	cluster, _, err := bs.cs.GetDynamicResource(constant.MysqlClusterYaml(bs.cs.GetNameSpace(), ""), dbCluster.K8sName)
	if err != nil {
		return err
	}
	var updateCluster bool
	clusterSpec := cluster.Object["spec"].(map[string]interface{})
	if updateStorage {
		updateCluster = true
		clusterSpec["backupSecretName"] = backupStorage.Name
		clusterSpec["backupURL"] = backupStorage.Bucket
		clusterSpec["backupRemoteDeletePolicy"] = "delete"
	}

	// 一个实例只能有一个周期任务, 如果已存在则更新
	var updateDB bool
	switch task.Type {
	case models.BackupTypeCycle:
		if task.KeepCopy == 0 {
			return errors.New("keep copy cannot be 0")
		}
		updateDB = existCycle
		clusterSpec["backupScheduleJobsHistoryLimit"] = task.KeepCopy
		clusterSpec["backupSchedule"] = task.Crontab
		updateCluster = true
	case models.BackupTypeOnce:
		if !utils.StringLength(task.Name, 1, 30) {
			return errors.New("the length of the backup name needs to be 1-30")
		}
		onceBackup := fmt.Sprintf(`apiVersion: mysql.presslabs.org/v1alpha1
kind: MysqlBackup
metadata:
  namespace: %v
  name: %v
spec:
  clusterName: %v`, bs.cs.GetNameSpace(), task.Name, dbCluster.K8sName)
		_, err = bs.cs.CreateDynamicResource(onceBackup)
		if err != nil {
			return
		}
	}

	if updateCluster {
		err = bs.cs.UpdateDynamicResource(constant.MysqlClusterYaml(bs.cs.GetNameSpace(), ""), cluster)
		if err != nil {
			return err
		}
	}

	if updateDB {
		_, err = bs.engine.ID(cycleTask.Id).MustCols("close", "set_type", "set_date", "set_time").Update(task)
	} else {
		task.UserId = dbCluster.UserId
		_, err = bs.engine.Insert(task)
	}
	return
}

func (bs *backupService) Delete(jobId int) error {
	job := models.BackupJob{Id: jobId}
	exist, err := bs.engine.Cols("job_name", "backup_task_id").Get(&job)
	if !exist {
		return fmt.Errorf("not found job %v, error: %v", jobId, err)
	}
	y := constant.MysqlBackupYaml(bs.cs.GetNameSpace(), job.JobName)
	err = bs.cs.DeleteDynamicResource(y)
	if err != nil && !utils.ErrorContains(err, "not found") {
		return err
	}
	_, err = bs.engine.ID(jobId).Delete(&job)
	if err != nil {
		return err
	}
	task := models.BackupTask{Id: job.BackupTaskId}
	exist, _ = bs.engine.Cols("type").Get(&task)
	if exist && task.Type == models.BackupTypeOnce {
		_, _ = bs.engine.ID(task.Id).Delete(&task)
	}
	return nil
}

func (bs *backupService) DeleteCycle(clusterId int) error {
	dbCluster := models.ClusterInstance{Id: clusterId}
	exist, err := bs.engine.Cols("k8s_name").Get(&dbCluster)
	if !exist {
		return fmt.Errorf("not found cluster %v, error: %v", clusterId, err)
	}
	cluster, _, err := bs.cs.GetDynamicResource(constant.MysqlClusterYaml(bs.cs.GetNameSpace(), ""), dbCluster.K8sName)
	if err != nil {
		return err
	}
	clusterSpec := cluster.Object["spec"].(map[string]interface{})
	delete(clusterSpec, "backupScheduleJobsHistoryLimit")
	delete(clusterSpec, "backupSchedule")
	err = bs.cs.UpdateDynamicResource(constant.MysqlClusterYaml(bs.cs.GetNameSpace(), ""), cluster)
	if err != nil && !utils.ErrorContains(err, "No content found to be updated") {
		return err
	}
	_, err = bs.engine.Where("cluster_id = ?", clusterId).And("type = ?", models.BackupTypeCycle).Cols("close").Update(&models.BackupTask{Close: true})
	return err
}
