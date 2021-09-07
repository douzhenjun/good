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
	"github.com/go-xorm/xorm"
	core1 "k8s.io/api/core/v1"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"strings"
	"time"
)

type NodeService interface {
	List(page int, pageSize int, key string) ([]models.Node, int64)
	AddLabel(id int, key string, value string) (bool, string)
	DeleteLabel(id int, key string) (bool, string)
	AsyncDbLabel()
	Operator(mode string, scName string) error
	OperatorPodList(page int, pageSize int, key string) ([]models.MysqlOperator, string, int64)
	OperatorLogList(filter string) ([]models.PodLog, int64)
	GetOperatorStatus() map[string]string
	ComputedSum() (int, int)
	GetOperatorMode() string
	OperatorImage() error
}

type nodeService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func (ns *nodeService) OperatorImage() error {
	imageList := make([]models.Images2, 0)
	err := ns.Engine.SQL("select i.*, it.category from images i inner join image_type it on it.id = i.image_type_id where it.type = 'Operator'").Find(&imageList)
	if err != nil {
		return err
	}
	for i := range imageList {
		if imageList[i].Status != "Valid" {
			return fmt.Errorf("%v:%v is not Valid", imageList[i].ImageName, imageList[i].Version)
		}
	}
	y := constant.Yaml("apps/v1", "StatefulSet", ns.cs.GetNameSpace(), "")
	operator, _, err := ns.cs.GetDynamicResource(y, "mysql-operator")
	if err != nil {
		return err
	}
	var ch bool
	for i := range imageList {
		spec, f, err := unstructured.NestedSlice(operator.Object, "spec", "template", "spec", "containers")
		if !f {
			return fmt.Errorf("not found spec -> template -> spec, error: %v", err)
		}
		dbImage := fmt.Sprintf("%v:%v", imageList[i].ImageName, imageList[i].Version)
		var s string
		switch imageList[i].Category {
		case "Operator":
			s = spec[0].(map[string]interface{})["image"].(string)
		case "Orchestrator":
			s = spec[1].(map[string]interface{})["image"].(string)
		case "sidecar5.7":
			s = spec[0].(map[string]interface{})["args"].([]interface{})[2].(string)
		case "sidecar8":
			s = spec[0].(map[string]interface{})["args"].([]interface{})[3].(string)
		}
		if dbImage != s[strings.Index(s, "/")+1:] {
			ch = true
			break
		}
	}
	if !ch {
		return response.ErrorOperatorNoChange
	}
	s, f, err := unstructured.NestedSlice(operator.Object, "spec", "volumeClaimTemplates")
	if !f || len(s) == 0 {
		return fmt.Errorf("not found spec -> volumeClaimTemplates, error: %v", err)
	}
	sc, f, err := unstructured.NestedString(s[0].(map[string]interface{}), "spec", "storageClassName")
	if !f {
		return fmt.Errorf("not found spec -> volumeClaimTemplates[0] -> spec -> storageClassName, error: %v", err)
	}
	replicas, err := models.GetConfigInt("operator@replicas", ns.Engine)
	if err != nil {
		return err
	}
	var mode string
	if replicas == 3 {
		mode = "multi"
	}
	return ns.Operator(mode, sc)
}

func NewNodeService(engine *xorm.Engine, cs CommonService) NodeService {
	return &nodeService{
		Engine: engine,
		cs:     cs,
	}
}

//label 同步数据库
func (ns *nodeService) AsyncDbLabel() {
	ns.cs.AsyncNodeInfo()
}

func (ns *nodeService) ComputedSum() (int, int) {
	node := make([]models.Node, 0)
	err := ns.Engine.Find(&node)
	utils.LoggerError(err)
	mgmtTag := 0
	computeTag := 0
	for i, m := range node {
		nodeLabelMap := map[string]interface{}{}
		err := json.Unmarshal([]byte(m.Label), &nodeLabelMap)
		if err != nil {
			utils.LoggerError(err)
		} else {
			if nodeLabelMap["iwhalecloud.dbassoperator"] == "mysqlha" {
				node[i].MgmtTag = "true"
				mgmtTag += 1
			} else {
				node[i].MgmtTag = "false"
			}

			if nodeLabelMap["iwhalecloud.dbassnode"] == "mysql" {
				node[i].ComputeTag = "true"
				computeTag += 1
			} else {
				node[i].ComputeTag = "false"
			}
		}

	}
	return mgmtTag, computeTag
}

