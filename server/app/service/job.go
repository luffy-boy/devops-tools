package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	gote "github.com/linxiaozhi/go-telnet"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
	"tools/server/app/common"
	"tools/server/app/models"
	"tools/server/app/util"
	"tools/server/app/util/robfig/cron"
)

type Job struct {
	Id         int                            // taskID
	LogId      interface{}                    // 日志记录ID
	ServerId   int                            // 执行器信息
	ServerName string                         // 执行器名称
	ServerType int8                           // 执行器类型， 0-ssh 1-telnet 2-agent
	Name       string                         // 任务名称
	Task       *models.Task                   // 任务对象
	RunFunc    func(time.Duration) *JobResult // 执行函数
	Status     int                            // 任务状态，大于0表示正在执行中
	Concurrent bool                           // 同一个任务是否允许并行执行
}

type JobResult struct {
	OutMsg    string
	ErrMsg    string
	IsOk      bool
	IsTimeout bool
}

type RpcResult struct {
	Status  int
	Message string
}

//调度计数器
var (
	Counter    sync.Map
	workRunNum chan bool
	mainCron   *cron.Cron
	lock       sync.Mutex
)

func GetCounter(key string) int {
	if v, ok := Counter.LoadOrStore(key, 0); ok {
		n := v.(int)
		return n
	}
	return 0
}

func SetCounter(key string) {
	if v, ok := Counter.Load(key); ok {
		n := v.(int)
		m := n + 1
		if n > 1000 {
			m = 0
		}
		Counter.Store(key, m)
	}
}

func AddJob(spec string, job *Job) bool {
	lock.Lock()
	defer lock.Unlock()

	if GetEntryById(job.Task.Id) {
		return false
	}
	enrtyId, err := mainCron.AddJob(spec, job)
	if err != nil {
		beego.Error("AddJob: ", err.Error())
		return false
	}
	job.Task.JobEntryId = int(enrtyId)
	job.Task.Create([]string{"job_entry_id"})
	return true
}

//增加定时方法
func AddJobFunc() {
	mainCron.AddFunc("0 */5 * * * *", TaskLogGather)
}

//初始化定时任务
func InitJobs() {
	var (
		filters []interface{}
		fields  []string
	)

	mainCron = cron.New(cron.WithSeconds())
	mainCron.Start()

	AddJobFunc()
	if size, _ := beego.AppConfig.Int("jobs.runnum"); size > 0 {
		workRunNum = make(chan bool, size)
	}

	filters = append(filters, "is_delete", 0)
	filters = append(filters, "is_audit", 1)
	filters = append(filters, "status", 1)
	fields = []string{"id", "group_id", "server_ids", "run_type", "description", "cron_spec", "concurrent", "command", "timeout", "is_notify", "notify_type", "notify_tpl_id", "notify_user_ids", "status"}
	list, _ := models.GetTaskList(1, 10000, filters, fields)
	for _, task := range list {
		jobs, err := NewJobFromTask(task)
		if err != nil {
			continue
		}

		for _, job := range jobs {
			AddJob(task.CronSpec, job)
		}
	}
}

func GetEntryById(taskId int) bool {
	entries := mainCron.Entries()
	for _, e := range entries {
		if v, ok := e.Job.(*Job); ok {
			if v.Task.Id == taskId {
				return true
			}
		}
	}
	return false
}

