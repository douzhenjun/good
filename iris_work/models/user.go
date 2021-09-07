package models

import "time"

type User struct {
	Id            int64     `xorm:"pk autoincr" json:"id"`         //主键用户ID
	UserName     string     `xorm:"varchar(12)" json:"username"`  //用户名称
	RegisterTime  time.Time  `json:"register_time"`                //用户注册时间
	Mobile        string     `xorm:"varchar(11)" json:"mobile"`    //用户的移动手机号
	Pwd          string     `json:"password"`                  //用户的账户密码
	CityName      string     `xorm:"varchar(24)" json:"city_name"` //用户所在城市的名称
}
