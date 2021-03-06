package service

import (
	"context"
	"errors"
	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"strings"
	"time"
	"tools/server/app/common"
	"tools/server/app/models"
	"tools/server/app/response"
	"tools/server/app/util"
)

//Admin 管理员-请求数据
type AdminRequest struct {
	UserId       int    `json:"user_id"`      //用户id
	UserName     string `json:"username"`     //用户名
	PassWord     string `json:"password"`     //密码
	RealName     string `json:"real_name"`    //姓名
	Birthday     string `json:"birthday"`     //出生日期
	Status       int8   `json:"status"`       //状态
	IsDelete     int8   `json:"is_delete"`    //是否删除
	Phone        string `json:"phone"`        //手机号码
	Email        string `json:"email"`        //邮箱
	Sex          int8   `json:"sex"`          //性别
	RoleId       int    `json:"role_id"`      //角色id
	Ip           string `json:"ip"`           //ip地址
	Introduction string `json:"introduction"` //描述
	RouteIds     string `json:"route_ids"`    //路由id
}

type AdminParams struct {
	UserId   int
	Page     int
	Limit    int
	RealName string
	RoleId   int
	Sort     string
}

//Admin 用户登录校验
type AdminLoginValid struct {
	UserName string `valid:"Required;MinSize(5);MaxSize(32)"`
	Password string `valid:"Required;Length(32)"`
}

//Admin 用户数据校验--注册
type AdminDataValid struct {
	UserName string `valid:"Required;MinSize(5);MaxSize(32)"`
	Password string `valid:"Required;Length(32)"`
	RealName string `valid:"Required;MinSize(2);MaxSize(10)"`
	Birthday string `valid:"Required"`
	Phone    string `valid:"Required;Phone"`
	Email    string `valid:"Required;Email"`
	Sex      int8   `valid:"Required;Range(0,2)"`
	Status   int8   `valid:"Required;Range(0,2)"`
	RoleId   int    `valid:"Required"`
	RouteIds string `valid:"MaxSize(200)"`
}

//Admin 用户数据校验--修改
type AdminUpdateDataValid struct {
	UserId   int    `valid:"Required"`
	RealName string `valid:"Required;MinSize(2);MaxSize(10)"`
	Birthday string `valid:"Required"`
	Phone    string `valid:"Required;Phone"`
	Email    string `valid:"Required;Email"`
	Sex      int8   `valid:"Required;Range(0,2)"`
	RoleId   int    `valid:"Required"`
	RouteIds string `valid:"MaxSize(200)"`
}

//Admin 用户登录返回值
type AdminLoginData struct {
	UserId       int    `json:"user_id"`      //用户id
	RealName     string `json:"real_name"`    //用户名称
	Sex          int8   `json:"sex"`          //性別
	Token        string `json:"token"`        //秘钥
	Introduction string `json:"introduction"` //描述
	Menu         *Menu  `json:"menu"`         // 菜单
}

type AdminListData struct {
	Total int64        `json:"total"` //用户数量
	List  []*AdminData `json:"list"`  // 用户列表
}

type AdminData struct {
	Id       int    `json:"user_id"`
	RealName string `json:"real_name"` //用户名称
	Sex      int8   `json:"sex"`       //性別
	Status   int8   `json:"status"`    //状态
	Birthday int    `json:"birthday"`  //生日
	Phone    string `json:"phone"`     //手机
	RoleId   int    `json:"role_id"`   //角色id
	RoleName string `json:"role_name"` //角色名称
}

type AdminBrief struct {
	Id       int    `json:"user_id"`
	RealName string `json:"real_name"` //用户名称
}

