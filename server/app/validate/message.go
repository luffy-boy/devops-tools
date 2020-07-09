package validate

import (
	"errors"
	"github.com/astaxie/beego/validation"
)

type ValidateMessage struct {
	TplId       int               `label:"模板id" valid:"Required;Numeric"` //模板id
	ToUser      []string          `label:"接收人" valid:"Required"`          //接收人
	ExtraParams map[string]string `label:"额外参数" valid:"Required"`         //额外参数
}

func (a *ValidateMessage) ValidMessage() error {
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
