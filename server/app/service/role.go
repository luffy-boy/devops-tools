package service

import (
	"errors"
	"strconv"
	"strings"
	"tools/server/app/models"
	"tools/server/app/response"
)

//Role 角色-请求数据
type RoleRequest struct {
	RoleId   int    `json:"role_id"`                                //角色id
	Role     string `json:"role" valid:"Required;MaxSize(32)"`      //角色
	RoleName string `json:"role_name" valid:"Required;MaxSize(32)"` //角色名称
	ParentId int    `json:"parent_id"`                              //父级id
	Status   int8   `json:"status"`                                 //状态
	IsDelete int8   `json:"is_delete"`                              //删除状态
}

//查询
type RoleParams struct {
	Page  int
	Limit int
}

//权限编辑
type RouteEditParams struct {
	RoleId   int
	RouteIds string
}

//查询返回数据
type RoleListData struct {
	Total int64           `json:"total"` //角色数量
	List  []*RoleDataTree `json:"list"`  // 角色列表
}

type RoleDataTree struct {
	Id       int             `json:"role_id"`
	ParentId int             `json:"parent_id"` //父id
	RoleName string          `json:"role_name"` //角色名称
	Role     string          `json:"role"`      //角色别名
	Status   int8            `json:"status"`    //状态
	Ctime    int             `json:"ctime"`     //创建时间
	Child    []*RoleDataTree `json:"child"`     //子节点
}

type RoleData struct {
	Id       int    `json:"role_id"`
	ParentId int    `json:"parent_id"` //父id
	RoleName string `json:"role_name"` //角色名称
	Role     string `json:"role"`      //角色别名
	Status   int8   `json:"status"`    //状态
	Ctime    int    `json:"ctime"`     //创建时间
}

type RoleDetail struct {
	Id       int    `json:"role_id"`
	RoleName string `json:"role_name"`
	Role     string `json:"role"`
	ParentId int    `json:"parent_id"`
	RouteIds string `json:"route_ids"`
	Status   int8   `json:"status"`
	IsDelete int8   `json:"is_delete"`
	Ctime    int    `json:"ctime,omitempty"`
	Utime    int    `json:"utime,omitempty"`
}

//查询全部管理员角色
func RoleAll() ([]*RoleData, error) {
	var (
		filters []interface{}
		fields  []string
	)

	filters = append(filters, "is_delete", 0)
	filters = append(filters, "status", 1)
	fields = []string{"id", "role_name", "role", "parent_id", "status", "ctime"}
	list, _ := models.GetRolesList(0, 0, filters, fields)

	var roleDataList []*RoleData
	for _, v := range list {
		arr := &RoleData{
			Id:       v.Id,
			ParentId: v.ParentId,
			RoleName: v.RoleName,
			Role:     v.Role,
			Status:   v.Status,
			Ctime:    v.Ctime,
		}
		roleDataList = append(roleDataList, arr)
	}
	return roleDataList, nil
}

//查询管理员列表
func RoleList(roleJ *RoleParams) (*RoleListData, error) {
	var (
		filters []interface{}
		fields  []string
	)

	filters = append(filters, "is_delete", 0)
	fields = []string{"id", "role", "role_name", "parent_id", "status", "ctime"}
	list, total := models.GetRolesList(roleJ.Page, roleJ.Limit, filters, fields)

	var roleDataList []*RoleData
	for _, v := range list {
		arr := &RoleData{
			Id:       v.Id,
			ParentId: v.ParentId,
			RoleName: v.RoleName,
			Role:     v.Role,
			Status:   v.Status,
			Ctime:    v.Ctime,
		}
		roleDataList = append(roleDataList, arr)
	}

	arr := superRoleTree(roleDataList, 0)
	data := &RoleListData{
		Total: 0,
		List:  nil,
	}
	data.List = arr
	data.Total = total
	return data, nil
}

//递归实现
func superCategory(allRoleData []*RoleData, pid int) []*RoleData {
	var arr []*RoleData
	for _, v := range allRoleData {
		if pid == v.ParentId {
			arr = append(arr, v)
			childRoleData := superCategory(allRoleData, v.Id)
			arr = append(arr, childRoleData...)
		}
	}
	return arr
}

