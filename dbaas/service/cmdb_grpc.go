package service

import (
	"DBaas/config"
	"DBaas/grpc/cmdb"
	"DBaas/utils"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func initCmdService() *grpc.ClientConn {
	// cmdb grpc 连接引擎
	c := config.GetConfig()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeout, c.CAddress, grpc.WithInsecure(), grpc.WithBlock())
	utils.LoggerError(err)
	return conn
}

type CmdService interface {
	GetCmdbFilterResources(User int32, ModelCode string, FilterCondition string) ([]map[string]interface{}, string, string, bool)
	GetCmdbResources(User int32, ModelCode string, InstId int32, FilterCondition string) (map[string]interface{}, string, string, bool)
	DeleteInstApply(Module string, Type string, ApplyId int32) (bool, string, string)
	AddInstApply(ModelId string, Module string, InstId string, Type string, ApplyId int32, Desc string) (bool, string, string)
	GetCmdbModelField(ModelId int32, selectType string, instId int32, attrId int32, attrName string, collectionIndex string, modelCode string) ([]interface{}, string, string)
	GetDeployHostinstRequest(sql string, user int32, status string) []map[string]interface{}
	GetCmdbObjectInst(Objectinst string) []map[string]interface{}
	DeleteInstanceDb(userId string, orgTag string, userTag string, instId int32) (bool, string, string)
	AddInstanceDb(dbName string, orgTag string, userTag string, modelCode string, moduleName string, field string) (map[string]interface{}, string, string)
	GetRelatedInst(instId int32, user int32, modelId int32, Type string, module string, modelCode string, filterCondition string) ([]map[string]interface{}, string, string)
	RelateInstance(instId int32, userId string, RalateInstId string, Operation string, OrgTag string, UserTag string, Role string, ApiPath string) (bool, string, string)
	UpdateInstanceDb(instId int32, dbName string, orgTag string, userTag string, modelCode string, moduleName string, field string, user string) (map[string]interface{}, string, string)
	GetFreeRelatedInst(user int32, modelId int32, module string, modelCode string, filterCondition string, relateModelId int32, relateModelCode string) ([]map[string]interface{}, string, string)
	JudgeSameClusterdbInst(port string, hostids string) (bool, string, string)
	WorkflowAddInstance(createInfo string, relateInfo string) (map[string]interface{}, string, string)
	WorkflowGetInstance(instId string, filterCondition string, user int32, module string, modelId int32, instRelInfo string, authTag string) (map[string]interface{}, string, string)
	WorkflowDeleteInstance(userId string, orgTag string, userTag string, instId string, authTag string, module string, deleteModel string) (map[string]interface{}, string, string)
	WorkflowGetInstRels(userId string, instId string, authTag string, module string, instRelInfo string) (map[string]interface{}, string, string)
	GetAreainfos(selecttype string, key string) ([]interface{}, string, string)
	AddUpdateDeploy(hostId string, name string, areaId string, orgTag string, userTag string, module string, deployUsername string, deployPassword string, deployIp string, deployPort string, field string, operation string, status string, initvalue string, userId string, authTag string) (map[string]interface{}, string, string)
}

type cmdService struct {
	conn *grpc.ClientConn
}

func NewCmdService() (CmdService, *grpc.ClientConn) {
	conn := initCmdService()
	return &cmdService{
		conn: conn,
	}, conn
}

//  获取主机的列表
func (cs *cmdService) GetDeployHostinstRequest(sql string, user int32, status string) []map[string]interface{} {
	c := cmdb.NewCmdbClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetDeployHostinst(ctx, &cmdb.GetDeployHostinstRequest{User: user, FilterCondition: sql, Status: status})
	var hostMap []map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return hostMap
	}
	err = json.Unmarshal([]byte(res.Data), &hostMap)
	utils.LoggerError(err)
	return hostMap
}

//  获取主机实例的详情
func (cs *cmdService) GetCmdbObjectInst(Objectinst string) []map[string]interface{} {
	c := cmdb.NewCmdbClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetCmdbObjectInst(ctx, &cmdb.GetObjectInstRequest{Objectinst: Objectinst})
	var hostDetailMap []map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return hostDetailMap
	}
	err = json.Unmarshal([]byte(res.Data), &hostDetailMap)
	utils.LoggerError(err)
	return hostDetailMap
}

