package models

import "encoding/json"

type Combo struct {
	Id        int `xorm:"notnull pk autoincr unique" json:"id"`
	ComboLite `xorm:"extends"`
	Qos       QosLite `xorm:"extends" json:"qos"`
	AssignAll bool    `xorm:"default false" json:"-"`

	UserList  json.RawMessage `xorm:"-" json:"userList,omitempty"`
	Available bool            `xorm:"-" json:"available"`
}

type ComboLite struct {
	Name    string `xorm:"varchar(20) unique" json:"name"`
	Cpu     int    `xorm:"notnull default 0" json:"cpu"`
	Mem     int    `xorm:"notnull default 0" json:"mem"`
	Storage int    `xorm:"notnull default 0" json:"storage"`
	Tags    string `xorm:"varchar(100)" json:"tags"`
	Remark  string `json:"remark"`
	Copy    int    `xorm:"notnull default 1" json:"copy"`
}

type ComboTag struct {
	Id     int    `xorm:"notnull pk autoincr unique" json:"id"`
	Name   string `xorm:"notnull varchar(20)" json:"name"`
	Preset bool   `xorm:"notnull default false" json:"preset"`
}
