<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.real_name" placeholder="管理员名称" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-select v-model="listQuery.role_id" placeholder="角色" clearable style="width: 90px" class="filter-item">
        <el-option v-for="item in rolesOptions" :key="item.role_id" :label="item.role_name" :value="item.role_id" />
      </el-select>
      <el-select v-model="listQuery.sort" style="width: 140px" class="filter-item" @change="handleFilter">
        <el-option v-for="item in sortOptions" :key="item.key" :label="item.label" :value="item.key" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <router-link :to="'/auth/admin/detail'">
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit">
        新增
      </el-button>
      </router-link>
      <el-button v-waves :loading="downloadLoading" class="filter-item" type="primary" icon="el-icon-download" @click="handleDownload">
        导出
      </el-button>
      <el-checkbox v-model="showReviewer" class="filter-item" style="margin-left:15px;" @change="tableKey=tableKey+1">
        reviewer
      </el-checkbox>
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange"
    >
      <el-table-column label="user_id" prop="user_id" sortable align="center" :class-name="getSortClass('user_id')">
        <template slot-scope="{row}">
          <span>{{ row.user_id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="生日" align="center">
        <template slot-scope="{row}">
          <span>{{ row.birthday | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
        </template>
      </el-table-column>
      <el-table-column label="昵称" prop="real_name" align="center"   >
        <template slot-scope="{row}">
          <span>{{row.real_name}}</span>
        </template>
      </el-table-column>
      <el-table-column label="性别" prop="sex" align="center" >
        <template slot-scope="{row}">
          <span>{{row.sex | sexFilter}}</span>
        </template>
      </el-table-column>
      <el-table-column label="角色组" prop="role_name" align="center"   >
        <template slot-scope="{row}">
          <span>{{row.role_name}}</span>
        </template>
      </el-table-column>
      <el-table-column label="联系方式" prop="phone" align="center" >
        <template slot-scope="{row}">
          <span>{{row.phone}}</span>
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
        <template slot-scope="{row,$index}">
          <router-link :to="'/auth/admin/detail/'+row.user_id">
            <el-button type="primary" size="small" icon="el-icon-edit">
              编辑
            </el-button>
          </router-link>
          <el-button v-if="row.status!=1" size="small" type="danger" icon="el-icon-delete" @click="handleDelete(row,$index)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

  </div>
</template>

<script>
import { fetchList, fetchPv, createAdmin, updateAdmin } from '@/api/admin'
import {roleAll} from '@/api/role'
import waves from '@/directive/waves' // waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // secondary package based on el-pagination

const sexTypeOptions = [
  { id: 1, name: '男' },
  { id: 2, name: '女' },
  { id: 3, name: '不详' },
]

const sexTypeKeyValue = sexTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})

const statusTypeOptions = [
  { id: 0, name: '无效' },
  { id: 1, name: '有效' },
]

const statusTypeKeyValue = statusTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name
  return acc
}, {})

export default {
  name: 'adminList',
  components: { Pagination },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        0: 'danger',
        1: 'success',
      }
      return statusMap[status]
    },
    sexFilter(type) {
      return sexTypeKeyValue[type]
    },
    statusInfoFilter(id){
      return statusTypeKeyValue[id]
    }
  },
  data() {
    return {
      tableKey: 0,
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 20,
        real_name: undefined,
        role_id: undefined,
        sort: '-id'
      },
      sexTypeOptions,
      rolesOptions:[],
      sortOptions: [{ label: 'ID 顺序', key: '+id' }, { label: 'ID 逆序', key: '-id' }],
      showReviewer: false,
      dialogStatus: '',
      rules: {
        type: [{ required: true, message: 'type is required', trigger: 'change' }],
        timestamp: [{ type: 'date', required: true, message: 'timestamp is required', trigger: 'change' }],
        title: [{ required: true, message: 'title is required', trigger: 'blur' }]
      },
      downloadLoading: false
    }
  },
  created() {
    this.getList()
    this.getRoleList()
  },
  methods: {
    getRoleList() {
        roleAll().then(response => {
        this.rolesOptions = response.data.list
      })
    },
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
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    sortChange(data) {
      const { prop, order } = data
      if (prop === 'id') {
        this.sortByID(order)
      }
    },
    sortByID(order) {
      if (order === 'ascending') {
        this.listQuery.sort = '+id'
      } else {
        this.listQuery.sort = '-id'
      }
      this.handleFilter()
    },
    handleDelete(row, index) {
      this.$notify({
        title: 'Success',
        message: 'Delete Successfully',
        type: 'success',
        duration: 2000
      })
      this.list.splice(index, 1)
    },
    handleDownload() {
      this.downloadLoading = true
      import('@/vendor/Export2Excel').then(excel => {
        const tHeader = ['管理员id', '昵称', '手机号', '角色组', '状态']
        const filterVal = ['user_id', 'real_name', 'phone', 'role_name', 'status']
        const data = this.formatJson(filterVal)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: '管理员'
        })
        this.downloadLoading = false
      })
    },
    formatJson(filterVal) {
      return this.list.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    },
    getSortClass: function(key) {
      const sort = this.listQuery.sort
      return sort === `+${key}` ? 'ascending' : 'descending'
    }
  }
}
</script>
