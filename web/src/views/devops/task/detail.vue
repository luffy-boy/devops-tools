<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container" label-width="120px">
        <el-form-item label="任务名称" prop="task_name" >
          <el-input v-model="postForm.task_name" type="text" />
        </el-form-item>
       <el-form-item label="任务分组"  prop="group_id" >
          <el-select v-model="postForm.group_id" class="filter-item" placeholder="请选择">
            <el-option v-for="item in groupList" :key="item.id" :label="item.group_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器列表" prop="server_ids">
          <el-select v-model="postForm.server_ids" filterable multiple collapse-tags placeholder="请选择">
            <el-option v-for="item in serverList" :key="item.server_id" :label="item.server_name" :value="item.server_id" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行方式"  prop="run_type" >
          <el-select v-model="postForm.run_type" class="filter-item" placeholder="请选择">
            <el-option v-for="item in runTypeList" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间表达式" prop="cron_spec" >
          <el-tooltip placement="top-start">
            <div slot="content">
              crontab表达式 <br/>
              ┌─────────────second 范围 (0 - 60)<br/>
              │ ┌───────────── min (0 - 59)<br/>
              │ │ ┌────────────── hour (0 - 23) <br/>
              │ │ │ ┌─────────────── day of month (1 - 31)<br/>
              │ │ │ │ ┌──────────────── month (1 - 12)<br/>
              │ │ │ │ │ ┌───────────────── day of week (0 - 6) (0 to 6 are Sunday to Saturday)<br/>
              │ │ │ │ │ │                 <br/>
              │ │ │ │ │ │<br/>
              │ │ │ │ │ │<br/>
              *&nbsp;&nbsp;*&nbsp;&nbsp;*&nbsp;&nbsp;*&nbsp;&nbsp;*&nbsp;&nbsp;*
              </div>
            <el-input v-model="postForm.cron_spec" type="text" placeholder="0/1 * * * * *">
              <el-button slot="append" icon="el-icon-alarm-clock" @click="runTime">执行时间</el-button>
            </el-input>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="执行命令" prop="command" >
          <el-input v-model="postForm.command" :autosize="{ minRows: 2, maxRows: 4}" type="textarea"/>
        </el-form-item>
        <el-form-item label="并发执行"  prop="concurrent" >
          <el-switch v-model="postForm.concurrent"></el-switch>
        </el-form-item>
        <el-form-item label="任务描述" prop="description" >
          <el-input v-model="postForm.description" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" />
        </el-form-item>
        <el-form-item label="超时时间" prop="timeout" >
          <el-input v-model="postForm.timeout" type="number" placeholder="单位(秒)" />
        </el-form-item>
        <el-form-item label="是否通知"  prop="is_notify" >
          <el-switch v-model="postForm.is_notify"></el-switch>
        </el-form-item>
        <el-form-item label="模板id"  prop="notify_tpl_id"  v-show="postForm.is_notify" >
          <el-select v-model="postForm.notify_tpl_id" class="filter-item" placeholder="请选择">
            <el-option v-for="item in notifyTplList" :key="item.id" :label="item.tpl_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户集合"  prop="notify_user_ids"  v-show="postForm.is_notify">
          <el-select v-model="postForm.notify_user_ids" filterable multiple collapse-tags placeholder="请选择">
            <el-option v-for="item in adminList" :key="item.user_id" :label="item.real_name" :value="item.user_id" />
          </el-select>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" :disabled="postForm.status" @click="onSubmit">提交</el-button>
            <el-button>取消</el-button>
          </el-form-item>
    </el-form>

    <el-drawer
      :title="cronTime.title"
      :visible.sync="cronTime.table"
      direction="rtl"
      size="20%">
      <el-table :data="runTimeLsit">
          <el-table-column property="run_date" label="执行时间" width="200"></el-table-column>
        </el-table>
    </el-drawer>

  </div>
</template>

<script>
import {create,getDetail,getGroupList,getNotifyData} from '@/api/task'
import {cronSpeNextTime} from '@/api/cron'
import {
  fetchList as serversFetchList
} from '@/api/servers'

const defaultForm = {
  task_id: 0,
  task_name: '',
  group_id: '',
  server_ids: [],
  run_type: '',
  cron_spec: '',
  command: '',
  concurrent: false,
  description: '',
  timeout: 60,
  is_notify: false,
  notify_type: '',
  notify_tpl_id: '',
  notify_user_ids: [],
  status: false
}

