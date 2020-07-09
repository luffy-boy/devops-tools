package admin

import (
	"time"
	"tools/server/app/response"
	"tools/server/app/service"
)

//Crontab 时间查询
type CronController struct {
	BaseController
}

// @Title cron
// @Description 根据时间表达式 返回下次执行执行
// @Success 200 {object}  Res {"code":1,"data":null,"msg":"ok"}
// @Param   cron_spec     string true       "任务表达式"
// @Param   cal_run_num       string false      "执行次数 默认5次"
// @Failure 400 no enough input
// @Failure 500 server error
// @router /next_runtime [get]
func (c *CronController) NextRunTime() {
	cronSpec := c.GetString("cron_spec")
	calRunNum,_:= c.GetInt("cal_run_num",5)
	if cronSpec == ""{
		c.ResponseToJson(response.ParamsErr,nil)
	}

	sched, err := service.CronParse.Parse(cronSpec)
	if err != nil{
		c.ResponseToJson(response.DevOpsCronSpecErr, nil,"时间表达式错误")
		return
	}

	list := service.RunTimeList{}
	nextTime := time.Now()
	for i:=0; i < calRunNum ;i++  {
		nextTime = sched.Next(nextTime)
		t := time.Unix(nextTime.Unix(),0)
		date := t.Format("2006-01-02 15:04:05")
		list.List = append(list.List,service.RunData{RunData:date})
	}

	c.ResponseToJson(response.Success, list)
}
