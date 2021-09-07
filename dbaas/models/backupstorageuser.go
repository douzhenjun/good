package models

type BackupStorageUser struct {
	Id        int `xorm:"notnull pk autoincr unique INT" json:"-"`
	StorageId int `xorm:"notnull INT" json:"-"`
	UserId    int `xorm:"notnull INT" json:"userId"`
}
