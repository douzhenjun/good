package models

type UserRegular struct {
	Id     int    `xorm:"not null pk autoincr unique INTEGER"`
	UserId int    `xorm:"INTEGER"`
	RoleId int    `xorm:"INTEGER"`
	Remake string `xorm:"VARCHAR(100)"`
}
