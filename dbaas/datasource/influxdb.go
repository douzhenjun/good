package datasource

import (
	"DBaas/config"
	"github.com/influxdata/influxdb/client/v2"
)

func NewInfluxdbEngine(c *config.AppConfig) (client.Client, error) {
	database := c.DataBase[1]
	conn, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + database.Host + ":" + database.Port,
		Username: database.User,
		Password: database.Pwd,
	})
	if err != nil {
		return nil, err
	}
	return conn, err
}

func NewBatchPoints() (client.BatchPoints, error) {
	dbConf := config.GetConfig().DataBase[1]
	return client.NewBatchPoints(client.BatchPointsConfig{
		Database:  dbConf.Database,
		Precision: "s",
	})
}
