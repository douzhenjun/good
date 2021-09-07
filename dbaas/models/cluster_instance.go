package models

import (
	"DBaas/utils"
	"encoding/json"
	"time"

	"github.com/go-xorm/xorm"
)

type ClusterInstance struct {
	Id                 int       `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Name               string    `xorm:"not null VARCHAR(40)" json:"name"`
	K8sName            string    `xorm:"VARCHAR(40)" json:"-"`
	SecretName         string    `xorm:"VARCHAR(40)" json:"-"`
	InnerConnectString string    `xorm:"VARCHAR(40)" json:"innerConnectString,omitempty"`
	ConnectString      string    `xorm:"VARCHAR(40)" json:"connectString"`
	ConsolePort        string    `xorm:"VARCHAR(40)" json:"consolePort,omitempty"`
	Status             string    `xorm:"VARCHAR(40)" json:"status"`
	UserId             int       `xorm:"INTEGER" json:"userId"`
	ImageId            int       `xorm:"INTEGER" json:"-"`
	ScName             string    `xorm:"VARCHAR(100)" json:"scName,omitempty"`
	LimitMem           int       `xorm:"INTEGER" json:"limitMem"`
	LimitCpu           int       `xorm:"INTEGER" json:"limitCpu"`
	Storage            int       `xorm:"INTEGER" json:"storage"`
	Remark             string    `xorm:"TEXT" json:"remark,omitempty"`
	YamlText           string    `xorm:"TEXT" json:"-"`
	Replicas           string    `xorm:"VARCHAR(40)" json:"replicas"`
	ActualReplicas     string    `xorm:"VARCHAR(40)" json:"actualReplicas"`
	Master             string    `xorm:"VARCHAR(100)" json:"master,omitempty"`
	Operator           string    `xorm:"VARCHAR(40)" json:"operator"`
	PodStatus          string    `xorm:"TEXT" json:"-"`
	OrgTag             string    `xorm:"VARCHAR(10)" json:"orgTag,omitempty"`
	UserTag            string    `xorm:"VARCHAR(10)" json:"userTag,omitempty"`
	IsDeploy           bool      `xorm:"BOOL default false" json:"-"` // 是否是通过deployment部署的（pv恢复）
	PvId               int       `xorm:"INTEGER" json:"-"`
	Secret             string    `xorm:"VARCHAR(50)" json:"-"`
	ComboId            int       `xorm:"default 0" json:"comboId"`
	DeletedAt          time.Time `xorm:"deleted default NULL" json:"-"`

	PodName      string          `xorm:"-" json:"podName,omitempty"`
	PvName       string          `xorm:"-" json:"pvName,omitempty"`
	UserName     string          `xorm:"-" json:"userName,omitempty"`
	ImageName    string          `xorm:"-" json:"imageName,omitempty"`
	ScNodes      int             `xorm:"-" json:"sc_nodes,omitempty"`
	ScType       string          `xorm:"-" json:"scType,omitempty"`
	Events       []PodLog        `xorm:"-" json:"events,omitempty"`
	Instance     []Instance      `xorm:"-" json:"instance,omitempty"`
	SecretMap    json.RawMessage `xorm:"-" json:"secretMap,omitempty"`
	PodStatusMap json.RawMessage `xorm:"-" json:"podStatusMap,omitempty"`

	CycleInfo *CycleInfo `xorm:"-" json:"crontabInfo,omitempty"`
	Qos       *QosLite   `xorm:"-" json:"qos,omitempty"`
}

const (
	ClusterStatusCreating = "Creating"
	ClusterStatusTrue     = "True"
	ClusterStatusFalse    = "False"
	ClusterStatusNotFound = "NotFound"
	ClusterStatusDisable  = "Disable"
)

func (ci *ClusterInstance) SetOperator(operator string, engine *xorm.Engine) {
	ci.Operator = operator
	_, err := engine.ID(ci.Id).Cols("operator").Update(ci)
	utils.LoggerError(err)
}

func (ci *ClusterInstance) GetSelf(engine *xorm.Engine) *ClusterInstance {
	current := ClusterInstance{Id: ci.Id}
	_, err := engine.Get(&current)
	utils.LoggerError(err)
	return &current
}
