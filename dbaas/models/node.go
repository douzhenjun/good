package models

type Node struct {
	Id         int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	NodeName   string `xorm:"unique VARCHAR(30)" json:"nodeName"`
	Status     string `xorm:"VARCHAR(30)" json:"status"`
	Age        string `xorm:"VARCHAR(20)" json:"age"`
	Label      string `xorm:"text" json:"label"`
	MgmtTag    string `xorm:"-" json:"mgmtTag"`
	ComputeTag string `xorm:"-"  json:"computeTag"`
	OrgTag     string `xorm:"VARCHAR(10)" json:"orgTag"`
	UserTag    string `xorm:"VARCHAR(10)" json:"userTag"`
}
