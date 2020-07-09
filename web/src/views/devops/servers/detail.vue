<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container" label-width="120px">
       <el-form-item label="服务器分组"  prop="group_id" >
          <el-select v-model="postForm.group_id" class="filter-item" placeholder="请选择">
            <el-option v-for="item in groupList" :key="item.id" :label="item.group_name" :value="item.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器名称" prop="server_name">
          <el-input v-model="postForm.server_name" type="text" />
        </el-form-item>
        <el-form-item label="连接方式"  prop="connection_type" >
          <el-select v-model="postForm.connection_type" class="filter-item" placeholder="请选择">
            <el-option v-for="item in connectionList" :key="item.id" :label="item.name" :value="item.id"
            :disabled="item.disabled" />
          </el-select>
        </el-form-item>
        <el-form-item label="服务器IP" prop="server_ip" >
          <el-input v-model="postForm.server_ip" type="text" />
        </el-form-item>
        <el-form-item label="端口" prop="port" >
          <el-input v-model="postForm.port" type="number"/>
        </el-form-item>
        <el-form-item label="登录方式"  prop="type" >
          <el-select v-model="postForm.type" class="filter-item" placeholder="请选择">
            <el-option v-for="item in typeList" :key="item.id" :label="item.name" :value="item.id"
            :disabled="item.disabled" />
          </el-select>
        </el-form-item>
        <el-form-item v-show="postForm.type == 0" label="登录账号" prop="server_account" >
          <el-input v-model="postForm.server_account" type="text"/>
        </el-form-item>
        <el-form-item v-show="postForm.type == 0" label="登录密码" prop="password" >
          <el-input v-model="postForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item v-show="postForm.type == 1" label="私钥路径" prop="private_key_src" >
          <el-input v-model="postForm.private_key_src" type="text"/>
        </el-form-item>
        <el-form-item v-show="postForm.type == 1"  label="公钥路径" prop="publickey_src" >
          <el-input v-model="postForm.publickey_src" type="text"/>
        </el-form-item>
        <el-form-item label="detail">
          <el-input v-model="postForm.detail" :autosize="{ minRows: 2, maxRows: 4}" type="textarea" placeholder="Please input" />
        </el-form-item>
        <el-form-item label="状态"  prop="status" >
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
import {create,getDetail,getGroupList} from '@/api/servers'

const defaultForm = {
  group_id:1,
  connection_type:0,
  server_name:"",
  type:1,
  server_account:"",
  password:"",
  server_ip:"",
  server_outer_ip:"",
  port:"",
  private_key_src:"",
  publickey_src:"",
  detail:"",
  status:true,
}

const connectionList = [
  {id:0,name:"SSH"},
  {id:1,name:"Telnet"},
]

const typeList = [
  {id:0,name:"密码"},
  {id:1,name:"秘钥"},
]

export default {
  data(){
    return {
      connectionList,
      typeList,
      groupList:[],
      postForm: Object.assign({},defaultForm),
      rules: {
        group_id: [{required: true, message: '请选择服务器分组', trigger: 'change' }],
        server_name: [{required: true, message: '请输入服务器名称', trigger: 'blur' }],
        server_ip: [{required: true, message: '请输入服务器IP', trigger: 'blur' }],
        port: [{required: true, message: '请输入服务器端口', trigger: 'blur' }],
        type: [{required: true, message: '请选择登录方式', trigger: 'change' }],
        connection_type: [{required: true, message: '请选择连接方式', trigger: 'change' }],
        status: [{required: true, message: '请选择状态', trigger: 'change' }],
      },
    }
  },
  created(){
    if (this.$route.params.id !== undefined){
         this.postForm.server_id = Number(this.$route.params.id)
         this.getDetail()
    }
    this.getGroupList()
  },
  methods:{
    getGroupList() {
        getGroupList().then(response => {
          if(response.code === 1){
            this.groupList = response.data.list
          }
      })
    },
    getDetail(){
      let params = {
        server_id: this.postForm.server_id
      }
      getDetail(params).then(response => {
        let self = this
        if(response.code !== 1){
           return
        }
        for (var index in response.data){
          self.postForm[index] = response.data[index]
          if(index === "status"){
            self.postForm[index] = self.postForm[index] === 1
          }
        }
      })
    },
    onSubmit() {
      let self = this
      this.$refs['postForm'].validate((valid) => {
        if (valid) {
          let form = Object.assign({},this.postForm)
          form.status = form.status ? 1 : 0
          form.port = Number(form.port)
          let method = form.route_id > 0 ? "post" : "put"
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
  }
}
</script>

