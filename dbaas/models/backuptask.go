package models

type BackupTask struct {
	Id        int    `xorm:"notnull pk autoincr unique" json:"-"`
	Type      string `xorm:"notnull varchar(10)" json:"type,omitempty"`
	KeepCopy  int    `json:"keepCopy,omitempty"`
	Crontab   string `xorm:"varchar(20)" json:"crontab,omitempty"`
	StorageId int    `xorm:"notnull" json:"-"`
	ClusterId int    `xorm:"notnull" json:"-"`
	UserId    int    `xorm:"notnull" json:"userId,omitempty"`
	Name      string `xorm:"varchar(30)" json:"-"` // 手动备份名称
	Close     bool   `xorm:"default false" json:"-"`

	SetType string `xorm:"varchar(10)" json:"setType,omitempty"`
	SetDate string `xorm:"varchar(10)" json:"setDate,omitempty"`
	SetTime string `xorm:"varchar(15)" json:"setTime,omitempty"`
}

type CycleInfo struct {
	BackupTask  `xorm:"extends"`
	StorageName string `json:"storageName,omitempty"`
}
