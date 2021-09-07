package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"sort"
	"strconv"
	"time"
)

type HomeController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//首页功能实体
	CommonService  service.CommonService
	NodeService    service.NodeService
	UserService    service.UserService
	ClusterService service.ClusterService
	Service        service.HomeService
}

//  获取首页基本信息
func (hc *HomeController) GetBaseinfo() mvc.Result {
	utils.LoggerInfo(" 获取首页基本信息 ")
	userTag := hc.Ctx.GetCookie("userTag")
	userName := hc.Ctx.GetCookie("userName")
	rootBaseinfoMap := make(map[string]interface{})
	commonBaseinfoMap := make(map[string]interface{})
	returnDataMap := make(map[string]interface{})
	if userTag == "AAAA" {
		userCount, userCounterr := hc.UserService.GetUserCount("")
		if userCounterr != nil {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": userCounterr,
					"error_msg_zh": userCounterr,
					"data":         returnDataMap,
				},
			}
		}
		clusterCount, clusterCounterr := hc.UserService.GetClusterInstanceCount(0)
		if clusterCounterr != nil {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": clusterCounterr,
					"error_msg_zh": clusterCounterr,
					"data":         returnDataMap,
				},
			}
		}

		var operateList = make([]models.MysqlOperator, 0)
		operatorStatus, _ := hc.Service.GetOperatorStatus()
		k8sStatus := ""
		if operatorStatus == "" {
			k8sStatus = "not_deployed"
		} else if operatorStatus == "true" {
			k8sStatus = "running"
			var errMsg string
			operateList, errMsg, _ = hc.NodeService.OperatorPodList(0, 0, "")
			if errMsg != "" {
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": errMsg,
						"error_msg_zh": errMsg,
						"data":         returnDataMap,
					},
				}
			}
		} else {
			k8sStatus = "error"
		}
		rootBaseinfoMap["status"] = k8sStatus
		rootBaseinfoMap["podList"] = operateList
		rootBaseinfoMap["userNum"] = userCount
		rootBaseinfoMap["clusterNum"] = clusterCount
	} else {
		commonBaseinfoMap, _ = hc.Service.CommonUserBaseInfo(userName)
	}
	returnDataMap["k8s"] = rootBaseinfoMap
	returnDataMap["baseInfo"] = commonBaseinfoMap

	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    returnDataMap,
		},
	}
}