//创建任务
func NewJobFromTask(task *models.Task) ([]*Job, error) {
	if task.Id < 1 {
		return nil, errors.New("缺少任务id")
	}

	if task.ServerIds == "" {
		return nil, errors.New("未配置执行服务器")
	}

	TaskServerIdsArr := strings.Split(task.ServerIds, ",")
	jobArr := make([]*Job, 0)
	for _, serverId := range TaskServerIdsArr {
		if serverId == "0" {
			//本地执行
			job := NewCommandJob(task.Id, 0, task.TaskName, task.Command)
			job.Task = task
			job.Concurrent = false
			if task.Concurrent == 1 {
				job.Concurrent = true
			}
			job.ServerId = 0
			job.ServerName = "本地服务器"
			jobArr = append(jobArr, job)
		} else {
			//远程执行
			serverId, _ := strconv.Atoi(serverId)
			var filters []interface{}
			filters = append(filters, "id", serverId)
			filters = append(filters, "status__in", []int{0, 1})
			fields := []string{"id", "group_id", "connection_type", "server_name", "server_ip", "server_account", "port", "password", "private_key_src", "public_key_src", "type", "detail", "status"}
			server := &models.TaskServers{}
			err := server.GetDetail(fields, filters)
			if err != nil {
				return nil, errors.New(err.Error())
			}
			if server.ConnectionType == 0 {
				if server.Type == 0 {
					//密码验证登录服务器
					job := RemoteCommandJobByPassword(task.Id, serverId, task.TaskName, task.Command, server)
					job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					job.ServerId = serverId
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				} else {
					job := RemoteCommandJob(task.Id, serverId, task.TaskName, task.Command, server)
					job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					job.ServerId = serverId
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				}
			} else if server.ConnectionType == 1 {
				if server.Type == 0 {
					//密码验证登录服务器
					job := RemoteCommandJobByTelnetPassword(task.Id, serverId, task.TaskName, task.Command, server)
					job.Task = task
					job.Concurrent = false
					if task.Concurrent == 1 {
						job.Concurrent = true
					}
					//job.Concurrent = task.Concurrent == 1
					job.ServerId = serverId
					job.ServerName = server.ServerName
					jobArr = append(jobArr, job)
				}
			} else if server.ConnectionType == 2 {
				//密码验证登录服务器
				job := RemoteCommandJobByAgentPassword(task.Id, serverId, task.TaskName, task.Command, server)
				job.Task = task
				job.Concurrent = false
				if task.Concurrent == 1 {
					job.Concurrent = true
				}
				//job.Concurrent = task.Concurrent == 1
				job.ServerId = serverId
				job.ServerName = server.ServerName
				jobArr = append(jobArr, job)

			}
		}
	}

	return jobArr, nil
}
func (j *Job) agentRun() (reply *JobResult) {

	var filters []interface{}
	filters = append(filters, "id", j.ServerId)
	fields := []string{"id", "status"}
	server := &models.TaskServers{}
	err := server.GetDetail(fields, filters)
	if err != nil {
		return
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.ServerIp, server.Port))
	reply = new(JobResult)
	if err != nil {
		logs.Error("Net error:", err)
		reply.IsOk = false
		reply.ErrMsg = "Net error:" + err.Error()
		reply.IsTimeout = false
		reply.OutMsg = ""
		return reply
	}

	defer conn.Close()
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	defer client.Close()
	reply = new(JobResult)

	task := j.Task
	err = client.Call("RpcTask.RunTask", task, &reply)
	if err != nil {
		reply.IsOk = false
		reply.ErrMsg = "Net error:" + err.Error()
		reply.IsTimeout = false
		reply.OutMsg = ""
		return reply
	}
	return
}

//测试服务器
func TestServer(server *models.TaskServers) error {
	if server.ConnectionType == 0 {
		switch server.Type {
		case 0:
			//密码登录
			return RemoteCommandByPassword(server)
		case 1:
			//密钥登录
			return RemoteCommandByKey(server)
		default:
			return errors.New("未知的登录方式")

		}
	} else if server.ConnectionType == 1 {
		if server.Type == 0 {
			//密码登录]
			return RemoteCommandByTelnetPassword(server)
		} else {
			return errors.New("Telnet方式暂不支持密钥登陆！")
		}

	} else if server.ConnectionType == 2 {
		return RemoteAgent(server)
	}

	return errors.New("未知错误")
}

func RunServer(j *Job) bool {
	//判断是否是当前执行器执行
	TaskServerIdsArr := strings.Split(j.Task.ServerIds, ",")
	num := len(TaskServerIdsArr)

	if num == 0 {
		return false
	}

	count := GetCounter(strconv.Itoa(j.Task.Id))
	index := count % num
	runServerId, _ := strconv.Atoi(TaskServerIdsArr[index])

	if j.ServerId != runServerId {
		return false
	}

	//本地服务器
	if runServerId == 0 {
		return true
	}

	//判断执行器或者服务器是否存活
	var filters []interface{}
	filters = append(filters, "id", runServerId)
	fields := []string{"id", "group_id", "connection_type", "server_name", "server_ip", "server_account", "port", "password", "private_key_src", "public_key_src", "type", "detail", "status"}
	server := &models.TaskServers{}
	err := server.GetDetail(fields, filters)
	if err != nil {
		return false
	}

	if server.Status != 1 {
		return false
	}

	if err := TestServer(server); err != nil {
		server.Status = 0
		server.Create([]string{"status"})
		return false
	}

	return true

}

