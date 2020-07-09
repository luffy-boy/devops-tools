package service

import (
	"github.com/astaxie/beego/validation"
	"strings"
	"tools/server/app/models"
)

//查询
type ServersParams struct {
	Page   int
	Limit  int
	Extend string
}

//查询返回数据
type ServersListData struct {
	Total int64          `json:"total"` //資源服務器总数
	List  []*ServersData `json:"list"`  // 資源服務器列表
}

type ServersData struct {
	Id             int    `json:"server_id"`
	GroupName      string `json:"group_name"`      //分組名名称
	ServerName     string `json:"server_name"`     //服务器名称
	ConnectionType int8   `json:"connection_type"` //连接方式
	ServerIp       string `json:"server_ip"`       //连接ip
	Port           int    `json:"port"`            //端口
	Detail         string `json:"detail"`          //备注
	Status         int8   `json:"status"`          //状态
	Ctime          int    `json:"ctime"`           //创建时间
	Utime          int    `json:"utime"`           //修改时间
}

type ServersDetail struct {
	ServerId       int    `json:"server_id"`
	GroupId        int    `json:"group_id" valid:"Required;"`         //服务器id
	ConnectionType int8   `json:"connection_type" valid:"Range(0,1)"` //连接方式
	ServerName     string `json:"server_name" valid:"Required;"`      //服务器名称
	Port           int    `json:"port" valid:"Required;"`             //端口
	ServerIp       string `json:"server_ip" valid:"Required;IP"`      //服务器id
	ServerAccount  string `json:"server_account"`                     //服务器账号
	Password       string `json:"password"`                           //密码
	PrivateKeySrc  string `json:"private_key_src"`                    //私钥路径
	PublicKeySrc   string `json:"public_key_src"`                     //公钥路径
	Type           int8   `json:"type" valid:"Range(0,1)"`            //登录类型
	Detail         string `json:"detail" valid:"Required;"`           //备注
	Status         int8   `json:"status" valid:"Required;Range(0,1)"` //状态
	Ctime          int    `json:"ctime,omitempty"`
	Utime          int    `json:"utime,omitempty"`
}

func (s *ServersDetail) Valid(v *validation.Validation) {
	if s.Type == 0 {
		if s.ServerAccount == "" {
			v.SetError("server_account", "server_account Can not be empty")
			return
		}
		if s.Password == "" {
			v.SetError("server_account", "password Can not be empty")
			return
		}
	} else if s.Type == 1 {
		if s.PrivateKeySrc == "" {
			v.SetError("private_key_src", "private_key_src Can not be empty")
			return
		}
		if s.PublicKeySrc == "" {
			v.SetError("public_key_src", "public_key_src Can not be empty")
			return
		}
	}
}

var ServersGroupList = map[int]map[string]interface{}{
	1: {"id": 1, "group_name": "内网服务器"},
	2: {"id": 2, "group_name": "外网服务器"},
}

//查询管理员列表
func ServersList(params *ServersParams) (error, *ServersListData) {
	var (
		filters []interface{}
		fields  []string
	)

	filters = append(filters, "is_delete", 0)
	fields = []string{"id", "group_id", "server_name", "connection_type", "server_ip", "port", "status", "detail", "ctime", "utime"}
	server := &models.TaskServers{}
	list, total := server.List(params.Page, params.Limit, filters, fields)

	data := &ServersListData{
		Total: 0,
		List:  []*ServersData{},
	}

	if params.Extend != "" {
		extend := strings.Split(params.Extend, ",")
		for _, v := range extend {
			switch v {
			case "local":
				server := &ServersData{
					Id:         0,
					GroupName:  "本地",
					ServerName: "本地服务器",
				}
				data.List = append(data.List, server)
			}
		}
	}

	for _, v := range list {
		groupName, ok := ServersGroupList[v.GroupId]["group_name"]
		if !ok {
			groupName = ""
		}
		server := &ServersData{
			Id:             v.Id,
			GroupName:      groupName.(string),
			ServerName:     v.ServerName,
			ConnectionType: v.ConnectionType,
			ServerIp:       v.ServerIp,
			Port:           v.Port,
			Detail:         v.Detail,
			Status:         v.Status,
			Ctime:          v.Ctime,
			Utime:          v.Utime,
		}
		data.List = append(data.List, server)
	}

	data.Total = total
	return nil, data
}

//创建或者新增
func ServersCreateOrUpdate(params *ServersDetail) error {
	var (
		Servers *models.TaskServers
		fields  []string
	)

	Servers = &models.TaskServers{
		Id:             params.ServerId,
		GroupId:        params.GroupId,
		ConnectionType: params.ConnectionType,
		ServerName:     params.ServerName,
		ServerAccount:  params.ServerAccount,
		ServerIp:       params.ServerIp,
		Port:           params.Port,
		Password:       params.Password,
		PrivateKeySrc:  params.PrivateKeySrc,
		PublicKeySrc:   params.PublicKeySrc,
		Type:           params.Type,
		Detail:         params.Detail,
		Status:         params.Status,
	}

	fields = []string{"group_id", "connection_type", "server_name", "server_account", "server_outer_ip", "server_ip", "port",
		"password", "private_key_src", "public_key_src", "type", "detail", "status", "utime"}
	if _, err := Servers.Create(fields); err != nil {
		return err
	}
	return nil
}

func ServersDetailById(serverId int) (*ServersDetail, error) {
	var (
		filters []interface{}
	)
	filters = append(filters, "id", serverId)
	filters = append(filters, "is_delete", 0)
	server := &models.TaskServers{}
	err := server.GetDetail([]string{"id", "group_id", "connection_type", "server_name", "server_account", "server_ip", "port", "password", "private_key_src", "public_key_src", "type", "detail", "status"}, filters)
	if err != nil {
		return nil, err
	}

	data := &ServersDetail{
		ServerId:       server.Id,
		GroupId:        server.GroupId,
		ConnectionType: server.ConnectionType,
		ServerName:     server.ServerName,
		Port:           server.Port,
		ServerIp:       server.ServerIp,
		ServerAccount:  server.ServerAccount,
		Password:       server.Password,
		PrivateKeySrc:  server.PrivateKeySrc,
		PublicKeySrc:   server.PublicKeySrc,
		Type:           server.Type,
		Detail:         server.Detail,
		Status:         server.Status,
	}

	return data, nil
}

func ServersDelete(routeId int) error {
	var (
		fields []string
	)
	route := &models.Routes{
		Id:       routeId,
		IsDelete: 1,
	}
	fields = append(fields, "is_delete")
	if _, err := route.Create(route, fields); err != nil {
		return err
	}
	return nil
}
