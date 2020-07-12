<template>
  <div class="app-container">
    <div class="filter-container">
      <router-link :to="'/devops/servers/add'">
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit">
        新增
      </el-button>
      </router-link>

    </div>

    <el-table
      :data="list"
      style="width: 100%;margin-bottom: 20px;"
      border>
      <el-table-column
        prop="server_id"
        label="id">
      </el-table-column>
      <el-table-column
        prop="group_name"
        label="分组信息">
      </el-table-column>
        <el-table-column
        prop="server_name"
        label="服务器名称">
      </el-table-column>
      <el-table-column  label="连接ip" width="180">
        <template slot-scope="{row}">
           {{ row.server_ip }}
        </template>
      </el-table-column>
      <el-table-column label="连接方式" width="180">
        <template slot-scope="{row}">
           {{ row.connection_type | connectionTypeFilter }}
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template slot-scope="{row}">
           {{ row.ctime | parseTime }}
        </template>
      </el-table-column>
      <el-table-column label="状态" class-name="status-col">
        <template slot-scope="{row}">
          <el-tag :type="row.status | statusFilter">
            {{ row.status | statusInfoFilter }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="230" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <router-link :to="'/devops/servers/edit/'+row.server_id">
            <el-button type="primary" size="small" icon="el-icon-edit">
              编辑
            </el-button>
          </router-link>
          <el-button v-if="row.status!=1" size="small" icon="el-icon-delete" type="danger" @click="handleDelete(row.server_id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />
  </div>
</template>

<script>
import { fetchList,del} from '@/api/servers'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

const statusTypeOptions = [
  { id: 0, name: '无效' },
  { id: 1, name: '有效' },
]

const statusTypeKeyValue = statusTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})

const connectionTypeOptions = [
  { id: 0, name: 'SSH' },
  { id: 1, name: 'Telnet' },
]

const connectionTypeKeyValue = connectionTypeOptions.reduce((acc, cur) => {
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
    statusInfoFilter(id){
      return statusTypeKeyValue[id]
    },
    connectionTypeFilter(type){
      return connectionTypeKeyValue[type]
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
      downloadLoading: false
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
        this.total = response.data.total
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    handleDelete(route_id) {
      let params = {
        route_id : route_id
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
    }
  }
}
</script>
