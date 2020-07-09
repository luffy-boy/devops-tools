package validate

import (
	"errors"
	"github.com/astaxie/beego/validation"
)

type ValidateEmail struct {
	Title    string   `label:"邮件标题" valid:"Required"`             //邮件标题
	Content  string   `label:"邮件内容" valid:"Required"`             //邮件内容
	ToUser   []string `label:"接收人" valid:"Required"`              //接收人
	FromUser string   `label:"发送人" valid:"Required"`              //发送人
	PassWord string   `label:"密码" valid:"Required"`               //密码
	Host     string   `label:"SMTP服务器地址" valid:"Required"`        //SMTP服务器地址
	Port     string   `label:"SMTP服务端口" valid:"Required;Numeric"` //SMTP服务端口
}

func (a *ValidateEmail) ValidEmail() error {
	valid := validation.Validation{}
	b, err := valid.Valid(a)
	if err != nil {
		return err
	}
	if !b {
		for _, err := range valid.Errors {
			msg := err.Key + ": " + err.Message
			return errors.New(msg)
		}
	}
	return nil
}
