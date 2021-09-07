package models

import (
	"encoding/json"
	"github.com/go-xorm/xorm"
)

const (
	ScTypeUnique = "unique-storage"
	ScTypeShared = "shared-storage"
)

type Sc struct {
	Id            int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Name          string `xorm:"VARCHAR(40) unique" json:"name"`
	ScType        string `xorm:"VARCHAR(60) not null" json:"scType"`
	NodeNum       int    `xorm:"INTEGER not null" json:"nodeNum"`
	Describe      string `xorm:"VARCHAR(100)" json:"describe"`
	ReclaimPolicy string `xorm:"VARCHAR(200)" json:"reclaimPolicy"`
	UserTag       string `xorm:"VARCHAR(10)" json:"userTag"`
	OrgTag        string `xorm:"VARCHAR(10)" json:"orgTag"`
	AssignAll     bool   `xorm:"notnull default false" json:"-"`
}

type ReturnSc struct {
	Sc
	Children []PersistentVolume `xorm:"-" json:"children"`
	Cluster  []ClusterInstance  `xorm:"-" json:"cluster"`
	ScUser   json.RawMessage    `xorm:"-" json:"scUser"`
}

func (sc *Sc) CheckNodeNum(engine *xorm.Engine) {
	if sc.ScType != "unique-storage" {
		return
	}
	count, err := engine.Where(" sc_id = ?", sc.Id).Count(new(PersistentVolume))
	if err != nil {
		return
	}
	sc.NodeNum = int(count)
}
