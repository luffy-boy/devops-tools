package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"tools/server/app/response"
	"tools/server/app/service"
)

//路由
type RouteController struct {
	BaseController
}

// @Title 获取路由列表
// @Description 获取角色列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *RouteController) List() {
	var (
		role *service.RouteParams //查询信息
		err  error
		data *service.RouteListData
	)

	page := 1
	limit := 100000

	role = &service.RouteParams{
		Page:  page,
		Limit: limit,
	}

	if err, data = service.RouteList(role); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 获取全部路由列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /all [get]
func (c *RouteController) All() {
	data := service.RouteAll()

	res := make(map[string]interface{})
	res["list"] = data
	c.ResponseToJson(response.Success, res)
}

// @Title 获取菜单列表 根据管理员角色
// @Description 获取菜单列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /menu [get]
func (c *RouteController) Menu() {
	var (
		err  error
		data *service.Menu
	)

	if data, err = service.MenuList(c.userId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 获取详细信息
// @Description
// @Param   route_id      body   int false    "路由id"
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *RouteController) Detail() {
	var (
		data *service.RouteDetail
		err  error
	)
	routeId, _ := c.GetInt("route_id")
	if routeId < 1 {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.RouteDetailById(routeId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 创建路由
// @Description 创建后台角色
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   route_name  body   string true       "路由名称"
// @Param   route   	body   string true       "路由url"
// @Param   request   	body   string false      "请求类型"
// @Param   parent_id   body   int true         "上级路由id"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *RouteController) Put() {
	var (
		route *service.RouteDetail //路由信息
		err   error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &route)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	b, errValid := valid.Valid(route)
	if errValid != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}
	if !b {
		for _, err := range valid.Errors {
			msg := err.Key + ": " + err.Message
			c.ResponseToJson(response.ParamsErr, nil, msg)
			return
		}
	}

	//创建或修改路由信息
	if err = service.RouteCreateOrUpdate(route); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改路由
// @Description 修改后台路由
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   route_id    body   int  true         "路由id"
// @Param   route_name  body   string false       "路由角色"
// @Param   route   	body   string false       "路由url"
// @Param   request   	body   string false       "请求类型"
// @Param   parent_id   body   int true          "上级路由id"
// @Param   status      body   int8 false         "状态 0.无效 1.有效"
// @Param   is_delete   body   int8 false         "删除状态  0.正常 1.删除"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *RouteController) Post() {
	var (
		route *service.RouteDetail
		err   error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &route)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	b, errValid := valid.Valid(route)
	if errValid != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}
	if !b {
		for _, err := range valid.Errors {
			msg := err.Key + ": " + err.Message
			c.ResponseToJson(response.ParamsErr, nil, msg)
			return
		}
	}

	//创建或修改路由信息
	if err = service.RouteCreateOrUpdate(route); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改路由排序
// @Description 修改路由
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   route_id   body   int  int       "任务id"
// @Param   sort       body   int16 int      "排序值"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /edit_sort [post]
func (c *RouteController) EditSort() {
	var (
		route service.RouteOperation
		err   error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &route)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	if route.RouteId <= 0 {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}
	if route.RouteId <= 0 {
		c.ResponseToJson(response.ParamsErr, nil, "缺少路由id")
		return
	}
	if route.Sort < 0 {
		c.ResponseToJson(response.ParamsErr, nil, "排序值填写错误")
		return
	}

	if err := service.SetRouteSort(route.RouteId, route.Sort); err != nil {
		c.ResponseToJson(response.DevOpsTaskStartErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, nil)
}

// @Title 删除路由
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   route_id    body   int  true         "路由id"
// @router / [delete]
func (c *RouteController) Delete() {
	routeId, _ := c.GetInt("route_id")

	if routeId < 1 {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	resCode := service.RouteDelete(routeId)
	c.ResponseToJson(resCode, nil)
}
