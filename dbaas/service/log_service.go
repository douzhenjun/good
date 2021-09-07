/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: ddh
 * @LastEditTime: 2021-02-07 16:32:07
 */

package service

import (
	"DBaas/models"
	"DBaas/utils"
	"github.com/go-xorm/xorm"
	"time"
)

type LogService interface {
	List(page int, pageSize int, key string, typeString string) ([]models.OperLog, int64, error)
	Add(level string, logSource string, people string, content string) error
}

type logService struct {
	Engine *xorm.Engine
	cs     CommonService
}

func NewLogService(engine *xorm.Engine) LogService {
	return &logService{
		Engine: engine,
	}
}

func (ls *logService) List(page int, pageSize int, key string, typeString string) ([]models.OperLog, int64, error) {
	opeLog := make([]models.OperLog, 0)
	err := ls.Engine.Where("log_source like ?", "%"+typeString+"%").And("content like ?", "%"+key+"%").Limit(pageSize, (page-1)*pageSize).Desc("oper_date").Find(&opeLog)
	if err != nil {
		return nil, 0, err
	}
	count, err := ls.Engine.Where("log_source like ?", "%"+typeString+"%").And("content like ?", "%"+key+"%").Count(&models.OperLog{})
	if err != nil {
		return nil, 0, err
	}
	for i := range opeLog {
		opeLog[i].ToResult()
	}
	return opeLog, count, nil
}

func (ls *logService) Add(level string, logSource string, people string, content string) error {
	log := models.OperLog{
		Content:    content,
		LogSource:  logSource,
		TypeLevel:  level,
		OperPeople: people,
		OperDate:   time.Now(),
	}
	_, err := ls.Engine.Insert(&log)
	utils.LoggerError(err)
	return err
}
