package models

type ScUser struct {
	Id     int `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	ScId   int `xorm:"INTEGER" json:"scId"`
	UserId int `xorm:"INTEGER" json:"userId"`

	UserName string `xorm:"-" json:"userName"`
}
