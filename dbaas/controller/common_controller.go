/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: ddh
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: ddh
 * @LastEditTime: 2021-02-07 16:32:07
 */

package controller

import (
	"DBaas/service"
	"github.com/kataras/iris/v12"
)

type CommonController struct {
	//iris框架自动为每个请求都绑定上下文对象
	Ctx iris.Context
	//common功能实体
	Service service.CommonService
	SysparamService service.ParameterService
	PodService       service.PodService

}


//  获取采集的性能数据
//func (cc *CommonController) GetCommonDetail() mvc.Result {
//	iris.New().Logger().Info(" 获取采集的性能数据  ")
//	attrId := cc.Ctx.URLParam("attrId")
//	selectType := cc.Ctx.URLParam("type")
//	modelId := cc.Ctx.URLParam("modelId")
//	condition := cc.Ctx.URLParam("condition")
//	conditionMap := map[string]interface{}{}
//	err := json.Unmarshal([]byte(condition), &conditionMap)
//	utils.LoggerError(err)
//	podId := cc.Ctx.URLParam("podId")
//	time := cc.Ctx.URLParam("time")
//	attrIdInt, _ := strconv.ParseInt(attrId, 10, 32)
//	podIdInt, _ := strconv.ParseInt(podId, 10, 32)
//	pod, errMsg := cc.PodService.SelectOne(int(podIdInt))
//	if errMsg != "" {
//			return mvc.Response{
//				Object: map[string]interface{}{
//					"errorno":      utils.RECODE_FAIL,
//					"data":         "{}",
//					"error_msg_en": errMsg,
//					"error_msg_zh": errMsg,
//				},
//			}
//	}
//	modelIdInt, _ := strconv.ParseInt(modelId, 10, 32)
//	timeInt, _ := strconv.ParseInt(time, 10, 32)
//	if timeInt == 0 {
//		timeInt = 3000
//	}
//	hostDetailInformation := service.GetPerformanceData(int32(modelIdInt), selectType, pod.Name, int32(attrIdInt), "", timeInt, conditionMap)
//	return mvc.Response{
//		Object: map[string]interface{}{
//			"errorno": utils.RECODE_OK,
//			"data":    hostDetailInformation,
//		},
//	}
//}
