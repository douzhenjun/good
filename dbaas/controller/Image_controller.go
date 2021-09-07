package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"DBaas/x/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type ImageController struct {
	Ctx           iris.Context
	Service       service.ImageService
	CommonService service.CommonService
}

func (ic *ImageController) GetTypes() mvc.Response {
	imageType := ic.Ctx.URLParam("type")
	ret := ic.Service.GetImageType(imageType)
	return mvc.Response{Object: utils.ResponseOk(ret)}
}

func (ic *ImageController) GetList() mvc.Result {
	page := ic.Ctx.URLParamIntDefault("page", 0)
	pageSize := ic.Ctx.URLParamIntDefault("pagesize", 0)
	key := ic.Ctx.URLParam("key")
	imagesList, count, err := ic.Service.List(page, pageSize, key)
	utils.LoggerError(err)
	result := make(map[string]interface{})
	result["all"] = count
	result["detail"] = imagesList
	result["page"] = page
	result["pagesize"] = pageSize
	return response.Success(result)
}

func (ic *ImageController) PostAdd() mvc.Result {
	var err error; defer utils.LoggerErrorP(&err)
	imageName := ic.Ctx.PostValue("imageName")
	version := ic.Ctx.PostValue("version")
	description := ic.Ctx.PostValue("description")
	imageType := ic.Ctx.PostValue("type")
	category := ic.Ctx.PostValue("category")
	image := models.Images{
		ImageName:   imageName,
		Version:     version,
		Description: description,
		Type:        imageType,
		Category:    category,
	}
	userName := ic.Ctx.GetCookie("userName")
	err = ic.Service.Add(image)
	if err == nil {
		ic.CommonService.AddLog("info", "system-image", userName, fmt.Sprintf("add image %v successful", imageName+":"+version))
		return response.Success(nil)
	}
	if utils.ErrorContains(err, "unique_name_version") {
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("init add image error %v", response.ErrorImageExist.En))
		return response.Fail(response.ErrorImageExist)
	} else {
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("init add image error %v", err.Error()))
		return response.Error(err)
	}
}

func (ic *ImageController) PostInitAdd() mvc.Result {
	var err error; defer utils.LoggerErrorP(&err)
	imageList := make([]models.Images, 0)
	param := ic.Ctx.PostValue("imageList")
	err = json.Unmarshal([]byte(param), &imageList)
	if err != nil {
		return response.Fail(response.ErrorParameter)
	}

	userName := ic.Ctx.GetCookie("userName")
	err = ic.Service.InitAdd(imageList)
	if err == nil {
		ic.CommonService.AddLog("info", "system-image", userName, fmt.Sprintf("init add images successful"))
		return response.Success(nil)
	}
	if utils.ErrorContains(err, "unique_name_version") {
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("init add image error %v", response.ErrorImageExist.En))
		return response.Fail(response.ErrorImageExist)
	} else {
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("init add image error %v", err.Error()))
		return response.Error(err)
	}
}

func (ic *ImageController) PostUpdate() mvc.Result {
	var err error; defer utils.LoggerErrorP(&err)
	imageId := ic.Ctx.PostValueIntDefault("id", -1)
	userName := ic.Ctx.GetCookie("userName")
	if imageId <= -1 {
		err = errors.New("invalid id")
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("update image %v error %s", imageId, err))
		return response.Error(err)
	}
	version := ic.Ctx.PostValue("version")
	description := ic.Ctx.PostValue("description")
	imageType := ic.Ctx.PostValue("type")
	category := ic.Ctx.PostValue("category")
	image := models.Images{
		Id:          imageId,
		Version:     version,
		Description: description,
		Type:        imageType,
		Category:    category,
	}
	err = ic.Service.Update(image)
	if err == nil {
		ic.CommonService.AddLog("info", "system-image", userName, fmt.Sprintf("update image %v successful", imageId))
		return response.Success(nil)
	}
	ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("update image %v error %s", imageId, err))
	return response.Error(err)
}

func (ic *ImageController) PostDelete() mvc.Result {
	var err error; defer utils.LoggerErrorP(&err)
	imageId := ic.Ctx.PostValueIntDefault("id", -1)
	userName := ic.Ctx.GetCookie("userName")
	if imageId <= -1 {
		err = errors.New("invalid id")
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("delete image %v error %s", imageId, err))
		return response.Error(err)
	}
	err = ic.Service.Delete(imageId)
	if err == nil {
		ic.CommonService.AddLog("info", "system-image", userName, fmt.Sprintf("update image %v successful", imageId))
		return response.Success(nil)
	}
	ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("update image %v error %s", imageId, err))
	if err.Error() == response.ErrorImageOccupied.En {
		return response.Fail(response.ErrorImageOccupied)
	}
	return response.Error(err)
}

func (ic *ImageController) PostParam() mvc.Result {
	imageId, _ := ic.Ctx.PostValueInt("id")
	paramList, err := ic.Service.Param(imageId)
	if err == nil {
		return response.Success(map[string]interface{}{"detail": paramList})
	}
	utils.LoggerError(err)
	return response.Error(err)
}
