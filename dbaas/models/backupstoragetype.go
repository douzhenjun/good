package models

type BackupStorageType struct {
	Id   int    `xorm:"notnull pk autoincr unique" json:"id"`
	Type string `xorm:"notnull varchar(20)" json:"type"`
}
