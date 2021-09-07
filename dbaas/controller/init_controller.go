/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date:  2021-02-22 09:44:07
 * @LastEditors: ddh
 * @LastEditTime: 2021-02-22 09:44:07
 */

package controller

import (
	"DBaas/models"
	"DBaas/service"
	"DBaas/utils"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strings"
	"time"
)

type InitController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context

	Service service.InitService

	CommonService service.CommonService

	ParamService service.ParameterService

	NodeService service.NodeService

	StorageService service.StorageService

	ImageService service.ImageService
}

func (ic *InitController) PostSave() mvc.Result {
	utils.LoggerInfo("初始化-保存")
	step := ic.Ctx.PostValue("step")
	storageStr := ic.Ctx.PostValue("storage")
	hostListStr := ic.Ctx.PostValue("hostList")
	imageListStr := ic.Ctx.PostValue("imageList")
	operatorStr := ic.Ctx.PostValue("operator")
	userName := ic.Ctx.GetCookie("userName")
	if step == "storage" {
		if storageStr == "" {
			ic.CommonService.AddLog("error", "system-init-checkout-storage", userName, utils.ERROR_PARAMETER_EN)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_PARAMETER_EN,
					"error_msg_zh": utils.ERROR_PARAMETER_ZH,
				},
			}
		} else {
			storageInitInfo := models.Initinfo{Name: "storage", Message: storageStr, Isaccess: "False"}
			ic.Service.AddOrModifyInitinfo("storage", storageInitInfo)
			storage := make(map[string]interface{})
			jsonerr := json.Unmarshal([]byte(storageStr), &storage)
			if jsonerr != nil {
				utils.LoggerError(jsonerr)
				ic.CommonService.AddLog("error", "system-init-checkout-storage", userName, jsonerr.Error())
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_DATA_TRANSFER_EN,
						"error_msg_zh": utils.ERROR_DATA_TRANSFER_ZH,
					},
				}
			}
			if _, ok := storage["scName"]; ok {
				scName := storage["scName"].(string)
				if scName == "" {
					//ic.CommonService.AddLog("error", "system-init-checkout-storage", userName, "Missing SC name")
					//return mvc.Response{
					//	Object: map[string]interface{}{
					//		"errorno":      utils.RECODE_FAIL,
					//		"error_msg_en": "Missing SC name",
					//		"error_msg_zh": "缺少SC名称",
					//	},
					//}
				} else {
					err, _ := ic.CommonService.GetResources("sc", scName, ic.CommonService.GetNameSpace(), meta1.GetOptions{})
					if err != nil {
						utils.LoggerError(err)
						ic.CommonService.AddLog("error", "system-init-checkout-storage", userName, "Storage parameter error"+err.Error())
						return mvc.Response{
							Object: map[string]interface{}{
								"errorno":      utils.RECODE_FAIL,
								"error_msg_en": "Storage parameter error" + err.Error(),
								"error_msg_zh": "存储参数错误：" + err.Error(),
							},
						}
					}
				}
			} else {
				//ic.CommonService.AddLog("error", "system-init-checkout-storage", userName, "Missing SC name")
				//return mvc.Response{
				//	Object: map[string]interface{}{
				//		"errorno":      utils.RECODE_FAIL,
				//		"error_msg_en": "Missing SC name",
				//		"error_msg_zh": "缺少SC名称",
				//	},
				//}
			}
			storageInitInfo.Isaccess = "True"
			ic.Service.AddOrModifyInitinfo("storage", storageInitInfo)
			ic.CommonService.AddLog("info", "system-init-checkout-storage", userName, "checkout storage successful")
		}
	}
	if step == "host" {
		if hostListStr == "" {
			ic.CommonService.AddLog("error", "system-init-checkout-host", userName, utils.ERROR_PARAMETER_EN)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_PARAMETER_EN,
					"error_msg_zh": utils.ERROR_PARAMETER_ZH,
				},
			}
		} else {
			hostInitInfo := models.Initinfo{Name: "host", Message: hostListStr, Isaccess: "False"}
			ic.Service.AddOrModifyInitinfo("host", hostInitInfo)
			hostList := make([]map[string]interface{}, 0)
			jsonerr := json.Unmarshal([]byte(hostListStr), &hostList)
			if jsonerr != nil {
				utils.LoggerError(jsonerr)
				ic.CommonService.AddLog("error", "system-init-checkout-host", userName, jsonerr.Error())
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_DATA_TRANSFER_EN,
						"error_msg_zh": utils.ERROR_DATA_TRANSFER_ZH,
					},
				}
			}

			for _, host := range hostList {
				node, err := ic.Service.GetNodeById(int(host["id"].(float64)))
				if err != "" {
					ic.CommonService.AddLog("error", "system-init-checkout-host", userName, "Parameter error, failed to get node information:"+err)
					return mvc.Response{
						Object: map[string]interface{}{
							"errorno":      utils.RECODE_FAIL,
							"error_msg_en": "Parameter error, failed to get node information:" + err,
							"error_msg_zh": "参数错误，获取节点信息失败：" + err,
						},
					}
				}
				if node.NodeName == "" {
					ic.CommonService.AddLog("error", "system-init-checkout-host", userName, "Parameter error, failed to get node information")
					return mvc.Response{
						Object: map[string]interface{}{
							"errorno":      utils.RECODE_FAIL,
							"error_msg_en": "Parameter error, failed to get node information",
							"error_msg_zh": "参数错误，获取节点信息失败",
						},
					}
				}
			}
			hostInitInfo.Isaccess = "True"
			ic.Service.AddOrModifyInitinfo("host", hostInitInfo)
			ic.CommonService.AddLog("info", "system-init-checkout-host", userName, "checkout host successful")

		}
	}
	if step == "image" {
		if imageListStr == "" {
			ic.CommonService.AddLog("error", "system-init-checkout-image", userName, utils.ERROR_PARAMETER_EN)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_PARAMETER_EN,
					"error_msg_zh": utils.ERROR_PARAMETER_ZH,
				},
			}
		} else {
			imageInitInfo := models.Initinfo{Name: "image", Message: imageListStr, Isaccess: "False"}
			ic.Service.AddOrModifyInitinfo("image", imageInitInfo)
			imageList := make([]map[string]interface{}, 0)
			jsonerr := json.Unmarshal([]byte(imageListStr), &imageList)
			if jsonerr != nil {
				utils.LoggerError(jsonerr)
				ic.CommonService.AddLog("error", "system-init-checkout-image", userName, jsonerr.Error())
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_DATA_TRANSFER_EN,
						"error_msg_zh": utils.ERROR_DATA_TRANSFER_ZH,
					},
				}
			}
			sameImageName := ""
			errorImageName := ""
			ImageNameMap := make(map[string]interface{})
			for _, image := range imageList {
				if _, ok := ImageNameMap[image["imageName"].(string)]; ok {
					sameImageName = image["imageName"].(string) + "," + ""
				}
				ImageNameMap[image["imageName"].(string)] = true
				harborAddressParam := ic.ParamService.SelectOneByKey("harbor_address")
				checkUrl := fmt.Sprintf("http://%v/api/repositories/%v/tags/%v", harborAddressParam.ParamValue, image["imageName"], image["version"])
				client := &http.Client{Timeout: 10 * time.Second}
				resp, err := client.Get(checkUrl)
				if err == nil {
					_ = resp.Body.Close()
					if resp.StatusCode == http.StatusNotFound {
						errorImageName = image["imageName"].(string) + "," + ""
					}
				} else {
					errorImageName = image["imageName"].(string) + "," + ""
				}
			}
			if sameImageName != "" {
				sameImageName := strings.TrimRight(sameImageName, ",")
				ic.CommonService.AddLog("error", "system-init-checkout-image", userName, fmt.Sprintf("Duplicate image name: %s", sameImageName))
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": fmt.Sprintf("Duplicate image name: %s", sameImageName),
						"error_msg_zh": fmt.Sprintf("镜像名称: %s 重复", sameImageName),
					},
				}
			}
			if errorImageName != "" {
				errorImageName := strings.TrimRight(errorImageName, ",")
				ic.CommonService.AddLog("error", "system-init-checkout-image", userName, fmt.Sprintf("No image was found according to the filled image parameters: %s", errorImageName))
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": fmt.Sprintf("No image was found according to the filled image parameters: %s", errorImageName),
						"error_msg_zh": fmt.Sprintf("根据填写的镜像参数未找到镜像：%s", errorImageName),
					},
				}
			}
			imageInitInfo.Isaccess = "True"
			ic.Service.AddOrModifyInitinfo("image", imageInitInfo)
			ic.CommonService.AddLog("info", "system-init-checkout-image", userName, "checkout image successful")

		}
	}
	if step == "operator" {
		if operatorStr == "" {
			ic.CommonService.AddLog("error", "system-init-checkout-operator", userName, utils.ERROR_PARAMETER_EN)
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": utils.ERROR_PARAMETER_EN,
					"error_msg_zh": utils.ERROR_PARAMETER_ZH,
				},
			}
		} else {
			operatorInitInfo := models.Initinfo{Name: "operator", Message: operatorStr, Isaccess: "False"}
			ic.Service.AddOrModifyInitinfo("operator", operatorInitInfo)
			operatorMap := make(map[string]interface{})
			jsonerr := json.Unmarshal([]byte(operatorStr), &operatorMap)
			if jsonerr != nil {
				utils.LoggerError(jsonerr)
				ic.CommonService.AddLog("error", "system-init-checkout-operator", userName, jsonerr.Error())
				return mvc.Response{
					Object: map[string]interface{}{
						"errorno":      utils.RECODE_FAIL,
						"error_msg_en": utils.ERROR_DATA_TRANSFER_EN,
						"error_msg_zh": utils.ERROR_DATA_TRANSFER_ZH,
					},
				}
			}

			if _, ok := operatorMap["scName"]; ok {
				scName := operatorMap["scName"].(string)
				if scName == "" {
					//ic.CommonService.AddLog("error", "system-init-checkout-operator", userName,"Missing SC name")
					//return mvc.Response{
					//	Object: map[string]interface{}{
					//		"errorno":      utils.RECODE_FAIL,
					//		"error_msg_en": "Missing SC name",
					//		"error_msg_zh": "缺少SC名称",
					//	},
					//}
				} else {
					err, _ := ic.CommonService.GetResources("sc", scName, ic.CommonService.GetNameSpace(), meta1.GetOptions{})
					if err != nil {
						utils.LoggerError(err)
						ic.CommonService.AddLog("error", "system-init-checkout-operator", userName, err.Error())
						return mvc.Response{
							Object: map[string]interface{}{
								"errorno":      utils.RECODE_FAIL,
								"error_msg_en": err.Error(),
								"error_msg_zh": err.Error(),
							},
						}
					}
				}

			} else {
				//ic.CommonService.AddLog("error", "system-init-checkout-operator", userName,"Missing SC name")
				//return mvc.Response{
				//	Object: map[string]interface{}{
				//		"errorno":      utils.RECODE_FAIL,
				//		"error_msg_en": "Missing SC name",
				//		"error_msg_zh": "缺少SC名称",
				//	},
				//}
			}
			operatorInitInfo.Isaccess = "True"
			ic.Service.AddOrModifyInitinfo("operator", operatorInitInfo)
			ic.CommonService.AddLog("info", "system-init-checkout-operator", userName, "checkout operator successful")

		}
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

