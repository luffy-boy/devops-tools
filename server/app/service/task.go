package service

import (
	"errors"
	"github.com/astaxie/beego/validation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
	"tools/server/app/common"
	"tools/server/app/models"
	"tools/server/app/util"
	"tools/server/app/util/robfig/cron"
)

//查询
type TaskParams struct {
	Page  int
	Limit int
	Audit int8
}

type TaskLogParams struct {
	Page   int
	Limit  int
	TaskId int
}

//查询返回数据
type TaskListData struct {
	Total int64       `json:"total"`
	List  []*TaskData `json:"list"`
}

type TaskData struct {
	Id           int    `json:"task_id"`
	GroupId      int    `json:"group_id"`
	GroupName    string `json:"group_name"`
	RunType      int8   `json:"run_type"`
	TaskName     string `json:"task_name"`
	CronSpec     string `json:"cron_spec"`
	Concurrent   int8   `json:"concurrent"`
	ExecuteTimes int    `json:"execute_times"`
	PrevTime     int    `json:"prev_time,omitempty"`
	NextTime     int    `json:"next_time,omitempty"`
	IsAudit      int8   `json:"is_audit"`
	Status       int8   `json:"status"`
	Ctime        int    `json:"ctime"`
	Utime        int    `json:"utime"`
}

type TaskLogListData struct {
	Total int64          `json:"total"`
	List  []*TaskLogData `json:"list"`
}

type TaskLogData struct {
	Id          string     `json:"id"`
	TaskId      int        `json:"task_id"`
	ServerId    int        `json:"server_id"`
	ServerName  string     `json:"server_name"`
	Output      string     `json:"output"`
	Error       string     `json:"error"`
	Status      int8       `json:"status"`
	ProcessTime int        `json:"process_time"`
	Ctime       int        `json:"ctime"`
	TaskInfo    TaskDetail `json:"task_info"`
}

//Task
type TaskDetail struct {
	TaskId        int    `json:"task_id"`
	GroupId       int    `json:"group_id" valid:"Required;"`
	ServerIds     string `json:"server_ids" valid:"Required;"`
	RunType       int8   `json:"run_type" valid:"Required;"`
	TaskName      string `json:"task_name" valid:"Required;"`
	Description   string `json:"description" valid:"Required;"`
	CronSpec      string `json:"cron_spec" valid:"Required;"`
	Concurrent    int8   `json:"concurrent"`
	Command       string `json:"command" valid:"Required;"`
	Timeout       int    `json:"timeout"`
	IsNotify      int8   `json:"is_notify"`
	NotifyType    int8   `json:"notify_type"`
	NotifyTplId   int    `json:"notify_tpl_id"`
	NotifyUserIds string `json:"notify_user_ids"`
	IsAudit       int8   `json:"is_audit"`
	UpdateId      int    `json:"update_id"`
	CreateId      int    `json:"create_id"`
	Status        int8   `json:"status"`
	Ctime         int    `json:"-"`
	Utime         int    `json:"-"`
}

type TaskOperation struct {
	TaskIds string `json:"task_ids"`
	Status  int8   `json:"status"`
	Audit   int8   `json:"audit"`
}

type BanListFind struct {
	Id      primitive.ObjectID `bson:"_id"`
	BanList []string           `bson:"ban_list"`
	Ctime   int                `bson:"ctime"`
	Utime   int                `bson:"utime"`
}

type BanListData struct {
	Id      string   `json:"_id"`
	BanList []string `json:"ban_list"`
	Ctime   int      `json:"ctime"`
	Utime   int      `json:"utime"`
}

type NotifyData struct {
	AdminList  []AdminBrief       `json:"admin_list,omitempty"`
	NotifyType []NotifyTypeDetail `json:"notify_type,omitempty"`
	NotifyTpl  []NotifyTplBrief   `json:"notify_tpl,omitempty"`
}

func (task *TaskDetail) Valid(v *validation.Validation) {
	_, err := CronParse.Parse(task.CronSpec)
	if err != nil {
		v.SetError("CronSpec", "时间表达式填写错误")
		return
	}

	if task.IsNotify == 1 && task.NotifyTplId == 0 {
		v.SetError("CronSpec", "请选择通知模板")
		return
	}

	if task.TaskId > 0 {
		fields := []string{"id", "status"}
		filters := make([]interface{}, 2)
		filters[0] = "id"
		filters[1] = task.TaskId
		task := &models.Task{}
		err := task.GetDetail(fields, filters)
		if err != nil {
			v.SetError("CronSpec", err.Error())
			return
		}
		if task.Status == 1 {
			v.SetError("CronSpec", "该任务处于运行状态,禁止修改")
			return
		}
	}
}

