package models

import (
	"strconv"
	"time"
)

type NotifyTpl struct {
	Id         int
	TplName    string //模板名称
	TplData    string //消息模板内容
	NotifyType int8   //通知类型
	CreateId   int    //创建者id
	UpdateId   int    //修改者id
	Status     int8   //状态 0.无效 1.有效
	IsAudit    int8   //审核状态 0.未审核 1.已审核 2.已拒绝
	IsDelete   int8   //删除状态 0.未删除 1.已删除
	Ctime      int    //创建时间
	Utime      int    //修改时间
}

func (c *NotifyTpl) TableName() string {
	return "notify_tpl"
}

func (self *NotifyTpl) List(page int, pageSize int, filters []interface{}, fields []string, orderBy ...string) ([]*NotifyTpl, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "tpl_name", "tpl_data", "is_audit", "status"}
	}
	if len(orderBy) == 0 {
		orderBy = append(orderBy, "id")
	}
	offset := (page - 1) * pageSize
	list := make([]*NotifyTpl, 0)
	model := o.QueryTable(self)
	l := len(filters)
	if l > 0 {
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := model.Count()
	model.Limit(pageSize, offset).OrderBy(orderBy...).All(&list, fields...)
	return list, total
}

func (self *NotifyTpl) Detail(fields []string, filters []interface{}) error {
	if len(fields) == 0 {
		fields = []string{"id", "tpl_name", "tpl_data", "is_audit", "status"}
	}
	model := o.QueryTable(self)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	if err := model.One(self, fields...); err != nil {
		return err
	}
	return nil
}

func (self *NotifyTpl) Create(fields []string) (num int64, err error) {
	ctime := time.Now().Unix()
	str := strconv.FormatInt(ctime, 10)
	if self.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"tpl_name", "tpl_data", "status", "create_id", "update_id", "utime"}
		}
		self.Utime, _ = strconv.Atoi(str)
		if num, err = o.Update(self, fields...); err != nil {
			return
		}
	} else {
		self.Ctime, _ = strconv.Atoi(str)
		if num, err = o.Insert(self); err != nil {
			return
		}
	}
	return
}