// 获取node节点
func (ns *nodeService) List(page int, pageSize int, key string) ([]models.Node, int64) {
	node := make([]models.Node, 0)
	err := ns.Engine.Where("node_name like ? or label like ? ", "%"+key+"%", "%"+key+"%").Limit(pageSize, pageSize*(page-1)).OrderBy("id").Find(&node)
	count, err := ns.Engine.Where("node_name like ? or label like ? ", "%"+key+"%", "%"+key+"%").Count(&models.Node{})
	for i, m := range node {
		nodeLabelMap := map[string]interface{}{}
		err := json.Unmarshal([]byte(m.Label), &nodeLabelMap)
		if err != nil {
			utils.LoggerError(err)
		} else {
			if nodeLabelMap["iwhalecloud.dbassoperator"] == "mysqlha" {
				node[i].MgmtTag = "true"
			} else {
				node[i].MgmtTag = "false"
			}

			if nodeLabelMap["iwhalecloud.dbassnode"] == "mysql" {
				node[i].ComputeTag = "true"
			} else {
				node[i].ComputeTag = "false"
			}
		}

	}
	utils.LoggerError(err)
	return node, count
}

// 添加node label
func (ns *nodeService) AddLabel(id int, key string, value string) (bool, string) {
	node := models.Node{Id: id}
	success, err := ns.Engine.Get(&node)
	if !success {
		if err != nil {
			return false, err.Error()
		} else {
			return false, ""
		}
	}
	err, k8sNode := ns.cs.GetResources("node", node.NodeName, "default", meta1.GetOptions{})
	if err != nil {
		utils.LoggerError(err)
		return false, err.Error()
	}
	if k8sNode, ok := (*k8sNode).(*core1.Node); ok {
		labels := k8sNode.Labels
		labels[key] = value
		patchData := map[string]interface{}{"metadata": map[string]map[string]string{"labels": labels}}
		playLoadBytes, err := json.Marshal(patchData)
		if err != nil {
			utils.LoggerError(err)
			return false, err.Error()
		}
		err = ns.cs.PatchOption("node", node.NodeName, "default", playLoadBytes, meta1.PatchOptions{}, types.StrategicMergePatchType)
		if err != nil {
			utils.LoggerError(err)
			return false, err.Error()
		}
	} else {
		return false, ""
	}
	return true, ""
}

// 删除node label
func (ns *nodeService) DeleteLabel(id int, key string) (bool, string) {
	node := models.Node{Id: id}
	success, err := ns.Engine.Get(&node)
	if !success {
		if err != nil {
			return false, err.Error()
		} else {
			return false, ""
		}
	}
	err, k8sNode := ns.cs.GetResources("node", node.NodeName, "default", meta1.GetOptions{})
	if err != nil {
		utils.LoggerError(err)
		return false, err.Error()
	}
	type RemoveStringValue struct {
		Op   string `json:"op"`
		Path string `json:"path"`
	}
	if k8sNode, ok := (*k8sNode).(*core1.Node); ok {
		payloads := make([]RemoveStringValue, 0)
		payloads = append(payloads, RemoveStringValue{
			Op:   "remove",
			Path: "/metadata/labels/" + key,
		})
		payloadBytes, _ := json.Marshal(payloads)
		err = ns.cs.PatchOption("node", k8sNode.Name, "default", payloadBytes, meta1.PatchOptions{}, types.JSONPatchType)
		if err != nil {
			utils.LoggerError(err)
			return false, err.Error()
		}
	} else {
		return false, ""
	}
	return true, ""
}