//  新增实例
func (cs *cmdService) AddInstanceDb(dbName string, orgTag string, userTag string, modelCode string, moduleName string, field string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.AddInstance(ctx, &cmdb.AddInstanceRequest{Name: dbName, ModelCode: modelCode, Field: field, ModuleName: moduleName, UserTag: userTag, OrgTag: orgTag})
	var dbBase map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return dbBase, "cmdb rpc error", "cmdb rpc error"
	}
	err = json.Unmarshal([]byte(res.Data), &dbBase)
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  修改数据库实例
func (cs *cmdService) UpdateInstanceDb(instId int32, dbName string, orgTag string, userTag string, modelCode string, moduleName string, field string, user string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.UpdateInstance(ctx, &cmdb.UpdateInstanceRequest{InstId: instId, Name: dbName, ModelCode: modelCode, Field: field, ModuleName: moduleName, UserTag: userTag, OrgTag: orgTag, Userid: user})
	var dbBase map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return dbBase, "failed to modify the CMDB instance information", "修改cmdb实例信息失败"
	}
	err = json.Unmarshal([]byte(res.Data), &dbBase)
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  删除数据库实例
func (cs *cmdService) DeleteInstanceDb(userId string, orgTag string, userTag string, instId int32) (bool, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.DeleteInstance(ctx, &cmdb.DeleteInstanceRequest{Userid: userId, OrgTag: orgTag, UserTag: userTag, InstId: instId})
	if err != nil {
		utils.LoggerError(err)
		//println(err.Error())
		return false, "", ""
	}
	if res.Errorno == 0 {
		return true, "", ""
	}
	return false, res.ErrorMsgEn, res.ErrorMsgZh
}

