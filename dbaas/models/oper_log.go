package models

import (
	"time"
)

type OperLog struct {
	Id             int       `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	OperDate       time.Time `xorm:"DATETIME" json:"-"`
	OperDateString string    `xorm:"-" json:"operDateString"`
	OperPeople     string    `xorm:"VARCHAR(40)" json:"operPeople"`
	Content        string    `xorm:"text" json:"content"`
	TypeLevel      string    `xorm:"VARCHAR(40)" json:"typeLevel"`
	LogSource      string    `xorm:"VARCHAR(100)" json:"logSource"`
}

func (o *OperLog) ToResult() {
	o.OperDateString = o.OperDate.Format("2006-01-02 15:04:05")
}
