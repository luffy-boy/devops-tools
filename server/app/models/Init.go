package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
)

var (
	o           orm.Ormer
	tablePrefix string
)

func Init() {
	dbHost := beego.AppConfig.String("db.host")
	dbPort := beego.AppConfig.String("db.port")
	dbUser := beego.AppConfig.String("db.user")
	dbPassword := beego.AppConfig.String("db.password")
	dbName := beego.AppConfig.String("db.name")
	timezone := beego.AppConfig.String("db.timezone")
	tablePrefix = beego.AppConfig.String("db.prefix")

	if dbPort == "" {
		dbPort = "3306"
	}
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8"
	if timezone != "" {
		dsn = dsn + "&loc=" + url.QueryEscape(timezone)
	}
	err := orm.RegisterDataBase("default", "mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}

	orm.RegisterModelWithPrefix(tablePrefix,
		new(Admin),
		new(Routes),
		new(Roles),
		new(Task),
		new(TaskServers),
		new(NotifyTpl),
	)

	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	o = orm.NewOrm()
}
