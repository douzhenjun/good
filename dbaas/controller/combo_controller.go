package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ComboController struct {
	Ctx     iris.Context
	Service service.ComboService
	Common  service.CommonService
}

func (cc *ComboController) GetList() mvc.Result {
	page := cc.Ctx.URLParamIntDefault("page", 0)
	pageSize := cc.Ctx.URLParamIntDefault("pagesize", 0)
	userId, _ := cc.Ctx.URLParamInt("userId")
	key := cc.Ctx.URLParam("key")
	clusterId, _ := cc.Ctx.URLParamInt("clusterId")
	list, count, err := cc.Service.List(page, pageSize, userId, clusterId, key)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(map[string]interface{}{
		"all":      count,
		"page":     page,
		"pagesize": pageSize,
		"detail":   list,
	})
}

func (cc *ComboController) PostAdd() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	combo := models.Combo{}
	qos := cc.Ctx.PostValue("qos")
	err = json.Unmarshal(utils.Str2bytes(qos), &combo.Qos)
	if err != nil {
		return response.Error(err)
	}
	delete(cc.Ctx.FormValues(), "qos")
	userIdStr := cc.Ctx.PostValue("userIds")
	delete(cc.Ctx.FormValues(), "userIds")
	err = cc.Ctx.ReadForm(&combo)
	if err != nil {
		return response.Error(err)
	}
	err = cc.Service.Add(combo, userIdStr)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (cc *ComboController) PostUser() mvc.Result {
	userIdStr := cc.Ctx.PostValue("userIds")
	comboId, _ := cc.Ctx.PostValueInt("comboId")
	err := cc.Service.User(comboId, userIdStr)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(nil)
}

func (cc *ComboController) PostEdit() mvc.Result {
	var err error
	defer utils.LoggerErrorP(&err)
	combo := models.Combo{}
	qos := cc.Ctx.PostValue("qos")
	if len(qos) != 0 {
		err = json.Unmarshal(utils.Str2bytes(qos), &combo.Qos)
		if err != nil {
			return response.Error(err)
		}
	}
	delete(cc.Ctx.FormValues(), "qos")
	err = cc.Ctx.ReadForm(&combo)
	if err != nil {
		return response.Error(err)
	}
	err = cc.Service.Edit(combo)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (cc *ComboController) PostDelete() mvc.Result {
	comboId, _ := cc.Ctx.PostValueInt("comboId")
	err := cc.Service.Delete(comboId)
	if err != nil {
		utils.LoggerError(err)
		return response.Error(err)
	}
	return response.Success(nil)
}

func (cc *ComboController) GetTagList() mvc.Result {
	list, err := cc.Service.TagList()
	if err != nil {
		return response.Error(err)
	}
	return response.Success(list)
}

func (cc *ComboController) PostTagAdd() mvc.Result {
	name := cc.Ctx.PostValue("name")
	err := cc.Service.TagAdd(name)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}

func (cc *ComboController) PostTagDelete() mvc.Result {
	tagId, _ := cc.Ctx.PostValueInt("id")
	err := cc.Service.TagDelete(tagId)
	if err != nil {
		return response.Error(err)
	}
	return response.Success(nil)
}
