package service

import (
	"DBaas/config"
	"DBaas/grpc/capricorn"
	"DBaas/utils"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func initCapricornService() *grpc.ClientConn {
	// capricorn grpc 连接引擎
	c := config.GetConfig()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeout, c.CAddress, grpc.WithInsecure(), grpc.WithBlock())
	utils.LoggerError(err)
	return conn
}

type CapricornService interface {
	GetUserResources(id string, username string, tag string) ([]map[string]interface{},string, string)
	AddUserResources(roleId string, operUsername string,username string, password string,organizationId string, isManager string) (map[string]interface{},string, string)
	UpdateUserResources(id string,roleId string, operUsername string,username string, password string,organizationId string, isManager string) (map[string]interface{},string, string)
	DeleteUserResources(id string, operUsername string) (string,string, string)
	OperateUserResources(id string, operUsername string) (string,string, string)
	GetRoleResources(id string, roleName string) ([]map[string]interface{},string, string)
	GetRandomPassword(userid string, operUsername string) (map[string]interface{},string, string)
}

type capricornService struct {
	conn *grpc.ClientConn
}

func NewCapricornService() (CapricornService, *grpc.ClientConn) {
	conn := initCapricornService()
	return &capricornService{
		conn: conn,
	}, conn
}

//获取用户信息
func (cs *capricornService) GetUserResources(id string, username string, tag string) ([]map[string]interface{},string, string ){
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var userMap []map[string]interface{}
	res, err := c.GetUserResources(ctx, &capricorn.GetUserResourcesRequest{Id: id, Username: username, Tag: tag})
	if err != nil {
		utils.LoggerError(err)
		return userMap,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &userMap)
	utils.LoggerError(err)
	return userMap,res.ErrorMsgEn, res.ErrorMsgZh
}

//新增用户信息
func (cs *capricornService) AddUserResources(roleId string, operUsername string,username string, password string,organizationId string, isManager string) (map[string]interface{},string, string) {
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	returndata := make(map[string]interface{})
	res, err := c.AddUserResources(ctx, &capricorn.AddUserResourcesRequest{RoleId: roleId, OperUsername: operUsername,Username:username,Password:password,OrganizationId:organizationId,IsManager:isManager})
	if err != nil {
		utils.LoggerError(err)
		return returndata,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &returndata)
	utils.LoggerError(err)
	return returndata,res.ErrorMsgEn, res.ErrorMsgZh
}

//修改用户信息
func (cs *capricornService) UpdateUserResources(id string,roleId string, operUsername string,username string, password string,organizationId string, isManager string) (map[string]interface{},string, string)  {
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	returndata := make(map[string]interface{})
	res, err := c.UpdateUserResources(ctx, &capricorn.UpdateUserResourcesRequest{Id:id,RoleId: roleId, OperUsername: operUsername,Username:username,Password:password,OrganizationId:organizationId,IsManager:isManager})
	if err != nil {
		utils.LoggerError(err)
		return returndata,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &returndata)
	utils.LoggerError(err)
	return returndata,res.ErrorMsgEn, res.ErrorMsgZh
}



//删除用户信息
func (cs *capricornService) DeleteUserResources(id string, operUsername string) (string,string, string) {
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var returndata string
	res, err := c.DeleteUserResources(ctx, &capricorn.DeleteUserResourcesRequest{Id: id, OperUsername: operUsername})
	if err != nil {
		utils.LoggerError(err)
		return returndata,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &returndata)
	utils.LoggerError(err)
	return returndata,res.ErrorMsgEn, res.ErrorMsgZh
}

//启用/禁用用户信息
func (cs *capricornService) OperateUserResources(id string, operUsername string) (string,string, string) {
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var returndata string
	res, err := c.OperateUserResources(ctx, &capricorn.OperateUserResourcesRequest{Id: id, OperUsername: operUsername})
	if err != nil {
		utils.LoggerError(err)
		return returndata,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &returndata)
	utils.LoggerError(err)
	return returndata,res.ErrorMsgEn, res.ErrorMsgZh
}

//获取角色信息
func (cs *capricornService) GetRoleResources(id string, roleName string) ([]map[string]interface{}, string, string) {
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var roleMap []map[string]interface{}
	res, err := c.GetRoleResources(ctx, &capricorn.GetRoleResourcesRequest{Id: id, RoleName: roleName})
	if err != nil {
		utils.LoggerError(err)
		return roleMap, "User module grpc communication error", "用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &roleMap)
	utils.LoggerError(err)
	return roleMap, res.ErrorMsgEn, res.ErrorMsgZh
}

//获取随机密码信息
func (cs *capricornService) GetRandomPassword(userid string, operUsername string) (map[string]interface{},string, string){
	c := capricorn.NewCapricornClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	passwordMap := make(map[string]interface{})
	res, err := c.GetRandomPassword(ctx, &capricorn.GetRandomPasswordRequest{Userid: userid, OperUsername: operUsername})
	if err != nil {
		utils.LoggerError(err)
		return passwordMap,"User module grpc communication error","用户模块grpc通信错误"
	}
	err = json.Unmarshal([]byte(res.Data), &passwordMap)
	utils.LoggerError(err)
	return passwordMap,res.ErrorMsgEn, res.ErrorMsgZh
}