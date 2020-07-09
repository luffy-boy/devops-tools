package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
	"tools/server/app/models"
	"tools/server/app/util"
)

/**
redis数据操作
*/

type RoleCache struct {
	Id        int    `json:"id"`
	Role      string `json:"role"`
	RoleName  string `json:"role_name"`
	ParentId  int    `json:"parent_id"`
	RoutesIds string `json:"routes_ids"`
}

type RouteCache struct {
	Id        int    `json:"id"`
	Route     string `json:"route"`
	RouteName string `json:"route_name"`
	Request   string `json:"request"`
}

type MessageData struct {
	TplId       int               `json:"tpl_id"`       //模板id
	ToUser      []string          `json:"to_user"`      //接收人
	ExtraParams map[string]string `json:"extra_params"` //额外参数
}

//定时任务任务日志
type TaskRunLog struct {
	TaskId      int    `json:"task_id"`
	ServerId    int    `json:"server_id"`
	ServerName  string `json:"server_name"`
	Output      string `json:"output"`
	Error       string `json:"error"`
	Status      int8   `json:"status"`
	ProcessTime int    `json:"process_time"`
	Ctime       int    `json:"ctime"`
}

//Admin 管理员-登录redis缓存数据
type AdminLoginRedisData struct {
	UserId       int    `json:"user_id"`      //用户id
	RealName     string `json:"real_name"`    //姓名
	Birthday     int    `json:"birthday"`     //出生日期时间戳
	Phone        string `json:"phone"`        //手机号码
	Sex          int8   `json:"sex"`          //性别
	RoleId       int    `json:"role_id"`      //角色id
	Role         string `json:"role"`         //角色
	RouteIds     string `json:"route_ids"`    //路由id
	Introduction string `json:"introduction"` //介绍
}

//redis数据重新加载
func RedisDataReload(funcList ...string) {
	if len(funcList) == 0 {
		SetRoleData()
		SetRouteData()
		return
	}
	for _, v := range funcList {
		switch v {
		case "SetRoleData":
			SetRoleData()
		case "SetRouteData":
			SetRouteData()
		}
	}
}

//设置管理员登录信息
func SetAdminLoginData(data *models.Admin, Role string) bool {
	if data.Id <= 0 {
		return false
	}
	var uData = &AdminLoginRedisData{
		UserId:       data.Id,
		RealName:     data.RealName,
		Birthday:     data.Birthday,
		Phone:        data.Phone,
		Sex:          data.Sex,
		RoleId:       data.RoleId,
		Role:         Role,
		RouteIds:     data.RouteIds,
		Introduction: data.Introduction,
	}

	tokenExpired, _ := beego.AppConfig.Int64("redis.tokenExpired")
	expire := time.Second * time.Duration(tokenExpired)
	str, err := json.Marshal(uData)
	if err != nil {
		return false
	}
	key := util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(data.Id)
	if err := util.RedisObj.SetEx(key, expire, string(str)); err != nil {
		return false
	}
	return true
}

//删除管理员登录信息
func DelAdminLoginData(userId int) bool {
	if userId == 0 {
		return false
	}
	if ok := util.RedisObj.Del(util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(userId)); ok != nil {
		return false
	}

	return true
}

//获取用户登录信息
func GetAdminLoginData(UserId int) AdminLoginRedisData {
	data := AdminLoginRedisData{}
	ret, err := util.RedisObj.Get(util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(UserId))
	if err != nil {
		return data
	}

	if ret == nil {
		return data
	}

	if err := json.Unmarshal(ret.([]byte), &data); err != nil {
		return data
	}
	return data
}

//获取角色组信息
func getRoleData() (map[int]*RoleCache, error) {
	var (
		RoleList map[int]*RoleCache
	)

	ret, err := util.RedisObj.Get(util.RedisKeyList[util.CmsRoleKey])
	if err != nil {
		return nil, err
	}

	if ret == nil {
		RoleList, _ = SetRoleData()
		return RoleList, nil
	}

	if err := json.Unmarshal(ret.([]byte), &RoleList); err != nil {
		return nil, err
	}

	return RoleList, nil
}

