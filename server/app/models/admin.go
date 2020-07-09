package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//管理员表
type Admin struct {
	Id           int
	Username     string //用户名
	Password     string //密码
	RealName     string //真实姓名
	Salt         string //加密盐
	Status       int8   //有效状态
	IsDelete     int8   //是否删除
	Birthday     int    //生日
	Phone        string //手机
	Email        string //邮箱
	Sex          int8   //性别 1.男  2.女 3.不详
	RoleId       int    //角色id
	RouteIds     string //角色id
	LoginIpAddr  string //登陆ip
	LoginTime    int    //登录时间
	Introduction string //介绍
	Ctime        int    //创建时间
	Utime        int    //修改时间
}

func (c *Admin) TableName() string {
	return "admin"
}

func GetAdminList(page int, pageSize int, filters []interface{}, fields []string, orderBy ...string) ([]*Admin, int64) {
	if len(fields) == 0 {
		fields = []string{"id", "real_name", "salt", "status", "role_id"}
	}
	if len(orderBy) == 0 {
		orderBy = append(orderBy, "id")
	}
	offset := (page - 1) * pageSize
	list := make([]*Admin, 0)
	admin := &Admin{}
	model := o.QueryTable(admin)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}
	total, _ := model.Count()
	model.Limit(pageSize, offset).OrderBy(orderBy...).All(&list, fields...)
	return list, total
}

/**
查询管理员信息
*/
func GetDetail(fields []string, filters []interface{}) (*Admin, error) {
	if len(fields) == 0 {
		fields = []string{"id", "username", "password", "real_name", "salt", "status", "role_id"}
	}

	admin := &Admin{}
	model := o.QueryTable(admin)
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			model = model.Filter(filters[k].(string), filters[k+1])
		}
	}

	if err := model.One(admin, fields...); err != nil && err != orm.ErrNoRows {
		return nil, err
	}
	return admin, nil
}

//创建或者修改管理员信息
func (self *Admin) Create(data *Admin, fields []string) (num int64, err error) {
	ctime := time.Now().Unix()
	str := strconv.FormatInt(ctime, 10)
	if data.Id == 1 {
		fields = []string{"login_ip_addr", "login_time"}
	}
	if data.Id > 0 {
		if len(fields) == 0 {
			fields = []string{"password", "real_name", "status", "is_delete", "birthday", "phone", "sex", "role_id", "route_ids", "utime"}
		}
		data.Utime, _ = strconv.Atoi(str)
		if num, err = o.Update(data, fields...); err != nil {
			return
		}
	} else {
		data.Ctime, _ = strconv.Atoi(str)
		if num, err = o.Insert(data); err != nil {
			return
		}
	}
	return
}