var TaskGroupList = map[int]map[string]interface{}{
	1: {"id": 1, "group_name": "服务器任务"},
	2: {"id": 2, "group_name": "业务层任务"},
}

func TaskList(params *TaskParams) (error, *TaskListData) {
	var (
		filters  []interface{}
		fields   []string
		sched    cron.Schedule
		nextTime time.Time
	)

	filters = append(filters, "is_delete", 0)
	if params.Audit == 1 {
		filters = append(filters, "is_audit", 1)
	} else if params.Audit == 2 {
		isAudit := []int{0, 2}
		filters = append(filters, "is_audit__in", isAudit)
	}
	fields = []string{"id", "group_id", "task_name", "run_type", "cron_spec", "concurrent", "execute_times", "prev_time", "is_audit", "status", "is_audit", "ctime", "utime"}
	list, total := models.GetTaskList(params.Page, params.Limit, filters, fields)

	data := &TaskListData{
		Total: 0,
		List:  []*TaskData{},
	}

	cTime := time.Now()
	for _, v := range list {
		groupName, _ := TaskGroupList[v.GroupId]["group_name"]
		sched, _ = CronParse.Parse(v.CronSpec)
		nextTime = sched.Next(cTime)

		task := &TaskData{
			Id:           v.Id,
			GroupId:      v.GroupId,
			GroupName:    groupName.(string),
			Concurrent:   v.Concurrent,
			RunType:      v.RunType,
			TaskName:     v.TaskName,
			CronSpec:     v.CronSpec,
			ExecuteTimes: v.ExecuteTimes,
			PrevTime:     v.PrevTime,
			NextTime:     int(nextTime.Unix()),
			IsAudit:      v.IsAudit,
			Status:       v.Status,
			Ctime:        v.Ctime,
			Utime:        v.Utime,
		}
		data.List = append(data.List, task)
	}

	data.Total = total
	return nil, data
}

func TaskDetailById(taskId int) (*TaskDetail, error) {
	var (
		filters []interface{}
	)
	filters = append(filters, "id", taskId)
	filters = append(filters, "is_delete", 0)
	task := &models.Task{}
	err := task.GetDetail([]string{"id", "group_id", "server_ids", "run_type", "task_name", "description", "cron_spec", "concurrent", "command",
		"timeout", "execute_times", "prev_time", "is_notify", "notify_tpl_id", "notify_user_ids", "status", "ctime", "utime", "create_id", "update_id"}, filters)
	if err != nil {
		return nil, err
	}

	data := &TaskDetail{
		TaskId:        task.Id,
		GroupId:       task.GroupId,
		ServerIds:     task.ServerIds,
		RunType:       task.RunType,
		TaskName:      task.TaskName,
		Description:   task.Description,
		CronSpec:      task.CronSpec,
		Concurrent:    task.Concurrent,
		Command:       task.Command,
		Timeout:       task.Timeout,
		IsNotify:      task.IsNotify,
		NotifyTplId:   task.NotifyTplId,
		NotifyUserIds: task.NotifyUserIds,
		Status:        task.Status,
	}

	return data, nil
}

//创建或者新增
func TaskCreateOrUpdate(params *TaskDetail, userId int) error {
	var (
		Task   *models.Task
		fields []string
	)

	if err := IsBan(params.Command); err != nil && err.Error() == "isBan" {
		return errors.New("存在禁用命令")
	}

	Task = &models.Task{
		Id:            params.TaskId,
		GroupId:       params.GroupId,
		ServerIds:     params.ServerIds,
		RunType:       params.RunType,
		TaskName:      params.TaskName,
		Description:   params.Description,
		CronSpec:      params.CronSpec,
		Concurrent:    params.Concurrent,
		Command:       params.Command,
		Timeout:       params.Timeout,
		IsNotify:      params.IsNotify,
		NotifyTplId:   params.NotifyTplId,
		NotifyUserIds: params.NotifyUserIds,
	}

	if Task.Id > 0 {
		Task.UpdateId = userId
	} else {
		Task.CreateId = userId
	}
	fields = []string{"group_id", "server_ids", "run_type", "task_name", "description", "cron_spec",
		"concurrent", "command", "timeout", "is_notify", "notify_tpl_id", "notify_user_ids", "status",
		"utime", "update_id"}
	if _, err := Task.Create(fields); err != nil {
		return err
	}
	return nil
}

