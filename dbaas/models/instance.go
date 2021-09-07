package models

import (
	"encoding/json"
	"time"
)

type Instance struct {
	Id                  int             `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	Name                string          `xorm:"VARCHAR(100)" json:"name"`
	DomainName          string          `xorm:"VARCHAR(40)" json:"domainName"`
	ClusterId           int             `xorm:"INTEGER" json:"clusterId"`
	Version             string          `xorm:"VARCHAR(20)" json:"version"`
	Status              string          `xorm:"VARCHAR(20)" json:"status"`
	Role                string          `xorm:"VARCHAR(20)" json:"role"`
	Volume              string          `xorm:"TEXT" json:"-"`
	BaseInfo            string          `xorm:"TEXT" json:"-"`
	InitContainer       string          `xorm:"TEXT" json:"-"`
	ContainerInfo       string          `xorm:"TEXT" json:"-"`
	Events              []PodLog        `xorm:"-" json:"events"`
	VolumeObject        json.RawMessage `xorm:"-" json:"volumeObject"`
	BaseInfoObject      json.RawMessage `xorm:"-" json:"baseInfoObject"`
	InitContainerObject json.RawMessage `xorm:"-" json:"initContainerObject"`
	ContainerInfoObject json.RawMessage `xorm:"-" json:"containerInfoObject"`
	DeletedAt           time.Time       `xorm:"deleted default NULL" json:"-"`
}
