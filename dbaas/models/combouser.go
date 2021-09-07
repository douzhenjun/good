package models

type ComboUser struct {
	Id      int `xorm:"notnull pk autoincr unique" json:"-"`
	ComboId int `xorm:"notnull" json:"-"`
	UserId  int `xorm:"notnull" json:"-"`
}
