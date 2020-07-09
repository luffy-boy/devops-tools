package common

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

//公共方法

/**
检查数组中是否存在某个值
*/
func InArray(needle interface{}, haystack []interface{}) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

//数组转字符串
func Implode(glue string, pieces []string) string {
	str := ""
	for _, v := range pieces {
		str += v + glue
	}
	str = strings.TrimRight(str, glue)
	return str
}

//生成MD5
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}

//生成指定长度随机字符串
func NonceStr(length int) string {
	var (
		chars string
		bytes []byte
		str   []byte
	)
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes = []byte(chars)
	str = []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		str = append(str, bytes[r.Intn(len(bytes))])
	}
	return string(str)
}

//记录debug日志
func DebugLog(fileName string, level string, data string) {
	var (
		log *logs.BeeLogger
	)
	log = logs.NewLogger(10000) // 创建一个日志记录器，参数为缓冲区的大小

	config := make(map[string]interface{})
	config["filename"] = "logs/" + fileName
	config["level"] = logs.LevelDebug
	//config["maxlines"] = 0
	//config["maxsize"] = 0
	//config["daily"] = 0
	//config["maxdays"] = 0
	//config["color"] = 0
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println("marshal failed, err:", err)
		return
	}
	log.SetLogger(logs.AdapterFile, string(configStr))
	log.SetLevel(logs.LevelDebug) // 设置日志写入缓冲区的等级
	log.EnableFuncCallDepth(true) // 输出log时能显示输出文件名和行号（非必须）

	switch level {
	case "Emergency":
		log.Emergency(data)
	case "Alert":
		log.Alert(data)
	case "Critical":
		log.Critical(data)
	case "Error":
		log.Error(data)
	case "Warning":
		log.Warning(data)
	case "Notice":
		log.Notice(data)
	case "Informational":
		log.Informational(data)
	case "Debug":
		log.Debug(data)
	default:
		log.Debug(data)
	}

	log.Flush() // 将日志从缓冲区读出，写入到文件
	defer log.Close()
}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string) (response string) {
	var ()
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

//发送POST请求
//url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
//content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

type youduMsg struct {
	Content     string   `json:"content"`     //推送内容
	ToUsers     []string `json:"toUsers"`     //推送管理员
	Application string   `json:"application"` //id
	Committer   string   `json:"committer"`   //标识
	NonceStr    string   `json:"nonceStr"`    //随机字符串
	Sign        string   `json:"sign"`        //签名
}

//发送有度短信
func SendYouDuMsg(content string, admins []string, committer string) bool {
	beego.AppConfig.String("youdu.msg.sms_domain")
	var (
		domain      string //有度短信发送域名
		nonceStr    string //随机字符串
		application string
		users       string //发送人员名单
		signStr     string //签名
		key         string //秘钥
		putData     youduMsg
	)

	users = Implode(",", admins)
	domain = beego.AppConfig.String("youdu.msg.sms_domain")
	nonceStr = NonceStr(16)
	application = beego.AppConfig.String("youdu.msg.application")
	key = beego.AppConfig.String("youdu.msg.key")

	committer = "域名报警"
	signStr = users + committer + content + key + nonceStr

	putData = youduMsg{
		Content:     content,
		ToUsers:     admins,
		Application: application,
		Committer:   committer,
		NonceStr:    nonceStr,
		Sign:        Md5(signStr),
	}

	result := Post(domain, putData, "application/json;charset=UTF-8")

	//记录日志
	date := time.Now().Format("2006-01-02")
	fileName := "common/sms/" + date + ".log"
	data := make(map[string]interface{})
	data["data"] = putData
	data["result"] = result
	lgoData, _ := json.Marshal(data)
	DebugLog(fileName, "Warning", string(lgoData))
	if result != "消息发送成功" {
		return false
	}
	return true
}

//日期转时间戳
func Strtotime(date string, formatType uint8) int64 {
	loc, _ := time.LoadLocation("Local")

	dateFormat := ""
	//设置格式化时间戳
	if formatType == 1 {
		dateFormat = "2006-01-02"
	} else {
		dateFormat = "2006-01-02 15:04:05"
	}
	tt, _ := time.ParseInLocation(dateFormat, date, loc)
	return tt.Unix()
}

func GbkAsUtf8(str string) string {
	srcDecoder := mahonia.NewDecoder("gbk")
	desDecoder := mahonia.NewDecoder("utf-8")
	resStr := srcDecoder.ConvertString(str)
	_, resBytes, _ := desDecoder.Translate([]byte(resStr), true)
	return string(resBytes)
}

//是否含义中文字符
func HasChineseChar(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Scripts["Han"], r) || (regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]").MatchString(string(r))) {
			return true
		}
	}
	return false
}

//判断数据是否存在  不存在:true 存在:false
func IsEmpty(value interface{}, array []interface{}) bool {
	if value == nil {
		return true
	}
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		v, ok := value.(int)
		if !ok {
			return true
		}
		if v > s.Len() {
			return true
		}
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return false
			}
		}
	}
	return true
}
