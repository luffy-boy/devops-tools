package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"tools/server/app/response"
	"tools/server/app/service"
)

type NotifyController struct {
	BaseController
}

// @Title 通知模板
// @Description 获取通知模板列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        query   int false         "当前页数"
// @Param   limit       query   int false         "查询条数"
// @Param   audit    	query   int false         "审核状态 -1查询所有"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *NotifyController) List() {
	var (
		params *service.NotifyTplParams //查询信息
		err    error
		data   service.NotifyTplListData
	)

	page, _ := c.GetInt("page", 1)
	limit, _ := c.GetInt("limit", 20)

	params = &service.NotifyTplParams{
		Page:  page,
		Limit: limit,
	}

	if err, data = service.NotifyTplList(params); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 创建通知模板
// @Description 创建通知模板
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   tpl_name  	body   string true       "通知模板名称"
// @Param   tpl_data   	body   string true       "通知模板内容"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *NotifyController) Put() {
	var (
		data *service.NotifyTplDetail
		err  error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &data)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	valid := validation.Validation{}
	b, errValid := valid.Valid(data)
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

	if err = service.NotifyCreateOrUpdate(data, c.userId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改通知模板
// @Description 修改通知模板
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   tpl_name  	body   string true       "通知模板名称"
// @Param   tpl_data   	body   string true       "通知模板内容"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *NotifyController) Post() {
	var (
		data *service.NotifyTplDetail
		err  error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &data)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	valid := validation.Validation{}
	b, errValid := valid.Valid(data)
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

	if err = service.NotifyCreateOrUpdate(data, c.userId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 获取通知模板详细
// @Description
// @Param   id      body   int false    "通知模板id"
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *NotifyController) Detail() {
	var (
		data service.NotifyTplDetail
		err  error
	)
	id, _ := c.GetInt("id")
	if id < 1 {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.NotifyDetailById(id); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 删除通知模板
// @Description 删除通知模板
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   notify_tpl_id   body   int  int       "通知模板id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [delete]
func (c *NotifyController) Delete() {
	notifyTplId, _ := c.GetInt("notify_tpl_id")
	if notifyTplId == 0 {
		c.ResponseToJson(response.ParamsErr, nil, "缺少task_id")
		return
	}

	if err := service.NotifyDelete(notifyTplId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, nil)
}
