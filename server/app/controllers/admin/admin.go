package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"tools/server/app/response"
	"tools/server/app/service"
	"tools/server/app/util"
)

//管理员
type AdminController struct {
	BaseController
}

// @Title 登录
// @Description 用户登录
// @Success 200 {object} service.AdminLoginData Res {"code":1,"data":null,"msg":"ok"}
// @Param   username    body   string true       "用户名"
// @Param   password    body   string true       "密码"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /login [post]
func (c *AdminController) Login() {
	var (
		admin    *service.AdminRequest   //登录信息
		data     *service.AdminLoginData //返回登录信息
		httpCode int                     //返回登录信息
		err      error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &admin)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	userLoginValid := service.AdminLoginValid{UserName: admin.UserName, Password: admin.PassWord}
	b, errValid := valid.Valid(&userLoginValid)
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
	//获取登录管理员信息
	admin.Ip = c.Ctx.Input.IP()
	if data, httpCode, err = service.AdminLogin(admin); err != nil {
		c.ResponseToJson(httpCode, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 退出登录
// @Description
// @Success 200 {object} service.AdminLoginData Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /logout [post]
func (c *AdminController) Logout() {
	token := c.Ctx.Request.Header.Get("Token")
	if token == "" {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}
	//校验token
	admin, err := util.ValidateToken(token)
	if err != nil || admin == nil {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}
	if admin.Id < 1 {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}

	if ok := service.DelAdminLoginData(admin.Id); !ok {
		c.ResponseToJson(response.LogoutError, nil)
		return
	}

	c.ResponseToJson(response.Success, nil)
	return
}

// @Title 重新登录
// @Description
// @Success 200  Res {"code":1,"data":service.AdminLoginData ,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /re_login [post]
func (c *AdminController) ReLogin() {
	var (
		data     *service.AdminLoginData //返回登录信息
		httpCode int
		err      error
	)
	token := c.Ctx.Request.Header.Get("Token")
	if token == "" {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}
	//校验token
	admin, err := util.ValidateToken(token)
	if err != nil || admin == nil {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}
	if admin.Id < 1 {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}

	if data, httpCode, err = service.AdminReLogin(admin.Id, c.Ctx.Input.IP(), token); err != nil {
		c.ResponseToJson(httpCode, nil)
		return
	}

	c.ResponseToJson(response.Success, data)
	return
}

// @Title 获取管理员列表
// @Description 获取管理员列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   user_id     body   uint true         "管理员id"
// @Param   page        body   uint false         "当前页数"
// @Param   limit       body   uint false         "查询条数"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *AdminController) List() {
	var (
		admin *service.AdminParams //查询信息
		err   error
		data  *service.AdminListData
	)

	admin = &service.AdminParams{}
	admin.UserId, _ = c.GetInt("user_id")
	admin.Page, _ = c.GetInt("page", 1)
	admin.Limit, _ = c.GetInt("limit", 20)
	admin.RoleId, _ = c.GetInt("role_id")
	admin.RealName = c.GetString("real_name", "")
	admin.Sort = c.GetString("sort", "")

	//获取管理员列表
	if err, data = service.AdminList(admin); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 获取管理员信息
// @Description
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [get]
func (c *AdminController) Get() {
	if c.userId < 1 {
		c.ResponseToJson(response.TokenInvalidErr, nil)
		return
	}
	var (
		role string
	)
	userData := make(map[string]interface{})
	data := service.GetAdminLoginData(c.userId)
	if c.userId == 1 {
		role = "Administrators"
	} else {
		role = data.Role
	}
	userData["real_name"] = data.RealName
	userData["introduction"] = data.Introduction
	userData["sex"] = data.Sex
	userData["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	userData["roles"] = role

	c.ResponseToJson(response.Success, userData)
}

// @Title 获取详细信息
// @Description
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *AdminController) Detail() {
	var (
		data service.AdminDetail
		err  error
	)
	userId, _ := c.GetInt("user_id")
	if userId < 1 {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.AdminDetailById(userId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 创建管理员账号
// @Description 创建管理员账号
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   username    	body   string true       "管理员账号"
// @Param   password    	body   string true       "管理员密码"
// @Param   realName    	body   string true       "真实姓名"
// @Param   Introduction    body   string false      "个人描述"
// @Param   phone       body   string true       "手机号码"
// @Param   sex         body   uint8 true        "性别"
// @Param   role_id     body   uint true         "角色id"
// @Param   birthday    body   int false         "生日"
// @Param   w   body   string false      "扩展路由id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *AdminController) Put() {
	var (
		admin *service.AdminRequest //登录信息
		err   error
	)

	//解析参数
	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &admin)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	//参数校验
	valid := validation.Validation{}
	userLoginValid := service.AdminDataValid{
		UserName: admin.UserName,
		Password: admin.PassWord,
		RealName: admin.RealName,
		Birthday: admin.Birthday,
		Phone:    admin.Phone,
		Email:    admin.Email,
		Sex:      admin.Sex,
		RoleId:   admin.RoleId,
		RouteIds: admin.RouteIds,
		Status:   admin.Status,
	}
	b, errValid := valid.Valid(&userLoginValid)
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
	if err = service.AdminCreateOrUpdate(admin); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改管理员账号
// @Description 修改管理员账号
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   user_id     body   uint   true       "用户id"
// @Param   password    body   string false       "管理员密码"
// @Param   real_name   body   string false       "真实姓名"
// @Param   phone       body   string false       "手机号码"
// @Param   sex         body   uint8 false        "性别"
// @Param   role_id     body   uint false         "角色id"
// @Param   birthday    body   string false         "生日 Y-m-d"
// @Param   route_ids   body   string false      "扩展路由id"
// @Param   is_delete   body   int8 false      "扩展路由id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *AdminController) Post() {
	var (
		admin *service.AdminRequest //登录信息
		err   error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &admin)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}
	admin.UserName = ""
	//参数校验
	valid := validation.Validation{}
	userRegValid := service.AdminUpdateDataValid{
		UserId:   admin.UserId,
		RealName: admin.RealName,
		Birthday: admin.Birthday,
		Phone:    admin.Phone,
		Email:    admin.Email,
		Sex:      admin.Sex,
		RoleId:   admin.RoleId,
		RouteIds: admin.RouteIds,
	}
	b, errValid := valid.Valid(&userRegValid)
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
	if err = service.AdminCreateOrUpdate(admin); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 获取首页展示信息
// @Description 获取首页展示信息
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /index_data [get]
func (c *AdminController) IndexData() {
	c.ResponseToJson(response.Success, nil)
}