//设置角色组信息
func SetRoleData() (map[int]*RoleCache, error) {
	var (
		RoleList map[int]*RoleCache
	)

	RoleList = make(map[int]*RoleCache)
	var filters []interface{}
	filters = append(filters, "status", 1)
	filters = append(filters, "is_delete", 0)
	list, _ := models.GetRolesList(0, 0, filters, []string{})
	if len(list) == 0 {
		expire := time.Second * 600
		str, _ := json.Marshal(&RoleList)
		if err := util.RedisObj.SetEx(util.RedisKeyList[util.CmsRoleKey], expire, string(str)); err != nil {
			return RoleList, err
		}
		return RoleList, nil
	}

	for _, v := range list {
		RoleList[v.Id] = &RoleCache{
			Id:        v.Id,
			Role:      v.Role,
			RoleName:  v.RoleName,
			ParentId:  v.ParentId,
			RoutesIds: v.RouteIds,
		}
	}

	str, err := json.Marshal(&RoleList)
	if err != nil {
		return RoleList, err
	}

	expire := time.Second * (3600*2 + 300)
	if err := util.RedisObj.SetEx(util.RedisKeyList[util.CmsRoleKey], expire, string(str)); err != nil {
		return RoleList, err
	}
	return RoleList, nil
}

//获取路由组信息
func GetRouteData() (map[int]*RouteCache, error) {
	var (
		RouteList map[int]*RouteCache
	)

	ret, err := util.RedisObj.Get(util.RedisKeyList[util.CmsRouteKey])
	if err != nil {
		return nil, err
	}

	if ret == nil {
		RouteList, _ = SetRouteData()
		return RouteList, nil
	}

	if err := json.Unmarshal(ret.([]byte), &RouteList); err != nil {
		return nil, err
	}
	return RouteList, nil
}

//设置路由信息
func SetRouteData() (map[int]*RouteCache, error) {
	var (
		RouteList map[int]*RouteCache
	)
	RouteList = make(map[int]*RouteCache)
	var filters []interface{}
	filters = append(filters, "status", 1)
	filters = append(filters, "is_delete", 0)
	list, _ := models.GetRoutesList(0, 0, filters, []string{})
	if len(list) == 0 {
		expire := time.Second * 600
		str, _ := json.Marshal(&RouteList)
		if err := util.RedisObj.SetEx(util.RedisKeyList[util.CmsRouteKey], expire, string(str)); err != nil {
			return RouteList, err
		}
		return RouteList, nil
	}
	for _, v := range list {
		RouteList[v.Id] = &RouteCache{
			Id:        v.Id,
			Route:     v.Route,
			RouteName: v.RouteName,
			Request:   v.Request,
		}
	}

	str, err := json.Marshal(&RouteList)
	if err != nil {
		return RouteList, err
	}
	expire := time.Second * 3600 * 2
	if err := util.RedisObj.SetEx(util.RedisKeyList[util.CmsRouteKey], expire, string(str)); err != nil {
		return RouteList, err
	}
	return RouteList, nil
}

func PushTakRunLog(log *TaskRunLog) error {
	logStr, _ := json.Marshal(&log)
	if err := util.RedisObj.LPush(util.RedisKeyList[util.JobTaskRunLog], string(logStr)); err != nil {
		return errors.New(fmt.Sprintf("定时任务日志写入Redis失败，错误原因:%s", err.Error()))
	}
	return nil
}

func PushMessage(message *MessageData) error {
	logStr, _ := json.Marshal(&message)
	if err := util.RedisObj.LPush(util.RedisKeyList[util.MessageQueue], string(logStr)); err != nil {
		return errors.New(fmt.Sprintf("消息通知写入Redis失败，错误原因:%s", err.Error()))
	}
	return nil
}

//拉取任务运行日志信息
func PullTaskRunLog() TaskRunLog {
	data := TaskRunLog{}
	ret, err := util.RedisObj.RPop(util.RedisKeyList[util.JobTaskRunLog])
	if err != nil {
		return data
	}
	if ret == nil {
		return data
	}
	if err := json.Unmarshal(ret.([]byte), &data); err != nil {
		return data
	}
	return data
}
