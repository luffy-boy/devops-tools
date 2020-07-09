package service

import (
	"errors"
	"strings"
	"tools/server/app/models"
)

//查询
type NotifyTplParams struct {
	Page  int
	Limit int
}

type NotifyTplOperation struct {
	TplIds string `json:"tpl_ids"`
	Status int8   `json:"status"`
	Audit  int8   `json:"audit"`
}

//查询返回数据
type NotifyTplListData struct {
	Total int64             `json:"total"`
	List  []NotifyTplDetail `json:"list"`
}

type NotifyTplDetail struct {
	Id             int    `json:"id"`
	TplName        string `json:"tpl_name" valid:"Required;"`
	NotifyType     int8   `json:"notify_type" valid:"Required;"`
	NotifyTypeName string `json:"notify_type_name"`
	TplData        string `json:"tpl_data,omitempty" valid:"Required;MaxSize(500)" `
	Status         int8   `json:"status"`
	Ctime          int    `json:"ctime,omitempty"`
	Utime          int    `json:"utime,omitempty"`
}

type NotifyTplBrief struct {
	Id      int    `json:"id"`
	TplName string `json:"tpl_name"`
}

type NotifyTypeDetail struct {
	Id         int8   `json:"id"`
	NotifyName string `json:"notify_name"`
}

var (
	NotifyTypeList = map[int8]NotifyTypeDetail{
		1: {Id: 1, NotifyName: "邮件"},
	}
)

func NotifyTplList(params *NotifyTplParams) (error, NotifyTplListData) {
	var (
		filters        []interface{}
		fields         []string
		notifyTypeName string
	)

	filters = append(filters, "is_delete", 0)
	fields = []string{"id", "tpl_name", "notify_type", "tpl_data", "status", "ctime", "utime"}
	notify := &models.NotifyTpl{}
	list, total := notify.List(params.Page, params.Limit, filters, fields)

	data := NotifyTplListData{
		Total: 0,
		List:  []NotifyTplDetail{},
	}

	for _, v := range list {
		notifyTypeName = ""
		if v, ok := NotifyTypeList[v.NotifyType]; ok {
			notifyTypeName = v.NotifyName
		}
		data.List = append(data.List, NotifyTplDetail{
			Id:             v.Id,
			TplName:        v.TplName,
			TplData:        v.TplData,
			NotifyType:     v.NotifyType,
			NotifyTypeName: notifyTypeName,
			Status:         v.Status,
			Ctime:          v.Ctime,
			Utime:          v.Utime,
		})
	}

	data.Total = total
	return nil, data
}

func NotifyDetailById(id int) (NotifyTplDetail, error) {
	var (
		filters []interface{}
		data    NotifyTplDetail
	)
	filters = append(filters, "id", id)
	filters = append(filters, "is_delete", 0)
	NotifyTpl := &models.NotifyTpl{
		Id:       id,
		IsDelete: 0,
	}
	err := NotifyTpl.Detail([]string{"id", "tpl_name", "notify_type", "tpl_data", "status"}, filters)
	if err != nil {
		return data, err
	}
	if NotifyTpl.Id == 0 {
		return data, errors.New("消息模板不存在")
	}

	data.Id = NotifyTpl.Id
	data.TplName = NotifyTpl.TplName
	data.TplData = NotifyTpl.TplData
	data.NotifyType = NotifyTpl.NotifyType
	data.Status = NotifyTpl.Status

	return data, nil
}

//创建或者新增
func NotifyCreateOrUpdate(data *NotifyTplDetail, adminId int) error {
	var (
		Notify *models.NotifyTpl
		fields []string
	)

	data.TplName = strings.TrimSpace(strings.Trim(strings.ToLower(data.TplName), " "))
	Notify = &models.NotifyTpl{
		Id:         data.Id,
		TplName:    data.TplName,
		TplData:    data.TplData,
		NotifyType: data.NotifyType,
		Status:     data.Status,
	}
	if Notify.Id > 0 {
		Notify.UpdateId = adminId
	} else {
		Notify.CreateId = adminId
	}

	fields = []string{"tpl_name", "tpl_data", "notify_type", "update_id", "status", "utime"}
	if _, err := Notify.Create(fields); err != nil {
		return err
	}
	return nil
}

func NotifyDelete(NotifyTplId int) error {
	var (
		filters []interface{}
		fields  []string
	)
	filters = append(filters, "id", NotifyTplId)
	filters = append(filters, "is_delete", 0)
	NotifyTpl := &models.NotifyTpl{
		Id:       NotifyTplId,
		IsDelete: 0,
	}
	err := NotifyTpl.Detail([]string{"id"}, filters)
	if err != nil {
		return err
	}
	if NotifyTpl.Id == 0 {
		return err
	}
	if NotifyTpl.IsAudit != 2 {
		return errors.New("该模板无法删除，只能设置无效")
	}
	if NotifyTpl.Status == 1 {
		return errors.New("请先设置无效")
	}
	NotifyTpl.IsDelete = 1
	fields = append(fields, "is_delete")
	if _, err := NotifyTpl.Create(fields); err != nil {
		return err
	}
	return nil
}
