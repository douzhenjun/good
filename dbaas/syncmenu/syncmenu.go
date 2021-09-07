package main

import (
	"DBaas/config"
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"xorm.io/core"
)

func pgEngine() *xorm.Engine {
	c := config.GetConfig()
	db := c.DataBase[0]
	pwd := os.Getenv("pg_passwd")
	if len(pwd) == 0 {
		pwd = db.Pwd
	}
	dbSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", db.User, pwd, db.Host, db.Port, db.Database)
	engine, err := xorm.NewEngine(db.Drive, dbSource)
	if err != nil {
		panic(err)
	}
	engine.SetMapper(core.GonicMapper{})
	engine.SetMaxOpenConns(10)
	return engine
}

func main() {
	menu := &menu{pgEngine()}
	mj := make(map[string]interface{})
	jf, err := ioutil.ReadFile("support-files/scripts/dbaas_menu_data.json")
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(jf, &mj)
	if err != nil {
		log.Fatalln(err)
	}
	// check menu
	sm := getSubArr(mj, "subMenu")
	smDB :=
	log.Println(mj)
}

func getSubArr(m map[string]interface{}, sub string) []map[string]interface{} {
	return m[sub].([]map[string]interface{})
}

type menu struct {
	engine *xorm.Engine
}

func (m *menu) getMenu() {

}