//  获取图表类模型字段接口
func (cs *cmdService) GetCmdbModelField(ModelId int32, selectType string, instId int32, attrId int32, attrName string, collectionIndex string, modelCode string) ([]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetCmdbModelField(ctx, &cmdb.GetModelFieldRequest{ModelId: ModelId, Type: selectType, InstId: instId, AttrId: attrId, AttrName: attrName, CollectionIndex: collectionIndex, ModelCode: modelCode})
	var modelFieldMap []interface{}
	var modelField map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return modelFieldMap, "", ""
	} else {
		if res.Errorno == 0 {
			err = json.Unmarshal([]byte(res.Data), &modelField)
			if err != nil {
				utils.LoggerError(err)
			}
			modelFieldMap = modelField["detail"].([]interface{})
			return modelFieldMap, res.ErrorMsgEn, res.ErrorMsgZh
		} else {
			return modelFieldMap, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

//
func (cs *cmdService) AddInstApply(ModelId string, Module string, InstId string, Type string, ApplyId int32, Desc string) (bool, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.AddInstApply(ctx, &cmdb.AddInstApplyRequest{ModelId: ModelId, Module: Module, InstId: InstId, Type: Type, ApplyId: ApplyId, Desc: Desc})
	if err != nil {
		return false, "", ""
	} else {
		if res.Errorno == 0 {
			return true, "", ""
		} else {
			fmt.Println(res.ErrorMsgZh)
			return false, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

func (cs *cmdService) DeleteInstApply(Module string, Type string, ApplyId int32) (bool, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	fmt.Println(Module)
	fmt.Println(Type)
	fmt.Println(ApplyId)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.DeleteInstApply(ctx, &cmdb.DeleteInstApplyRequest{Module: Module, Type: Type, ApplyId: ApplyId})
	fmt.Println(res.ErrorMsgZh)
	if err != nil {
		utils.LoggerError(err)
		return false, "", ""
	} else {
		if res.Errorno == 0 {
			return true, "", ""
		} else {
			return false, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

//获取实例
func (cs *cmdService) GetCmdbResources(User int32, ModelCode string, InstId int32, FilterCondition string) (map[string]interface{}, string, string, bool) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetCmdbResources(ctx, &cmdb.GetResourcesRequest{User: User, ModelCode: ModelCode, InstId: InstId, FilterCondition: FilterCondition})
	var cmdbResourcesMap map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return cmdbResourcesMap, "", "", false
	} else {
		if res.Errorno == 0 {
			err = json.Unmarshal([]byte(res.Data), &cmdbResourcesMap)
			if err != nil {
				utils.LoggerError(err)
				println(err)
			}
			return cmdbResourcesMap, "", "", true
		} else {
			return cmdbResourcesMap, res.ErrorMsgEn, res.ErrorMsgZh, false
		}
	}
}

//获取实例
func (cs *cmdService) GetCmdbFilterResources(User int32, ModelCode string, FilterCondition string) ([]map[string]interface{}, string, string, bool) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetCmdbResources(ctx, &cmdb.GetResourcesRequest{User: User, ModelCode: ModelCode, FilterCondition: FilterCondition})
	var cmdbResourcesMap []map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return cmdbResourcesMap, "", "", false
	} else {
		if res.Errorno == 0 {
			err = json.Unmarshal([]byte(res.Data), &cmdbResourcesMap)
			if err != nil {
				println(err)
			}
			return cmdbResourcesMap, "", "", true
		} else {
			return cmdbResourcesMap, res.ErrorMsgEn, res.ErrorMsgZh, false
		}
	}
}

//  获取实例关联的模型下实例接口
func (cs *cmdService) GetRelatedInst(instId int32, user int32, modelId int32, Type string, module string, modelCode string, filterCondition string) ([]map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetRelatedInst(ctx, &cmdb.GetRelatedInstRequest{ModelId: modelId, InstId: instId, User: user, Type: Type, ModelCode: modelCode, Module: module, FilterCondition: filterCondition})
	var relatedInstMap []map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return relatedInstMap, "", ""
	} else {
		if res.Errorno == 0 {
			err = json.Unmarshal([]byte(res.Data), &relatedInstMap)
			if err != nil {
				fmt.Println(err)
			}
			return relatedInstMap, res.ErrorMsgEn, res.ErrorMsgZh
		} else {
			return relatedInstMap, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

//  实例关联的模型下实例接口
func (cs *cmdService) RelateInstance(instId int32, userId string, RalateInstId string, Operation string, OrgTag string, UserTag string, Role string, ApiPath string) (bool, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.RelateInstance(ctx, &cmdb.RelateInstanceRequest{InstId: instId, Userid: userId, RalateInstId: RalateInstId, Operation: Operation, OrgTag: OrgTag, UserTag: UserTag, Role: Role, ApiPath: ApiPath})
	if err != nil {
		utils.LoggerError(err)
		return false, "", ""
	} else {
		if res.Errorno == 0 {
			return true, res.ErrorMsgEn, res.ErrorMsgZh
		} else {
			return false, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

//  获取实例关联的模型下实例接口
func (cs *cmdService) GetFreeRelatedInst(user int32, modelId int32, module string, modelCode string, filterCondition string, relateModelId int32, relateModelCode string) ([]map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.GetRelatedFreeInst(ctx, &cmdb.GetRelatedFreeInstRequest{ModelId: modelId, User: user, ModelCode: modelCode, Module: module, FilterCondition: filterCondition, RelateModelCode: relateModelCode, RelateModelId: relateModelId})
	var relatedFreeInstMap []map[string]interface{}
	if err != nil {
		utils.LoggerError(err)
		return relatedFreeInstMap, "", ""
	} else {
		if res.Errorno == 0 {
			err = json.Unmarshal([]byte(res.Data), &relatedFreeInstMap)
			if err != nil {
				fmt.Println(err)
			}
			return relatedFreeInstMap, res.ErrorMsgEn, res.ErrorMsgZh
		} else {
			return relatedFreeInstMap, res.ErrorMsgEn, res.ErrorMsgZh
		}
	}
}

//  判断是否拥有相同ip和port的mysql实例接口
func (cs *cmdService) JudgeSameClusterdbInst(port string, hostids string) (bool, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.JudgeSameClusterdbInst(ctx, &cmdb.JudgeSameClusterdbInstRequest{Port: port, Hostids: hostids})
	if err != nil {
		utils.LoggerError(err)
		return false, "CMDB grpc communication error, unable to determine whether there is the same database instance", "cmdb grpc通信错误，无法判断是否有相同数据库实例"
	} else {
		if res.Result {
			return res.Result, "a database instance with the same port already exists on the selected host", "目前所选主机上已拥有相同端口的数据库实例"
		} else {
			return res.Result, "", ""
		}
	}
}

//  以工作流形式新增实例
func (cs *cmdService) WorkflowAddInstance(createInfo string, relateInfo string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.WorkflowAdd(ctx, &cmdb.WorkflowAddRequest{CreateInfo: createInfo, RelateInfo: relateInfo})
	var dbBase map[string]interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return dbBase, "", ""
	}
	err = json.Unmarshal([]byte(res.Data), &dbBase)
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  以工作流形式查询实例信息
func (cs *cmdService) WorkflowGetInstance(instId string, filterCondition string, user int32, module string, modelId int32, instRelInfo string, authTag string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.WorkflowGetInst(ctx, &cmdb.WorkflowGetInstRequest{InstId: instId, FilterCondition: filterCondition, User: user, Module: module, ModelId: modelId, InstRelInfo: instRelInfo, AuthTag: authTag})
	var dbBase map[string]interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return dbBase, "", ""
	}
	err = json.Unmarshal([]byte(res.Data), &dbBase)
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  以工作流形式删除实例
func (cs *cmdService) WorkflowDeleteInstance(userId string, orgTag string, userTag string, instId string, authTag string, module string, deleteModel string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.WorkflowDeleteInst(ctx, &cmdb.WorkflowDeleteInstRequest{Userid: userId, OrgTag: orgTag, UserTag: userTag, InstId: instId, AuthTag: authTag, Module: module, DeleteModel: deleteModel})
	var dbBase map[string]interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return dbBase, "", ""
	}
	returnData := res.Data
	jsonerr := json.Unmarshal([]byte(returnData), &dbBase)
	if jsonerr != nil {
		fmt.Println(jsonerr)
	}
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  以工作流形式获取实例下关联实例id
func (cs *cmdService) WorkflowGetInstRels(userId string, instId string, authTag string, module string, instRelInfo string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.WorkflowGetInstRels(ctx, &cmdb.WorkflowGetInstRelsRequest{Userid: userId, InstId: instId, AuthTag: authTag, Module: module, InstRelInfo: instRelInfo})
	var dbBase map[string]interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return dbBase, "", ""
	}
	err = json.Unmarshal([]byte(res.Data), &dbBase)
	return dbBase, res.ErrorMsgEn, res.ErrorMsgZh
}

//  以工作流形式获取实例下关联实例id
func (cs *cmdService) GetAreainfos(selecttype string, key string) ([]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.GetAreainfos(ctx, &cmdb.GetAreainfosRequest{Type: selecttype, Key: key})
	var Areainfos []interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return Areainfos, "", ""
	}
	err = json.Unmarshal([]byte(res.Data), &Areainfos)
	return Areainfos, res.ErrorMsgEn, res.ErrorMsgZh
}

//  以工作流形式获取实例下关联实例id
func (cs *cmdService) AddUpdateDeploy(hostId string, name string, areaId string, orgTag string, userTag string, module string, deployUsername string, deployPassword string, deployIp string, deployPort string, field string, operation string, status string, initvalue string, userId string, authTag string) (map[string]interface{}, string, string) {
	c := cmdb.NewCmdbClient(cs.conn)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := c.AddUpdateDeploy(ctx, &cmdb.AddUpdateDeployHostRequest{HostId: hostId, Name: name, AreaId: areaId, OrgTag: orgTag, UserTag: userTag, Module: module, DeployUsername: deployUsername, DeployPassword: deployPassword, DeployIp: deployIp, DeployPort: deployPort, Field: field, Operation: operation, DeployStatus: status, InitValue: initvalue, UserId: userId, AuthTag: authTag})
	var hostmap map[string]interface{}
	//fmt.Println(res.Errorno)
	if err != nil {
		utils.LoggerError(err)
		//fmt.Println(err.Error())
		return hostmap, "", ""
	}
	err = json.Unmarshal([]byte(res.Data), &hostmap)
	return hostmap, res.ErrorMsgEn, res.ErrorMsgZh
}
