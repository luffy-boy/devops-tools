package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//用户角色配置表
type Roles struct {
	Id       int
	Role	 string	//角色标识
	RoleName string //角色名称
	ParentId int   //上级id
	Status   int8   //有效状态
	IsDelete int8   //是否删除
	RouteIds string //有效路由id集合
	Ctime    int    //创建时间
	Utime    int    //修改时间
}

func (self *Roles) TableName() string {
	return "roles"
}

func GetRolesList(page, pageSize int, filters []interface{}, fields []string) ([]*Roles, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "role", "role_name", "parent_id", "status", "is_delete", "route_ids", "ctime", "utime"}
	}
	offset := (page - 1) * pageSize
	list := make([]*Roles, 0)
	model := o.QueryTable("app_roles")
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
func (self *Roles) Create(data *Roles, fields []string) (num int64, err error) {
	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)

	if data.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"role","role_name", "parent_id", "status", "is_delete", "route_ids"}
		}
		data.Utime, _ = strconv.Atoi(sTime)
		fields = append(fields, "utime")
		if num, err = o.Update(data, fields...); err != nil {
			return
		}
	} else {
		data.Ctime, _ = strconv.Atoi(sTime)
		if num, err = o.Insert(data); err != nil {
			return
		}
	}
	return
}

//详情
func (self *Roles) Detail(fields []string, filters []interface{}) error{
	if len(fields) == 0 {
		fields = []string{"id", "role_name", "role", "parent_id", "status"}
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
