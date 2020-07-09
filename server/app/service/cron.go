package service

import "tools/server/app/util/robfig/cron"

var (
	options = cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month |cron.Dow
	CronParse = cron.NewParser(options)
)

//执行时间列表
type RunTimeList struct {
	List  []RunData `json:"list"`  //执行列表
}

type RunData struct {
	RunData       string     `json:"run_date"`
}
