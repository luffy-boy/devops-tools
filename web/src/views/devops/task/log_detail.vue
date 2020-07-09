<template>
  <el-container>
    <el-header>
    <el-page-header @back="goBack" title="" content="日志内容">
    </el-page-header>
    </el-header>
    <el-main>
    <el-form  class="form-container" label-width="120px">
      <el-divider content-position="left"><span>执行结果</span></el-divider>
        <el-form-item label="日志id" >
          <span>{{Detail.id}}</span>
        </el-form-item>
        <el-form-item label="运行服务器" >
          <span>{{Detail.server_name}}</span>
        </el-form-item>
        <el-form-item label="执行时间" >
          <span>{{Detail.ctime | parseTime}}</span>
        </el-form-item>
        <el-form-item label="执行耗时" >
          <span>{{Detail.process_time}}ms</span>
        </el-form-item>
        <el-form-item label="运行结果" >
          <el-tag :type="Detail.status | statusFilter">
            {{ Detail.status | statusInfoFilter }}
          </el-tag>
        </el-form-item>

        <el-divider content-position="left"><span>执行命令</span></el-divider>
         <el-form-item label="执行耗时" >
          <span>{{Detail.task_info.command}}ms</span>
        </el-form-item>

        <el-divider content-position="left"><span>执行输出</span></el-divider>
        <el-form-item label="" >
          <el-input v-model="Detail.output" :autosize="{ minRows: 2, maxRows: 4}" readonly="readonly" type="textarea" />
        </el-form-item>

        <el-divider content-position="left"><span>错误输出</span></el-divider>
        <el-form-item label="" >
          <el-input v-model="Detail.error" :autosize="{ minRows: 2, maxRows: 4}" readonly="readonly" type="textarea" />
        </el-form-item>
    </el-form>
    </el-main>
 </el-container>
</template>

<script>
import {logDetail} from '@/api/task'

const statusTypeOptions = [
  { id: 3, name: '执行异常' },
  { id: 2, name: '执行超时' },
  { id: 1, name: '执行成功' },
]

const statusTypeKeyValue = statusTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})

export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        3: 'warning',
        2: 'danger',
        1: 'success',
      }
      return statusMap[status]
    },
    statusInfoFilter(id){
      return statusTypeKeyValue[id]
    }
  },
  data() {
    return {
      Detail: {
        log_id: '',
        task_id: 0,
        server_id: 0,
        server_name: '',
        output: '',
        error: '',
        status: 0,
        process_time: 0,
        ctime: 0
      }
    }
  },
  created(){
    if (this.$route.params.id !== undefined) {
         this.Detail.log_id = this.$route.params.id
         this.getDetail()
    }
  },
  methods:{
    goBack() {
      this.$router.go(-1)
    },
    getDetail() {
      const params = {
        log_id: this.Detail.log_id
      }
      logDetail(params).then(response => {
        const self = this
        if(response.code !== 1){
          return
        }
        self.Detail = response.data
      })
    }
  }
}
</script>