func (ic *InitController) GetList() mvc.Result {
	//userTag := ic.Ctx.GetCookie("userTag")
	utils.LoggerInfo("初始化-查询步骤信息")
	step := ic.Ctx.URLParam("step")
	data := make(map[string]interface{})
	returnStep := ""
	err := ""
	returnStep, err = ic.Service.GetLastDeployStep()
	if err != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": err,
				"error_msg_zh": err,
			},
		}
	}
	data["step"] = returnStep
	if step == "storage" {
		storage, errMsg := GetStorageInitinfo(ic)
		if errMsg != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		data["storage"] = storage
	} else if step == "host" {
		hostList, err := GetHostInitinfo(ic, data)
		if err != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": err,
					"error_msg_zh": err,
				},
			}
		}
		if len(hostList) > 0 {
			data["hostList"] = hostList
		}
	} else if step == "image" {
		imageList, imageErr := GetImageInitinfo(ic)
		if imageErr != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": imageErr,
					"error_msg_zh": imageErr,
				},
			}
		}
		data["imageList"] = imageList
	} else if step == "operator" {
		operatorMap, err := GetOperatorInitinfo(ic)
		if err != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": err,
					"error_msg_zh": err,
				},
			}
		}
		data["operator"] = operatorMap
		hostList, err := GetHostInitinfo(ic, data)
		if err != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": err,
					"error_msg_zh": err,
				},
			}
		}
		if len(hostList) > 0 {
			data["hostList"] = hostList
		}
		storage, errMsg := GetStorageInitinfo(ic)
		if errMsg != "" {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		data["storage"] = storage
	}

	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    data,
		},
	}
}