func (j *Job) Run() {
	//执行策略 轮询
	if j.Task.RunType == 1 {
		if !RunServer(j) {
			return
		} else {
			SetCounter(strconv.Itoa(j.Task.Id))
		}
	}

	if !j.Concurrent && j.Status > 0 {
		fileName := "task.log"
		common.DebugLog(fileName, "Warning", fmt.Sprintf("任务[%d]上一次执行尚未结束，本次被忽略。", j.Task.Id))
		return
	}

	fileName := "task.log"
	defer func() {
		if err := recover(); err != nil {
			common.DebugLog(fileName, "Emergency", string(debug.Stack()))
		}
	}()

	if workRunNum != nil {
		workRunNum <- true
		defer func() {
			<-workRunNum
		}()
	}

	common.DebugLog(fileName, "Debug", fmt.Sprintf("开始执行任务: %d", j.Task.Id))

	j.Status++
	defer func() {
		j.Status--
	}()

	t := time.Now()
	timeout := time.Minute * 10
	if j.Task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.Task.Timeout)
	}

	var jobResult = new(JobResult)
	//anget
	if j.ServerType == 2 {
		jobResult = j.agentRun()
	} else {
		jobResult = j.RunFunc(timeout)
	}

	//任务执行耗时
	ut := time.Now().Sub(t) / time.Millisecond

	var logStatus int8
	logStatus = 1
	if jobResult.IsTimeout {
		logStatus = 2
	} else if !jobResult.IsOk {
		logStatus = 3
	}
	//插入日志
	log := TaskRunLog{
		TaskId:      j.Id,
		ServerId:    j.ServerId,
		ServerName:  j.ServerName,
		Output:      jobResult.OutMsg,
		Error:       jobResult.ErrMsg,
		Status:      logStatus,
		ProcessTime: int(ut),
		Ctime:       int(t.Unix()),
	}
	if err := PushTakRunLog(&log); err != nil {
		//写入错误日志
		common.DebugLog("task.log", "Error", err.Error())
	}

	if log.Status < 0 && j.Task.IsNotify == 1 {
		var toUser []string
		toUser = strings.Split(j.Task.NotifyUserIds, ",")
		//写入消息通知队列
		extraParams := make(map[string]string)
		extraParams["task_name"] = j.Task.TaskName
		extraParams["server_id"] = strconv.Itoa(j.ServerId)
		extraParams["command"] = j.Task.Command
		extraParams["execute_time"] = beego.Date(time.Unix(t.Unix(), 0), "Y-m-d H:i:s")
		extraParams["out_msg"] = jobResult.OutMsg
		extraParams["err_msg"] = jobResult.ErrMsg
		message := MessageData{
			TplId:       j.Task.NotifyTplId,
			ToUser:      toUser,
			ExtraParams: extraParams,
		}
		if err := PushMessage(&message); err != nil {
			//写入错误日志
			common.DebugLog("message.log", "Error", err.Error())
		}
	}

	// 更新上次执行时间
	dateTime := time.Now().Unix()
	sTime := strconv.FormatInt(dateTime, 10)
	j.Task.PrevTime, _ = strconv.Atoi(sTime)
	j.Task.ExecuteTimes++
	if _, err := j.Task.Create([]string{"execute_times", "prev_time"}); err != nil {
		//错误log日志
		common.DebugLog(fileName, "Error", err.Error())
	}
}

//生成一个指定任务
func NewCommandJob(id int, serverId int, name string, command string) *Job {
	job := &Job{
		Id:   id,
		Name: name,
	}

	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		//cmd := exec.Command("/bin/bash", "-c", command)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("CMD", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		cmd.Start()
		err, isTimeout := runCmdWithTimeout(cmd, timeout)
		jobresult = new(JobResult)
		jobresult.OutMsg = bufOut.String()
		jobresult.ErrMsg = bufErr.String()

		jobresult.IsOk = true
		if err != nil {
			jobresult.IsOk = false
		}

		jobresult.IsTimeout = isTimeout

		return jobresult
	}
	return job
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration) (error, bool) {
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	var err error
	select {
	case <-time.After(timeout):
		beego.Warn(fmt.Sprintf("任务执行时间超过%d秒，进程将被强制杀掉: %d", int(timeout/time.Second), cmd.Process.Pid))
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil {
			beego.Error(fmt.Sprintf("进程无法杀掉: %d, 错误信息: %s", cmd.Process.Pid, err))
		}
		return err, true
	case err = <-done:
		return err, false
	}
}

