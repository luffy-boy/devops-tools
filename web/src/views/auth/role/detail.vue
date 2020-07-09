<template>
  <div class="createPost-container">
    <el-form
      ref="postForm"
      :model="postForm"
      :rules="rules"
      class="form-container"
      label-width="120px"
    >
      <el-form-item label="角色名称" prop="role_name">
        <el-input v-model="postForm.role_name" type="text" />
      </el-form-item>
      <el-form-item label="角色别名" prop="role">
        <el-input v-model="postForm.role" type="text" />
      </el-form-item>
      <el-form-item label="上级角色" prop="parent_id">
        <selectTree
          v-if="isFinished"
          :data="roleList"
          :defaultProps="{children:'child',label:'role_name',id:'role_id'}"
          v-model="postForm.parent_id"
          :filterable="postForm.role_id"
        ></selectTree>
      </el-form-item>
      <el-form-item label="状态" prop="status">
        <el-switch v-model="postForm.status"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">提交</el-button>
        <el-button @click="$router.back(-1)">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { create, getDetail, fetchList } from '@/api/role'
import selectTree from '@/components/SelectTree'

const defaultForm = {
  role_id: 0,
  role_name: '',
  role: '',
  parent_id: '',
  status: true,
  is_delete: 0
}

export default {
  name: 'RoleDetail',
  components: { selectTree },
  data() {
    return {
      roleList: [],
      isFinished: false,
      listQuery: {
        page: 0,
        limit: 0
      },
      postForm: Object.assign({}, defaultForm),
      rules: {
        role_name: [
          { required: true, message: '请输入角色名称', trigger: 'blur' }
        ],
        role: [{ required: true, message: '请输入角色别名', trigger: 'blur' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getRoleList()
    if (this.$route.params.id !== undefined) {
      this.postForm.role_id = Number(this.$route.params.id)
      this.getDetail()
    }
  },
  methods: {
    getRoleList() {
      this.listLoading = true
      fetchList(this.listQuery).then(response => {
        this.roleList = response.data.list
        this.isFinished = true
        this.listLoading = false
      })
    },
    getDetail() {
      const params = {
        role_id: this.postForm.role_id
      }
      getDetail(params).then(response => {
        const self = this
        if (response.code !== 1) {
          return
        }
        for (var index in response.data) {
          self.postForm[index] = response.data[index]
          if (index === 'status') {
            self.postForm[index] = self.postForm[index] === 1
          }
          if (index === 'parent_id') {
            self.postForm[index] = Number(self.postForm[index])
          }
        }
      })
    },
    onSubmit() {
      let self = this;
      this.$refs['postForm'].validate(valid => {
        if (valid) {
          const form = Object.assign({}, this.postForm)
          form.status = (form.status === true) ? 1 : 0
          form.parent_id = Number(form.parent_id)
          const method = form.role_id > 0 ? 'post' : 'put'
          create(form, method).then(response => {
            if (response.code !== 1) {
              return
            }
            this.$message({
              message: '成功',
              type: 'success',
              onClose: function() {
                self.$router.back(-1)
              }
            })
          })
        }
      })
    }
  }
}
</script>

