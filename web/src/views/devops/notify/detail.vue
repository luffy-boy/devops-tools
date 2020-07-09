<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container" label-width="120px">
        <el-form-item label="模板名称" prop="tpl_name">
          <el-input v-model="postForm.tpl_name" type="text" />
        </el-form-item>
        <el-form-item label="通知方式">
          <el-radio-group v-model="postForm.notify_type" size="medium">
            <el-radio v-for="item in notify_type_list"  border :key="item.id" :label="item.id">{{item.notify_name}}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="通知内容">
          <el-input v-model="postForm.tpl_data" :autosize="{ minRows: 8, maxRows: 12}" type="textarea" placeholder="Please input" />
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
import {create,getDetail} from '@/api/notify'
import {getNotifyData} from '@/api/task'

const defaultForm = {
  id: 0,
  tpl_name:"",
  notify_type: 0,
  tpl_data: "",
  status: true,
}

export default {
  data(){
    return {
      notify_type_list:[],
      postForm: Object.assign({},defaultForm),
      rules: {
        notify_type: [{required: true, message: '请选择通知方式', trigger: 'change' }],
        tpl_name: [{required: true, message: '请输入模板名称', trigger: 'blur' }],
        tpl_data: [{required: true, message: '请输入模板内容', trigger: 'blur' }],
        status: [{required: true, message: '请选择状态', trigger: 'change' }],
      },
    }
  },
  created(){
    if (this.$route.params.id !== undefined){
         this.postForm.id = Number(this.$route.params.id)
         this.getDetail()
    }
    this.getNotifyTypeList()
  },
  methods:{
    getNotifyTypeList() {
      const params = {
        expand: "notify_type"
      }
        getNotifyData(params).then(response => {
          if(response.code === 1){
            this.notify_type_list = response.data.notify_type
          }
      })
    },
    getDetail(){
      let params = {
        id: this.postForm.id
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
          let method = form.id > 0 ? "post" : "put"
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