//递归实现(返回树状结果得数据)
func superRoleTree(allRoleData []*RoleData, pid int) []*RoleDataTree {
	var arr []*RoleDataTree
	for _, v := range allRoleData {
		if pid == v.ParentId {
			cTree := &RoleDataTree{
				Id:       v.Id,
				ParentId: v.ParentId,
				RoleName: v.RoleName,
				Role:     v.Role,
				Status:   v.Status,
				Ctime:    v.Ctime,
				Child:    nil,
			}
			chileRole := superRoleTree(allRoleData, v.Id)
			cTree.Child = chileRole
			arr = append(arr, cTree)
		}
	}
	return arr
}

//创建或者新增
func RoleCreateOrUpdate(roleJ *RoleRequest) error {
	var (
		role   *models.Roles
		fields []string
	)

	role = &models.Roles{
		Id:       roleJ.RoleId,
		Role:     roleJ.Role,
		RoleName: roleJ.RoleName,
		ParentId: roleJ.ParentId,
		Status:   roleJ.Status,
		IsDelete: roleJ.IsDelete,
	}

	fields = []string{"role", "role_name", "parent_id", "status", "route_ids"}
	if _, err := role.Create(role, fields); err != nil {
		return err
	}
	RedisDataReload([]string{"SetRoleData"}...) //刷新角色信息
	return nil
}

func RoleDetailById(roleId int) (*RoleDetail, error) {
	var (
		filters []interface{}
	)
	filters = append(filters, "id", roleId)
	model := &models.Roles{}
	err := model.Detail([]string{"id", "role_name", "route_ids", "role", "parent_id", "status", "is_delete"}, filters)
	if err != nil {
		return nil, err
	}

	data := &RoleDetail{
		Id:       model.Id,
		RoleName: model.RoleName,
		Role:     model.Role,
		ParentId: model.ParentId,
		RouteIds: model.RouteIds,
		Status:   model.Status,
		IsDelete: model.IsDelete,
		Ctime:    model.Ctime,
		Utime:    model.Utime,
	}

	return data, nil
}

func RoleDelete(roleId int) error {

	//查询是否有子角色

	var (
		filters []interface{}
	)
	filters = append(filters, "parent_id", roleId)
	filters = append(filters, "is_delete", 0)
	model := &models.Roles{}
	err := model.Detail([]string{"id"}, filters)
	if err != nil {
		return err
	}
	if model.Id > 0 {
		return errors.New("hasChild")
	}

	//执行删除
	var (
		fields []string
	)

	role := &models.Roles{
		Id:       roleId,
		IsDelete: 1,
	}
	fields = append(fields, "is_delete")
	if _, err := role.Create(role, fields); err != nil {
		return err
	}
	RedisDataReload([]string{"SetRoleData"}...) //刷新路由信息
	return nil
}

func RouteEdit(routeParams *RouteEditParams) int {
	var (
		fields []string
	)

	role := &models.Roles{
		Id:       routeParams.RoleId,
		RouteIds: routeParams.RouteIds,
	}
	fields = append(fields, "route_ids")
	if _, err := role.Create(role, fields); err != nil {
		return response.DbReadErr
	}
	RedisDataReload([]string{"SetRoleData"}...) //刷新路由信息
	return response.Success
}

//获取角色组 路由id
func getRoleCacheInfoById(roleId, adminId int) (RoleData *RoleCache) {
	var (
		ok bool
	)
	if adminId == 1 {
		list, _ := GetRouteData()
		routesIds := ""
		if list != nil {
			for k, _ := range list {
				routesIds += strconv.Itoa(k) + ","
			}
			routesIds = strings.TrimRight(routesIds, ",")
		}
		return &RoleCache{
			Id:        0,
			Role:      "Administrators",
			RoleName:  "超级管理员",
			ParentId:  0,
			RoutesIds: routesIds,
		}
	}
	//组装路由id集合
	roleList, err := getRoleData()
	if err != nil {
		return
	}
	if RoleData, ok = roleList[roleId]; !ok {
		return &RoleCache{}
	}
	return
}
