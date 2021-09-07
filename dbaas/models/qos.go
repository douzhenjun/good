package models

import disegrpc "DBaas/grpc/dise"

type Qos struct {
	Id        int `xorm:"notnull pk autoincr unique" json:"-"`
	ClusterId int `xorm:"notnull" json:"-"`

	*QosLite `xorm:"extends" json:"-"`
}

type QosLite struct {
	ReadBps   int64 `xorm:"notnull" json:"qosReadBPS,omitempty"`
	WriteBps  int64 `xorm:"notnull" json:"qosWriteBPS,omitempty"`
	ReadIops  int64 `xorm:"notnull" json:"qosReadIops,omitempty"`
	WriteIops int64 `xorm:"notnull" json:"qosWriteIops,omitempty"`
}

// TG QosLite转换为grpc中的类型
func (q *QosLite) TG() *disegrpc.QoS {
	return &disegrpc.QoS{
		ReadBps:   q.ReadBps,
		WriteBps:  q.WriteBps,
		ReadIops:  q.ReadIops,
		WriteIops: q.WriteIops,
	}
}

// FG 由grpc中的类型转换为QosLite
func (q *QosLite) FG(qos *disegrpc.QoS) *QosLite {
	q.ReadBps = qos.ReadBps
	q.WriteBps = qos.ReadBps
	q.ReadIops = qos.ReadIops
	q.WriteIops = qos.WriteIops
	return q
}
