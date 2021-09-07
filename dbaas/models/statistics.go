package models

import (
	"time"
)

// StatisticsCluster 集群统计数据
type StatisticsCluster struct {
	Id            int       `xorm:"notnull pk autoincr unique" json:"id"`
	Cluster       string    `xorm:"notnull varchar(50) index" json:"cluster"`
	DeployTimeout bool      `xorm:"default false" json:"deployTimeout"`
	Type          string    `xorm:"notnull varchar(10)" json:"type"`
	Replicas      int       `json:"replicas"`
	DeployStart   time.Time `json:"-"`
	DeployEnd     time.Time `json:"-"`
	DeleteAt      time.Time `json:"-"`

	DeployStartF   string `xorm:"-" json:"deployStart"`
	DeployEndF     string `xorm:"-" json:"deployEnd"`
	DeleteAtF      string `xorm:"-"  json:"deleteAt"`
	DeployDuration int    `xorm:"-" json:"deployDuration"`
	UseDuration    int    `xorm:"-" json:"useDuration"`
}

func (sc *StatisticsCluster) Format() {
	f := "2006-01-02 15:04:05"
	sc.DeployStartF = sc.DeployStart.Format(f)
	sc.DeployEndF = sc.DeployEnd.Format(f)
	if !sc.DeleteAt.IsZero() {
		sc.DeleteAtF = sc.DeleteAt.Format(f)
	}
	sc.DeployDuration = int(sc.DeployEnd.Sub(sc.DeployStart).Seconds())
	var deleteEnd time.Time
	if sc.DeleteAt.IsZero() {
		deleteEnd = time.Now()
	} else {
		deleteEnd = sc.DeleteAt
	}
	sc.UseDuration = int(deleteEnd.Sub(sc.DeployStart).Seconds())
}
