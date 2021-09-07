package models

type RegularRole struct {
	Id      int    `xorm:"not null pk autoincr unique INTEGER"`
	Role    string `xorm:"VARCHAR(16)"`
	ResType string `xorm:"not null VARCHAR(20)"`
	ResName string `xorm:"not null VARCHAR(100)"`
	Level   int    `xorm:"not null INTEGER"`
}