//远程执行任务 密钥验证
func RemoteCommandJob(id int, serverId int, name string, command string, servers *models.TaskServers) *Job {
	job := &Job{
		Id:       id,
		Name:     name,
		ServerId: serverId,
	}

	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false

		key, err := ioutil.ReadFile(servers.PrivateKeySrc)
		if err != nil {
			jobresult.IsOk = false
			return
		}
		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			jobresult.IsOk = false
			return
		}
		addr := fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)
		config := &ssh.ClientConfig{
			User: servers.ServerAccount,
			Auth: []ssh.AuthMethod{
				// Use the PublicKeys method for remote authentication.
				ssh.PublicKeys(signer),
			},
			//HostKeyCallback: ssh.FixedHostKey(hostKey),
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		// Connect to the remote server and perform the SSH handshake.47.93.220.5
		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			jobresult.IsOk = false
			return
		}

		defer client.Close()

		session, err := client.NewSession()
		if err != nil {
			jobresult.IsOk = false
			return
		}
		defer session.Close()

		// Once a Session is created, you can execute a single command on
		// the remote side using the Run method.

		var b bytes.Buffer
		var c bytes.Buffer
		session.Stdout = &b
		session.Stderr = &c

		//session.Output(command)
		if err := session.Run(command); err != nil {
			jobresult.IsOk = false
			return
		}
		jobresult.OutMsg = b.String()
		jobresult.ErrMsg = c.String()
		jobresult.IsOk = true
		jobresult.IsTimeout = false
		return
	}
	return job
}

func RemoteCommandJobByPassword(id int, serverId int, name string, command string, servers *models.TaskServers) *Job {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	job := &Job{
		Id:         id,
		Name:       name,
		ServerId:   serverId,
		ServerType: servers.ConnectionType,
	}
	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false

		// get auth method
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(servers.Password))

		clientConfig = &ssh.ClientConfig{
			User: servers.ServerAccount,
			Auth: auth,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			//Timeout: 1000 * time.Second,
		}

		// connet to ssh
		addr = fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)

		if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
			jobresult.IsOk = false
			return
		}
		defer client.Close()
		// create session
		if session, err = client.NewSession(); err != nil {
			jobresult.IsOk = false
			return
		}

		var b bytes.Buffer
		var c bytes.Buffer
		session.Stdout = &b
		session.Stderr = &c
		//session.Output(command)
		if err := session.Run(command); err != nil {
			jobresult.IsOk = false
			return
		}
		jobresult.OutMsg = b.String()
		jobresult.ErrMsg = c.String()
		jobresult.IsOk = true
		jobresult.IsTimeout = false
		return
	}

	return job
}

func RemoteCommandJobByTelnetPassword(id int, serverId int, name string, command string, servers *models.TaskServers) *Job {

	job := &Job{
		Id:       id,
		Name:     name,
		ServerId: serverId,
	}

	job.RunFunc = func(timeout time.Duration) (jobresult *JobResult) {
		jobresult = new(JobResult)
		jobresult.OutMsg = ""
		jobresult.ErrMsg = ""
		jobresult.IsTimeout = false

		addr := fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)
		conn, err := gote.DialTimeout("tcp", addr, timeout)
		if err != nil {
			jobresult.IsOk = false
			return
		}

		defer conn.Close()

		buf := make([]byte, 4096)

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			return
		}

		if _, err = conn.Write([]byte(servers.ServerAccount + "\r\n")); err != nil {
			jobresult.IsOk = false
			return
		}

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			return
		}

		if _, err = conn.Write([]byte(servers.Password + "\r\n")); err != nil {
			jobresult.IsOk = false
			return
		}

		if _, err = conn.Read(buf); err != nil {
			jobresult.IsOk = false
			return
		}

		loginStr := common.GbkAsUtf8(string(buf[:]))
		if !strings.Contains(loginStr, ">") {
			jobresult.ErrMsg = jobresult.ErrMsg + "Login failed!"
			jobresult.IsOk = false
			return
		}

		commandArr := strings.Split(command, "\n")

		out, n := "", 0
		for _, c := range commandArr {
			_, err = conn.Write([]byte(c + "\r\n"))
			if err != nil {
				jobresult.IsOk = false
				return
			}

			n, err = conn.Read(buf)

			out = out + common.GbkAsUtf8(string(buf[0:n]))
			if err != nil ||
				strings.Contains(out, "'"+c+"' is not recognized as an internal or external command") ||
				strings.Contains(out, "'"+c+"' 不是内部或外部命令，也不是可运行的程序") {
				jobresult.ErrMsg = jobresult.ErrMsg + " " + common.GbkAsUtf8(string(buf[0:n]))
				jobresult.IsOk = false
				jobresult.OutMsg = out
				return
			}
		}
		jobresult.IsOk = true
		jobresult.OutMsg = out
		return
	}

	return job
}

