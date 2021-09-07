package dbaasgrpcservice

import (
	"DBaas/config"
	"DBaas/service"
	"DBaas/utils"
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type Service struct {
	ps service.PodService
}

func (s *Service) GetPodInfoForAlarm(ctx context.Context, in *Empty) (*Response, error) {
	returnData, err := s.ps.SelectPodsForAlarm()
	if err != nil {
		errMsg := err.Error()
		return &Response{Errorno: -1, ErrorMsgEn: errMsg, ErrorMsgZh: errMsg, Data: "{}"}, nil
	}
	j, err := json.Marshal(returnData)
	utils.LoggerError(err)
	return &Response{Errorno: 0, ErrorMsgEn: "", ErrorMsgZh: "", Data: string(j)}, nil
}

func RungGRPCServer(ps service.PodService, c *config.AppConfig) {
	// 启动一个grpc server
	grpcServer := grpc.NewServer()
	// 绑定服务实现 RegisterHelloWorldServiceServer
	RegisterDbaasgrpcserviceServer(grpcServer, &Service{ps})
	// 监听端口
	listen, err := net.Listen("tcp", ":"+c.GrpcServerPort)
	utils.LoggerError(err)
	// 绑定监听端口
	utils.LoggerInfo("serve gRPC server: 127.0.0.1:" + c.GrpcServerPort)
	if err = grpcServer.Serve(listen); err != nil {
		utils.LoggerInfo("failed to serve:", err)
		return
	}
}
