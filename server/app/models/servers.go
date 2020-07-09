package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type TaskServers struct {
	Id             int
	GroupId        int
	ConnectionType int8
	ServerName     string
	ServerAccount  string
	ServerIp       string
	Port           int
	Password       string
	PrivateKeySrc  string
	PublicKeySrc   string
	Type           int8
	Detail         string
	Status         int8
	IsDelete       int8
	Ctime          int
	Utime          int
}

func (self *TaskServers) TableName() string {
	return "task_servers"
}

func (self *TaskServers) List(page, pageSize int, filters []interface{}, fields []string) ([]*TaskServers, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "route_name", "route", "request", "parent_id", "status"}
	}
	offset := (page - 1) * pageSize
	list := make([]*TaskServers, 0)
	model := o.QueryTable(self)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := model.Count()
	model.OrderBy("id").Limit(pageSize, offset).All(&list, fields...)
	return list, total
}

//创建或者修改 角色信息
func (self *TaskServers) Create(fields []string) (num int64, err error) {
	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)

	if self.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"group_id", "connection_type", "server_name", "server_account", "server_outer_ip", "server_ip", "port",
				"password", "private_key_src", "public_key_src", "type", "detail", "status", "utime"}
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

func (self *TaskServers) GetDetail(fields []string, filters []interface{}) error {
	if len(fields) == 0 {
		fields = []string{"group_id", "connection_type", "server_name", "server_account", "server_outer_ip", "server_ip"}
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