// 部署Operator
func (ns *nodeService) Operator(mode string, scName string) error {
	_, err := ns.Engine.Where(" param_key = ? ", "kubernetes_namespace").
		Cols("is_modifiable").Update(&models.Sysparameter{IsModifiable: false})
	utils.LoggerError(err)
	err, _ = ns.cs.GetResources("sc", scName, ns.cs.GetNameSpace(), meta1.GetOptions{})
	if err != nil {
		return err
	}

	joinSql := "select images.* from images left join image_type it on images.image_type_id = it.id where it.type = ? and it.category = ?"
	var operatorImage, orchestratorImage, sidecar5Image, sidecar8Image models.Images
	var exist bool
	checkImage := func(t string, i models.Images) error {
		if !exist || err != nil {
			return fmt.Errorf("not found need images %v error: %v", t, err)
		}
		if i.Status == "Invalid" {
			return fmt.Errorf("%v image is invalid", t)
		}
		return nil
	}
	exist, err = ns.Engine.SQL(joinSql, "Operator", "Operator").Get(&operatorImage)
	if err := checkImage("operator", operatorImage); err != nil {
		return err
	}
	exist, err = ns.Engine.SQL(joinSql, "Operator", "Orchestrator").Get(&orchestratorImage)
	if err := checkImage("orchestrator", orchestratorImage); err != nil {
		return err
	}
	exist, err = ns.Engine.SQL(joinSql, "Operator", "sidecar5.7").Get(&sidecar5Image)
	if err := checkImage("sidecar5.7", sidecar5Image); err != nil {
		return err
	}
	exist, err = ns.Engine.SQL(joinSql, "Operator", "sidecar8").Get(&sidecar8Image)
	if err := checkImage("sidecar8", sidecar8Image); err != nil {
		return err
	}
	address := models.Sysparameter{ParamKey: "harbor_address"}
	exist, err = ns.Engine.Get(&address)
	if !exist || err != nil {
		return errors.New("not found need images address")
	}

	var deploymentYAML string
	nodeNum := 1
	if mode == "multi" {
		nodeNum = 3
	}
	envTZ := "Asia/Shanghai"
	sysTZ := models.Sysparameter{ParamKey: "operator_tz"}
	_, _ = ns.Engine.Cols("param_value").Get(&sysTZ)
	if len(sysTZ.ParamValue) != 0 {
		envTZ = sysTZ.ParamValue
	}
	deploymentYAML = fmt.Sprintf(`
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
    helm.sh/hook: crd-install
  name: mysqlbackups.mysql.presslabs.org
  namespace: %v
  labels:
    app: mysql-operator
spec:
  group: mysql.presslabs.org
  names:
    kind: MysqlBackup
    listKind: MysqlBackupList
    plural: mysqlbackups
    singular: mysqlbackup
  scope: Namespaced
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
    helm.sh/hook: crd-install
  name: mysqlclusters.mysql.presslabs.org
  namespace: %v
  labels:
    app: mysql-operator
spec:
  additionalPrinterColumns:
    - JSONPath: .status.conditions[?(@.type == 'Ready')].status
      description: The cluster status
      name: Ready
      type: string
    - JSONPath: .spec.replicas
      description: The number of desired nodes
      name: Replicas
      type: integer
    - JSONPath: .metadata.creationTimestamp
      name: Age
      type: date
  group: mysql.presslabs.org
  names:
    kind: MysqlCluster
    listKind: MysqlClusterList
    plural: mysqlclusters
    shortNames:
      - mysql
    singular: mysqlcluster
  scope: Namespaced
  subresources:
    scale:
      specReplicasPath: .spec.replicas
      statusReplicasPath: .status.readyNodes
    status: {}
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
    helm.sh/hook: crd-install
  name: mysqldatabases.mysql.presslabs.org
  namespace: %v
  labels:
    app: mysql-operator
spec:
  additionalPrinterColumns:
    - JSONPath: .status.conditions[?(@.type == 'Ready')].status
      description: The database status
      name: Ready
      type: string
    - JSONPath: .spec.clusterRef.name
      name: Cluster
      type: string
    - JSONPath: .spec.database
      name: Database
      type: string
    - JSONPath: .metadata.creationTimestamp
      name: Age
      type: date
  group: mysql.presslabs.org
  names:
    kind: MysqlDatabase
    listKind: MysqlDatabaseList
    plural: mysqldatabases
    singular: mysqldatabase
  scope: Namespaced
  subresources:
    status: {}
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.3.0
    helm.sh/hook: crd-install
  name: mysqlusers.mysql.presslabs.org
  namespace: %v
  labels:
    app: mysql-operator
spec:
  additionalPrinterColumns:
    - JSONPath: .status.conditions[?(@.type == 'Ready')].status
      description: The user status
      name: Ready
      type: string
    - JSONPath: .spec.clusterRef.name
      name: Cluster
      type: string
    - JSONPath: .spec.user
      name: UserName
      type: string
    - JSONPath: .metadata.creationTimestamp
      name: Age
      type: date
  group: mysql.presslabs.org
  names:
    kind: MysqlUser
    listKind: MysqlUserList
    plural: mysqlusers
    singular: mysqluser
  scope: Namespaced
  subresources:
    status: {}
  version: v1alpha1
  versions:
    - name: v1alpha1
      served: true
      storage: true
---
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: mysql-operator
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
spec:
  maxUnavailable: 1
  selector:
    matchLabels:
      app: mysql-operator
---
# Source: mysql-operator/templates/serviceaccount.yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: mysql-operator
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
---
apiVersion: v1
kind: Secret
metadata:
  name: mysql-operator-orc
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
data:
  TOPOLOGY_USER: "b3JjaGVzdHJhdG9y"
  TOPOLOGY_PASSWORD: "NUJGUU01UjdiVQ=="
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: mysql-operator-orc
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
data:
  orchestrator.conf.json: "{\n  \"ApplyMySQLPromotionAfterMasterFailover\": false,\n  \"BackendDB\": \"sqlite\",\n  \"Debug\": false,\n  \"DetachLostReplicasAfterMasterFailover\": true,\n  \"DetectClusterAliasQuery\": \"SELECT CONCAT(SUBSTRING(@@hostname, 1, LENGTH(@@hostname) - 1 - LENGTH(SUBSTRING_INDEX(@@hostname,'-',-2))),'.',SUBSTRING_INDEX(@@report_host,'.',-1))\",\n  \"DetectInstanceAliasQuery\": \"SELECT @@hostname\",\n  \"DiscoverByShowSlaveHosts\": false,\n  \"FailMasterPromotionIfSQLThreadNotUpToDate\": true,\n  \"HTTPAdvertise\": \"http://{{ .Env.HOSTNAME }}-svc:80\",\n  \"HostnameResolveMethod\": \"none\",\n  \"InstancePollSeconds\": 5,\n  \"ListenAddress\": \":3000\",\n  \"MasterFailoverDetachReplicaMasterHost\": true,\n  \"MySQLHostnameResolveMethod\": \"@@report_host\",\n  \"MySQLTopologyCredentialsConfigFile\": \"/etc/orchestrator/orc-topology.cnf\",\n  \"OnFailureDetectionProcesses\": [\n    \"/usr/local/bin/orc-helper event -w '{failureClusterAlias}' 'OrcFailureDetection' 'Failure: {failureType}, failed host: {failedHost}, lost replcas: {lostReplicas}' || true\",\n    \"/usr/local/bin/orc-helper failover-in-progress '{failureClusterAlias}' '{failureDescription}' || true\"\n  ],\n  \"PostIntermediateMasterFailoverProcesses\": [\n    \"/usr/local/bin/orc-helper event '{failureClusterAlias}' 'OrcPostIntermediateMasterFailover' 'Failure type: {failureType}, failed hosts: {failedHost}, slaves: {countSlaves}' || true\"\n  ],\n  \"PostMasterFailoverProcesses\": [\n    \"/usr/local/bin/orc-helper event '{failureClusterAlias}' 'OrcPostMasterFailover' 'Failure type: {failureType}, new master: {successorHost}, slaves: {slaveHosts}' || true\"\n  ],\n  \"PostUnsuccessfulFailoverProcesses\": [\n    \"/usr/local/bin/orc-helper event -w '{failureClusterAlias}' 'OrcPostUnsuccessfulFailover' 'Failure: {failureType}, failed host: {failedHost} with {countSlaves} slaves' || true\"\n  ],\n  \"PreFailoverProcesses\": [\n    \"/usr/local/bin/orc-helper failover-in-progress '{failureClusterAlias}' '{failureDescription}' || true\"\n  ],\n  \"ProcessesShellCommand\": \"sh\",\n  \"RaftAdvertise\": \"{{ .Env.HOSTNAME }}-svc\",\n  \"RaftBind\": \"{{ .Env.HOSTNAME }}\",\n  \"RaftDataDir\": \"/var/lib/orchestrator\",\n  \"RaftEnabled\": true,\n  \"RaftNodes\": [\n    \"mysql-operator-0-svc\",\n    \"mysql-operator-1-svc\",\n    \"mysql-operator-2-svc\"\n  ],\n  \"RecoverIntermediateMasterClusterFilters\": [\n    \".*\"\n  ],\n  \"RecoverMasterClusterFilters\": [\n    \".*\"\n  ],\n  \"RecoveryIgnoreHostnameFilters\": [],\n  \"RecoveryPeriodBlockSeconds\": 300,\n  \"RemoveTextFromHostnameDisplay\": \":3306\",\n  \"SQLite3DataFile\": \"/var/lib/orchestrator/orc.db\",\n  \"SlaveLagQuery\": \"SELECT TIMESTAMPDIFF(SECOND,ts,NOW()) as drift FROM sys_operator.heartbeat ORDER BY drift ASC LIMIT 1\",\n  \"UnseenInstanceForgetHours\": 1\n}"
  orc-topology.cnf: |
    [client]
    user = {{ .Env.ORC_TOPOLOGY_USER }}
    password = {{ .Env.ORC_TOPOLOGY_PASSWORD }}
---
# Source: mysql-operator/templates/clusterrole.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: 'mysql-operator'
  namespace: %v
  labels:
    app: 'mysql-operator'
    chart: 'mysql-operator-0.1.1_master'
    release: 'mysql-operator'
    heritage: 'Helm'
rules:
  - apiGroups:
      - apps
    resources:
      - statefulsets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - batch
    resources:
      - jobs
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - configmaps
      - events
      - jobs
      - persistentvolumeclaims
      - pods
      - secrets
      - services
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - ""
    resources:
      - pods/status
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mysql.presslabs.org
    resources:
      - mysqlbackups
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mysql.presslabs.org
    resources:
      - mysqlclusters
      - mysqlclusters/status
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mysql.presslabs.org
    resources:
      - mysqldatabases
      - mysqldatabases/status
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - mysql.presslabs.org
    resources:
      - mysqlusers
      - mysqlusers/status
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
  - apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
---
# Source: mysql-operator/templates/clusterrolebinding.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: mysql-operator
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: mysql-operator
subjects:
  - name: mysql-operator
    namespace: "default"
    kind: ServiceAccount
---
# Source: mysql-operator/templates/orc-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-operator-0-svc
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  type: ClusterIP
  ports:
  - name: web
    port: 80
    targetPort: 3000
  - name: raft
    port: 10008
    targetPort: 10008
  selector:
    statefulset.kubernetes.io/pod-name: mysql-operator-0
---
# Source: mysql-operator/templates/orc-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-operator-1-svc
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  type: ClusterIP
  ports:
  - name: web
    port: 80
    targetPort: 3000
  - name: raft
    port: 10008
    targetPort: 10008
  selector:
    statefulset.kubernetes.io/pod-name: mysql-operator-1
---
# Source: mysql-operator/templates/orc-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-operator-2-svc
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
  annotations:
    service.alpha.kubernetes.io/tolerate-unready-endpoints: "true"
spec:
  type: ClusterIP
  ports:
  - name: web
    port: 80
    targetPort: 3000
  - name: raft
    port: 10008
    targetPort: 10008
  selector:
    statefulset.kubernetes.io/pod-name: mysql-operator-2
---
# Source: mysql-operator/templates/orc-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-operator
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
spec:
  type: NodePort
  selector:
    app: mysql-operator
  ports:
  - name: http
    port: 80
    protocol: TCP
    targetPort: 3000
---
# Source: mysql-operator/templates/statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql-operator
  namespace: %v
  labels:
    app: mysql-operator
    chart: mysql-operator-0.1.1_master
    release: mysql-operator
    heritage: Helm
spec:
  replicas: %v
  serviceName: mysql-operator-orc
  podManagementPolicy: Parallel
  selector:
    matchLabels:
      app: mysql-operator
      release: mysql-operator
  template:
    metadata:
      labels:
        app: mysql-operator
        release: mysql-operator
      annotations:
        checksum/configs: b3e75caf30cfd2c28865d7366bb43815f657685dfe41f3958d0f3f4a94fb8e6f
        checksum/secret: 4522ddf748e2edf59c244434f5af103d44f888560c9d24705b419e6fa3b3f3ec
    spec:
      serviceAccountName: mysql-operator
      containers:
        - name: operator
          image: "%v"
          imagePullPolicy: IfNotPresent
          env:
            - name: TZ
              value: %v
            - name: ORC_TOPOLOGY_USER
              valueFrom:
                secretKeyRef:
                  name: mysql-operator-orc
                  key: TOPOLOGY_USER
            - name: ORC_TOPOLOGY_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: mysql-operator-orc
                  key: TOPOLOGY_PASSWORD
          args:
            - --leader-election-namespace=default
            # connect to orchestrator on localhost
            - --orchestrator-uri=http://127.0.0.1:3000/api
            - --sidecar-image=%v
            - --sidecar-mysql8-image=%v
            - --failover-before-shutdown=false
          resources:
            {}
          # TODO: add livenessProbe to controller
          # livenessProbe:
          #   httpGet:
          #     path: /health
          #     port: 80
        - name: orchestrator
          image: %v
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 3000
              name: web
              protocol: TCP
            - containerPort: 10008
              name: raft
              protocol: TCP
          env:
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
          envFrom:
            - prefix: ORC_
              secretRef:
                name: mysql-operator-orc
          volumeMounts:
            - name: data
              mountPath: /var/lib/orchestrator/
            - name: configs
              mountPath: /templates/
          livenessProbe:
            timeoutSeconds: 10
            initialDelaySeconds: 200
            httpGet:
              path: /api/lb-check
              port: 3000
          # https://github.com/github/orchestrator/blob/master/docs/raft.md#proxy-healthy-raft-nodes
          readinessProbe:
            timeoutSeconds: 10
            httpGet:
              path: /api/raft-health
              port: 3000
          resources:
            null

      volumes:
        - name: configs
          configMap:
            name: mysql-operator-orc

      # security context to mount corectly the volume for orc
      securityContext:
        fsGroup: 777
      nodeSelector:
        iwhalecloud.dbassoperator: mysqlha

      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - topologyKey: kubernetes.io/hostname
            labelSelector:
              matchLabels:
                app: mysql-operator
      tolerations:
        - effect: NoSchedule
          key: node-role.kubernetes.io/master
          operator: Equal
          value: ""
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: [ ReadWriteOnce ]
        resources:
          requests:
            storage: 10Gi
        storageClassName: "%v"
---
apiVersion: v1
kind: Secret
metadata:
  name: my-secret #密码对象名称
  namespace: %v
type: Opaque
data:
  # root password is required to be specified
  ROOT_PASSWORD: bm90LXNvLXNlY3VyZQ==   #base64编码
  # a user name to be created, not required
  USER: dXNlcm5hbWU= #base64编码
  # a password for user, not required
  PASSWORD: dXNlcnBhc3Nz
  # a name for database that will be created, not required
  DATABASE: dXNlcmRi #base64编码

`, ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), ns.cs.GetNameSpace(), nodeNum, fmt.Sprintf("%v/%v:%v", address.ParamValue, operatorImage.ImageName, operatorImage.Version), envTZ, fmt.Sprintf("%v/%v:%v", address.ParamValue, sidecar5Image.ImageName, sidecar5Image.Version), fmt.Sprintf("%v/%v:%v", address.ParamValue, sidecar8Image.ImageName, sidecar8Image.Version), fmt.Sprintf("%v/%v:%v", address.ParamValue, orchestratorImage.ImageName, orchestratorImage.Version), scName, ns.cs.GetNameSpace())
	flagIndex := 0
	yamlArray := strings.Split(deploymentYAML, "---")

	_, err = ns.Engine.Where("key like ?", "operator@%").Delete(new(models.MiscConfig))
	utils.LoggerError(err)
	deleteCount, err := ns.Engine.Delete(new(models.MysqlOperator))
	utils.LoggerError(err)
	for i := 0; i < int(deleteCount); i++ {
		ns.cs.ClearEvent(fmt.Sprintf("mysql-operator-%v", i))
	}

	_ = models.SetConfigInt("operator@step", len(yamlArray), ns.Engine)
	_ = models.SetConfigInt("operator@replicas", nodeNum, ns.Engine)
	for i, s := range yamlArray {
		if i >= 4 {
			_ = ns.cs.DeleteDynamicResource(s)
		}
	}
	var operatorStatus = true
	var operatorStepped int
	for i, value := range yamlArray {
		operatorStepped = i + 1
		if i == 1 {
			_, _, err = ns.cs.GetDynamicResource(value, "mysqlclusters.mysql.presslabs.org")
			if err != nil {
				utils.LoggerError(err)
				_, err = ns.cs.CreateDynamicResource(value)
			}
		} else if i == 0 {
			_, _, err = ns.cs.GetDynamicResource(value, "mysqlbackups.mysql.presslabs.org")
			if err != nil {
				utils.LoggerError(err)
				_, err = ns.cs.CreateDynamicResource(value)
			}
		} else if i == 2 {
			_, _, err = ns.cs.GetDynamicResource(value, "mysqldatabases.mysql.presslabs.org")
			if err != nil {
				utils.LoggerError(err)
				_, err = ns.cs.CreateDynamicResource(value)
			}
		} else if i == 3 {
			_, _, err = ns.cs.GetDynamicResource(value, "mysqlusers.mysql.presslabs.org")
			if err != nil {
				utils.LoggerError(err)
				_, err = ns.cs.CreateDynamicResource(value)
			}
		} else {
			_, err = ns.cs.CreateDynamicResource(value)
		}
		if err != nil {
			utils.LoggerError(err)
			operatorStatus = false
			flagIndex = i
			_ = models.SetConfig("operator@reason", err.Error(), ns.Engine)
			break
		}
	}
	_ = models.SetConfigInt("operator@stepped", operatorStepped, ns.Engine)
	_ = models.SetConfigBool("operator@status", operatorStatus, ns.Engine)

	// 部署失败, 回退
	if !operatorStatus {
		for i := 4; i <= flagIndex; i++ {
			_ = ns.cs.DeleteDynamicResource(yamlArray[i])
		}
		return err
	}

	_ = models.SetConfigBool("operator@creating", true, ns.Engine)
	// 设置部署超时10分钟
	utils.Timer(func() { _ = models.SetConfigBool("operator@creating", false, ns.Engine) }, 10*time.Minute)
	// 部署成功之后写入部署日志
	go ns.cs.AsyncOperatorLog()
	return nil
}