//  获取首页基本信息
func (hc *HomeController) GetChart() mvc.Result {
	utils.LoggerInfo(" 获取首页图表数据 ")
	userName := hc.Ctx.GetCookie("userName")
	chartType := hc.Ctx.URLParam("type")
	model := hc.Ctx.URLParam("model")
	//returnUserDetailData := make(map[string]interface{})
	//returnClusterDetailData := make(map[string]interface{})
	//returnPodDetailData := make(map[string]interface{})
	//returnDetailDataList := make([]map[string]interface{},0)
	var err error
	returnDataList := make([]map[string]interface{}, 0)
	if model == "cluster" {
		if userName == "root" {
			//获取用户信息
			userList := make([]models.User, 0)
			var err error
			userList, err = hc.Service.User3DInfo("")
			if err != nil {
				utils.LoggerError(err)
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_MSG_EN,
						"error_msg_zh": utils.ERROR_MSG_ZH,
					},
				}
			}

			userDetailMap := make(map[string]interface{})
			userSortMap := make(map[int]interface{})
			userDetailList := make([]map[string]interface{}, 0)
			if len(userList) > 0 {
				signaluserch := make(chan map[string]interface{}, len(userList)+1)
				for _, user := range userList {
					go GetSignalUserForHomeChart(user, signaluserch, hc)
				}
				for {
					time.Sleep(10 * time.Millisecond)
					userDetailMap = <-signaluserch
					userSortMap[userDetailMap["id"].(int)] = userDetailMap
					userDetailList = append(userDetailList, userDetailMap)
					if len(userDetailList) == len(userList) {
						break
					}
				}
				close(signaluserch)
				if chartType != "" {
					sort.Slice(userDetailList, func(i, j int) bool {
						return userDetailList[i][chartType].(float64) > userDetailList[j][chartType].(float64)
					})
				}
				for _, signauser := range userDetailList {
					returnDataList = append(returnDataList, (userSortMap[signauser["id"].(int)].(map[string]interface{})["cpuUsage"]).(map[string]interface{}))
					returnDataList = append(returnDataList, (userSortMap[signauser["id"].(int)].(map[string]interface{})["memUsage"]).(map[string]interface{}))
					returnDataList = append(returnDataList, (userSortMap[signauser["id"].(int)].(map[string]interface{})["storUsage"]).(map[string]interface{}))
				}
				if len(returnDataList) > 60{
					returnDataList = returnDataList[:60]
				}
			}
		} else {
			//获取Database信息
			clusterList := make([]models.ClusterInstance, 0)
			clusterList, err = hc.Service.Cluster3DInfo(userName)
			clusterDetailMap := make(map[string]interface{})
			clusterDetailList := make([]map[string]interface{}, 0)
			clusterSortMap := make(map[int]interface{})
			if err != nil {
				utils.LoggerError(err)
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_MSG_EN,
						"error_msg_zh": utils.ERROR_MSG_ZH,
					},
				}
			}
			if len(clusterList) > 0 {
				signalclusterch := make(chan map[string]interface{}, len(clusterList)+1)
				for _, cluster := range clusterList {
					go GetSignalClusterForHomeChart(cluster, signalclusterch, hc)
				}
				for {
					time.Sleep(10 * time.Millisecond)
					clusterDetailMap = <-signalclusterch
					clusterSortMap[clusterDetailMap["id"].(int)] = clusterDetailMap
					clusterDetailList = append(clusterDetailList, clusterDetailMap)
					if len(clusterDetailList) == len(clusterList) {
						break
					}
				}
				close(signalclusterch)
				if chartType != "" {
					sort.Slice(clusterDetailList, func(i, j int) bool {
						return clusterDetailList[i][chartType].(float64) > clusterDetailList[j][chartType].(float64)
					})
				}
				for _, signacluser := range clusterDetailList {
					returnDataList = append(returnDataList, (clusterSortMap[signacluser["id"].(int)].(map[string]interface{})["cpuUsage"]).(map[string]interface{}))
					returnDataList = append(returnDataList, (clusterSortMap[signacluser["id"].(int)].(map[string]interface{})["memUsage"]).(map[string]interface{}))
					returnDataList = append(returnDataList, (clusterSortMap[signacluser["id"].(int)].(map[string]interface{})["storUsage"]).(map[string]interface{}))
				}
				if len(returnDataList) > 60{
					returnDataList = returnDataList[:60]
				}
			}
		}

	} else if model == "pod" {
		//获取pod信息
		podList := make([]models.Instance, 0)
		if userName == "root" {
			podList, err = hc.Service.Pod3DInfo("")
		} else {
			podList, err = hc.Service.Pod3DInfo(userName)
		}

		podDetailMap := make(map[string]interface{})
		podSortMap := make(map[int]interface{})
		podDetailList := make([]map[string]interface{}, 0)
		if err != nil {
			utils.LoggerError(err)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_MSG_EN,
					"error_msg_zh": utils.ERROR_MSG_ZH,
				},
			}
		}
		if len(podList) > 0 {
			signalpodch := make(chan map[string]interface{}, len(podList)+1)
			for _, pod := range podList {
				go GetSignalPodForHomeChart(pod, signalpodch, hc)
			}
			for {
				time.Sleep(10 * time.Millisecond)
				podDetailMap = <-signalpodch
				podSortMap[podDetailMap["id"].(int)] = podDetailMap
				podDetailList = append(podDetailList, podDetailMap)
				if len(podDetailList) == len(podList) {
					break
				}
			}
			close(signalpodch)
			if chartType != "" {
				sort.Slice(podDetailList, func(i, j int) bool {
					return podDetailList[i][chartType].(float64) > podDetailList[j][chartType].(float64)
				})
			}
			for _, signapod := range podDetailList {
				if signapod["cpu"] != float64(0) || signapod["mem"] != float64(0) {
					if len((podSortMap[signapod["id"].(int)].(map[string]interface{})["cpuUsage"]).(map[string]interface{})) > 0 {
						returnDataList = append(returnDataList, (podSortMap[signapod["id"].(int)].(map[string]interface{})["cpuUsage"]).(map[string]interface{}))

					}
					if len((podSortMap[signapod["id"].(int)].(map[string]interface{})["memUsage"]).(map[string]interface{})) > 0 {
						returnDataList = append(returnDataList, (podSortMap[signapod["id"].(int)].(map[string]interface{})["memUsage"]).(map[string]interface{}))
					}
				}
			}
			if len(returnDataList) > 40{
				returnDataList = returnDataList[:40]
			}
		}
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    returnDataList,
		},
	}
}

