package models

type MysqlOperator struct {
	Id               int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	NodeName         string `xorm:"VARCHAR(40)" json:"nodeName"`
	Name             string `xorm:"VARCHAR(40) unique" json:"name"`
	Ready            string `xorm:"VARCHAR(20)" json:"ready"`
	ContainerStatus  string `xorm:"VARCHAR(20)" json:"containerStatus""`
	Ip               string `xorm:"VARCHAR(40)" json:"ip"`
	Replicas         string `xorm:"VARCHAR(40)" json:"replicas""`
	ServiceIP        string `xorm:"VARCHAR(40)" json:"serviceIP"`
	Status           string `xorm:"VARCHAR(20)" json:"status"`
	DeploymentStatus string `xorm:"VARCHAR(40)" json:"deploymentStatus"`
	OrgTag           string `xorm:"VARCHAR(10)" json:"orgTag"`
	UserTag          string `xorm:"VARCHAR(10)" json:"userTag"`
}
