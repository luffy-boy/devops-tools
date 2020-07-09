<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container" label-width="120px">
      <el-form-item label='账号' prop="username" >
        <el-input v-model="postForm.username" type="text" :readonly="postForm.user_id > 0" />
      </el-form-item>
      <el-form-item label='密码' prop="password" >
        <el-input v-model="postForm.password" type="password" auto-complete="new-password" show-password />
      </el-form-item>
      <el-form-item label='确认密码' prop="confirm_password">
        <el-input v-model="postForm.confirm_password" type="password" show-password />
      </el-form-item>
        <el-form-item label='昵称' prop="real_name" >
        <el-input v-model="postForm.real_name" type="text" />
      </el-form-item>
      <el-form-item label='手机' prop="phone" >
        <el-input v-model="postForm.phone" type="text" />
      </el-form-item>
      <el-form-item label='Email' prop="email" >
        <el-input v-model="postForm.email" type="text" />
      </el-form-item>
      <el-form-item label='性别'  prop="sex" >
        <el-select v-model="postForm.sex" class="filter-item" placeholder="Please select">
          <el-option v-for="item in sexTypeOptions" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
      </el-form-item>
        <el-form-item label='出生日期'  prop="birthday" >
        <el-date-picker v-model="postForm.birthday" type="datetime" placeholder="Please pick a birthday" />
      </el-form-item>
        <el-form-item label='角色id'  prop="role_id" >
        <el-select v-model="postForm.role_id" class="filter-item" placeholder="Please select">
          <el-option v-for="item in rolesOptions" :key="item.role_id" :label="item.role_name" :value="item.role_id" />
        </el-select>
      </el-form-item>
      <el-form-item label='扩展权限'  prop="route_ids" >
        <el-input v-model="postForm.route_ids" type="text" />
      </el-form-item>
      <el-form-item label='个人说明'>
        <el-input v-model="postForm.introduction" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="Please input" />
      </el-form-item>
      <el-form-item label='状态'  prop="status" >
        <el-switch v-model="postForm.status"></el-switch>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">提交</el-button>
        <el-button>取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { roleAll } from '@/api/role'
import { create, getDetail } from '@/api/admin'
import { parseTime } from '@/utils'
import { validUsername, validPassword } from '@/utils/validate'

const defaultForm = {
  user_id: 0,
  username: '',
  password: '',
  confirm_password: '',
  real_name: '',
  phone: '',
  email: '',
  sex: '',
  birthday: '',
  role_id: '',
  route_ids: '',
  introduction: '',
  status: true
}

const sexTypeOptions = [
  { id: 1, name: '男' },
  { id: 2, name: '女' },
  { id: 3, name: '不详' }
]

export default {
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error('长度在 5 到 32 个字符'))
      }
      callback()
    }
    const validatePassword = (rule, value, callback) => {
      if (this.postForm.user_id > 1 && value.length === 0){
        callback()
      }
      if (!validPassword(value)) {
        callback(new Error('长度在 5 到 32 个字符'))
      }
      callback()
    }
    const validateConfirmPassword = (rule, value, callback) => {
      if (this.postForm.password.length === 0){
        callback()
      }
      if (!validPassword(value)) {
        callback(new Error('长度在 5 到 32 个字符'))
      }
      if (value != this.postForm.password) {
        callback(new Error('两次输入密码不一致'))
      }
      callback()
    }
    return {
      sexTypeOptions,
      rolesOptions:[],
      postForm: Object.assign({}, defaultForm),
      rules: {
        username: [{ validator: validateUsername, trigger: 'blur' }],
        password: [{ validator: validatePassword, trigger: 'blur' }],
        confirm_password: [{ validator: validateConfirmPassword, trigger: 'blur' }],
        real_name: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
        phone: [{ required: true, message: '请输入手机', trigger: 'blur' }],
        email: [
          { required: true, message: '请输入邮箱', trigger: 'blur' },
          { type: "email", message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
         ],
        sex: [{ required: true, message: '请选择性别', trigger: 'change' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }],
        birthday: [{ type: "date",required: true, message: '请输入出生日期', trigger: 'blur' }],
        role_id: [{ required: true, message: '请选择角色', trigger: 'blur' }],
      },
    }
  },
  created() {
    if (this.$route.params.id !== undefined) {
      this.postForm.user_id = Number(this.$route.params.id)
      this.getDetail()
    }
    this.getRoleList()
  },
  methods: {
    getRoleList() {
      roleAll().then(response => {
      this.rolesOptions = response.data.list
      })
    },
    getDetail() {
      const params = {
        user_id: this.postForm.user_id
      }
      getDetail(params).then(response => {
        this.postForm.user_id = response.data.user_id
        this.postForm.real_name = response.data.real_name
        this.postForm.username = response.data.username
        this.postForm.sex = response.data.sex
        this.postForm.phone = response.data.phone
        this.postForm.email = response.data.email
        this.postForm.role_id = response.data.role_id
        this.postForm.route_ids = response.data.route_ids
        this.postForm.introduction = response.data.introduction
        this.postForm.password = ''
        this.postForm.confirm_password = ''
        this.postForm.birthday = new Date(response.data.birthday)
        this.postForm.status = response.data.status === 1
      })
    },
    onSubmit() {
      const self = this
      this.$refs['postForm'].validate((valid) => {
        if (valid) {
          const form = {
            user_id: this.postForm.user_id,
            username: this.postForm.username,
            password: this.postForm.password.length > 0 ? this.$md5(this.postForm.password) : "",
            real_name: this.postForm.real_name,
            phone: this.postForm.phone,
            email: this.postForm.email,
            sex: this.postForm.sex,
            birthday: parseTime(this.postForm.birthday),
            role_id: this.postForm.role_id,
            route_ids: this.postForm.route_ids,
            introduction: this.postForm.introduction,
            status: this.postForm.status ? 1 : 0,
          }
          let method = form.user_id > 0 ? 'post' : 'put'
          create(form,method).then(response => {
            if (response.code != 1) {
              return
            }
            this.$message({
              message: '成功',
              type: 'success',
              onClose: function(){
                 self.$router.go(0)
              }
            })
          })
        }
      })
    }
  }
}
</script>

