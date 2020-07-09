package models

import (
	"github.com/astaxie/beego/orm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

type Task struct {
	Id            int
	GroupId       int    //任务分组id
	ServerIds     string //服务器id
	RunType       int8   //执行策略
	TaskName      string //任务名称
	Description   string //任务描述
	CronSpec      string //时间表达式
	Concurrent    int8   //是否允许并发
	Command       string //执行命令
	Timeout       int    //超时设置
	JobEntryId    int    //Job id
	ExecuteTimes  int    //累计执行次数
	PrevTime      int    //上次执行时间
	IsNotify      int8   //是否通知
	NotifyTplId   int    //通知模板id
	NotifyUserIds string //通知用户集合
	CreateId      int    //创建者id
	UpdateId      int    //修改者id
	IsAudit       int8   //审核状态
	Status        int8   //状态
	IsDelete      int8   //状态
	Ctime         int    //创建时间
	Utime         int    //修改时间
}

type TaskRunLog struct {
	Id          primitive.ObjectID `bson:"_id"`
	TaskId      int                `bson:"task_id"`
	ServerId    int                `bson:"server_id"`
	ServerName  string             `bson:"server_name"`
	Output      string             `bson:"output"`
	Error       string             `bson:"error"`
	Status      int8               `bson:"status"` //任务执行状态 1.成功 2.超时 3.失败
	ProcessTime int                `bson:"process_time"`
	Ctime       int                `bson:"ctime"`
}

func (self *Task) TableName() string {
	return "task"
}

func GetTaskList(page int, pageSize int, filters []interface{}, fields []string) ([]*Task, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "group_id", "server_ids", "run_type", "description", "cron_spec", "concurrent", "command", "timeout", "is_notify", "notify_type", "notify_tpl_id", "notify_user_ids", "status"}
	}
	offset := (page - 1) * pageSize
	list := make([]*Task, 0)
	admin := &Task{}
	model := o.QueryTable(admin)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := model.Count()
	model.Limit(pageSize, offset).OrderBy("-id").All(&list, fields...)
	return list, total
}

//创建或者修改 角色信息
func (self *Task) Create(fields []string) (num int64, err error) {
	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)

	if self.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"group_id", "server_ids", "run_type", "task_name", "description", "cron_spec",
				"concurrent", "command", "timeout", "is_notify", "notify_tpl_id", "notify_user_ids",
				"status", "utime", "update_id"}
		}
		self.Utime, _ = strconv.Atoi(sTime)
		if num, err = o.Update(self, fields...); err != nil {
			return
		}
	} else {
		self.Ctime, _ = strconv.Atoi(sTime)
		if num, err = o.Insert(self); err != nil {
			return
		}
	}
	return
}

func (self *Task) GetDetail(fields []string, filters []interface{}) error {
	if len(fields) == 0 {
		fields = []string{"id", "group_id", "server_ids", "run_type", "task_name", "description", "cron_spec", "concurrent", "command",
			"timeout", "execute_times", "prev_time", "is_notify", "notify_type", "notify_tpl_id", "notify_user_ids", "status", "ctime", "utime", "create_id", "update_id"}
	}
	model := o.QueryTable(self)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}

	if err := model.One(self, fields...); err != nil && err != orm.ErrNoRows {
		return err
	}
	return nil
}