func GetSignalUserForHomeChart(user models.User, signaluserch chan map[string]interface{}, hc *HomeController) {
	signalUserChartMap := make(map[string]interface{})
	//returnClusterChartMap:= make(map[int]interface{})

	signalCpuChart := make(map[string]interface{})
	signalMemChart := make(map[string]interface{})
	signalStorageChart := make(map[string]interface{})
	xValue := user.UserName
	cpuUsed, memUsed, storageUsed, err := hc.UserService.GetClusterCpuMemStorage(user.Id)
	utils.LoggerError(err)
	signalUserChartMap["id"] = user.Id
	signalCpuChart["label"] = "CPU"
	signalCpuChart["xValue"] = xValue
	signalCpuChart["used"] = cpuUsed
	signalCpuChart["total"] = user.CpuAll
	cpuUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(cpuUsed)/float64(user.CpuAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signaluserch <- signalUserChartMap
	}
	signalCpuChart["yValue"] = cpuUsage
	signalUserChartMap["cpuUsage"] = signalCpuChart

	signalMemChart["label"] = "MEM"
	signalMemChart["xValue"] = xValue
	signalMemChart["used"] = memUsed
	signalMemChart["total"] = user.MemAll
	memUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(memUsed)/float64(user.MemAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signaluserch <- signalUserChartMap
	}
	signalMemChart["yValue"] = memUsage
	signalUserChartMap["memUsage"] = signalMemChart

	signalStorageChart["label"] = "STOR"
	signalStorageChart["xValue"] = xValue
	signalStorageChart["used"] = storageUsed
	signalStorageChart["total"] = user.StorageAll
	storUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(storageUsed)/float64(user.StorageAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signaluserch <- signalUserChartMap
	}
	signalStorageChart["yValue"] = storUsage
	signalUserChartMap["storUsage"] = signalStorageChart
	signalUserChartMap["cpu"] = cpuUsage
	signalUserChartMap["mem"] = memUsage
	signalUserChartMap["stor"] = storUsage
	signaluserch <- signalUserChartMap
}

func GetSignalClusterForHomeChart(cluster models.ClusterInstance, signalclusterch chan map[string]interface{}, hc *HomeController) {
	signalClusterChartMap := make(map[string]interface{})
	//returnClusterChartMap:= make(map[int]interface{})
	user, errM := hc.UserService.SelectOne(cluster.UserId)
	if errM != "" {
		signalclusterch <- signalClusterChartMap
	}
	signalCpuChart := make(map[string]interface{})
	signalMemChart := make(map[string]interface{})
	signalStorageChart := make(map[string]interface{})
	xValue := cluster.Name
	signalClusterChartMap["id"] = cluster.Id
	signalCpuChart["label"] = "CPU"
	signalCpuChart["xValue"] = xValue
	signalCpuChart["used"] = cluster.LimitCpu
	signalCpuChart["total"] = user.CpuAll
	cpuUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(cluster.LimitCpu)/float64(user.CpuAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signalclusterch <- signalClusterChartMap
	}
	signalCpuChart["yValue"] = cpuUsage
	signalClusterChartMap["cpuUsage"] = signalCpuChart

	signalMemChart["label"] = "MEM"
	signalMemChart["xValue"] = xValue
	signalMemChart["used"] = cluster.LimitMem
	signalMemChart["total"] = user.MemAll
	memUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(cluster.LimitMem)/float64(user.MemAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signalclusterch <- signalClusterChartMap
	}
	signalMemChart["yValue"] = memUsage
	signalClusterChartMap["memUsage"] = signalMemChart

	signalStorageChart["label"] = "STOR"
	signalStorageChart["xValue"] = xValue
	signalStorageChart["used"] = cluster.Storage
	signalStorageChart["total"] = user.StorageAll
	storUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(cluster.Storage)/float64(user.StorageAll))*100), 64)
	if err != nil {
		utils.LoggerError(err)
		signalclusterch <- signalClusterChartMap
	}
	signalStorageChart["yValue"] = storUsage
	signalClusterChartMap["storUsage"] = signalStorageChart
	signalClusterChartMap["cpu"] = cpuUsage
	signalClusterChartMap["mem"] = memUsage
	signalClusterChartMap["stor"] = storUsage
	signalclusterch <- signalClusterChartMap
}