func GetOperatorInitinfo(ic *InitController) (map[string]interface{}, string) {
	operatorMap := make(map[string]interface{})
	operatorInfo, err := ic.Service.SelectOneByName("operator")
	if err != "" {
		return operatorMap, err
	}

	if operatorInfo.Message == "" {
		return operatorMap, ""
	}
	jsonerr := json.Unmarshal([]byte(operatorInfo.Message), &operatorMap)
	if jsonerr != nil {
		utils.LoggerError(jsonerr)
		return operatorMap, jsonerr.Error()
	}
	return operatorMap, ""
}

func GetImageInitinfo(ic *InitController) ([]map[string]interface{}, string) {
	imageList := make([]map[string]interface{}, 0)
	imageInfo, err := ic.Service.SelectOneByName("image")
	if err != "" {
		return imageList, err
	}

	if imageInfo.Message == "" {
		return imageList, ""
	}
	jsonerr := json.Unmarshal([]byte(imageInfo.Message), &imageList)
	if jsonerr != nil {
		utils.LoggerError(jsonerr)
		return imageList, jsonerr.Error()
	}
	return imageList, ""
}

func GetHostInitinfo(ic *InitController, data map[string]interface{}) ([]map[string]interface{}, string) {
	hostList := make([]map[string]interface{}, 0)
	hostInfo, err := ic.Service.SelectOneByName("host")
	if err != "" {
		return hostList, err
	}

	if hostInfo.Message == "" || hostInfo.Message == "[]" {
		truehostList, _ := ic.NodeService.List(0, 0, "")
		data["hostList"] = truehostList
		return hostList, ""

	}
	jsonerr := json.Unmarshal([]byte(hostInfo.Message), &hostList)
	if jsonerr != nil {
		utils.LoggerError(jsonerr)
		return hostList, jsonerr.Error()
	}
	falseHostTag := 0
	for _, host := range hostList {
		node, err := ic.Service.GetNodeById(int(host["id"].(float64)))
		if err != "" {
			return hostList, err
		}
		if node.NodeName != "" {
			host["nodeName"] = node.NodeName
			host["status"] = node.Status
			host["age"] = node.Age
		} else {
			falseHostTag += 1
		}
	}
	if falseHostTag > 0 {
		truehostList, _ := ic.NodeService.List(0, 0, "")
		data["hostList"] = truehostList
		return make([]map[string]interface{}, 0), ""
	}
	return hostList, ""
}

