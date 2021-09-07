package models

import (
	"github.com/go-xorm/xorm"
	"strconv"
)

type MiscConfig struct {
	Key   string `xorm:"VARCHAR(50) not null pk unique"`
	Value string `xorm:"VARCHAR(100)"`
}

/*
模糊搜索配置, 传入like模糊查询参数
*/
func SearchConfig(like string, engine *xorm.Engine) ([]MiscConfig, error) {
	list := make([]MiscConfig, 0)
	err := engine.Where("key like ?", like).Find(&list)
	return list, err
}

func GetConfig(key string, engine *xorm.Engine) (string, error) {
	mc := MiscConfig{Key: key}
	_, err := engine.Get(&mc)
	return mc.Value, err
}

func GetConfigInt(key string, engine *xorm.Engine) (int, error) {
	value, err := GetConfig(key, engine)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(value)
}

func GetConfigBool(key string, engine *xorm.Engine) (bool, error) {
	value, err := GetConfig(key, engine)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(value)
}

func SetConfig(key string, value string, engine *xorm.Engine) error {
	config := MiscConfig{Key: key}
	exist, err := engine.Get(&config)
	if err != nil {
		return err
	}
	oldValue := config.Value
	config.Value = value
	if !exist {
		_, err = engine.Insert(&config)
		return err
	}
	if oldValue != value {
		_, err = engine.Where("key = ?", config.Key).Cols("value").Update(&config)
	}
	return err
}

func SetConfigInt(key string, value int, engine *xorm.Engine) error {
	return SetConfig(key, strconv.Itoa(value), engine)
}

func SetConfigBool(key string, value bool, engine *xorm.Engine) error {
	return SetConfig(key, strconv.FormatBool(value), engine)
}
