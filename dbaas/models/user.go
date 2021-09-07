package models

type User struct {
	Id         int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	ZdcpId     int    `xorm:"INTEGER" json:"-"`
	UserName   string `xorm:"not null VARCHAR(100)" json:"name"`
	Password   string `xorm:"VARCHAR(32)" json:"-"`
	MemAll     int64  `xorm:"BIGINT" json:"-"`
	CpuAll     int    `xorm:"INTEGER" json:"-"`
	StorageAll int    `xorm:"INTEGER" json:"-"`
	Remarks    string `xorm:"text" json:"-"`
	UserTag    string `xorm:"VARCHAR(100)" json:"-"`
	AutoCreate bool   `xorm:"BOOL default false" json:"-"`
	BackupMax  int    `json:"-"`
}
