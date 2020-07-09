<template>
  <div class="app-container">
    <div class="filter-container">
      <router-link :to="'/devops/task/detail'">
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit">
        新增
      </el-button>
      </router-link>
      <el-button class="filter-item" style="margin-left: 10px;"  @click="auditTask(1)" type="success" icon="el-icon-success">
        审核通过
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" @click="auditTask(2)" type="danger" icon="el-icon-error">
        审核不通过
      </el-button>
    </div>

    <el-table
      :data="list"
      style="width: 100%;margin-bottom: 20px;"
      @selection-change="handleSelectionChange"
      border>
      <el-table-column
        type="selection"
        :selectable='selectable'
        width="55">
      </el-table-column>
      <el-table-column
        prop="task_id"
        label="id">
      </el-table-column>
      <el-table-column
        prop="group_name"
        label="分组信息">
      </el-table-column>
      <el-table-column
        prop="task_name"
        label="任务名称">
      </el-table-column>
      <el-table-column label="执行策略" width="180">
        <template slot-scope="{row}">
           {{ row.run_type | runTypeFilter }}
        </template>
      </el-table-column>
      <el-table-column
        prop="cron_spec"
        label="时间表达式">
      </el-table-column>
      <el-table-column prop="concurrent" label="并发执行">
         <template slot-scope="{row}">
           {{ row.concurrent | concurrentFilter }}
        </template>
      </el-table-column>
      <el-table-column
        prop="execute_times"
        label="累计执行次数">
      </el-table-column>
      <el-table-column label="下次执行时间">
        <template slot-scope="{row}">
           {{ row.next_time | parseTime }}
        </template>
      </el-table-column>
      <el-table-column label="审核状态" class-name="status-col">
        <template slot-scope="{row}">
          <el-tag :type="row.is_audit | auditFilter">
            {{ row.is_audit | auditInfoFilter }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="运行状态" class-name="status-col">
        <template slot-scope="{row}">
          <el-switch
            v-model="row.status"
            :disabled="row.is_audit !== 1"
            @change="runningTask($event,row)"
            active-color="#13ce66"
            inactive-color="#ff4949">
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="400" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <router-link :to="'/devops/task/detail/'+row.task_id">
            <el-button type="primary" size="small" icon="el-icon-edit">
              查看
            </el-button>
          </router-link>
          <router-link :to="'/devops/task/log/'+row.task_id">
            <el-button size="small" icon="el-icon-reading">
              日志
            </el-button>
          </router-link>
          <el-button v-if="row.is_audit !== 2" size="small" icon="el-icon-video-play" type="success" @click="executeTask(row.task_id)">
            测试
          </el-button>
          <el-button v-if="row.is_audit === 2" size="small" icon="el-icon-delete" type="danger" @click="handleDelete(row.task_id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { fetchList, del, auditTask, executeTask, runningTask} from '@/api/task'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination


const auditTypeOptions = [
  { id: 0, name: '待审核' },
  { id: 1, name: '已通过' },
  { id: 2, name: '未通过' },
]

const auditTypeKeyValue = auditTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})


export default {
  name: 'routeList',
  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        0: 'danger',
        1: 'success',
      }
      return statusMap[status]
    },
    auditFilter(status) {
      const auditMap = {
        0: 'warning',
        1: 'success',
        2: 'danger',
      }
      return auditMap[status]
    },
    auditInfoFilter(id){
      return auditTypeKeyValue[id]
    },
    concurrentFilter(value){
      const concurrentMap = {
        0: '否',
        1: '是'
      }
      return concurrentMap[value]
    },
    runTypeFilter(value){
      const runTypetMap = {
        0: '同时执行',
        1: '轮询执行'
      }
      return runTypetMap[value]
    }
  },
  data() { 
    return {
      list: [],
      total: 0,
      listQuery: {
        page: 1,
        limit: 20,
      },
      listLoading: true,
      downloadLoading: false,
      multipleSelection:[],
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.list = response.data.list
        // console.log(this.list)
        for (let v of this.list){
          v.status = v.status === 1
        }
        this.total = response.data.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    handleDelete(task_id) {
      const params = {
        task_id : task_id
      }
      this.$confirm('确认是否删除？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
      }).then(() => {
        del(params).then(response => {
          if (response.code === 1){
              this.getList()
              this.$message({
                message: '删除成功',
                type: 'success'
              });
          }else{
            this.$message({
              message: '删除失败',
              type: 'warning'
            });
          }
        }) 
      }).catch(() => {         
      });
    },
    auditTask(audit) {
      //审核任务
      if(this.multipleSelection.length === 0){
        this.$message({
          showClose: true,
          message: '请选择任务~',
          type: 'error'
        });
        return
      }

      const msg = audit === 1 ? '通过' : '不通过'
      this.$confirm('确认审核'+msg+'?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let taskIds = ""
        for(let val of this.multipleSelection){
          taskIds = taskIds + val.id + ","
        }
        taskIds =  taskIds.substring(0, taskIds.lastIndexOf(','));
        const form ={
          audit: audit,
          task_ids: taskIds
        }
        auditTask(form).then(response => {
          if(response.code != 1){
            return
          }
          this.$message({
            message: '成功',
            type: 'success'
          });
          this.getList()
        })
      }).catch(() => {});
    },
    executeTask(taskId) {
      //执行一次
      this.$confirm('确定执行该任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const query ={
          task_id: taskId
        }
        executeTask(query).then(response => {
          if(response.code != 1){
            return
          }
          this.$message({
            message: '成功',
            type: 'success'
          });
          this.getList()
        })
      }).catch(() => {});
    },
    runningTask(runningType,task) {
      //修改任务运行状态
      const msg = runningType ? '启动' : '停止'
      this.$confirm('确认'+msg+'任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const form ={
          status: runningType ? 1 : 0,
          task_ids: task.task_id.toString()
        }
        runningTask(form).then(response => {
          if(response.code != 1){
            task.status = !runningType
            return
          }
          this.$message({
            message: '成功',
            type: 'success'
          });
          this.getList()
        })
      }).catch(() => {
        task.status = !runningType
      });
    },
    stopTask() {
      //暂停任务
      this.$confirm('确认审核任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        let taskIds = ""
        for(let val of this.multipleSelection){
          taskIds = taskIds + val.id + ","
        }
        taskIds =  taskIds.substring(0, taskIds.lastIndexOf(','));
        const form ={
          audit: audit,
          task_ids: taskIds
        }
        auditTask(form).then(response => {
          if(response.code != 1){
            return
          }
          this.$message({
            message: '成功',
            type: 'success'
          });
          this.getList()
        })
      }).catch(() => {});
    },
    selectable(row,index){
      return row.is_audit === 0
    },
    handleSelectionChange(val) {
      this.multipleSelection = [];
      for(let i=0; i < val.length; i++){
        this.multipleSelection.push({id:val[i]['task_id']})
      }
    }
  }
}
</script>
