package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"tools/server/app/response"
	"tools/server/app/service"
)

//资源服务器管理
type ServersController struct {
	BaseController
}

// @Title 资源服务器列表
// @Description 获取服务器列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        body   int false         "当前页数"
// @Param   limit       body   int false         "查询条数"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /group_list [get]
func (c *ServersController)GroupList() {
	var (
		list  []map[string]interface{}
		result  map[string]interface{}
	)
	result = make(map[string]interface{},1)
	for _,value := range service.ServersGroupList{
		list = append(list,value)
	}

	result["list"] = list
	c.ResponseToJson(response.Success, result)
}

// @Title 资源服务器列表
// @Description 获取服务器列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        body   int false         "当前页数"
// @Param   limit       body   int false         "查询条数"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *ServersController) List() {
	var (
		params *service.ServersParams //查询信息
		err    error
		data   *service.ServersListData
	)

	page,_   := c.GetInt("page", 1)
	limit,_  := c.GetInt("limit", 20)
	extend   := c.GetString("extend")

	params = &service.ServersParams{
		Page:  page,
		Limit: limit,
		Extend: extend,
	}
	//解析参数
	err = c.ParseForm(params)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//获取管理员列表
	if err, data = service.ServersList(params); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 创建资源服务器由
// @Description 创建资源服务器
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   Servers_name  body   string true       "路由名称"
// @Param   Servers   	body   string true       "路由url"
// @Param   request   	body   string false      "请求类型"
// @Param   parent_id   body   uint true         "上级路由id"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *ServersController) Put() {
	var (
		Servers *service.ServersDetail //路由信息
		err     error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &Servers)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	valid := validation.Validation{}
	b, errValid := valid.Valid(Servers)
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

	if err = service.ServersCreateOrUpdate(Servers); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改资源服务器
// @Description 修改资源服务器
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   Servers_id    body   int  true         "路由id"
// @Param   Servers_name  body   string false       "路由角色"
// @Param   Servers   	body   string false       "路由url"
// @Param   request   	body   string false       "请求类型"
// @Param   parent_id   body   int true          "上级路由id"
// @Param   status      body   int8 false         "状态 0.无效 1.有效"
// @Param   is_delete   body   int8 false         "删除状态  0.正常 1.删除"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *ServersController) Post() {
	var (
		Servers *service.ServersDetail
		err     error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &Servers)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	b, errValid := valid.Valid(&Servers)
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
	if err = service.ServersCreateOrUpdate(Servers); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 获取详细信息
// @Description
// @Param   server_id      body   int false    "服务器id"
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *ServersController)Detail(){
	var (
		data *service.ServersDetail
		err error
	)
	serverId,_ := c.GetInt("server_id")
	if serverId < 1{
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if  data,err = service.ServersDetailById(serverId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 删除路由
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   route_id    body   int  true         "路由id"
// @router / [delete]
func (c *ServersController) Delete() {
	serverId,_ := c.GetInt("server_id")

	if serverId < 1 {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//创建或修改路由信息
	if err := service.ServersDelete(serverId); err != nil {
		c.ResponseToJson(response.DbDeleteErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}