func GetSignalPodForHomeChart(pod models.Instance, signalpodch chan map[string]interface{}, hc *HomeController) {
	collectService, conn := service.NewCollectService()
	defer service.CloseGrpc(conn)
	podBaseInformation := make(map[string]interface{})
	podBaseSortCpuMap := make(map[string]interface{})
	podBaseSortMemMap := make(map[string]interface{})
	podBaseInformation["id"] = pod.Id
	podBaseInformation["name"] = pod.Name
	podBaseInformation["cpuUsage"] = podBaseSortCpuMap
	podBaseInformation["memUsage"] = podBaseSortMemMap
	podBaseInformation["cpu"] = float64(0.00)
	podBaseInformation["mem"] = float64(0.00)
	cluster, err := hc.Service.SelectOneCluster(pod.ClusterId)
	if err != nil {
		signalpodch <- podBaseInformation
	}
	//获取pod性能数据
	//cpuUsage
	cpuSql := fmt.Sprintf(`SELECT  * FROM "metrics_pod_cpu" WHERE "pod"='%v' AND time>=now()-300s order by time desc limit 1;`, pod.Name)
	cpuResult := collectService.GetInfluxDbData(cpuSql, "1")
	if len(cpuResult) > 0 {
		cpuResult := cpuResult[len(cpuResult)-1 : len(cpuResult)]
		if valueFloat, ok := cpuResult[0]["cpu"].(float64); ok {
			cpuUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(valueFloat)/float64(cluster.LimitCpu))*100), 64)
			utils.LoggerError(err)
			podBaseInformation["cpu"] = cpuUsage
			podBaseSortCpuMap["label"] = "CPU"
			podBaseSortCpuMap["xValue"] = cluster.Name + "/" + pod.Name
			podBaseSortCpuMap["yValue"] = cpuUsage
			podBaseInformation["cpuUsage"] = podBaseSortCpuMap
		}
	}

	//memUsage
	memSql := fmt.Sprintf(`SELECT  * FROM "metrics_pod_mem" WHERE "pod"='%v' AND time>=now()-300s order by time desc limit 1;`, pod.Name)
	memResult := collectService.GetInfluxDbData(memSql, "1")
	if len(memResult) > 0 {
		memResult := memResult[len(memResult)-1 : len(memResult)]
		if valueFloat, ok := memResult[0]["mem"].(float64); ok {
			memUsage, err := strconv.ParseFloat(fmt.Sprintf("%.2f", (float64(valueFloat)/float64(cluster.LimitMem))*100), 64)
			utils.LoggerError(err)
			podBaseInformation["mem"] = memUsage
			podBaseSortMemMap["label"] = "MEM"
			podBaseSortMemMap["xValue"] = cluster.Name + "/" + pod.Name
			podBaseSortMemMap["yValue"] = memUsage
			podBaseInformation["memUsage"] = podBaseSortMemMap
		}
	}
	signalpodch <- podBaseInformation
}