func GetStorageInitinfo(ic *InitController) (map[string]interface{}, string) {
	storage := make(map[string]interface{})
	storageInfo, err := ic.Service.SelectOneByName("storage")
	if err != "" {
		return storage, err
	}

	if storageInfo.Message == "" {
		return storage, ""
	}
	jsonerr := json.Unmarshal([]byte(storageInfo.Message), &storage)
	if jsonerr != nil {
		utils.LoggerError(jsonerr)
		return storage, utils.ERROR_DATA_TRANSFER_EN
	}
	return storage, ""
}

func (ic *InitController) PostDeploy() mvc.Result {
	utils.LoggerInfo("初始化-部署")
	//userId := ic.Ctx.GetCookie("userId")
	//userIdInt, err := strconv.ParseInt(userId, 10, 64)
	//utils.LoggerError(err)
	//if err != nil {
	//	return mvc.Response{
	//		Object: map[string]interface{}{
	//			"errorno":      utils.RECODE_FAIL,
	//			"error_msg_en": err.Error(),
	//			"error_msg_zh": err.Error(),
	//		},
	//	}
	//}
	//userIdSlice = append(userIdSlice, int(userIdInt))
	orgTag := ic.Ctx.GetCookie("orgTag")
	userTag := ic.Ctx.GetCookie("userTag")
	userName := ic.Ctx.GetCookie("userName")

	returnStep, errMsg := ic.Service.GetLastDeployStep()
	if errMsg != "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": errMsg,
				"error_msg_zh": errMsg,
			},
		}
	}
	if returnStep == "parameter" || returnStep == "storage" {
		//初始化sc
		result, errMsg := InitSc(ic, userName, orgTag, userTag)
		if result == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化node
		noderesult, errMsg := InitNode(ic, userName)
		if noderesult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化image
		imageresult, errMsg := InitImage(ic, userName)
		if imageresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化operator
		operatorresult, errMsg := InitOperator(ic, userName)
		if operatorresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
	} else if returnStep == "host" {
		//初始化node
		noderesult, errMsg := InitNode(ic, userName)
		if noderesult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化image
		imageresult, errMsg := InitImage(ic, userName)
		if imageresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化operator
		operatorresult, errMsg := InitOperator(ic, userName)
		if operatorresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}

	} else if returnStep == "image" {
		//初始化image
		imageresult, errMsg := InitImage(ic, userName)
		if imageresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
		//初始化operator
		operatorresult, errMsg := InitOperator(ic, userName)
		if operatorresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
	} else if returnStep == "operator" {
		//初始化operator
		operatorresult, errMsg := InitOperator(ic, userName)
		if operatorresult == false {
			return mvc.Response{
				Object: map[string]interface{}{
					"errorno":      utils.RECODE_FAIL,
					"error_msg_en": errMsg,
					"error_msg_zh": errMsg,
				},
			}
		}
	}
	skipStstusInitInfo := models.Initinfo{Name: "skipStstus", Isaccess: "True"}
	ic.Service.AddOrModifyInitinfo("skipStstus", skipStstusInitInfo)
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

