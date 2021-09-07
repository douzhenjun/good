package configs

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// 服务端配置
type AppConfig struct {
	Name            string     `json:"name" yaml:"name"`
	Port            string     `json:"port" yaml:"port"`
	DataBase        []DataBase `json:"data_base" yaml:"data_base"`
}

/**
 * pg配置
 */
type DataBase struct {
	Drive    string `json:"drive" yaml:"drive"`
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Pwd      string `json:"pwd" yaml:"pwd"`
	Database string `json:"database" yaml:"database"`
}

// 读取服务器配置
func ReadConfig() error {
	content, err := ioutil.ReadFile("configs.yaml")
	if err != nil {
		return err
	}
	conf = &AppConfig{}
	return yaml.Unmarshal(content, conf)
}

func init() {
	err := ReadConfig()
	if err != nil {
		panic(err)
	}
}

// 返回主机数据库配置
var conf *AppConfig

func GetConfig() *AppConfig {
	return conf
}