//  获取首页基本信息
func (hc *HomeController) GetThreed() mvc.Result {
	utils.LoggerInfo(" 获取首页3d数据 ")
	//userTag := hc.Ctx.GetCookie("userTag")
	userName := hc.Ctx.GetCookie("userName")
	//rootBaseinfoMap := make(map[string]interface{})
	//commonBaseinfoMap := make(map[string]interface{})
	//returnDataMap := make(map[string]interface{})
	returnUserDetailData := make(map[string]interface{})
	returnClusterDetailData := make(map[string]interface{})
	returnPodDetailData := make(map[string]interface{})
	returnDetailDataList := make([]map[string]interface{}, 0)
	//获取用户信息
	userList := make([]models.User, 0)
	var err error
	if userName == "root" {
		userList, err = hc.Service.User3DInfo("")
	} else {
		userList, err = hc.Service.User3DInfo(userName)
	}
	if err != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_MSG_EN,
				"error_msg_zh": utils.ERROR_MSG_ZH,
			},
		}
	}

	userDetailMap := make(map[string]interface{})
	userMap := make(map[int]interface{})
	returnDataList := make([]interface{}, 0)

	var userDetailList []interface{}
	if len(userList) > 0 {
		signaluserch := make(chan map[string]interface{}, len(userList)+1)
		for _, user := range userList {
			go GetSignalUserForHome(user, signaluserch, hc)
		}
		for {
			//time.Sleep(10 * time.Millisecond)
			userDetailMap = <-signaluserch
			signalUserForHome := make(map[string]interface{})
			signalUserForHome["type"] = "User"
			signalUserForHome["belong"] = ""
			signalUserForHome["data"] = userDetailMap
			userMap[userDetailMap["id"].(int)] = signalUserForHome
			userDetailList = append(userDetailList, userDetailMap)
			if len(userDetailList) == len(userList) {
				break
			}
		}
		close(signaluserch)
		for _, signaluser := range userList {
			returnDataList = append(returnDataList, userMap[signaluser.Id])
		}
	}
	returnUserDetailData["children"] = returnDataList
	returnUserDetailData["type"] = "User"
	returnDetailDataList = append(returnDetailDataList, returnUserDetailData)
	//获取Database信息
	clusterList := make([]models.ClusterInstance, 0)
	if userName == "root" {
		clusterList, err = hc.Service.Cluster3DInfo("")
	} else {
		clusterList, err = hc.Service.Cluster3DInfo(userName)
	}
	clusterDetailMap := make(map[string]interface{})
	clusterDetailList := make([]map[string]interface{}, 0)
	if err != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_MSG_EN,
				"error_msg_zh": utils.ERROR_MSG_ZH,
			},
		}
	}
	if len(clusterList) > 0 {
		signalclusterch := make(chan map[string]interface{}, len(clusterList)+1)
		for _, cluster := range clusterList {
			go GetSignalClusterForHome(cluster, signalclusterch, hc)
		}
		for {
			time.Sleep(10 * time.Millisecond)
			clusterDetailMap = <-signalclusterch
			clusterDetailList = append(clusterDetailList, clusterDetailMap)
			if len(clusterDetailList) == len(clusterList) {
				break
			}
		}
		close(signalclusterch)
	}
	returnClusterDetailData["children"] = clusterDetailList
	returnClusterDetailData["type"] = "Database"
	returnDetailDataList = append(returnDetailDataList, returnClusterDetailData)

	//获取pod信息
	podList := make([]models.Instance, 0)
	if userName == "root" {
		podList, err = hc.Service.Pod3DInfo("")
	} else {
		podList, err = hc.Service.Pod3DInfo(userName)
	}

	podDetailMap := make(map[string]interface{})
	podDetailList := make([]map[string]interface{}, 0)
	if err != nil {
		utils.LoggerError(err)
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_MSG_EN,
				"error_msg_zh": utils.ERROR_MSG_ZH,
			},
		}
	}
	if len(podList) > 0 {
		signalpodch := make(chan map[string]interface{}, len(podList)+1)
		for _, pod := range podList {
			go GetSignalPodForHome(pod, signalpodch, hc)
		}
		for {
			time.Sleep(10 * time.Millisecond)
			podDetailMap = <-signalpodch
			podDetailList = append(podDetailList, podDetailMap)
			if len(podDetailList) == len(podList) {
				break
			}
		}
		close(signalpodch)
	}
	returnPodDetailData["children"] = podDetailList
	returnPodDetailData["type"] = "Pod"
	returnDetailDataList = append(returnDetailDataList, returnPodDetailData)
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    returnDetailDataList,
		},
	}
}

