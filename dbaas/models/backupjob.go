package models

import "time"

const (
	BackupStatusCompleted = "Completed"
	BackupStatusRunning   = "Running"
	BackupStatusError     = "Error"
)

const (
	BackupTypeCycle = "cycle"
	BackupTypeOnce  = "once"
)

type BackupJob struct {
	Id           int       `xorm:"notnull pk autoincr unique" json:"id"`
	JobName      string    `xorm:"varchar(100) unique" json:"name"`
	PodName      string    `xorm:"varchar(100)" json:"-"`
	Status       string    `xorm:"varchar(35)" json:"status"`
	CreateTime   time.Time `xorm:"notnull" json:"-"`
	Duration     int       `json:"duration"`
	BackupSet    string    `xorm:"varchar(100)" json:"backupSet"`
	BackupTaskId int       `xorm:"notnull" json:"-"`

	CreateTimeF string `xorm:"-" json:"createTime"`
}

type BackupJobView struct {
	BackupJob     `xorm:"extends"`
	BackupTask    `xorm:"extends"`
	StorageName   string `json:"storageName"`
	ClusterName   string `json:"clusterName"`
	ImageId       int    `json:"imageId"`
	OldStorage    int    `json:"oldStorage"`
	ConnectString int    `json:"nodeport"`
}

func (bjv *BackupJobView) FormatDate() {
	bjv.CreateTimeF = bjv.CreateTime.Format("2006-01-02 15:04:05")
}