const runTypeList = [
  { id:0, name:'同时执行' },
  { id:1, name:'轮询执行' },
]

const typeList = [
  { id:0, name:'密码' },
  { id:1, name:'秘钥' },
]

//需要计算的指定次数
const calRunNum = 5

export default {
  data() {
    return {
      cronTime:{
        title:'接下来'+calRunNum+'次的执行时间',
        table: false,
      },
      runTypeList,
      typeList,
      calRunNum,
      groupList:[],
      serverList:[],
      notifyTplList:[],
      adminList:[],
      runTimeLsit:[],
      postForm: Object.assign({},defaultForm),
      rules: {
        group_id: [{ required: true, message: '请选择任务分组', trigger: 'change' }],
        task_name: [{ required: true, message: '请填写任务名称', trigger: 'blur' }],
        run_type: [{ required: true, message: '请选择执行方式', trigger: 'change' }],
        cron_spec: [{ required: true, message: '请填写时间表达式', trigger: 'blur' }],
        command: [{ required: true, message: '请填写执行命令', trigger: 'blur' }],
        concurrent: [{ required: true, message: '请选择并发模式', trigger: 'change' }],
        description: [{ required: true, message: '请填写任务描述信息', trigger: 'change' }],
      },
    }
  },
  created(){
    if (this.$route.params.id !== undefined) {
         this.postForm.task_id = Number(this.$route.params.id)
         this.getDetail()
    }
    this.getServerList()
    this.getGroupList()
    this.getNotifyData()
  },
  methods:{
    getGroupList() {
      getGroupList().then(response => {
        if(response.code === 1){
          this.groupList = response.data.list
        }
      })
    },
    getServerList() {
      const params = {
        page: 1,
        limit: 10000,
        extend: 'local'
      }
      serversFetchList(params).then(response => {
        if(response.code === 1){
          this.serverList = response.data.list
        }
      })
    },
    getDetail() {
      const params = {
        task_id: this.postForm.task_id
      }
      getDetail(params).then(response => {
        const self = this
        if(response.code !== 1){
          return
        }
        for (var index in response.data){
          self.postForm[index] = response.data[index]
          if (index === 'status') {
            self.postForm[index] = self.postForm[index] === 1
          }
          if (index === 'is_notify') {
            self.postForm[index] = self.postForm[index] === 1
          }
          if (index === 'server_ids' && self.postForm[index] != '') {
            self.postForm[index] = self.postForm[index].split(",").map(Number)
          }
          if (index === 'notify_user_ids') {
            if (self.postForm[index] !== '') {
              self.postForm[index] = self.postForm[index].split(",").map(Number)
            }else{
              self.postForm[index] = []
            }
          }
        }
      })
    },
    getNotifyData() {
      const self = this
      const params = {
        expand: "admin,notify_type,notify_tpl"
      }
      getNotifyData(params).then(response => {
        if(response.code === 1){
          response.data.admin_list && (self.adminList = response.data.admin_list)
          response.data.notify_tpl && (self.notifyTplList = response.data.notify_tpl)
        }
      })
    },
    onSubmit() {
      const self = this
      this.$refs['postForm'].validate((valid) => {
        if (valid) {
          const form = Object.assign({},this.postForm)
          form.status = form.status ? 1 : 0
          form.is_notify = form.is_notify ? 1 : 0
          form.concurrent = form.concurrent ? 1 : 0
          form.notify_tpl_id = Number(form.notify_tpl_id)
          form.notify_user_ids = form.notify_user_ids.join(',')
          form.server_ids = form.server_ids.join(',')
          
          const method = form.route_id > 0 ? 'post' : 'put'
          create(form,method).then(response => {
            if(response.code != 1){
              return
            }
            this.$message({
              message: '成功',
              type: 'success',
              onClose:function(){
                 self.$router.go(0);
              }
            });
          })
        }
      })
    },
    runTime() {
      let params = {
        cron_spec: this.postForm.cron_spec,
        cal_run_num: this.calRunNum,
      }
      if(params.cron_spec === ''){
        this.$message.error('请填写时间表达式');
          return
      }
      this.loading = true;
      cronSpeNextTime(params).then(response => {
        const self = this
        if(response.code !== 1){
          return
        }
        self.runTimeLsit = response.data.list
        self.cronTime.table = true
        this.loading = false
      })
    }
  }
}
</script>