type AdminDetail struct {
	UserId       int    `json:"user_id"`
	Username     string `json:"username"`
	RealName     string `json:"real_name"`
	Status       int8   `json:"status"`
	Birthday     string `json:"birthday"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Sex          int8   `json:"sex"`
	RoleId       int    `json:"role_id"`
	RouteIds     string `json:"route_ids"`
	Introduction string `json:"introduction"`
}

//首页统计数据
type IndexData struct {
	UserCount   DataCount `json:"user_count"`
	AuditTask   DataCount `json:"audit_task"`
	SuccessTask DataCount `json:"success_task"`
	FailTask    DataCount `json:"fail_task"`
}

//总用户数
type DataCount struct {
	ActualData []int64  `json:"actual_data"`
	DateData   []string `json:"date_data"`
	Total      int64    `json:"total"`
}

//管理员登录
func AdminLogin(params *AdminRequest) (*AdminLoginData, int, error) {
	var (
		LoginData *AdminLoginData
		filters   []interface{}
		routeIds  string
	)
	filters = append(filters, "username", params.UserName)
	filters = append(filters, "is_delete", 0)
	admin, err := models.GetDetail([]string{"id", "username", "password", "real_name", "salt", "status", "sex", "role_id", "route_ids"}, filters)
	if err != nil {
		return nil, response.DbQueryErr, err
	}
	if admin.Id == 0 {
		return nil, response.UserAccountErr, errors.New("用户不存在")
	}
	//校验密码
	if common.Md5(params.PassWord+admin.Salt) != admin.Password {
		return nil, response.UserAccountErr, errors.New("密码错误")
	}
	if admin.Status != 1 {
		return nil, response.UserStatusErr, errors.New("状态错误")
	}

	LoginData = &AdminLoginData{
		RealName:     admin.RealName,
		Sex:          admin.Sex,
		UserId:       admin.Id,
		Introduction: admin.Introduction,
	}

	adminJwt := &util.AdminJwtData{
		Id:       admin.Id,
		RealName: admin.RealName,
		Sex:      admin.Sex,
	}
	tokenExpired, _ := beego.AppConfig.Int("redis.tokenExpired")
	LoginData.Token, err = util.GenerateToken(adminJwt, tokenExpired)
	if err != nil {
		return nil, response.TokenCreateErr, err
	}

	ctime := time.Now().Unix()
	str := strconv.FormatInt(ctime, 10)
	admin.LoginIpAddr = params.Ip
	admin.LoginTime, _ = strconv.Atoi(str)

	fields := []string{"login_ip_addr", "login_time"}
	if _, err := admin.Create(admin, fields); err != nil {
		return nil, response.DbReadErr, err
	}

	//组装路由id集合
	RoleData := getRoleCacheInfoById(admin.RoleId, admin.Id)
	routeIds += RoleData.RoutesIds
	role := RoleData.Role
	if admin.RouteIds != "" {
		routeIds += "," + admin.RouteIds
	}
	routeIds = strings.Trim(routeIds, ",")
	admin.RouteIds = routeIds
	SetAdminLoginData(admin, role) //设置登录缓存信息

	//设置菜单
	if menu, err1 := MenuList(admin.Id); err1 == nil {
		LoginData.Menu = menu
	}
	return LoginData, 1, nil
}

//重新登录
func AdminReLogin(userId int, ip string, token string) (*AdminLoginData, int, error) {
	var (
		LoginData *AdminLoginData
		filters   []interface{}
		routeIds  string
	)
	filters = append(filters, "id", userId)
	filters = append(filters, "is_delete", 0)
	admin, err := models.GetDetail([]string{"id", "username", "password", "real_name", "salt", "status", "sex", "role_id", "route_ids"}, filters)
	if err != nil {
		return nil, response.DbQueryErr, err
	}
	if admin.Id == 0 {
		return nil, response.UserAccountErr, errors.New("用户不存在")
	}
	if admin.Status != 1 {
		return nil, response.UserStatusErr, errors.New("状态错误")
	}

	LoginData = &AdminLoginData{
		RealName:     admin.RealName,
		Sex:          admin.Sex,
		UserId:       admin.Id,
		Introduction: admin.Introduction,
	}

	tokenExpired, _ := beego.AppConfig.Int("redis.tokenExpired")
	LoginData.Token, err = util.RefreshToken(token, tokenExpired)
	if err != nil {
		return nil, response.TokenRefreshErr, err
	}

	admin.LoginIpAddr = ip
	admin.LoginTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))

	fields := []string{"login_ip_addr", "login_time"}
	if _, err := admin.Create(admin, fields); err != nil {
		return nil, response.DbReadErr, err
	}

	RoleData := getRoleCacheInfoById(admin.RoleId, admin.Id)

	routeIds += RoleData.RoutesIds
	role := RoleData.Role
	if admin.RouteIds != "" {
		routeIds += "," + admin.RouteIds
	}
	routeIds = strings.Trim(routeIds, ",")
	admin.RouteIds = routeIds
	SetAdminLoginData(admin, role) //设置登录缓存信息
	return LoginData, 1, nil
}

//刷新登陆信息
func RefreshLogin(userId int) error {
	var (
		filters  []interface{}
		routeIds string
	)

	userKey := util.RedisKeyList[util.CmsUserLoginKey] + strconv.Itoa(userId)
	if !util.RedisObj.Exists(userKey) {
		return errors.New("管理员登陆状态异常")
	}
	filters = append(filters, "id", userId)
	filters = append(filters, "is_delete", 0)
	admin, err := models.GetDetail([]string{"id", "username", "password", "real_name", "salt", "status", "sex", "role_id", "route_ids"}, filters)
	if err != nil {
		return err
	}
	if admin.Id == 0 {
		return errors.New("用户不存在")
	}
	if admin.Status != 1 {
		return errors.New("状态错误")
	}
	//组装路由id集合
	RoleData := getRoleCacheInfoById(admin.RoleId, userId)
	routeIds += RoleData.RoutesIds
	role := RoleData.Role
	if admin.RouteIds != "" {
		routeIds += "," + admin.RouteIds
	}
	routeIds = strings.Trim(routeIds, ",")
	admin.RouteIds = routeIds
	SetAdminLoginData(admin, role) //设置登录缓存信息
	return nil
}

//查询管理员列表
func AdminList(adminJ *AdminParams) (error, *AdminListData) {
	var (
		filters   []interface{}
		fields    []string
		adminList []*AdminData
	)

	orderBy := "-id"
	filters = append(filters, "id__gt", 1)
	filters = append(filters, "is_delete", 0)
	if adminJ.RealName != "" {
		filters = append(filters, "real_name__icontains", adminJ.RealName)
	}
	if adminJ.RoleId > 0 {
		filters = append(filters, "role_id", adminJ.RoleId)
	}
	if adminJ.Sort != "" {
		orderArr := map[string]string{"user_id": "id"}
		sort, ok := orderArr[adminJ.Sort[1:]]
		if !ok {
			sort = "id"
		}
		prefix := adminJ.Sort[0]
		switch prefix {
		case '+':
			orderBy = sort
		case '-':
			orderBy = "-" + sort
		}
	}
	fields = []string{"id", "real_name", "birthday", "phone", "status", "sex", "role_id", "route_ids"}
	list, total := models.GetAdminList(adminJ.Page, adminJ.Limit, filters, fields, orderBy)

	result, err := getRoleData()
	if err != nil {
		return err, nil
	}

	for _, v := range list {
		RoleName := ""
		RoleData, ok := result[v.RoleId]
		if ok {
			RoleName = RoleData.RoleName
		}
		arr := &AdminData{
			Id:       v.Id,
			RealName: v.RealName,
			Sex:      v.Sex,
			Status:   v.Status,
			Birthday: v.Birthday,
			Phone:    v.Phone,
			RoleId:   v.RoleId,
			RoleName: RoleName,
		}
		adminList = append(adminList, arr)
	}

	data := &AdminListData{
		Total: 0,
		List:  nil,
	}
	data.List = adminList
	data.Total = total
	return nil, data
}

//查询用户数据
func AdminDetailById(userId int) (AdminDetail, error) {
	var (
		filters []interface{}
	)
	filters = append(filters, "id", userId)
	admin, err := models.GetDetail([]string{"id", "username", "real_name", "birthday", "phone", "email", "status", "sex", "role_id", "route_ids", "introduction"}, filters)
	if err != nil {
		return AdminDetail{}, err
	}

	t := time.Unix(int64(admin.Birthday), 0)
	birthday := t.Format("2006-01-02 15:04:05")

	data := AdminDetail{
		UserId:       admin.Id,
		Username:     admin.Username,
		RealName:     admin.RealName,
		Status:       admin.Status,
		Birthday:     birthday,
		Phone:        admin.Phone,
		Email:        admin.Email,
		Sex:          admin.Sex,
		RoleId:       admin.RoleId,
		RouteIds:     admin.RouteIds,
		Introduction: admin.Introduction,
	}

	return data, nil
}

//创建或者新增
func AdminCreateOrUpdate(adminJ *AdminRequest) error {
	var (
		admin    *models.Admin
		err      error
		filters  []interface{}
		fields   []string
		salt     string
		password string
	)

	//查询用户是否重复
	if adminJ.UserId == 0 {
		filters = append(filters, "username", adminJ.UserName)
	} else {
		filters = append(filters, "id", adminJ.UserId)
	}
	admin, err = models.GetDetail([]string{"id", "username", "password", "real_name", "salt", "status", "sex", "role_id", "route_ids"}, filters)
	if err != nil {
		return err
	}
	if adminJ.UserId > 0 {
		if admin.Id != adminJ.UserId {
			return errors.New("用户不存在")
		}
		adminJ.UserName = ""
	} else {
		if admin.Id > 0 {
			return errors.New("用户已存在")
		}
	}

	if adminJ.PassWord != "" {
		if admin.Id == 0 {
			salt = common.NonceStr(4)
		} else {
			salt = admin.Salt
		}
		password = common.Md5(adminJ.PassWord + salt)
	}

	birthday := common.Strtotime(adminJ.Birthday, 2)
	str := strconv.FormatInt(birthday, 10)
	birthday1, _ := strconv.Atoi(str)

	admin = &models.Admin{
		Id:           adminJ.UserId,
		Username:     adminJ.UserName,
		Password:     password,
		RealName:     adminJ.RealName,
		Salt:         salt,
		Status:       adminJ.Status,
		IsDelete:     adminJ.IsDelete,
		Birthday:     birthday1,
		Phone:        adminJ.Phone,
		Email:        adminJ.Email,
		Sex:          adminJ.Sex,
		RoleId:       adminJ.RoleId,
		RouteIds:     adminJ.RouteIds,
		Introduction: adminJ.Introduction,
	}

	fields = []string{"real_name", "status", "is_delete", "birthday", "phone",
		"email", "sex", "role_id", "route_ids", "introduction", "utime"}
	if admin.Password != "" {
		fields = append(fields, "password")
	}
	if _, err := admin.Create(admin, fields); err != nil {
		return err
	}

	return nil
}

//设置更新管理员字段
func setAdminUpdateFields(fields []string, admin *models.Admin) []string {
	if admin.Password != "" {
		fields = append(fields, "password")
	}
	if admin.RealName != "" {
		fields = append(fields, "real_name")
	}
	if admin.Status >= 0 {
		fields = append(fields, "status")
	}
	if admin.IsDelete >= 0 {
		fields = append(fields, "is_delete")
	}
	if admin.Birthday != 0 {
		fields = append(fields, "birthday")
	}
	if admin.Phone != "" {
		fields = append(fields, "phone")
	}
	if admin.Email != "" {
		fields = append(fields, "email")
	}
	if admin.Sex >= 0 {
		fields = append(fields, "sex")
	}
	if admin.RoleId >= 0 {
		fields = append(fields, "role_id")
	}
	if admin.RouteIds != "" {
		fields = append(fields, "route_ids")
	}
	if admin.Introduction != "" {
		fields = append(fields, "introduction")
	}
	fields = append(fields, "utime")
	return fields
}

func GetIndexData() (IndexData, error) {

	timestamp := getNowZeroTime()

	//用户数
	uerr, userCount := countUser(timestamp)
	if uerr != nil {
		return IndexData{}, uerr
	}

	//待审核任务
	aerr, auditTask := countAuditTask(timestamp)
	if aerr != nil {
		return IndexData{}, aerr
	}

	//成功任务
	serr, successTask := countTaskStatus(timestamp, 1)
	if serr != nil {
		return IndexData{}, serr
	}

	//失败任务
	ferr, failTask := countTaskStatus(timestamp, 3)
	if ferr != nil {
		return IndexData{}, ferr
	}

	data := IndexData{
		UserCount:   userCount,
		AuditTask:   auditTask,
		SuccessTask: successTask,
		FailTask:    failTask,
	}

	return data, nil
}

func countUser(timestamp int64) (error, DataCount) {
	var (
		filters    []interface{}
		actualData []int64
		dateData   []string
		total      int64
		dateKey    map[string]int64
		ndate      string
		ntime      int64
	)
	//用户总数
	filters = append(filters, "is_delete", 0)
	data, err := models.CountUserAndTask(filters, "app_admin")
	if err != nil {
		return err, DataCount{}
	}
	for _, v := range data {
		total = v.Dcount
	}

	//最近七天的数据
	stime := timestamp - 6*86400
	list, gerr := models.GroupUserAndTask(stime, "app_admin")
	if gerr != nil {
		return gerr, DataCount{}
	}

	//日期map
	dateKey = map[string]int64{}
	for _, v := range list {
		dateKey[v.Cdate] = v.Dcount
	}

	for i := 6; i >= 0; i-- {
		ntime = timestamp - int64(i*86400)
		ndate = time.Unix(ntime, 0).Format("01-02")

		dateData = append(dateData, ndate)
		if _, ok := dateKey[ndate]; ok {
			actualData = append(actualData, dateKey[ndate])
		} else {
			actualData = append(actualData, 0)
		}
	}

	res := DataCount{
		ActualData: actualData,
		DateData:   dateData,
		Total:      total,
	}

	return nil, res
}

func countAuditTask(timestamp int64) (error, DataCount) {
	var (
		filters    []interface{}
		actualData []int64
		dateData   []string
		total      int64
		dateKey    map[string]int64
		ndate      string
		ntime      int64
	)
	//待审核任务总数
	filters = append(filters, "is_audit", 0)
	filters = append(filters, "is_delete", 0)
	data, err := models.CountUserAndTask(filters, "app_task")
	if err != nil {
		return err, DataCount{}
	}
	for _, v := range data {
		total = v.Dcount
	}

	//待审核任务最近七天的数据
	stime := timestamp - 7*86400
	list, gerr := models.GroupUserAndTask(stime, "app_task")
	if gerr != nil {
		return gerr, DataCount{}
	}

	//日期map
	dateKey = map[string]int64{}
	for _, v := range list {
		dateKey[v.Cdate] = v.Dcount
	}

	for i := 6; i >= 0; i-- {
		ntime = timestamp - int64(i*86400)
		ndate = time.Unix(ntime, 0).Format("01-02")

		dateData = append(dateData, ndate)
		if _, ok := dateKey[ndate]; ok {
			actualData = append(actualData, dateKey[ndate])
		} else {
			actualData = append(actualData, 0)
		}
	}

	res := DataCount{
		ActualData: actualData,
		DateData:   dateData,
		Total:      total,
	}

	return nil, res
}

func countTaskStatus(timestamp int64, status int8) (error, DataCount) {
	var (
		filters    interface{}
		actualData []int64
		dateData   []string
		total      int64
		dateKey    map[string]int64
		ndate      string
		ntime      int64
	)
	//任务总数
	filters = bson.D{{"status", status}}
	data, err := taskStatus(filters, 1)
	if err != nil {
		return err, DataCount{}
	}
	for _, v := range data {
		total = v.Dcount
	}

	//任务最近七天的数据
	timeString := "1970-01-01T08:00:00Z"
	t, _ := time.Parse("2006-01-02T15:04:05Z", timeString)
	match := bson.D{{"$match", bson.D{{"status", status}, {"ctime", bson.D{{"$gte", timestamp - 7*84600}}}}}}
	project := bson.D{
		{"$project",
			bson.D{
				{"cdate",
					bson.D{
						{"$dateToString",
							bson.D{
								{"format", "%m-%d"},
								{"date",
									bson.D{
										{"$add",
											bson.A{t,
												bson.D{
													{"$multiply", bson.A{"$ctime", 1000}}}}}}}}}}}}}}
	group1 := bson.D{{"$group", bson.D{{"_id", bson.D{{"cdate", "$cdate"}}}, {"dcount", bson.D{{"$sum", 1}}}}}}
	sort := bson.D{{"$sort", bson.D{{"_id.cdate", 1}}}}
	filters = mongo.Pipeline{match, project, group1, sort}
	list, gerr := taskStatus(filters, 2)
	if gerr != nil {
		return gerr, DataCount{}
	}

	//日期map
	dateKey = map[string]int64{}
	for _, v := range list {
		dateKey[v.Cdate] = v.Dcount
	}

	for i := 6; i >= 0; i-- {
		ntime = timestamp - int64(i*86400)
		ndate = time.Unix(ntime, 0).Format("01-02")
		dateData = append(dateData, ndate)
		if _, ok := dateKey[ndate]; ok {
			actualData = append(actualData, dateKey[ndate])
		} else {
			actualData = append(actualData, 0)
		}
	}

	res := DataCount{
		ActualData: actualData,
		DateData:   dateData,
		Total:      total,
	}

	return nil, res
}

func getNowZeroTime() int64 {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	day := time.Now().Format("02")
	datetime := year + "-" + month + "-" + day + " 00:00:00"
	timestamp := common.Strtotime(datetime, 2)

	return timestamp
}

func taskStatus(filters interface{}, qtype int8) ([]*models.IndexData, error) {

	countData := make([]*models.IndexData, 0)

	if qtype == 1 {
		//全部查询
		total, err := util.MgoClient.CountDocuments("log", "task_run_log", filters)
		if err != nil {
			return nil, err
		}
		countData = []*models.IndexData{{
			Cdate:  "",
			Dcount: total,
		}}
	} else {
		//聚合查询
		cursor, err := util.MgoClient.Aggregate("log", "task_run_log", filters, &options2.AggregateOptions{})
		if err != nil {
			return nil, err
		}

		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
		if err := cursor.Close(context.TODO()); err != nil {
			return nil, err
		}
		for _, result := range results {
			cid := result["_id"].(primitive.M)
			cdate := cid["cdate"]
			dcount := result["dcount"]
			countData = append(countData, &models.IndexData{Cdate: cdate.(string), Dcount: int64(dcount.(int32))})
		}
	}

	return countData, nil
}
