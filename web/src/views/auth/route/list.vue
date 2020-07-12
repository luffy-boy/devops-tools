<template>
  <div class="app-container">
    <div class="filter-container">
      <router-link :to="'/auth/route/add'">
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit">
        新增
      </el-button>
      </router-link>

    </div>

    <el-table
      :data="list"
      style="width: 100%;margin-bottom: 20px;"
      row-key="route_id"
      border
      default-expand-all
      :tree-props="{children: 'children', hasChildren: 'hasChildren'}">
      <el-table-column
        prop="route_name"
        label="路由名称"
        width="200">
      </el-table-column>
      <el-table-column
        prop="route"
        label="路由"
        width="180">
      </el-table-column>
       <el-table-column label="请求方式" width="180">
         <template slot-scope="{row}">
          <div slot="reference" class="name-wrapper">
              <el-tag size="medium">{{ row.request | requestInfoFilter}}</el-tag>
          </div>
         </template>
      </el-table-column>
      <el-table-column  
        label="排序"
        width="100">
        <template slot-scope="{row}">
          <el-input
            size="mini"
            type="number"
            placeholder="顺序排序"
            v-model="row.sort"
            @blur="editSort(row.route_id,row.sort)"
            >
          </el-input>
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
          <router-link :to="'/auth/route/edit/'+row.route_id">
            <el-button type="primary" size="small" icon="el-icon-edit">
              编辑
            </el-button>
          </router-link>
          <el-button v-if="row.status!=1" size="small" type="danger" icon="el-icon-delete" @click="handleDelete(row.route_id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
  </div>
</template>

<script>
import { fetchList,del,editSort} from '@/api/route'
import { deepClone } from '@/utils'

const statusTypeOptions = [
  { id: 0, name: '无效' },
  { id: 1, name: '有效' },
]

const statusTypeKeyValue = statusTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})

export default {
  name: 'routeList',
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
    requestFilter(request){
      const requestMap = {
        "GET": 'success',
        "POST": 'success',
        "PUT": 'success',
        "DELETE": 'danger',
      }
       return request == "" ? "as" : requestMap[request]
    },
    requestInfoFilter(request){
      return request == "" ? "全部" : request
    }
  },
  data() { 
    return {
      list: [],
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
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
    },
    editSort(id,sort) {
      const params = {
        route_id : id,
        sort : Number(sort),
      }
      editSort(params).then(response => {
        let type = 'warning'
        let message = '修改失败'
        if(response.code === 1){
          type = 'success'
          message = '修改成功'
        }
        this.$message({
            message: message,
            type: type
        });
        this.getList()
      })
    },
    handleDelete(route_id) {
      const params = {
        route_id : route_id
      }
      this.$confirm('确认是否删除？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
      }).then(() => {
        del(params).then(response => {
          let type = 'warning'
          let message = '删除失败'
          if(response.code === 1){
            type = 'success'
            message = '删除成功'
            this.getList()
          }
          this.$message({
              message: type,
              type: message
          });
        }) 
      }).catch(() => {         
      });
    }
  }
}
</script>
