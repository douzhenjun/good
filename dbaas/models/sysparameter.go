package models

type Sysparameter struct {
	Id           int    `xorm:"not null pk autoincr unique INTEGER"`
	ParamKey     string `xorm:"not null unique(unq_name_key) VARCHAR(16)"`
	ParamValue   string `xorm:"not null text"`
	DefaultValue string `xorm:"not null text"`
	IsModifiable bool   `xorm:"BOOL" json:"isModifiable"`
}
