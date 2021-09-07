package models

import "encoding/json"

const (
	BackupStorageStatusTrue    = "True"
	BackupStorageStatusFalse   = "False"
	BackupStorageStatusUnknown = "Unknown"
)

type BackupStorage struct {
	Id        int    `xorm:"notnull pk autoincr unique INT" json:"id"`
	Name      string `xorm:"notnull unique VARCHAR(50)" json:"name"`
	Type      string `xorm:"notnull VARCHAR(10)" json:"type"`
	EndPoint  string `xorm:"notnull VARCHAR(50)" json:"endPoint"`
	Bucket    string `xorm:"notnull VARCHAR(50)" json:"bucket"`
	AccessKey string `xorm:"notnull VARCHAR(50)" json:"-"`
	SecretKey string `xorm:"notnull VARCHAR(100)" json:"-"`
	Status    string `xorm:"VARCHAR(10)" json:"status"`
	AssignAll bool   `xorm:"default false" json:"-"`

	UserIds json.RawMessage `xorm:"-" json:"userIds"`
	Serial  int             `xorm:"-" json:"serial"` //序号
}

func (bs *BackupStorage) SetStatus() {
	// TODO: S3存储状态验证
	bs.Status = BackupStorageStatusUnknown
}
