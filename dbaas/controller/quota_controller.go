package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type QuotaController struct {
	Ctx     iris.Context
	Service service.QuotaService
}

func (qc *QuotaController) GetList() mvc.Result {
	data, err := qc.Service.ApiList()
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"detail": data})
}

func (qc *QuotaController) GetUsageList() mvc.Result {
	data, err := qc.Service.ApiUsageList()
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{"detail": data})
}

func (qc *QuotaController) PostEdit() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	var quota = new(models.ApiQuota)
	err = qc.Ctx.ReadForm(quota)
	if err != nil {
		return response.Error(err)
	}
	err = qc.Service.EditQuota(quota)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}
