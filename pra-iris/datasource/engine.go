package datasource

import (
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"pra-iris/configs"
	"pra-iris/models"
	"xorm.io/core"
)

func NewPgEngine(c *configs.AppConfig) *xorm.Engine {
	database := c.DataBase[0]
	dataSourceName := "postgres://" + database.User + ":" + database.Pwd + "@" + database.Host + ":" + database.Port + "/" + database.Database + "?sslmode=disable"
	engine, err := xorm.NewEngine(database.Drive, dataSourceName)
	if err != nil {
		panic(err)
	}
	engine.SetMapper(core.GonicMapper{})
	// 默认创建表
	err = engine.Sync2(
		new(models.User),
	)
	if err != nil {
		panic(err)
	}
	return engine
}