func RemoteCommandJobByAgentPassword(id int, serverId int, name string, command string, servers *models.TaskServers) *Job {

	job := &Job{
		Id:         id,
		Name:       name,
		ServerType: servers.ConnectionType,
	}

	job.RunFunc = func(timeout time.Duration) *JobResult {
		return new(JobResult)
	}
	return job

}

func RemoteCommandByTelnetPassword(servers *models.TaskServers) error {

	addr := fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)
	conn, err := gote.DialTimeout("tcp", addr, time.Second*10)
	if err != nil {
		return err
	}

	defer conn.Close()

	buf := make([]byte, 4096)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(servers.ServerAccount + "\r\n"))
	if err != nil {
		return err
	}

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(servers.Password + "\r\n"))
	if err != nil {
		return err
	}

	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	str := common.GbkAsUtf8(string(buf[:]))

	if strings.Contains(str, ">") {
		return nil
	}

	return errors.New("连接失败!")
}

func RemoteCommandByPassword(servers *models.TaskServers) error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(servers.Password))

	clientConfig = &ssh.ClientConfig{
		User: servers.ServerAccount,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	addr = fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)
	client, err := ssh.Dial("tcp", addr, clientConfig)
	if err == nil {
		defer client.Close()
	}
	return err
}

func RemoteCommandByKey(servers *models.TaskServers) error {
	key, err := ioutil.ReadFile(servers.PrivateKeySrc)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%d", servers.ServerIp, servers.Port)
	config := &ssh.ClientConfig{
		User: servers.ServerAccount,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err == nil {
		client.Close()
	}
	return err
}

func RemoteAgent(servers *models.TaskServers) error {

	conn, err := net.Dial("tcp", servers.ServerIp+":"+strconv.Itoa(servers.Port))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	var reply *RpcResult
	defer client.Close()

	ping := "ping"
	err = client.Call("RpcTask.HeartBeat", ping, &reply)
	if err != nil {
		return err
	}
	if reply.Status == 200 {
		return nil
	} else {
		return fmt.Errorf("链接错误：%v", reply.Message)
	}
}

//定时任务日志收集
func TaskLogGather() {
	var (
		takCache   TaskRunLog
		taskRunLog models.TaskRunLog
		err        error
	)
	maxCount := 300
	msNum := 2
	for i := 0; i < maxCount; i++ {
		takCache = PullTaskRunLog()
		if takCache.TaskId <= 0 {
			break
		}
		taskRunLog = models.TaskRunLog{
			Id:          primitive.NewObjectID(),
			TaskId:      takCache.TaskId,
			ServerId:    takCache.ServerId,
			ServerName:  takCache.ServerName,
			Output:      takCache.Output,
			Error:       takCache.Error,
			Status:      takCache.Status,
			ProcessTime: takCache.ProcessTime,
			Ctime:       takCache.Ctime,
		}
		_, err = util.MgoClient.InsertOne("log", "task_run_log", taskRunLog)
		if err != nil {
			logStr, _ := json.Marshal(&takCache)
			util.RedisObj.LPush(util.RedisKeyList[util.JobTaskRunLog], string(logStr))
			common.DebugLog("mongo.log", "Error", fmt.Sprintf("日志记录失败:%s。", err.Error()))
			break
		}
		time.Sleep(time.Duration(msNum) * time.Millisecond)
	}
}
