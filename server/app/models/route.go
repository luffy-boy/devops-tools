package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

//路由配置表
type Routes struct {
	Id        int
	RouteName string //路由名称
	Route     string //路由
	Request   string //请求方式
	ParentId  int    //父id
	Status    int8   //有效状态
	IsDelete  int8   //是否删除
	Path      string //路径
	Component string //组成部分 view地址
	Name      string //路由唯一标识
	Redirect  string //重定向地址
	Hidden    int8   //是否显示  0.不显示 1.显示
	Icon      string //图标
	Extra     string //额外参数
	Sort      int16  //排序
	IsRoute   int8   //前端路由 1.是 0.不是
	Ctime     int    //创建时间
	Utime     int    //修改时间
}

func (self *Routes) TableName() string {
	return "routes"
}

func GetRoutesList(page, pageSize int, filters []interface{}, fields []string, orderBy ...string) ([]*Routes, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "route_name", "route", "request", "parent_id", "status"}
	}
	if len(orderBy) == 0 {
		orderBy = append(orderBy, "id")
	}
	offset := (page - 1) * pageSize
	list := make([]*Routes, 0)
	model := o.QueryTable("app_routes")
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := model.Count()
	model.OrderBy(orderBy...).Limit(pageSize, offset).All(&list, fields...)
	return list, total
}

func (self *Routes) Create(data *Routes, fields []string) (num int64, err error) {
	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)

	data.Route = strings.Trim(strings.ToLower(data.Route), "/")
	if data.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"route_name", "route", "request", "status", "is_delete"}
		}
		data.Utime, _ = strconv.Atoi(sTime)
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

func (self *Routes) Detail(fields []string, filters []interface{}) error {
	if len(fields) == 0 {
		fields = []string{"id", "route_name", "route", "request", "parent_id", "status"}
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
