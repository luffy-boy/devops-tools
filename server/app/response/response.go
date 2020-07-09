package response

//成功
const (
	Success = 1
)

//系统级别错误
const (
	SystemErr              = iota + 10000 //系统繁忙
	DbReadErr                             //10001  数据写入失败
	DbQueryErr                            //10002  数据查询失败
	DbDeleteErr                           //10003  数据删除失败
	UserAccessForbiddenErr                //10004  权限不足  访问被禁止
	TokenMessingErr                       //10005  缺失token
	TokenInvalidErr                       //10006  token非法
	TokenExpiredErr                       //10007  token过期
	TokenCreateErr                        //10008  Token创建失败
	TokenRefreshErr                       //10009  Token刷新失败
	ParamsErr                             //10010  参数错误
	LogoutError                           //10011  退出登录失败
	RedisReadErr                          //10012  Redis数据写入失败
	RedisQueryErr                         //10013  Redis数据读取失败
	MgoQueryErr                           //10014  Mgo数据查询失败
)

//服务级错误
const (
	ServiceErr = 200000000
)

//服务级错误代码
const (
	AdminCmsApiErr = ServiceErr + 1000000 //后台Api模块错误码  201000000
)

//模块错误码
const (
	AdminAuthErr   = AdminCmsApiErr + 1000 //201001000 权限管理模块
	AdminDevOpsErr = AdminCmsApiErr + 2000 //201002000 运维模块错误码
)

//权限管理模块
const (
	AdminErr         = AdminAuthErr + iota
	UserStatusErr    //201001001   用户状态被禁止
	UserAccountErr   //201001002   账号不存在或者密码错误
	RoleHasChild     //201001003   角色拥有下级权限，无法删除
	RouteHasChildErr //201001004   路由拥有下级配置，无法删除
)

//服务器运维
const (
	DevOpsErr          = AdminDevOpsErr + iota
	DevOpsTaskAuditErr //201002001   任务审核失败
	DevOpsTaskRunErr   //201002002   任务执行失败
	DevOpsTaskStartErr //201002003   任务启动失败
	DevOpsCronSpecErr  //201002004   时间表达式错误
	DevBanHasChinese   //201002005   命令含有中文字符，操作失败
	DevIsBan           //201002006   为禁用命令
)

var ResponseList = map[int]map[string]interface{}{
	//成功
	Success: {"code": Success, "info": "ok"},

	//系统类
	SystemErr:              {"code": SystemErr, "info": "系统繁忙"},
	DbReadErr:              {"code": DbReadErr, "info": "暂停服务"},
	DbQueryErr:             {"code": DbQueryErr, "info": "暂停服务"},
	DbDeleteErr:            {"code": DbDeleteErr, "info": "暂停服务"},
	RedisReadErr:           {"code": RedisReadErr, "info": "暂停服务"},
	RedisQueryErr:          {"code": RedisQueryErr, "info": "暂停服务"},
	TokenMessingErr:        {"code": TokenMessingErr, "info": "Token验证失败，参数为空"},
	TokenInvalidErr:        {"code": TokenInvalidErr, "info": "Token非法~"},
	TokenExpiredErr:        {"code": TokenExpiredErr, "info": "Token过期~"},
	UserAccessForbiddenErr: {"code": UserAccessForbiddenErr, "info": "权限不足"},
	ParamsErr:              {"code": ParamsErr, "info": "参数错误"},
	LogoutError:            {"code": LogoutError, "info": "退出登录失败"},
	MgoQueryErr:            {"code": MgoQueryErr, "info": "数据查询失败"},

	//权限管理
	UserStatusErr:    {"code": UserStatusErr, "info": "该用户已被禁止登录"},
	UserAccountErr:   {"code": UserAccountErr, "info": "账号不存在或者密码错误"},
	RoleHasChild:     {"code": RoleHasChild, "info": "角色拥有下级权限，无法删除"},
	RouteHasChildErr: {"code": RoleHasChild, "info": "路由拥有下级配置，无法删除"},

	//运维
	DevOpsTaskAuditErr: {"code": DevOpsTaskAuditErr, "info": "任务审核失败"},
	DevOpsTaskRunErr:   {"code": DevOpsTaskRunErr, "info": "任务执行失败"},
	DevOpsTaskStartErr: {"code": DevOpsTaskStartErr, "info": "任务运行状态修改失败"},
	DevOpsCronSpecErr:  {"code": DevOpsCronSpecErr, "info": "时间表达式错误"},
	DevBanHasChinese:   {"code": DevBanHasChinese, "info": "命令含有中文字符，操作失败"},
	DevIsBan:           {"code": DevIsBan, "info": "含有禁用命令"},
}

type Json struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
