package service

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"regexp"
	"strings"
	"tools/server/app/validate"
)

type MessageRes struct {
	Status bool
	Info   string
}

func SendMsg(data *MessageData) (res *MessageRes, err error) {
	mess := &validate.ValidateMessage{
		TplId:       data.TplId,
		ToUser:      data.ToUser,
		ExtraParams: data.ExtraParams,
	}

	//参数验证
	err = mess.ValidMessage()
	if err != nil {
		return nil, err
	}

	//查询模板消息
	notifyDetail, err := NotifyDetailById(mess.TplId)
	if err != nil {
		return nil, err
	}

	//模板消息状态
	notifyStatus := notifyDetail.Status
	if notifyStatus != 1 {
		return nil, errors.New("模板id状态异常")
	}

	//模板消息类型及内容
	notifyType := notifyDetail.NotifyType
	notifyContent := notifyDetail.TplData

	switch notifyType {
	case 1:
		err = sendEmail(data, notifyContent)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("模板消息类型不存在")
	}

	res = &MessageRes{
		Status: true,
		Info:   "发送成功",
	}

	return res, nil
}

/*
发送邮件
*/
func sendEmail(data *MessageData, notifyContent string) error {

	_, isExist := data.ExtraParams["title"]
	if !isExist {
		return errors.New("邮件标题不能为空")
	}

	//检查替换content
	content, err := checkExtraContent(notifyContent, data.ExtraParams)
	if err != nil {
		return err
	}

	username := beego.AppConfig.String("mail.username")
	password := beego.AppConfig.String("mail.password")
	host := beego.AppConfig.String("mail.host")
	port := beego.AppConfig.String("mail.port")

	mail := &validate.ValidateEmail{
		Title:    data.ExtraParams["title"],
		ToUser:   data.ToUser,
		Content:  content,
		FromUser: username,
		PassWord: password,
		Host:     host,
		Port:     port,
	}
	//参数验证
	err = mail.ValidEmail()
	if err != nil {
		return err
	}

	config := `{"username":` + mail.FromUser + `,"password":` + mail.PassWord + `,"host":` + mail.Host + `,"port":` + mail.Port + `}`
	email := utils.NewEMail(config)
	email.To = mail.ToUser
	email.From = mail.FromUser
	email.Subject = mail.Title
	email.Text = mail.Content
	err = email.Send()

	if err != nil {
		return err
	}

	return nil
}

//校验并替换消息内容
func checkExtraContent(content string, extra map[string]string) (string, error) {

	reg := regexp.MustCompile("{{.*?\\.DATA}}")
	arrContent := reg.FindAllString(content, -1)

	if len(arrContent) >0 {
		//校验内容中需替换字段是否为空
		f := func(c rune) bool {
			if c == '{' || c == '.' {
				return true
			} else {
				return false
			}
		}

		for _,v := range arrContent{
			result := strings.FieldsFunc(v, f)
			ekey := result[0]
			_, isExist := extra[ekey]
			if !isExist {
				return "", errors.New(ekey + "不能为空")
			}

			//替换字符串
			content = strings.Replace(content, v, extra[ekey], -1 )
		}

	}

	return content, nil
}
