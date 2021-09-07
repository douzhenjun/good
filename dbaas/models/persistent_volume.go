package models

import "time"

type PersistentVolume struct {
	Id            int       `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Capacity      string    `xorm:"VARCHAR(20)" json:"capacity"`
	Name          string    `xorm:"VARCHAR(50)" json:"name"`
	IpAddr        string    `xorm:"VARCHAR(20)" json:"ipAddr,omitempty"`
	Port          string    `xorm:"VARCHAR(20)" json:"port,omitempty"`
	Iqn           string    `xorm:"VARCHAR(40)" json:"iqn,omitempty"`
	Lun           int       `xorm:"INTEGER" json:"lun,omitempty"`
	OrgTag        string    `xorm:"VARCHAR(10)" json:"orgTag"`
	UserTag       string    `xorm:"VARCHAR(10)" json:"userTag"`
	ReclaimPolicy string    `xorm:"VARCHAR(10)" json:"reclaimPolicy"`
	Status        string    `xorm:"VARCHAR(10)" json:"status"`
	PvcName       string    `xorm:"VARCHAR(50)" json:"pvcName"`
	ScId          int       `xorm:"INTEGER" json:"-"`
	PodId         int       `xorm:"INTEGER" json:"-"`
	DeletedAt     time.Time `xorm:"deleted default NULL" json:"-"`

	PodName string `xorm:"-" json:"podName"`
	SCName  string `xorm:"-" json:"scName"`

	// 用户信息
	Tenant   string `xorm:"-" json:"tenant"`
	UserId   int    `xorm:"-" json:"userId"`
	CpuTotal int    `xorm:"-" json:"cpuTotal"`
	MemTotal int    `xorm:"-" json:"memTotal"`
}
