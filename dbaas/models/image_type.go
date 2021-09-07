package models

type ImageType struct {
	Id       int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Type     string `xorm:"VARCHAR(100)" json:"type"`
	Category string `xorm:"VARCHAR(100)" json:"category"`
}