//任务删除
func TaskDelete(taskId int) error {
	filters := make([]interface{}, 2)
	filters[0] = "id"
	filters[1] = taskId
	task := &models.Task{}
	err := task.GetDetail([]string{"id", "status", "is_audit"}, filters)
	if err != nil {
		return err
	}
	if task.Id == 0 {
		return errors.New("任务不存在")
	}
	if task.IsAudit != 2 {
		return errors.New("该任务无法删除，只能设置无效")
	}
	if task.Status == 1 {
		return errors.New("运行中任务无法删除，请先终止")
	}
	task.IsDelete = 1
	editFields := []string{"is_delete"}
	if _, err := task.Create(editFields); err != nil {
		return err
	}
	return nil
}

//任务审核
func TaskAudit(taskId int, audit int8) error {
	filters := make([]interface{}, 2)
	filters[0], filters[1] = "id", taskId
	task := &models.Task{}
	err := task.GetDetail([]string{"id", "is_audit"}, filters)
	if err != nil {
		return err
	}
	if task.Id == 0 {
		return errors.New("任务不存在")
	}
	if task.IsAudit != 0 {
		return errors.New("该任务已被审核，无法操作")
	}
	task.IsAudit = audit
	if _, err := task.Create([]string{"is_audit"}); err != nil {
		return err
	}
	return nil
}

//任务执行
func ExecuteTask(taskId int) error {
	filters := make([]interface{}, 2)
	filters[0], filters[1] = "id", taskId
	task := &models.Task{}
	err := task.GetDetail([]string{"id", "server_ids", "task_name", "command", "is_audit"}, filters)
	if err != nil {
		return err
	}
	if task.Id == 0 {
		return errors.New("任务不存在")
	}
	jobs, err := NewJobFromTask(task)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		job.Run()
	}
	return nil
}

//设置任务状态
func SetTaskState(taskId int, status int8) error {
	fields := []string{"id", "group_id", "is_audit", "server_ids", "run_type", "task_name", "description", "cron_spec", "concurrent", "command",
		"timeout", "job_entry_id", "execute_times", "prev_time", "is_notify", "notify_type", "notify_tpl_id", "notify_user_ids", "status", "utime", "update_id"}
	filters := make([]interface{}, 2)
	filters[0] = "id"
	filters[1] = taskId
	task := &models.Task{}
	err := task.GetDetail(fields, filters)
	if err != nil {
		return err
	}
	if task.Id == 0 {
		return errors.New("任务不存在")
	}
	if task.IsAudit != 1 {
		return errors.New("任务未通过审核")
	}

	fields1 := []string{"status"}
	if status == 1 {
		jobs, err := NewJobFromTask(task)
		if err != nil {
			return err
		}
		for _, job := range jobs {
			AddJob(task.CronSpec, job)
		}
	} else if status == 0 {
		mainCron.Remove(cron.EntryID(task.JobEntryId))
		fields1 = append(fields1, "job_entry_id")
		task.JobEntryId = 0
	}
	task.Status = status
	if _, err := task.Create(fields1); err != nil {
		return err
	}
	return nil
}

func TaskLogList(params *TaskLogParams) (error, *TaskLogListData) {
	var (
		err  error
		opts *options2.FindOptions
	)
	data := &TaskLogListData{
		Total: 0,
		List:  []*TaskLogData{},
	}

	list := make([]*models.TaskRunLog, 0)

	opts = &options2.FindOptions{}
	opts.SetSort(bson.M{"ctime": -1})
	opts.SetLimit(int64(params.Limit))
	opts.SetSkip(int64((params.Page - 1) * params.Limit))
	filters := bson.D{{"task_id", params.TaskId}}

	data.Total, err = util.MgoClient.Find("log", "task_run_log", &list, filters, opts)
	if err != nil {
		return err, nil
	}
	for _, v := range list {
		server := &TaskLogData{
			Id:          v.Id.Hex(),
			TaskId:      v.TaskId,
			ServerId:    v.ServerId,
			ServerName:  v.ServerName,
			Output:      v.Output,
			Error:       v.Error,
			Status:      v.Status,
			ProcessTime: v.ProcessTime,
			Ctime:       v.Ctime,
		}
		data.List = append(data.List, server)
	}
	return nil, data
}

