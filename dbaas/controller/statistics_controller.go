package controller

import (
	"DBaas/service"
	"DBaas/x/response"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"sort"
)

type StatisticsController struct {
	Ctx iris.Context
}

func (sc *StatisticsController) GetCluster() mvc.Result {
	page, _ := sc.Ctx.URLParamInt("page")
	pageSize, _ := sc.Ctx.URLParamInt("pagesize")
	search := sc.Ctx.URLParam("search")
	t := sc.Ctx.URLParam("type")
	replicas, _ := sc.Ctx.URLParamInt("replicas")
	sortF := sc.Ctx.URLParam("sortField")
	var dSort string
	switch sortF {
	case "deployStart":
		sortF = "deploy_start"
	case "deployEnd":
		sortF = "deploy_end"
	case "deleteAt":
		sortF = "delete_at"
	default:
		dSort = sortF
		sortF = ""
	}
	sortT := sc.Ctx.URLParam("sortType")
	desc := sortT == "descend"
	list, summary, err := service.Statistics.ClusterList(page, pageSize, search, t, replicas, sortF, desc)
	switch dSort {
	case "useDuration":
		sort.Slice(list, func(i, j int) bool {
			if desc {
				return list[i].UseDuration > list[j].UseDuration
			} else {
				return list[i].UseDuration < list[j].UseDuration
			}
		})
	case "deployDuration":
		sort.Slice(list, func(i, j int) bool {
			if desc {
				return list[i].DeployDuration > list[j].DeployDuration
			} else {
				return list[i].DeployDuration < list[j].DeployDuration
			}
		})
	}
	if err != nil {
		return response.Error(err)
	}
	data := map[string]interface{}{}
	data["all"] = summary["count"]
	delete(summary, "count")
	data["summary"] = summary
	data["detail"] = list
	data["page"] = page
	data["pagesize"] = pageSize
	return response.Success(data)
}
