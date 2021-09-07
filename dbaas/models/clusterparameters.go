package models

type Clusterparameters struct {
	Id             int    `xorm:"not null pk autoincr unique INTEGER" json:"id"`
	ParameterName  string `xorm:"VARCHAR(100)" json:"parameterName"`
	ParameterValue string `xorm:"VARCHAR(200)" json:"parameterValue"`
	ClusterId      int    `xorm:"INTEGER" json:"-"`
}
