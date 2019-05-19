package services

import (
	"github.com/go-xorm/xorm"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var DBEngine *xorm.Engine

func init() {
	driverName := "mysql"
	dataSourceName := "helocard:helocardpwd@(127.0.0.1:3306)/go_chat?charset=utf8"
	DBEngine, err := xorm.NewEngine(driverName, dataSourceName)
	if nil != err {
		log.Fatal(err.Error())
	}
	// 是否显示 SQL 语句
	DBEngine.ShowSQL(true)
	// 数据库最大打开的连接数
	DBEngine.SetMaxOpenConns(2)

	// 自动建表
	//DBEngine.Sync2(new(User))

	fmt.Println("init data base ok")
}
