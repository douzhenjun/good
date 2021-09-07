package models

type Initinfo struct {
	Id       int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Name     string `xorm:"VARCHAR(20) unique" json:"name"`
	Message  string `xorm:"VARCHAR" json:"message"`
	Isaccess string `xorm:"VARCHAR(20)" json:"isaccess"`
	Isdeploy string `xorm:"VARCHAR(20)" json:"isdeploy"`
}
