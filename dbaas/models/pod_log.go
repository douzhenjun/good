package models

import (
	"time"
)

type PodLog struct {
	Id             int       `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	OperDate       time.Time `xorm:"not null DATETIME" json:"-"`
	OperDateString string    `xorm:"-" json:"operDateString"`
	Type           string    `xorm:"not null VARCHAR(60)" json:"type"`
	Reason         string    `xorm:"VARCHAR(60)" json:"reason"`
	Name           string    `xorm:"VARCHAR(100) unique" json:"name"`
	From           string    `xorm:"VARCHAR(60)" json:"from"`
	Message        string    `xorm:"TEXT" json:"message"`
}

func (p *PodLog) DateFormat() {
	p.OperDateString = p.OperDate.Format("2006-01-02 15:04:05")
}
