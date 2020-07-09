package admin

import (
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
	"tools/server/app/common"
	"tools/server/app/response"
	"tools/server/app/service"
)

//任务管理
type TaskController struct {
	BaseController
}

// @Title 任务列表
// @Description 获取任务列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        query   int false         "当前页数"
// @Param   limit       query   int false         "查询条数"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /group_list [get]
func (c *TaskController) GroupList() {
	var (
		list   []map[string]interface{}
		result map[string]interface{}
	)
	result = make(map[string]interface{}, 1)
	for _, value := range service.TaskGroupList {
		list = append(list, value)
	}

	result["list"] = list
	c.ResponseToJson(response.Success, result)
}

// @Title 任务列表
// @Description 获取服务器列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        query   int false         "当前页数"
// @Param   limit       query   int false         "查询条数"
// @Param   audit    	query   int false         "审核状态"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /list [get]
func (c *TaskController) List() {
	var (
		params *service.TaskParams //查询信息
		err    error
		data   *service.TaskListData
	)

	page, _ := c.GetInt("page", 1)
	limit, _ := c.GetInt("limit", 20)
	audit, _ := c.GetInt8("audit", 20)

	params = &service.TaskParams{
		Page:  page,
		Limit: limit,
		Audit: audit,
	}

	//获取管理员列表
	if err, data = service.TaskList(params); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 创建任务
// @Description 创建任务
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   Task_name  body   string true       "路由名称"
// @Param   Task   	body   string true       "路由url"
// @Param   request   	body   string false      "请求类型"
// @Param   parent_id   body   uint true         "上级路由id"
// @Param   status      body   int8 true         "状态 0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [put]
func (c *TaskController) Put() {
	var (
		Task *service.TaskDetail
		err  error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &Task)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	valid := validation.Validation{}
	b, errValid := valid.Valid(Task)
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

	if err = service.TaskCreateOrUpdate(Task, c.userId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 修改任务
// @Description 修改任务
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   Task_id    body   uint  true         "路由id"
// @Param   Task_name  body   string false       "路由角色"
// @Param   Task   	body   string false       "路由url"
// @Param   request   	body   string false       "请求类型"
// @Param   parent_id   body   uint true          "上级路由id"
// @Param   status      body   int8 false         "状态 0.无效 1.有效"
// @Param   is_delete   body   int8 false         "删除状态  0.正常 1.删除"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [post]
func (c *TaskController) Post() {
	var (
		Task *service.TaskDetail
		err  error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &Task)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	valid := validation.Validation{}
	b, errValid := valid.Valid(Task)
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

	if err = service.TaskCreateOrUpdate(Task, c.userId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 获取详细信息
// @Description
// @Param   task_id      body   int false    "任务id"
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /detail [get]
func (c *TaskController) Detail() {
	var (
		data *service.TaskDetail
		err  error
	)
	taskId, _ := c.GetInt("task_id")
	if taskId < 1 {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.TaskDetailById(taskId); err != nil {
		c.ResponseToJson(response.DbQueryErr, nil)
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 审核任务
// @Description 审核任务
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   task_id   body   int  int       "任务id"
// @Param   audit     body   int8 int         "审核状态  1.通过 2.不通过"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /audit [post]
func (c *TaskController) Audit() {
	var (
		TaskOperation service.TaskOperation
		taskId        int
		err           error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &TaskOperation)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	checkArray := make([]interface{}, 2)
	checkArray[0], checkArray[1] = int8(1), int8(2)

	if !common.InArray(TaskOperation.Audit, checkArray) {
		c.ResponseToJson(response.ParamsErr, nil, "audit不在范围内")
		return
	}
	if TaskOperation.TaskIds == "" {
		c.ResponseToJson(response.ParamsErr, nil, "缺少task_id")
		return
	}

	taskIdArr := strings.Split(TaskOperation.TaskIds, ",")
	for _, val := range taskIdArr {
		taskId, _ = strconv.Atoi(val)
		if err := service.TaskAudit(taskId, TaskOperation.Audit); err != nil {
			c.ResponseToJson(response.DevOpsTaskAuditErr, nil, err.Error())
			return
		}
	}
	c.ResponseToJson(response.Success, nil)
}

// @Title 任务 执行一次
// @Description 任务执行
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   task_id   body   int  int       "任务id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /execute [get]
func (c *TaskController) ExecuteTask() {
	taskId, _ := c.GetInt("task_id")
	if taskId == 0 {
		c.ResponseToJson(response.ParamsErr, nil, "缺少task_id")
		return
	}

	if err := service.ExecuteTask(taskId); err != nil {
		c.ResponseToJson(response.DevOpsTaskRunErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, nil)
}

// @Title 任务 启动-停止
// @Description 审核任务
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   task_id   body   int  int       "任务id"
// @Param   status    body   int8 int       "任务状态  0.无效 1.有效"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /running [post]
func (c *TaskController) TaskRunning() {
	var (
		TaskOperation service.TaskOperation
		taskId        int
		err           error
	)

	form := c.Ctx.Input.RequestBody
	err = json.Unmarshal(form, &TaskOperation)
	if err != nil {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	checkArray := make([]interface{}, 2)
	checkArray[0], checkArray[1] = int8(0), int8(1)

	if !common.InArray(TaskOperation.Status, checkArray) {
		c.ResponseToJson(response.ParamsErr, nil, "status不在范围内")
		return
	}
	if TaskOperation.TaskIds == "" {
		c.ResponseToJson(response.ParamsErr, nil, "缺少task_id")
		return
	}

	taskIdArr := strings.Split(TaskOperation.TaskIds, ",")
	for _, val := range taskIdArr {
		taskId, _ = strconv.Atoi(val)
		if err := service.SetTaskState(taskId, TaskOperation.Status); err != nil {
			c.ResponseToJson(response.DevOpsTaskStartErr, nil, err.Error())
			return
		}
	}
	c.ResponseToJson(response.Success, nil)
}

// @Title 删除任务
// @Description 删除任务
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   task_id   body   int  int       "任务id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router / [delete]
func (c *TaskController) Delete() {
	taskId, _ := c.GetInt("task_id")
	if taskId == 0 {
		c.ResponseToJson(response.ParamsErr, nil, "缺少task_id")
		return
	}

	if err := service.TaskDelete(taskId); err != nil {
		c.ResponseToJson(response.DbReadErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, nil)
}

// @Title 任务日志
// @Description 获取任务日志列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        query   int false         "当前页数"
// @Param   limit       query   int false         "查询条数"
// @Param   audit    	query   int false         "审核状态"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /log [get]
func (c *TaskController) Log() {
	var (
		params *service.TaskLogParams //查询信息
		err    error
		data   *service.TaskLogListData
	)

	page, _ := c.GetInt("page", 1)
	limit, _ := c.GetInt("limit", 20)
	TaskId, _ := c.GetInt("task_id")

	if TaskId <= 0 {
		c.ResponseToJson(response.ParamsErr, nil, "task_id必传")
		return
	}

	params = &service.TaskLogParams{
		Page:   page,
		Limit:  limit,
		TaskId: TaskId,
	}

	if err, data = service.TaskLogList(params); err != nil {
		c.ResponseToJson(response.MgoQueryErr, nil, err.Error())
		return
	}

	c.ResponseToJson(response.Success, data)
}

// @Title 任务日志详细
// @Description 获取任务日志详细
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   page        query   int false         "当前页数"
// @Param   limit       query   int false         "查询条数"
// @Param   audit    	query   int false         "审核状态"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /log_detail [get]
func (c *TaskController) LogDetail() {
	var (
		data *service.TaskLogData
		err  error
	)
	logId := c.GetString("log_id")
	if logId == "" {
		c.ResponseToJson(response.ParamsErr, data)
		return
	}

	if data, err = service.TaskLogDetailById(logId); err != nil {
		c.ResponseToJson(response.MgoQueryErr, err.Error())
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 获取禁用命令列表
// @Description 获取禁用命令列表
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Failure 400 no enough input
// @Failure 500 server error
// @router /ban_list [get]
func (c *TaskController) BanList() {
	var (
		err  error
		data interface{}
	)

	if data, err = service.BanList(); err != nil && (err.Error() != "mongo: no documents in result") {
		c.ResponseToJson(response.MgoQueryErr, nil, err.Error())
		return
	}
	c.ResponseToJson(response.Success, data)
}

// @Title 编辑禁用命令
// @Description 编辑禁用命令
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   ban_list   body  array  true   "禁用命令"
// @Param   id   	   body  string false   "id"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /edit_ban [post]
func (c *TaskController) EditBan() {
	var err error

	data := c.GetStrings("ban_list[]")
	id := c.GetString("id", "")
	if len(data) > 0 {
		if err = service.EditBan(data, id); err != nil {
			if err.Error() == "hasChinese" {
				c.ResponseToJson(response.DevBanHasChinese, nil, err.Error())
			} else {
				c.ResponseToJson(response.DbReadErr, nil, err.Error())
			}
			return
		}
	} else {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 是否为禁用命令
// @Description 是否为禁用命令
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   ban_order   body  string  true   "禁用命令"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /is_ban [get]
func (c *TaskController) IsBan() {
	var err error

	order := c.GetString("ban_order", "")
	if order != "" {
		if err = service.IsBan(order); err != nil && (err.Error() != "mongo: no documents in result") {
			if err.Error() == "isBan" {
				c.ResponseToJson(response.DevIsBan, nil)
			} else {
				c.ResponseToJson(response.MgoQueryErr, nil, err.Error())
			}
			return
		}
	} else {
		c.ResponseToJson(response.ParamsErr, nil)
		return
	}

	c.ResponseToJson(response.Success, nil)
}

// @Title 通知数据集合
// @Description 通知所需要的数据信息
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   expand   body  string  false   "扩展信息、 admin,notify_type,notify_tpl"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /notify_data [get]
func (c *TaskController) NotifyData() {
	expand := c.GetString("expand", "admin,notify_type,notify_tpl")
	data := service.GetNotifyData(expand)
	c.ResponseToJson(response.Success, data)
}
