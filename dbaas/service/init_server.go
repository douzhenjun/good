/**
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author:  zhangwei
 * @Date: 2020/11/16 10:30
 * @LastEditors: zhangwei
 * @LastEditTime: 2020/11/16 10:30
 **/
package service

import (
	"DBaas/models"
	"DBaas/utils"
	"fmt"
	"github.com/go-xorm/xorm"
)

//  初始化服务接口定义
type InitService interface {
	AddOrModifyInitinfo(name string, initinfo models.Initinfo) (bool, string)
	GetLastStep() (string, string)
	SelectOneByName(name string) (models.Initinfo, string)
	GetNodeById(id int) (models.Node, string)
	GetLastDeployStep() (string, string)
}

//  创建主机服务的接口
func NewInitService(db *xorm.Engine) InitService {
	return &initService{
		Engine: db,
	}
}

//  主机服务结构体
type initService struct {
	Engine *xorm.Engine
}

//  保存用户信息  新增用户
func (is *initService) AddOrModifyInitinfo(name string, initinfo models.Initinfo) (bool, string) {
	var init models.Initinfo
	_, err := is.Engine.Where(" name = ? ", name).Get(&init)
	if err != nil {
		utils.LoggerError(err)
		return false, err.Error()
	}

	if init.Name != "" {
		_, err := is.Engine.Id(init.Id).Update(&initinfo)
		if err != nil {
			utils.LoggerError(err)
			return false, err.Error()
		}
	} else {
		_, err := is.Engine.Insert(initinfo)
		if err != nil {
			utils.LoggerError(err)
			return false, err.Error()
		}
	}

	return true, ""
}

func (is *initService) GetLastStep() (string, string) {
	LastStep := "parameter"
	stepList := []string{"storage", "host", "image", "operator"}
	for _, step := range stepList {
		var Initinfo models.Initinfo
		_, err := is.Engine.Where(" name = ? ", step).Get(&Initinfo)
		if err != nil {
			utils.LoggerError(err)
			return step, err.Error()
		}
		if Initinfo.Name != "" {
			LastStep = Initinfo.Name
		}
		if Initinfo.Isaccess == "False" {
			break
		}
	}
	return LastStep, ""
}

func (is *initService) GetLastDeployStep() (string, string) {
	deployLastStep := "operator"
	LastStep, err := is.GetLastStep()
	if err != "" {
		fmt.Println(err)
	}
	if LastStep != "operator" {
		deployLastStep = LastStep
	} else {
		var operatorInitinfo models.Initinfo
		_, err := is.Engine.Where(" name = ? ", "operator").Get(&operatorInitinfo)
		if err != nil {
			utils.LoggerError(err)
			return deployLastStep, err.Error()
		}
		if operatorInitinfo.Isaccess == "True" {
			deployLastStep = "storage"
		}
		deployStepList := []string{"storage", "host", "image", "operator"}
		for _, step := range deployStepList {
			var Initinfo models.Initinfo
			_, err := is.Engine.Where(" name = ? ", step).Get(&Initinfo)
			if err != nil {
				utils.LoggerError(err)
				return step, err.Error()
			}
			if Initinfo.Name != "" && Initinfo.Isdeploy == "False" {
				deployLastStep = Initinfo.Name
				break
			}
		}
	}

	return deployLastStep, ""
}

func (is *initService) SelectOneByName(name string) (models.Initinfo, string) {
	var init models.Initinfo
	_, err := is.Engine.Where(" name = ? ", name).Get(&init)
	if err != nil {
		return init, err.Error()
	}
	return init, ""
}

func (is *initService) GetNodeById(id int) (models.Node, string) {
	var node models.Node
	_, err := is.Engine.Where(" id = ? ", id).Get(&node)
	if err != nil {
		return node, err.Error()
	}
	return node, ""
}
