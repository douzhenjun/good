package models

import (
	"fmt"
)

/*
ApiQuota 接口配额
*/
type ApiQuota struct {
	Id      int    `xorm:"notnull pk autoincr unique" json:"id"`
	Path    string `xorm:"notnull varchar(50)" json:"path"`
	Cpu     int    `xorm:"notnull default 0" json:"cpu"`
	Memory  int    `xorm:"notnull default 0" json:"memory"`
	Storage int    `xorm:"notnull default 0" json:"storage"`
}

type ApiQuotaView struct {
	Id      int    `json:"id"`
	Path    string `json:"path"`
	Cpu     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

func (aq *ApiQuota) ToUsageView(usage *ApiQuota) *ApiQuotaView {
	return &ApiQuotaView{
		Id:      aq.Id,
		Path:    aq.Path,
		Cpu:     fmt.Sprintf("%v/%v", usage.Cpu, aq.Cpu),
		Memory:  fmt.Sprintf("%vG/%vG", usage.Memory, aq.Memory),
		Storage: fmt.Sprintf("%vG/%.2fTB", usage.Storage, float64(aq.Storage)/1024),
	}
}
