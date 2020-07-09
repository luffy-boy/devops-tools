package admin

import (
	"tools/server/app/service"
)

/*
消息发送
*/
func SendMsg(data *service.MessageData) *service.MessageRes {
	var res  = &service.MessageRes{
		Status : false,
		Info : "",
	}
	sendRes, err := service.SendMsg(data)
	if err != nil{
		res.Status = false
		res.Info = err.Error()
	}

	if sendRes != nil {
		res = sendRes
	}

	return res
}
