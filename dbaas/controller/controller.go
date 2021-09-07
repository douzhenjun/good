package controller

import (
	"DBaas/models"
	"DBaas/utils"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
)

func ReadQos(ctx iris.Context) (*models.Qos, error) {
	qosStr := ctx.PostValue("qos")
	qos := new(models.Qos)
	var err error
	if len(qosStr) != 0 {
		err = json.Unmarshal(utils.Str2bytes(qosStr), &qos.QosLite)
		if err != nil {
			err = fmt.Errorf("parse qos error: %v", err)
		}
	} else {
		qos = nil
	}
	return qos, err
}
