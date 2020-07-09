package service

import (
	"tools/server/app/models"
	"tools/server/app/util"
)

func Init() {
	models.Init()
	util.RedisInit()        //初始化redis
	util.MongoInit()        //初始化mongodb
	RedisDataReload()       //加载所有缓存数据
	InitJobs()              //初始化任务模块
	//validate.InitValidate() //初始化验证模块
}