//初始化sc
func InitSc(ic *InitController, userName string, orgTag string, userTag string) (bool, string) {
	storageInitInfo := models.Initinfo{Name: "storage", Isdeploy: "False"}
	ic.Service.AddOrModifyInitinfo("storage", storageInitInfo)
	storage, errMsg := GetStorageInitinfo(ic)
	if errMsg != "" {
		return false, "fail to configure storage: " + errMsg
	}
	if len(storage) > 0 {
		scName := ""
		scType := ""
		nodeNum := 0
		if _, ok := storage["scName"]; ok {
			scName = storage["scName"].(string)
		}
		if _, ok := storage["scType"]; ok {
			scType = storage["scType"].(string)
		}
		if _, ok := storage["nodeNum"]; ok {
			nodeNum = int(storage["nodeNum"].(float64))
		}

		if scName == "" || scType == "" || nodeNum == 0 {
			ic.CommonService.AddLog("error", "system-init-storage", userName, "add sc "+scName+" error: "+utils.ERROR_PARAMETER_EN)
			return false, "fail to configure storage: " + "The storage cache information was not obtained"
		}
		_, err := ic.StorageService.Add(scName, "", "", orgTag, userTag, 0, scType, nodeNum, "")
		if err != nil {
			ic.CommonService.AddLog("error", "system-init-storage", userName, "add sc "+scName+" error: "+errMsg)
			return false, "fail to configure storage: " + errMsg
		}
		ic.CommonService.AddLog("info", "system-init-storage", userName, "add sc "+scName+" successful")

	}
	storageInitInfo.Isdeploy = "True"
	ic.Service.AddOrModifyInitinfo("storage", storageInitInfo)
	return true, ""
}

//初始化node
func InitNode(ic *InitController, userName string) (bool, string) {
	hostInitInfo := models.Initinfo{Name: "host", Isdeploy: "False"}
	ic.Service.AddOrModifyInitinfo("host", hostInitInfo)
	data := make(map[string]interface{})
	hostList, errMsg := GetHostInitinfo(ic, data)
	if errMsg != "" {
		return false, "fail to configure node: " + errMsg
	}
	if len(hostList) == 0 {
		return false, "fail to configure node: " + "The node cache information was not obtained"
	}
	for _, m := range hostList {
		success := true
		errMsg := ""
		if mgmtTag, ok := m["mgmtTag"]; ok {
			if mgmtTag.(bool) {
				success, errMsg = ic.NodeService.AddLabel(int(m["id"].(float64)), "iwhalecloud.dbassoperator", "mysqlha")
			} else {
				success, errMsg = ic.NodeService.DeleteLabel(int(m["id"].(float64)), "iwhalecloud.dbassoperator")
			}
		}

		if !success {
			ic.CommonService.AddLog("error", "system-init-node", userName, "node add or delete label error: "+errMsg)
			return false, "fail to configure node: " + errMsg
		}

		if mgmtTag, ok := m["computeTag"]; ok {
			if mgmtTag.(bool) {
				success, errMsg = ic.NodeService.AddLabel(int(m["id"].(float64)), "iwhalecloud.dbassnode", "mysql")
			} else {
				success, errMsg = ic.NodeService.DeleteLabel(int(m["id"].(float64)), "iwhalecloud.dbassnode")
			}
		}

		if !success {
			ic.CommonService.AddLog("error", "system-init-node", userName, "node add or delete label error: "+errMsg)
			return false, "fail to configure node: " + errMsg
		}
	}
	ic.NodeService.AsyncDbLabel()
	ic.CommonService.AddLog("info", "system-init-node", userName, "node add or delete label successful")
	hostInitInfo.Isdeploy = "True"
	ic.Service.AddOrModifyInitinfo("host", hostInitInfo)
	return true, ""
}

