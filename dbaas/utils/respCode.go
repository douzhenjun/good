/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: Dou
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: Dou
 * @LastEditTime: 2021-02-07 16:32:07
 */

package utils

//请求状态码
const (
	RECODE_OK      = 0  //请求成功 正常
	RECODE_FAIL    = 1  //失败
	RECODE_UNLOGIN = -1 //未登录 没有权限
)

func ResponseOk(data interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"errorno":    RECODE_OK,
		"error_msg_en": "",
		"error_msg_zh": "",
	}
	if data != nil {
		result["data"] = data
	}
	return result
}

func ResponseFail(errMsgEn string, errMsgZh string) map[string]interface{} {
	return map[string]interface{}{
		"errorno":    RECODE_FAIL,
		"error_msg_en": errMsgEn,
		"error_msg_zh": errMsgZh,
	}
}

// 响应信息国际化
const (
	ERROR_MSG_EN = "fail"
	ERROR_MSG_ZH = "请求失败"

	ERROR_PARAMETER_EN = "Parameter error"
	ERROR_PARAMETER_ZH = "参数错误"

	ERROR_SWITCH_EN = "MySQL instance switch failed"
	ERROR_SWITCH_ZH = "MySQL实例切换失败"

	ERROR_PASSWORD_EN = "Password error"
	ERROR_PASSWORD_ZH = "密码错误"

	ERROR_DATAQUERY_EN = "Data query error"
	ERROR_DATAQUERY_ZH = "数据查询出错"

	ERROR_DELETE_EN = "Delete error"
	ERROR_DELETE_ZH = "删除失败"

	ERROR_ADD_EN = "Add error"
	ERROR_ADD_ZH = "新增失败"

	ERROR_UPDATE_EN = "Update error"
	ERROR_UPDATE_ZH = "修改失败"

	ERROR_DEPLOY_EN = "Deploy error"
	ERROR_DEPLOY_ZH = "部署失败"

	ERROR_DATA_TRANSFER_EN = "Data transfer error"
	ERROR_DATA_TRANSFER_ZH = "数据转换格式时出错"

	ERROR_DELETE_HOST_MSG_EN = "There is a PostgreSQLClusterMember under this host. Please delete the PostgreSQLClusterMember first"
	ERROR_DELETE_HOST_MSG_ZH = "该主机下有PostgreSQL集群成员，请先删除PostgreSQL集群成员"
)