func GetSignalUserForHome(user models.User, signaluserch chan map[string]interface{}, hc *HomeController) {
	capricornService, conn := service.NewCapricornService()
	defer service.CloseGrpc(conn)
	cpuUsed, memUsed, storageUsed, err := hc.UserService.GetClusterCpuMemStorage(user.Id)
	utils.LoggerError(err)
	userBaseInformation := make(map[string]interface{})
	userBaseInformation["id"] = user.Id
	userBaseInformation["cpuTotal"] = user.CpuAll
	userBaseInformation["cpuUsed"] = cpuUsed
	userBaseInformation["memTotal"] = user.MemAll
	userBaseInformation["memUsed"] = memUsed
	userBaseInformation["storTotal"] = user.StorageAll
	userBaseInformation["storUsed"] = storageUsed
	userIdString := strconv.Itoa(user.ZdcpId)
	hostInfo := make(map[string]interface{}, 0)
	userList, ErrorMsgEn, ErrorMsgZh := capricornService.GetUserResources(userIdString, "", "")
	if ErrorMsgEn != "" && ErrorMsgZh != "" {
		iris.New().Logger().Error(ErrorMsgEn)
	}
	if len(userList) > 0 {
		hostInfo = userList[0]
		if _, ok := hostInfo["username"]; ok {
			userBaseInformation["username"] = hostInfo["username"]
		}
		if _, ok := hostInfo["status"]; ok {
			if hostInfo["status"] == "inactive" {
				userBaseInformation["status"] = "disable"
			} else {
				userBaseInformation["status"] = "enable"
			}
		}
	}
	signaluserch <- userBaseInformation
}

