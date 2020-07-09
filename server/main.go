package main

import (
	"github.com/astaxie/beego/plugins/cors"
	"os"
	"tools/server/app/service"
	_ "tools/server/routers"

	"github.com/astaxie/beego"
)

func init() {
	os.Setenv("ZONEINFO", "./zoneinfo.zip")
	service.Init()
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AllowHeaders:     []string{"Origin", "X-Requested-With", "Token", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}
	beego.Run()
}

func main() {
	beego.Run()
}