func TaskLogDetailById(LogId string) (*TaskLogData, error) {
	var (
		filters bson.M
		opts    *options2.FindOneOptions
	)
	opts = &options2.FindOneOptions{}
	objectId, err := primitive.ObjectIDFromHex(LogId)
	if err != nil {
		return nil, err
	}
	filters = bson.M{"_id": objectId}
	runLog := &models.TaskRunLog{}
	if err := util.MgoClient.FindOne("log", "task_run_log", filters, runLog, opts); err != nil {
		return nil, err
	}

	taskInfo := TaskDetail{}
	if runLog.TaskId > 0 {
		filtersTask := make([]interface{}, 2)
		filtersTask[0], filtersTask[1] = "id", runLog.TaskId
		task := &models.Task{}
		if err := task.GetDetail([]string{"id", "group_id", "server_ids", "task_name", "description", "cron_spec", "concurrent", "command",
			"timeout", "execute_times", "status"}, filtersTask); err == nil {
			taskInfo.TaskId = task.Id
			taskInfo.GroupId = task.GroupId
			taskInfo.ServerIds = task.ServerIds
			taskInfo.TaskName = task.TaskName
			taskInfo.Description = task.Description
			taskInfo.CronSpec = task.CronSpec
			taskInfo.Concurrent = task.Concurrent
			taskInfo.Command = task.Command
			taskInfo.Timeout = task.Timeout
			taskInfo.IsNotify = task.IsNotify
			taskInfo.Status = task.Status
		}
	}
	data := &TaskLogData{
		Id:          runLog.Id.Hex(),
		TaskId:      runLog.TaskId,
		ServerId:    runLog.ServerId,
		ServerName:  runLog.ServerName,
		Output:      runLog.Output,
		Error:       runLog.Error,
		Status:      runLog.Status,
		ProcessTime: runLog.ProcessTime,
		Ctime:       runLog.Ctime,
		TaskInfo:    taskInfo,
	}
	return data, nil
}

func BanList() (interface{}, error) {
	var (
		err error
	)
	banList := &BanListFind{}
	filters := bson.D{{"ban_list", bson.D{{"$ne", nil}}}}
	err = util.MgoClient.FindOne("config", "ban_list", filters, banList, &options2.FindOneOptions{})
	if err != nil {
		return nil, err
	}

	data := &BanListData{
		Id:      banList.Id.Hex(),
		BanList: banList.BanList,
		Ctime:   banList.Ctime,
	}

	return data, nil
}

func EditBan(data []string, id string) error {
	var err error

	for _, str := range data {
		if common.HasChineseChar(str) {
			return errors.New("hasChinese")
		}
	}

	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)
	ntime, _ := strconv.Atoi(sTime)

	if id == "" {
		//插入
		banList := bson.M{"_id": primitive.NewObjectID(), "ctime": ntime, "utime": 0, "ban_list": data}
		_, err = util.MgoClient.InsertOne("config", "ban_list", banList)
	} else {
		//更新
		banList := bson.M{"$set": bson.M{"utime": ntime, "ban_list": data}}
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		filters := bson.M{"_id": objectId}

		_, err = util.MgoClient.UpdateOne("config", "ban_list", filters, banList, &options2.UpdateOptions{})
	}

	if err != nil {
		return err
	}

	return nil
}

func IsBan(order string) error {
	var (
		err error
	)
	banList := &BanListFind{}
	filters := bson.D{{"ban_list", bson.D{{"$ne", nil}}}}
	err = util.MgoClient.FindOne("config", "ban_list", filters, banList, &options2.FindOneOptions{})
	if err != nil {
		return err
	}

	for _, v := range banList.BanList {
		if strings.Contains(order, v) {
			return errors.New("isBan")
		}
	}

	return nil
}

func GetNotifyData(expand string) NotifyData {
	var (
		fields  []string
		filters []interface{}
	)
	expandArr := strings.Split(expand, ",")
	data := NotifyData{
		AdminList:  nil,
		NotifyType: nil,
		NotifyTpl:  nil,
	}
	for _, v := range expandArr {
		switch v {
		case "admin":
			fields = []string{"id", "real_name"}
			filters = filters[0:0]
			filters = append(filters, "is_delete", 0)
			filters = append(filters, "status", 1)
			adminList, _ := models.GetAdminList(0, 0, filters, fields)
			for _, v := range adminList {
				data.AdminList = append(data.AdminList, AdminBrief{
					Id:       v.Id,
					RealName: v.RealName,
				})
			}
		case "notify_type":
			for _, v := range NotifyTypeList {
				data.NotifyType = append(data.NotifyType, NotifyTypeDetail{
					Id:         v.Id,
					NotifyName: v.NotifyName,
				})
			}
		case "notify_tpl":
			filters = append(filters, "status", 1)
			filters = append(filters, "is_delete", 0)
			fields = []string{"id", "tpl_name"}
			notify := &models.NotifyTpl{}
			notifyList, _ := notify.List(0, 0, filters, fields)
			for _, v := range notifyList {
				data.NotifyTpl = append(data.NotifyTpl, NotifyTplBrief{
					Id:      v.Id,
					TplName: v.TplName,
				})
			}
		}
	}
	return data
}