func GetSignalClusterForHome(cluster models.ClusterInstance, signalclusterch chan map[string]interface{}, hc *HomeController) {
	collectService, conn := service.NewCollectService()
	defer service.CloseGrpc(conn)
	clusterBaseInformation := make(map[string]interface{})
	signalclusterForHome := make(map[string]interface{})
	clusterBaseInformation["id"] = cluster.Id
	clusterBaseInformation["name"] = cluster.Name
	clusterBaseInformation["status"] = cluster.Status
	clusterBaseInformation["cpu"] = cluster.LimitCpu
	clusterBaseInformation["mem"] = cluster.LimitMem
	clusterBaseInformation["stor"] = cluster.Storage
	clusterBaseInformation["qps"] = ""
	clusterBaseInformation["session"] = ""
	clusterBaseInformation["runTime"] = ""
	//获取pod性能数据
	MasterPod, err := hc.Service.GetMasterPodByCluster(cluster.Id)
	if err != nil {
		signalclusterch <- signalclusterForHome
	}
	//session
	if MasterPod.Name != "" {
		sessionSql := fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs order by time desc limit 1;`, "value", "mysql_global_status_threads_connected", MasterPod.Name, 300)
		sessionResult := collectService.GetInfluxDbData(sessionSql, "1")
		if len(sessionResult) > 0 {
			sessionResult := sessionResult[len(sessionResult)-1 : len(sessionResult)]
			if valueFloat, ok := sessionResult[0]["value"].(float64); ok {
				clusterBaseInformation["session"] = valueFloat
			}
		}
		//qps
		qpsSql := fmt.Sprintf(`SELECT  DIFFERENCE( %v ) as %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-60s order by time desc`, "value", "value", "mysql_global_status_questions", MasterPod.Name)
		qpsResult := collectService.GetInfluxDbData(qpsSql, "1")
		if len(qpsResult) > 0 {
			qpsResult := qpsResult[len(qpsResult)-1 : len(qpsResult)]
			if valueFloat, ok := qpsResult[0]["value"].(float64); ok {
				clusterBaseInformation["qps"] = valueFloat
			}
		}

		//uptime
		uptimeSql := fmt.Sprintf(`SELECT  %v FROM zdcp.autogen.%v WHERE "kubernetes_pod_name"='%v' AND time>=now()-%vs order by time desc limit 1;`, "value", "mysql_global_status_uptime", MasterPod.Name, 300)
		uptimeResult := collectService.GetInfluxDbData(uptimeSql, "1")
		if len(uptimeResult) > 0 {
			uptimeResult := uptimeResult[len(uptimeResult)-1 : len(uptimeResult)]
			if valueFloat, ok := uptimeResult[0]["value"].(float64); ok {
				//valueInt, _ := strconv.ParseInt(valueString, 10, 32)
				clusterBaseInformation["runTime"] = utils.ResolveTime(int(valueFloat))
			}
		}
	}

	signalclusterForHome["type"] = "Database"
	user, errM := hc.UserService.SelectOne(cluster.UserId)
	if errM != "" {
		signalclusterch <- signalclusterForHome
	}
	signalclusterForHome["belong"] = user.UserName
	signalclusterForHome["data"] = clusterBaseInformation
	signalclusterch <- signalclusterForHome
}

func GetSignalPodForHome(pod models.Instance, signalpodch chan map[string]interface{}, hc *HomeController) {
	collectService, conn := service.NewCollectService()
	defer service.CloseGrpc(conn)
	podBaseInformation := make(map[string]interface{})
	signalPodForHome := make(map[string]interface{})
	podBaseInformation["id"] = pod.Id
	podBaseInformation["name"] = pod.Name
	podBaseInformation["status"] = pod.Status
	podBaseInformation["cpuUsage"] = ""
	podBaseInformation["memUsage"] = ""
	cluster, err := hc.Service.SelectOneCluster(pod.ClusterId)
	if err != nil {
		signalpodch <- podBaseInformation
	}
	//获取pod性能数据
	//cpuUsage
	cpuSql := fmt.Sprintf(`SELECT  * FROM "metrics_pod_cpu" WHERE "pod"='%v' AND time>=now()-60s order by time desc limit 1;`, pod.Name)
	cpuResult := collectService.GetInfluxDbData(cpuSql, "1")
	if len(cpuResult) > 0 {
		cpuResult := cpuResult[len(cpuResult)-1 : len(cpuResult)]
		if valueFloat, ok := cpuResult[0]["cpu"].(float64); ok {
			podBaseInformation["cpuUsed"] = float64(valueFloat)
			podBaseInformation["cpuTotal"] = float64(cluster.LimitCpu)
			cpuUsage := fmt.Sprintf("%.2f", (float64(valueFloat)/float64(cluster.LimitCpu))*100)
			podBaseInformation["cpuUsage"] = cpuUsage + "%"
		}
	}

	//memUsage
	memSql := fmt.Sprintf(`SELECT  * FROM "metrics_pod_mem" WHERE "pod"='%v' AND time>=now()-60s order by time desc limit 1;`, pod.Name)
	memResult := collectService.GetInfluxDbData(memSql, "1")
	if len(memResult) > 0 {
		memResult := memResult[len(memResult)-1 : len(memResult)]
		if valueFloat, ok := memResult[0]["mem"].(float64); ok {
			podBaseInformation["memUsed"] = float64(valueFloat)
			podBaseInformation["memTotal"] = float64(cluster.LimitMem)
			memUsage := fmt.Sprintf("%.2f", (float64(valueFloat)/float64(cluster.LimitMem))*100)
			podBaseInformation["memUsage"] = memUsage + "%"
		}
	}

	signalPodForHome["type"] = "Pod"
	signalPodForHome["belong"] = cluster.Name
	signalPodForHome["data"] = podBaseInformation
	signalpodch <- signalPodForHome
}
