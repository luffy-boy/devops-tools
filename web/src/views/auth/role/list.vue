<template>
  <div class="app-container">
    <div class="filter-container">
      <router-link :to="'/auth/role/detail'">
        <el-button class="filter-item" type="primary">新增角色</el-button>
      </router-link>
    </div>
    <div>
      <el-table
        :data="list"
        style="width: 100%;margin-bottom: 20px;"
        row-key="role_id"
        border
        default-expand-all
        :tree-props="{children: 'child'}"
      >
        <el-table-column prop="role_name" label="角色名称" width="200"></el-table-column>
        <el-table-column prop="role" label="别名" width="180"></el-table-column>
        <el-table-column prop="status" label="状态" sortable width="100">
          <template slot-scope="{row}">
            <el-tag :type="row.status | statusFilter">{{ row.status | statusInfoFilter }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ctime" label="创建时间" sortable>
          <template slot-scope="{row}">
            <span>{{ row.ctime | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
          </template>
        </el-table-column>
        <el-table-column label="管理操作">
          <template slot-scope="{row}">
            <router-link :to="'/auth/role/route_edit/'+row.role_id">
              <el-button type="warning" icon="el-icon-star-off" size="small">权限配置</el-button>
            </router-link>
            <router-link :to="'/auth/role/detail/'+row.role_id">
              <el-button type="primary" icon="el-icon-edit" size="small">编辑</el-button>
            </router-link>
            <el-button
              type="danger"
              icon="el-icon-delete"
              size="small"
              @click="handleDelete(row.role_id)"
            >删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
<script>
import { fetchList, del } from "@/api/role";
import { parseTime } from "@/utils";

const statusTypeOptions = [
  { id: 0, name: "无效" },
  { id: 1, name: "有效" }
];

const statusTypeKeyValue = statusTypeOptions.reduce((acc, cur) => {
  acc[cur.id] = cur.name;
  return acc;
}, {});

export default {
  name: "roleList",
  filters: {
    statusFilter(status) {
      const statusMap = {
        0: "danger",
        1: "success"
      };
      return statusMap[status];
    },
    statusInfoFilter(id) {
      return statusTypeKeyValue[id];
    }
  },
  data() {
    return {
      list: [],
      listQuery: {
        page: 0,
        limit: 0
      },
      listLoading: true,
      downloadLoading: false
    };
  },
  created() {
    this.getList();
  },
  methods: {
    getList() {
      this.listLoading = true;
      fetchList(this.listQuery).then(response => {
        this.list = response.data.list;
        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false;
        }, 1.5 * 1000);
      });
    },
    handleDelete(role_id) {
      let params = {
        role_id: role_id
      };
      this.$confirm("确认是否删除？", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          del(params).then(response => {
            if (response.code !== 1) {
              return
            }
            this.$message({
               message: '成功',
               type: 'success'
            });
            this.getList()
          });
        })
        .catch(() => {});
    }
  }
};
</script>
