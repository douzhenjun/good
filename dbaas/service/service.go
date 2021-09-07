package service

import (
	"DBaas/utils"
	"github.com/go-xorm/xorm"
	"google.golang.org/grpc"
	core1 "k8s.io/api/core/v1"
)

/*
分页查询
slicePtr: 切片指针,
session: 查询会话(内部会clone一份,不影响外部使用),
t: 要查询的类型 new(models.xxx)
*/
func pageFind(page, pageSize int, slicePtr interface{}, session *xorm.Session, t interface{}) (count int64, err error) {
	findSession := session.Clone()
	countSession := session.Clone()
	if utils.MustInt(page, pageSize) {
		findSession.Limit(pageSize, pageSize*(page-1))
	}
	err = findSession.Find(slicePtr)
	if err != nil {
		return
	}
	count, _ = countSession.Count(t)
	return
}

/*
ContainerStatusF 容器状态格式化
*/
func ContainerStatusF(s core1.ContainerStatus) string {
	switch {
	case s.State.Waiting != nil:
		return s.State.Waiting.Reason
	case s.State.Running != nil:
		return "Running"
	case s.State.Terminated != nil:
		return s.State.Terminated.Reason
	}
	return "FError"
}

var eventSelector = map[string]struct{}{"Pod": {}, "MysqlCluster": {}}

func CloseGrpc(conn *grpc.ClientConn) {
	if conn != nil {
		utils.LoggerError(conn.Close())
	}
}
