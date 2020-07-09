// @APIVersion 1.0.0
// @Title Tools API
// @Description 后台api
// @Contact 645536551@qq.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"tools/server/app/controllers/admin"
)

func init() {
	var ns *beego.Namespace
	ns = beego.NewNamespace("auth",
		beego.NSNamespace("admin",
			beego.NSInclude(
				&admin.AdminController{},
			),
		),
		beego.NSNamespace("route",
			beego.NSInclude(
				&admin.RouteController{},
			),
		),
		beego.NSNamespace("role",
			beego.NSInclude(
				&admin.RoleController{},
			),
		),
	)
	beego.AddNamespace(ns)

	ns = beego.NewNamespace("/devops",
		beego.NSNamespace("/servers",
			beego.NSInclude(
				&admin.ServersController{},
			),
		),
		beego.NSNamespace("/task",
			beego.NSInclude(
				&admin.TaskController{},
			),
		),
		beego.NSNamespace("/cron",
			beego.NSInclude(
				&admin.CronController{},
			),
		),
		beego.NSNamespace("/notify",
			beego.NSInclude(
				&admin.NotifyController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
