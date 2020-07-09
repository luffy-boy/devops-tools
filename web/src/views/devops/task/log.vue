<template>
  <div class="app-container">
    <div class="filter-container">
    </div>

    <el-table
      :data="list"
      style="width: 100%;margin-bottom: 20px;"
      border>
      <el-table-column
        prop="server_name"
        label="运行服务器">
      </el-table-column>
      <el-table-column label="执行结果" class-name="status-col">
        <template slot-scope="{row}">
          <el-tag :type="row.status | statusFilter">
            {{ row.status | statusInfoFilter }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column
        label="执行耗时">
        <template slot-scope="{row}">
          {{row.process_time}}ms
        </template>
      </el-table-column>
      <el-table-column label="执行时间">
        <template slot-scope="{row}">
           {{ row.ctime | parseTime }}
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="400" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <router-link :to="'/devops/task/log_detail/'+row.id">
            <el-button type="primary" size="small" icon="el-icon-edit">
              查看
            </el-button>
          </router-link>
        </template>
      </el-table-column>
    </el-table>
    
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { fetchLogList} from '@/api/task'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination


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
  name: 'routeList',
  components: { Pagination },
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
      list: [],
      total: 0,
      listQuery: {
        task_ld:0,
        page: 1,
        limit: 20,
      },
      listLoading: true,
    }
  },
  created() {
     if (this.$route.params.id !== undefined) {
       this.listQuery.task_id = Number(this.$route.params.id)
       this.getList()
    }
  },
  methods: {
    getList() {
      this.listLoading = true
      fetchLogList(this.listQuery).then(response => {
        this.list = response.data.list
        this.total = response.data.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    }
  }
}
</script>