//初始化image
func InitImage(ic *InitController, userName string) (bool, string) {
	imageInitInfo := models.Initinfo{Name: "image", Isdeploy: "False"}
	ic.Service.AddOrModifyInitinfo("image", imageInitInfo)
	imageInfo, err := ic.Service.SelectOneByName("image")
	if err != "" {
		return false, "fail to configure image: " + err
	}

	if imageInfo.Message == "" {
		return false, "fail to configure image: " + "The image cache information was not obtained"
	}

	imageList := make([]models.Images, 0)
	jsonerr := json.Unmarshal([]byte(imageInfo.Message), &imageList)
	if jsonerr != nil {
		return false, "fail to configure image: " + jsonerr.Error()
	}

	deployerr := ic.ImageService.InitAdd(imageList)
	if deployerr == nil {
		imageInitInfo.Isdeploy = "True"
		ic.Service.AddOrModifyInitinfo("image", imageInitInfo)
		ic.CommonService.AddLog("info", "system-init-image", userName, fmt.Sprintf("init add images successful"))
		return true, ""
	}
	var errMsgEn = deployerr.Error()
	if strings.Contains(errMsgEn, "unique_name_version") {
		ic.CommonService.AddLog("error", "system-image", userName, fmt.Sprintf("init add image error %v", errMsgEn))
		errMsgEn = "The image already exists!"
		return false, "fail to configure image: " + errMsgEn
	}
	return true, ""
}

//初始化operator
func InitOperator(ic *InitController, userName string) (bool, string) {
	operatorInitInfo := models.Initinfo{Name: "operator", Isdeploy: "False"}
	ic.Service.AddOrModifyInitinfo("operator", operatorInitInfo)
	operatorMap, err := GetOperatorInitinfo(ic)
	if err != "" {
		return false, "fail to deploy operator: " + "The Operator cache information was not obtained"

	}
	scName := ""
	mode := ""
	if _, ok := operatorMap["scName"]; ok {
		scName = operatorMap["scName"].(string)
	}
	if _, ok := operatorMap["mode"]; ok {
		mode = operatorMap["mode"].(string)
	}
	if scName == "" || mode == "" {
		ic.CommonService.AddLog("error", "system-init-operator", userName, "The operator cache information was not obtained")
		return false, "fail to deploy operator: " + "The operator cache information was not obtained"
	}
	operatorErr := ic.NodeService.Operator(mode, scName)
	if operatorErr != nil {
		ic.CommonService.AddLog("error", "system-init-operator", userName, "operator error: "+operatorErr.Error())
		return false, "fail to deploy operator: " + operatorErr.Error()
	}
	ic.CommonService.AddLog("info", "system-init-operator", userName, "operator successful")
	operatorInitInfo.Isdeploy = "True"
	ic.Service.AddOrModifyInitinfo("operator", operatorInitInfo)
	return true, ""
}

//初始化-跳过信息录入
func (ic *InitController) PostSkip() mvc.Result {
	utils.LoggerInfo("初始化-跳过信息录入")
	isSkip := ic.Ctx.PostValue("isSkip")
	if isSkip == "" {
		return mvc.Response{
			Object: map[string]interface{}{
				"errorno":      utils.RECODE_FAIL,
				"error_msg_en": utils.ERROR_PARAMETER_EN,
				"error_msg_zh": utils.ERROR_PARAMETER_ZH,
			},
		}
	}
	skipStstusInitInfo := models.Initinfo{Name: "skipStstus", Isaccess: "True"}
	ic.Service.AddOrModifyInitinfo("skipStstus", skipStstusInitInfo)
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
		},
	}
}

//初始化-跳过信息查询
func (ic *InitController) GetStatus() mvc.Result {
	utils.LoggerInfo(" 初始化-跳过信息查询 ")
	skipStstusMap := make(map[string]interface{})
	skipStstusInfo, err := ic.Service.SelectOneByName("skipStstus")
	if err != "" {
		skipStstusMap["isSkip"] = "True"
	} else {
		skipStstusMap["isSkip"] = skipStstusInfo.Isaccess
	}
	return mvc.Response{
		Object: map[string]interface{}{
			"errorno": utils.RECODE_OK,
			"data":    skipStstusMap,
		},
	}
}
