/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: Dou
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: Dou
 * @LastEditTime: 2021-02-07 16:32:07
 */

package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// 服务端配置
type AppConfig struct {
	Name            string     `json:"name" yaml:"name"`
	Port            string     `json:"port" yaml:"port"`
	GrpcServerPort  string     `json:"grpc_server_port" yaml:"grpc_server_port"`
	Mode            string     `json:"mode" yaml:"mode"`
	DataBase        []DataBase `json:"data_base" yaml:"data_base"`
	CAddress        string     `json:"cmdb_address" yaml:"cmdb_address"`
	ReportPath      string     `json:"report_path" yaml:"report_path"`
	RabbitmqAddress []Rabbitmq `json:"rabbitmq" yaml:"rabbitmq"`
}

/*
配额限制
*/
type Quota struct {
	Cpu     int `json:"cpu" yaml:"cpu"`
	Mem     int `json:"mem" yaml:"mem"`
	Storage int `json:"storage" yaml:"storage"`
}

/**
 * Oracle配置
 */
type DataBase struct {
	Drive    string `json:"drive" yaml:"drive"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Pwd      string `json:"pwd" yaml:"pwd"`
	Host     string `json:"host" yaml:"host"`
	Database string `json:"database" yaml:"database"`
}

/**
 * Rabbitmq配置
 */
type Rabbitmq struct {
	VirtualHost string `json:"virtual_host" yaml:"virtual_host"`
	Port        string `json:"port" yaml:"port"`
	User        string `json:"user" yaml:"user"`
	Pwd         string `json:"pwd" yaml:"pwd"`
	Host        string `json:"host" yaml:"host"`
}

// 读取服务器配置
func ReadConfig() error {
	content, err := ioutil.ReadFile("./configs.yaml")
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
	// 每隔10分钟读取一次配置
	//utils.LoopTask(func() { utils.LoggerError(ReadConfig()) }, time.Minute*10)
}

var conf *AppConfig

func GetConfig() *AppConfig {
	return conf
}
