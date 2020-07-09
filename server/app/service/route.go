package service

import (
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
	"tools/server/app/common"
	"tools/server/app/models"
	"tools/server/app/response"
)

//查询
type RouteParams struct {
	Page  int
	Limit int
}

//查询返回数据
type RouteListData struct {
	Total int64            `json:"total"` //路由总数
	List  []*RouteDataTree `json:"list"`  // 路由列表
}

type RouteDataTree struct {
	Id        int              `json:"route_id"`
	ParentId  int              `json:"parent_id"`          //父id
	RouteName string           `json:"route_name"`         //路由名称
	Route     string           `json:"route"`              //路由
	Request   string           `json:"request"`            //请求方式
	Status    int8             `json:"status"`             //状态
	Sort      int16            `json:"sort"`               //排序
	Child     []*RouteDataTree `json:"children,omitempty"` //子节点
}

type RouteDetail struct {
	Id        int    `json:"route_id"`
	RouteName string `json:"route_name" valid:"Required;MaxSize(32)"`
	Route     string `json:"route" valid:"Required;MaxSize(50)"`
	Request   string `json:"request" valid:"MaxSize(10)"`
	ParentId  int    `json:"parent_id,omitempty"`
	Component string `json:"component"`
	Title     string `json:"title"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Redirect  string `json:"redirect"`
	Hidden    int8   `json:"hidden"`
	Icon      string `json:"icon"`
	Extra     string `json:"extra"`
	Sort      int16  `json:"sort"`
	IsRoute   int8   `json:"is_route"`
	Status    int8   `json:"status"`
	Ctime     int    `json:"ctime,omitempty"`
	Utime     int    `json:"utime,omitempty"`
}

//前端菜单
type Menu struct {
	List []*MenuData `json:"list"` // 路由列表
}

type Meta struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
}
type MenuData struct {
	Id        int         `json:"-"`
	ParentId  int         `json:"-"`                  //父id
	Path      string      `json:"path"`               //路径
	Redirect  string      `json:"redirect"`           //重定向
	Name      string      `json:"name"`               //唯一标识
	Meta      Meta        `json:"meta,omitempty"`     //meta信息
	Component string      `json:"component"`          //view地址
	Child     []*MenuData `json:"children,omitempty"` //子节点
}

type RouteOperation struct {
	RouteId int   `json:"route_id"`
	Sort    int16 `json:"sort"`
}

func (data *RouteDetail) Valid(v *validation.Validation) {
	var chineseErr = "禁止输入中文字符"
	if data.Path != "" {
		if common.HasChineseChar(data.Path) {
			v.SetError("path", chineseErr)
			return
		}
	}
	if data.Icon != "" {
		if common.HasChineseChar(data.Icon) {
			v.SetError("icon", chineseErr)
			return
		}
	}
	if data.Redirect != "" {
		if common.HasChineseChar(data.Icon) {
			v.SetError("redirect", chineseErr)
			return
		}
	}
	if data.Route != "" {
		if common.HasChineseChar(data.Icon) {
			v.SetError("icon", chineseErr)
			return
		}
	}
}

func RouteList(routeJ *RouteParams) (error, *RouteListData) {
	var (
		filters []interface{}
		fields  []string
	)

	orderBy := "sort"
	filters = append(filters, "is_delete", 0)
	fields = []string{"id", "route_name", "request", "route", "parent_id", "sort", "status", "ctime", "utime"}
	list, total := models.GetRoutesList(routeJ.Page, routeJ.Limit, filters, fields, orderBy)

	var routeDataList []*RouteDataTree
	for _, v := range list {
		arr := &RouteDataTree{
			Id:        v.Id,
			ParentId:  v.ParentId,
			Route:     v.Route,
			RouteName: v.RouteName,
			Request:   v.Request,
			Status:    v.Status,
			Sort:      v.Sort,
			Child:     nil,
		}
		routeDataList = append(routeDataList, arr)
	}

	arr := superRouteTree(routeDataList, 0)
	data := &RouteListData{
		Total: 0,
		List:  nil,
	}
	data.List = arr
	data.Total = total
	return nil, data
}

func RouteAll() []*RouteDetail {
	var (
		filters   []interface{}
		fields    []string
		routeList []*RouteDetail
	)

	filters = append(filters, "is_delete", 0)
	filters = append(filters, "status", 1)
	fields = []string{"id", "route_name", "route", "parent_id", "status", "ctime", "utime"}
	list, _ := models.GetRoutesList(0, 0, filters, fields)

	for _, v := range list {
		arr := &RouteDetail{
			Id:        v.Id,
			ParentId:  v.ParentId,
			Route:     v.Route,
			RouteName: v.RouteName,
			Status:    v.Status,
		}
		routeList = append(routeList, arr)
	}
	return routeList
}

func MenuList(adminId int) (error, *Menu) {
	var (
		filters []interface{}
		fields  []string
	)
	orderBy := "sort"
	filters = append(filters, "is_delete", 0)
	filters = append(filters, "status", 1)
	filters = append(filters, "is_route", 1)
	fields = []string{"id", "route_name", "route", "parent_id", "path", "name", "redirect",
		"hidden", "icon", "extra", "component"}
	list, _ := models.GetRoutesList(0, 0, filters, fields, orderBy)

	routeList, err := GetUserRouteById(adminId)
	if err != nil {
		return err, nil
	}

	var menuData []*MenuData
	for _, v := range list {
		if !common.InArray(v.Id, routeList[:]) {
			continue
		}
		arr := &MenuData{
			Id:        v.Id,
			ParentId:  v.ParentId,
			Path:      v.Path + v.Extra,
			Component: v.Component,
			Redirect:  v.Redirect,
			Name:      v.Name,
			Meta: Meta{
				Title: v.RouteName,
				Icon:  v.Icon,
			},
			Child: nil,
		}
		menuData = append(menuData, arr)
	}

	arr := superMenuTree(menuData, 0)
	data := &Menu{
		List: arr,
	}
	data.List = arr
	return nil, data
}

func RouteDetailById(routeId int) (*RouteDetail, error) {
	var (
		filters []interface{}
	)
	filters = append(filters, "id", routeId)
	filters = append(filters, "is_delete", 0)
	route := &models.Routes{}
	err := route.Detail([]string{"id", "route_name", "route", "request", "parent_id",
		"path", "component", "name", "is_route", "redirect", "hidden", "icon", "extra", "sort", "status"}, filters)
	if err != nil {
		return nil, err
	}

	data := &RouteDetail{
		Id:        route.Id,
		RouteName: route.RouteName,
		Route:     route.Route,
		Request:   route.Request,
		ParentId:  route.ParentId,
		Path:      route.Path,
		Component: route.Component,
		Name:      route.Name,
		Redirect:  route.Redirect,
		Hidden:    route.Hidden,
		Icon:      route.Icon,
		Extra:     route.Extra,
		Sort:      route.Sort,
		Status:    route.Status,
		IsRoute:   route.IsRoute,
		Ctime:     route.Ctime,
		Utime:     route.Utime,
	}

	return data, nil
}

//递归实现(返回树状结果得数据)
func superRouteTree(allRouteData []*RouteDataTree, pid int) []*RouteDataTree {
	var arr []*RouteDataTree
	for _, v := range allRouteData {
		if pid == v.ParentId {
			chileRoute := superRouteTree(allRouteData, v.Id)
			v.Child = chileRoute
			arr = append(arr, v)
		}
	}
	return arr
}

//递归实现(返回树状结果得数据)
func superMenuTree(list []*MenuData, pid int) []*MenuData {
	var arr []*MenuData
	for _, v := range list {
		if pid == v.ParentId {
			chileRoute := superMenuTree(list, v.Id)
			v.Child = chileRoute
			arr = append(arr, v)
		}
	}
	return arr
}

//创建或者新增
func RouteCreateOrUpdate(params *RouteDetail) error {
	var (
		route  *models.Routes
		fields []string
	)

	params.RouteName = strings.TrimSpace(strings.Trim(strings.ToLower(params.RouteName), "/"))
	params.Path = strings.TrimSpace(strings.Trim(strings.ToLower(params.Path), "/"))
	params.Route = strings.TrimSpace(strings.ToLower(params.Route))
	route = &models.Routes{
		Id:        params.Id,
		RouteName: params.RouteName,
		Route:     params.Route,
		Request:   params.Request,
		ParentId:  params.ParentId,
		Status:    params.Status,
		Component: params.Component,
		Name:      params.Name,
		Path:      params.Path,
		Redirect:  params.Redirect,
		Hidden:    params.Hidden,
		Icon:      params.Icon,
		Extra:     params.Extra,
		Sort:      params.Sort,
		IsRoute:   params.IsRoute,
	}

	fields = []string{"route_name", "route", "request", "status", "parent_id", "utime", "component",
		"name", "path", "redirect", "hidden", "icon", "extra", "sort", "sort", "is_route"}
	if _, err := route.Create(route, fields); err != nil {
		return err
	}
	RedisDataReload([]string{"SetRouteData"}...) //刷新路由信息
	return nil
}

func RouteDelete(routeId int) int {
	var (
		filters []interface{}
		fields  []string
	)
	filters = append(filters, "parent_id", routeId)
	filters = append(filters, "is_delete", 0)
	route := &models.Routes{}
	err := route.Detail([]string{"id"}, filters)
	if err != nil {
		return response.DbQueryErr
	}
	if route.Id > 0 {
		return response.RouteHasChildErr
	}
	routeUpdate := &models.Routes{
		Id:       routeId,
		IsDelete: 1,
	}
	fields = append(fields, "is_delete")
	if _, err := routeUpdate.Create(routeUpdate, fields); err != nil {
		return response.DbReadErr
	}
	RedisDataReload([]string{"SetRouteData"}...) //刷新路由信息
	return response.Success
}

//设置排序顺序
func SetRouteSort(routeId int, sort int16) error {
	fields := []string{"id"}
	filters := make([]interface{}, 2)
	filters[0] = "id"
	filters[1] = routeId
	model := &models.Routes{}
	err := model.Detail(fields, filters)
	if err != nil {
		return err
	}
	model.Sort = sort
	fields = []string{"sort"}
	if _, err := model.Create(model, fields); err != nil {

	}
	return nil
}

//根据用户id获取路由信息
func GetUserRouteById(userId int) ([]interface{}, error) {
	data := GetAdminLoginData(userId)
	arr := strings.Split(data.RouteIds, ",")
	routeList := make([]interface{}, len(arr))
	for k, v := range arr {
		id, _ := strconv.Atoi(v)
		routeList[k] = id
	}

	return routeList, nil
}
