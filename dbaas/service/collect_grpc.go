package service

import (
	"DBaas/config"
	"DBaas/grpc/collect"
	"DBaas/utils"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

func initCollectService() *grpc.ClientConn {
	// collect grpc 连接引擎
	c := config.GetConfig()
	timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeout, c.CAddress, grpc.WithInsecure(), grpc.WithBlock())
	utils.LoggerError(err)
	return conn
}

type CollectService interface {
	GetInstItemList(userTag string, modelId int32, instId int32) map[string][]map[string]interface{}
	GetInfluxDbData(sql string, sec string) []map[string]interface{}
	AddInsts(modelId int32, instId int32, collectUser string, isHost bool, hostIp string, hostSn string) (bool, string)
	ChangeHostStatus(hostId int32, hostIp string, hostSn string, status string) (bool, string)
}

type collectService struct {
	conn *grpc.ClientConn
}

func NewCollectService() (CollectService, *grpc.ClientConn) {
	conn := initCollectService()
	return &collectService{
		conn: conn,
	}, conn
}

func (cs *collectService) GetInstItemList(userTag string, modelId int32, instId int32) map[string][]map[string]interface{} {
	c := collect.NewCollectClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var collectTable map[string][]map[string]interface{}
	res, err := c.GetInstItemList(ctx, &collect.GetInstItemListRequest{UserTag: userTag, ModelId: modelId, InstId: instId})
	if err != nil {
		return collectTable
	}
	err = json.Unmarshal([]byte(res.Data), &collectTable)
	utils.LoggerError(err)
	return collectTable
}

func (cs *collectService) GetInfluxDbData(sql string, sec string) []map[string]interface{} {
	c := collect.NewCollectClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var collectTable []map[string]interface{}
	res, err := c.GetInfluxdbData(ctx, &collect.GetInfluxdbDataRequest{Sql: sql, Epoch: sec})
	if err != nil {
		utils.LoggerInfo("influxdb 查询出错:", err)
		return collectTable
	}
	err = json.Unmarshal([]byte(res.Result), &collectTable)
	utils.LoggerError(err)
	return collectTable
}

// 新增实例,支持批量新增
func (cs *collectService) AddInsts(modelId int32, instId int32, collectUser string, isHost bool, hostIp string, hostSn string) (bool, string) {
	c := collect.NewCollectClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cc := &collect.Inst{ModelId: modelId, InstId: instId, CollectUser: collectUser, IsHost: isHost, HostIp: hostIp, HostSn: hostSn}
	fc := make([]*collect.Inst, 0)
	fc = append(fc, cc)
	res, err := c.AddInsts(ctx, &collect.AddInstsRequest{Inst: fc})
	if err != nil {
		utils.LoggerError(err)
		return false, err.Error()
	}

	return res.Result, res.Msg
}

// 新增实例,支持批量新增
func (cs *collectService) ChangeHostStatus(hostId int32, hostIp string, hostSn string, status string) (bool, string) {
	c := collect.NewCollectClient(cs.conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := c.ChangeHostStatus(ctx, &collect.ChangeHostStatusRequest{HostId: hostId, HostIp: hostIp, HostSn: hostSn, ChangedStatus: status})
	if err != nil {
		utils.LoggerError(err)
		return false, err.Error()
	}

	return res.Result, res.Msg
}
