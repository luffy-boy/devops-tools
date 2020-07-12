package admin

import (
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"strings"
	"time"
	"tools/server/app/common"
	"tools/server/app/response"
	"tools/server/app/service"
	"tools/server/app/util"
)

/**
后台基类
*/
type BaseController struct {
	beego.Controller
	//当前环境
	ENV    string
	userId int
}

//过滤不校验的Token的路由
var skipRouter = []string{
	"auth/admin/login",
	"auth/admin/logout",
	"auth/admin/refresh_login",
}

func (c *BaseController) Prepare() {
	c.ENV = beego.AppConfig.String("runmode")
	c.CheckLogin()
}

// @Title 测试接口
// @Description 测试数据
// @Success 200
// @router / [get]
func (c *BaseController) Get() {
	fmt.Println(c.userId)
	str := "time:" + strconv.FormatInt(time.Now().Unix(), 10)
	c.Ctx.WriteString(str)
}

/**
校验登录
*/
func (c *BaseController) CheckLogin() {
	var (
		routes  = strings.Trim(strings.ToLower(c.Ctx.Request.URL.Path), "/")
		request = strings.ToUpper(c.Ctx.Request.Method)
		routeId int
		admin   *util.AdminJwtData
		err     error
	)

	arr := make([]interface{}, len(skipRouter))
	for k, v := range skipRouter {
		arr[k] = v
	}
	admin = &util.AdminJwtData{}
	token := c.Ctx.Request.Header.Get("Token")
	if token != "" {
		//校验token
		admin, err = util.ValidateToken(token)
		if err != nil || admin == nil {
			c.ResponseToJson(response.TokenInvalidErr, nil)
			return
		}
		c.userId = admin.Id
	}

	if !common.InArray(routes, arr[:]) {
		if admin.Id < 1 {
			c.ResponseToJson(response.TokenInvalidErr, nil)
			return
		}

		userKey := util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(admin.Id)
		if !util.RedisObj.Exists(userKey) {
			c.ResponseToJson(response.TokenExpiredErr, nil)
			return
		}

		//重新刷新token时间
		refTokenExpired, _ := beego.AppConfig.Int64("redis.refTokenExpired")
		tokenExpired, _ := beego.AppConfig.Int("redis.tokenExpired")
		if expired := util.RedisObj.Ttl(util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(admin.Id)); expired > 0 && expired <= refTokenExpired {
			util.RefreshToken(token, tokenExpired)
			util.RedisObj.Expire(userKey, tokenExpired)
		}

		if admin.Id == 1 {
			return
		}

		//校验路由
		RouteList, err := service.GetRouteData()
		if err != nil {
			c.ResponseToJson(response.UserAccessForbiddenErr, nil)
			return
		}

		for _, v := range RouteList {
			if v.Route == routes {
				if v.Request == "" {
					routeId = v.Id
					break
				} else if request == v.Request {
					routeId = v.Id
					break
				}
			}
		}

		if routeId == 0 {
			c.ResponseToJson(response.UserAccessForbiddenErr, nil)
			return
		}

		routeList := service.GetUserRouteById(admin.Id)
		if len(routeList) == 0 {
			c.ResponseToJson(response.UserAccessForbiddenErr, nil)
			return
		}
		if !common.InArray(routeId, routeList[:]) {
			c.ResponseToJson(response.UserAccessForbiddenErr, nil)
			return
		}
	}
}

//返回结果  json
func (c *BaseController) ResponseToJson(httpCode int, data interface{}, info ...string) {
	rcode := response.ResponseList[httpCode]["code"].(int)
	rinfo := response.ResponseList[httpCode]["info"].(string)
	if rcode == 1 {
		rinfo = "ok"
	}
	if len(info) > 0 && info[0] != "" {
		rinfo = info[0]
	}
	result := response.Json{Code: rcode, Data: data, Message: rinfo}
	c.Data["json"] = &result
	c.ServeJSON()
}
