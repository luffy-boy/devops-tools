package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"tools/server/app/response"
	"tools/server/app/service"
)

//后台角色
type RoleController struct {
	BaseController
}

// @Title 获取全部角色列表
// @Description 获取全部角色列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /all [get]
func (c *RoleController) All() {
	var (
		err  error
		data []*service.RoleData
	)
	if data, err = service.RoleAll(); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}
	result := make(map[string]interface{})
	result["list"] = data
	c.ResponseToJson(response.Success, result)
}

// @Title 获取角色列表
// @Description 获取角色列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        body   uint false         "当前页数"
// @Param   limit       body   uint false         "查询条数"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *RoleController) List() {
	var (
		role *service.RoleParams //查询信息
		err  error
		data *service.RoleListData
	)

	page, _ := c.GetInt("page", 1)
	limit, _ := c.GetInt("limit", 20)

	role = &service.RoleParams{
		Page:  page,
		Limit: limit,
	}
	//解析参数
	err = c.ParseForm(role)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//获取管理员列表
	if data, err = service.RoleList(role); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 创建角色
// @Description 创建后台角色
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   role_name   body   string true       "角色名称"
// @Param   parent_id   body   uint true         "上级角色id"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Param   route_ids   body   string false      "扩展路由id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *RoleController) Put() {
	var (
		role *service.RoleRequest //角色信息
		err  error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &role)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	b, errValid := valid.Valid(role)
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

	//创建或修改管理员信息
	if err = service.RoleCreateOrUpdate(role); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改角色
// @Description 修改后台角色
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   role_id     body   uint  true         "角色id"
// @Param   role_name   body   string false       "角色名称"
// @Param   parent_id   body   uint false         "上级id"
// @Param   status      body   int8 false         "状态 0.无效 1.有效"
// @Param   is_delete   body   int8 false         "删除状态  0.正常 1.删除"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *RoleController) Post() {
	var (
		role *service.RoleRequest
		err  error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &role)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	b, errValid := valid.Valid(role)
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
	if err = service.RoleCreateOrUpdate(role); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 获取详细信息
// @Description
// @Param   role_id      body   int false    "角色id"
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *RoleController) Detail() {
	var (
		data *service.RoleDetail
		err  error
	)
	roleId, _ := c.GetInt("role_id")
	if roleId < 1 {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.RoleDetailById(roleId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 删除角色
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   role_id    body   int  true         "角色id"
// @router / [delete]
func (c *RoleController) Delete() {
	roleId, _ := c.GetInt("role_id")

	if roleId < 1 {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//删除角色信息
	if err := service.RoleDelete(roleId); err != nil {
		if err.Error() == "hasChild" {
			c.ResponseToJson(response.RoleHasChild, nil)
		} else {
			c.ResponseToJson(response.DbDeleteErr, nil, err.Error())
		}

		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 角色权限编辑
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   role_id    body   int     true         "角色id"
// @Param   route_ids  body   string  true      "路由ids"
// @router /route_edit [post]
func (c *RoleController) RouteEdit() {
	var (
		params  *service.RouteEditParams
		resCode int
		err     error
	)

	roleId, _ := c.GetInt("role_id", 0)
	routeIds  := c.GetString("route_ids", "")

	params = &service.RouteEditParams{
		RoleId:   roleId,
		RouteIds: routeIds,
	}
	//解析参数
	err = c.ParseForm(params)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	if params.RoleId < 1 {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//角色权限
	resCode = service.RouteEdit(params)
	c.ResponseToJson(resCode, nil)
}
