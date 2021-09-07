package service

import (
	disegrpc "DBaas/grpc/dise"
	"DBaas/models"
	"DBaas/utils"
	"fmt"
	"github.com/go-xorm/xorm"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type DiseService interface {
	SetQoSOfVolume(pv string, qos *models.QosLite) error
	GetQosOfVolume(pv string) (*models.QosLite, error)
}

type diseService struct {
	conn *grpc.ClientConn
}

func (ds *diseService) GetQosOfVolume(pv string) (*models.QosLite, error) {
	in := new(disegrpc.GetQoSOfVolumeRequest)
	in.VolumeName = pv
	out := new(disegrpc.GetQoSOfVolumeResponse)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := ds.conn.Invoke(ctx, "/disegrpc.DISE/GetQoSOfVolume", in, out)
	if err != nil {
		return nil, err
	}
	return new(models.QosLite).FG(out.RWQoS), nil
}

func (ds *diseService) SetQoSOfVolume(pv string, qos *models.QosLite) error {
	in := new(disegrpc.SetQoSOfVolumeRequest)
	in.VolumeName = pv
	in.RWQoS = qos.TG()
	out := new(disegrpc.SetQoSOfVolumeResponse)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return ds.conn.Invoke(ctx, "/disegrpc.DISE/SetQoSOfVolume", in, out)
}

func NewDiseService(engine *xorm.Engine) (DiseService, *grpc.ClientConn, error) {
	c, err := initDiseService(engine)
	if err != nil {
		return nil, nil, err
	}
	return &diseService{conn: c}, c, nil
}

func initDiseService(engine *xorm.Engine) (*grpc.ClientConn, error) {
	address := models.Sysparameter{ParamKey: "dise_grpc"}
	exist, err := engine.Cols("param_value").Get(&address)
	addressList := strings.Split(address.ParamValue, ",")
	if !exist || len(addressList) == 0 {
		return nil, fmt.Errorf("not found dise_grpc address, error: %v", err)
	}
	var conn *grpc.ClientConn
	for i := range addressList {
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn, err = grpc.DialContext(timeout, addressList[i], grpc.WithInsecure(), grpc.WithBlock())
		if err == nil && conn != nil {
			cancel()
			break
		}
		cancel()
		utils.LoggerError(fmt.Errorf("dise grpc connect %v fail, error: %v", addressList[i], err))
	}
	return conn, err
}
