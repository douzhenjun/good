/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: Dou
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: Dou
 * @LastEditTime: 2021-02-07 16:32:07
 */

package datasource

import (
	"DBaas/config"
	"DBaas/models"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq" // pg数据库 导入这个
	"os"
	"xorm.io/core"
)

/**
 * 实例化数据库引擎方法：pg的数据引擎
 */
func NewPgEngine(c *config.AppConfig) *xorm.Engine {
	database := c.DataBase[0]
	var pgPassword string
	pgPassword = os.Getenv("pg_passwd")
	if pgPassword == "" {
		pgPassword = database.Pwd
	}
	dataSourceName := "postgres://" + database.User + ":" + pgPassword + "@" + database.Host + ":" + database.Port + "/" + database.Database + "?sslmode=disable"
	engine, err := xorm.NewEngine(database.Drive, dataSourceName)
	if err != nil {
		panic(err)
	}
	//同步数据库结构：主要负责对数据结构实体同步更新到数据库表
	/**
	 * 自动检测和创建表，这个检测是根据表的名字
	 * 自动检测和新增表中的字段，这个检测是根据字段名，同时对表中多余的字段给出警告信息
	 * 自动检测，创建和删除索引和唯一索引，这个检测是根据索引的一个或多个字段名，而不根据索引名称。因此这里需要注意，如果在一个有大量数据的表中引入新的索引，数据库可能需要一定的时间来建立索引。
	 * 自动转换varchar字段类型到text字段类型，自动警告其它字段类型在模型和数据库之间不一致的情况。
	 * 自动警告字段的默认值，是否为空信息在模型和数据库之间不匹配的情况
	 */
	//Sync2是Sync的基础上优化的方法
	engine.SetMapper(core.GonicMapper{})
	err = engine.Sync2(
		new(models.StatisticsCluster),
	)
	if err != nil {
		panic(err)
	}
	//设置显示sql语句
	engine.ShowSQL(false)
	engine.SetMaxOpenConns(100)
	return engine
}