func (ns *nodeService) GetOperatorStatus() map[string]string {
	result := map[string]string{}
	list, _ := models.SearchConfig("operator@%", ns.Engine)
	for _, c := range list {
		result[c.Key[9:]] = c.Value
	}
	return result
}

func (ns *nodeService) OperatorPodList(page int, pageSize int, key string) ([]models.MysqlOperator, string, int64) {
	mysqlOperatorList := make([]models.MysqlOperator, 0)
	err := ns.Engine.Limit(pageSize, (page-1)*pageSize).Find(&mysqlOperatorList)
	if err != nil {
		utils.LoggerError(err)
		return mysqlOperatorList, err.Error(), 0
	}
	count, err := ns.Engine.Limit(pageSize, (page-1)*pageSize).Count(&models.MysqlOperator{})
	if err != nil {
		utils.LoggerError(err)
		return mysqlOperatorList, err.Error(), count
	}
	return mysqlOperatorList, "", count
}

func (ns *nodeService) OperatorLogList(filter string) ([]models.PodLog, int64) {
	return ns.cs.FilterEvent(filter + ".")
}

func (ns *nodeService) GetOperatorMode() string {
	replicas, _ := models.GetConfig("operator@replicas", ns.Engine)
	if replicas == "3" {
		return "multi"
	} else {
		return "single"
	}
}
