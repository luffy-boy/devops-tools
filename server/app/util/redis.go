package util

import (
	"errors"
	"github.com/astaxie/beego"
	"time"
	"tools/server/app/util/redigo/redis"
)

type RedisOp struct {
	p *redis.Pool
}

var (
	RedisObj = &RedisOp{}
	err      error
)

const (
	CmsRoleKey      = iota //0 后台角色 key
	CmsRouteKey            // 1 后台路由 key
	CmsUserLoginKey        // 2 后台用户登录 key
	JobTaskRunLog          // 4 定时任务执行日志
	MessageQueue           // 4 消息通知和队列
)

//redis key管理集合
var RedisKeyList = map[int]string{
	//成功
	CmsRoleKey:      "CmsRoleInfo",
	CmsRouteKey:     "CmsRouteInfo",
	CmsUserLoginKey: "CmsUserLoginInfo_",
	JobTaskRunLog:   "JobTaskRunLog",
	MessageQueue:    "MessageQueue",
}

//初始化连接redis
func RedisInit() {
	//redisHost := beego.AppConfig.String("redis.host")
	redisConn := beego.AppConfig.String("redis.conn")
	redisDbNum := beego.AppConfig.String("redis.redisDbNum")
	redisPassword := beego.AppConfig.String("redis.password")

	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", redisConn)
		if err != nil {
			return nil, err
		}

		if redisPassword != "" {
			if _, err := c.Do("AUTH", redisPassword); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", redisDbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	RedisObj.p = &redis.Pool{
		MaxIdle:     256,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}

func (r *RedisOp) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	c := r.p.Get()
	if err := c.Err(); err != nil {
		return nil, err
	}
	defer c.Close()

	return c.Do(commandName, args...)
}

func (c *RedisOp) Set(key string, value interface{}) error {
	if _, err := c.do("SET", key, value); err == nil {
		return nil
	}
	return err
}

func (c *RedisOp) Get(key string) (result interface{}, err error) {
	if result, err = c.do("GET", key); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *RedisOp) Del(key string) error {
	if _, err = c.do("Del", key); err != nil {
		return err
	}
	return nil
}

func (c *RedisOp) SetEx(key string, timeout time.Duration, value interface{}) error {
	if _, err := c.do("SETEX", key, int64(timeout/time.Second), value); err != nil {
		return err
	}
	return nil
}

func (c *RedisOp) LPush(key string, value interface{}) error {
	if _, err := c.do("LPUSH", key, value); err != nil {
		return err
	}
	return nil
}

func (c *RedisOp) RPop(key string) (result interface{}, err error) {
	if result, err = c.do("RPOP", key); err != nil {
		return nil, err
	}
	return result, nil
}

//返回key是否存在
func (c *RedisOp) Exists(key string) bool {
	result, err := c.do("Exists", key)
	result1 := result.(int64)
	if err != nil || result1 == 0 {
		return false
	}
	return true
}

//设置过期时间 单位秒
func (c *RedisOp) Expire(key string, expired int) bool {
	result, err := c.do("EXPIRE", key, expired)
	result1 := result.(int64)
	if err != nil || result1 == 0 {
		return false
	}
	return true
}

//返回过期时间  单位:秒    -2不存在或已过期 -1永不过期
func (c *RedisOp) Ttl(key string) int64 {
	result, err := c.do("TTL", key)
	expired := result.(int64)
	if err != nil {
		return 0
	}
	return expired
}
