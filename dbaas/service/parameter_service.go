/**
* @Description:
* @version:
* @Company: iwhalecloud
* @Author:  zhangwei
* @Date: 2021/02/24 10:30
* @LastEditors: zhangwei
* @LastEditTime: 2021/02/26 13:30
**/
package service

import (
	"DBaas/utils"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	//"github.com/kataras/iris/v12"
	"DBaas/models"
	//"DBaas/utils"
)

//  用户服务接口定义
type ParameterService interface {
	SelectOne(id int) models.Sysparameter
	SelectOneByKey(key string) models.Sysparameter
	ResetParameter(id int) bool
	ModifyParameter(paramsList []map[string]interface{}) bool
	ListParameter(limit int, offset int, key string) ([]models.Sysparameter, error)
	ListParameterAll(key string) ([]models.Sysparameter, error)
	GetParameterCount(key string) (int64, error)
	ModifyNamespaceIsModifiable()
}

//  创建主机服务的接口
func NewParameterService(db *xorm.Engine) ParameterService {
	return &parameterService{
		Engine: db,
	}
}

//  主机服务结构体
type parameterService struct {
	Engine *xorm.Engine
}

func (ps *parameterService) SelectOne(id int) models.Sysparameter {
	var param models.Sysparameter
	_, err := ps.Engine.Where(" id = ? ", id).Get(&param)
	utils.LoggerError(err)
	return param
}

func (ps *parameterService) SelectOneByKey(key string) models.Sysparameter {
	var param models.Sysparameter
	_, err := ps.Engine.Where(" param_key = ? ", key).Get(&param)
	utils.LoggerError(err)
	return param
}

//
func (ps *parameterService) ResetParameter(id int) bool {
	var parameter models.Sysparameter
	_, err := ps.Engine.Where(" id = ? ", id).Get(&parameter)
	utils.LoggerError(err)
	pa := models.Sysparameter{ParamValue: parameter.DefaultValue}
	_, err = ps.Engine.Where(" id = ? ", id).Update(&pa)
	utils.LoggerError(err)
	return err == nil
}

func (ps *parameterService) ModifyParameter(paramsList []map[string]interface{}) bool {
	session := ps.Engine.NewSession()
	if err := session.Begin(); err != nil {
		iris.New().Logger().Info(err.Error())
	}
	if len(paramsList) > 0 {
		for _, param := range paramsList {
			key := param["key"].(string)
			value := param["value"].(string)
			pa := models.Sysparameter{
				ParamValue: value,
			}
			var oldParam models.Sysparameter
			_, err := ps.Engine.Where(" param_key = ? ", key).Get(&oldParam)
			if err != nil {
				utils.LoggerError(err)
				return false
			}
			// 初始化无默认值时设置默认值
			if oldParam.DefaultValue == "" {
				oldParam.DefaultValue = value
				_, err := session.Where(" param_key = ? ", key).Update(&oldParam)
				if err != nil {
					session.Rollback()
					utils.LoggerError(err)
					return false
				}
			}
			if value == "" {
				_, err := session.Where(" param_key = ? ", key).Cols("param_value").Update(&pa)
				if err != nil {
					session.Rollback()
					utils.LoggerError(err)
					return false
				}
			} else {
				_, err := session.Where(" param_key = ? ", key).Update(&pa)
				if err != nil {
					session.Rollback()
					utils.LoggerError(err)
					return false
				}
			}
		}
	}
	session.Commit()
	return true
}

func (ps *parameterService) ListParameter(limit int, offset int, key string) ([]models.Sysparameter, error) {
	parameterList := make([]models.Sysparameter, 0)
	err := ps.Engine.Where(" param_key like ? ", "%"+key+"%").Or(" param_value like ? ", "%"+key+"%").Or(" default_value like ? ", "%"+key+"%").Limit(limit, offset).OrderBy("id").Find(&parameterList)
	utils.LoggerError(err)
	return parameterList, err
}

//
func (ps *parameterService) ListParameterAll(key string) ([]models.Sysparameter, error) {
	parameterList := make([]models.Sysparameter, 0)
	err := ps.Engine.Where(" param_key like ? ", "%"+key+"%").Or(" param_value like ? ", "%"+key+"%").Or(" default_value like ? ", "%"+key+"%").OrderBy("id").Find(&parameterList)
	utils.LoggerError(err)
	return parameterList, err
}

/**
 * 获取系统参数总数量
 */
func (ps *parameterService) GetParameterCount(key string) (int64, error) {
	count, err := ps.Engine.Where(" param_key like ? ", "%"+key+"%").Or(" param_value like ? ", "%"+key+"%").Or(" default_value like ? ", "%"+key+"%").Count(new(models.Sysparameter))
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ps *parameterService) ModifyNamespaceIsModifiable() {
	var pa models.Sysparameter
	_, err := ps.Engine.Where(" param_key = ? ", "kubernetes_namespace").Get(&pa)
	utils.LoggerError(err)
	pa.IsModifiable = false
	_, uperr := ps.Engine.Where(" param_key = ? ", "kubernetes_namespace").Cols("is_modifiable").Update(&pa)
	if uperr != nil {
		iris.New().Logger().Error(uperr.Error())
	}
}
