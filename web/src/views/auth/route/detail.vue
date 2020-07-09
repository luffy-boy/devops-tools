<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container" label-width="120px">
      <el-form-item label="路由名称" prop="route_name" >
        <el-input v-model="postForm.route_name" type="text" />
      </el-form-item>
      <el-form-item label="路由" prop="route" >
        <el-input v-model="postForm.route" type="text" />
      </el-form-item>
      <el-form-item label="上级路由" prop="parent_id">
        <selectTree
          v-if="isFinished"
          :data="routeList"
          :default_expand_all="true"
          :defaultProps="{children:'children',label:'route_name',id:'route_id'}"
          :filterable="postForm.route_id"
          v-model="postForm.parent_id"
        ></selectTree>
      </el-form-item>
      <el-form-item label="请求方式"  prop="request" >
        <el-select v-model="postForm.request" class="filter-item" placeholder="Please select">
          <el-option v-for="item in requestList" :key="item.key" :label="item.title" :value="item.key" />
        </el-select>
      </el-form-item>

      <el-form-item label="前端路由"  prop="is_route" >
        <el-switch v-model="postForm.is_route"></el-switch>
      </el-form-item>
      <div v-if="postForm.is_route">
      <el-form-item label="路径" prop="path" >
        <el-input v-model="postForm.path" placeholder="路由路径" type="text" />
      </el-form-item>
      <el-form-item label="路由标识" prop="name" >
        <el-input v-model="postForm.name" placeholder="路由唯一标识" type="text" />
      </el-form-item>
      <el-form-item label="view路径" prop="component" >
        <el-input v-model="postForm.component" placeholder="view路径" type="text" />
      </el-form-item>
      <el-form-item label="路由重定向" prop="redirect" >
        <el-input v-model="postForm.redirect" type="text" placeholder="重定向地址，在面包屑中点击会重定向去的地址" />
      </el-form-item>
      <el-form-item label="显示菜单"  prop="hidden" >
        <el-switch v-model="postForm.hidden" ></el-switch>
      </el-form-item>
      <el-form-item label="图标"  prop="icon" >
        <el-input v-model="postForm.icon" placeholder="icon图标" type="text" />
      </el-form-item>
      <el-form-item label="额外参数"  prop="extra" >
        <el-input v-model="postForm.extra" type="text" />
      </el-form-item>
      </div>
      <el-form-item label="排序"  prop="sort" >
        <el-input v-model="postForm.sort" type="number" />
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
import { create,getDetail,fetchList } from '@/api/route'
import selectTree from "@/components/SelectTree";

const defaultForm = {
  route_id: 0,
  route_name: '',
  route: '',
  parent_id: '',
  request: '',
  status: true,
  is_route: true,
  path: '',
  component: '',
  name: '',
  redirect: '',
  hidden: false,
  icon: '',
  extra: '',
  sort: 50
}

const requestList = [
  { title: '全部',key: '' },
  { title: 'GET',key: 'GET' },
  { title: 'POST',key: 'POST' },
  { title: 'PUT',key: 'PUT' }
]

export default {
  components: { selectTree },
  data() {
    return {
      requestList,
      isFinished: false,
      routeList: [],
      checkedKeys: [3],
      postForm: Object.assign({}, defaultForm),
      rules: {
        route_name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
        route: [{ required: true, message: '请输入路由', trigger: 'blur' }],
        status: [{ required: true, message: '请选择状态', trigger: 'change' }]
      }
    }
  },
  created() {
    if (this.$route.params.id !== undefined) {
      this.postForm.route_id = Number(this.$route.params.id)
      this.getDetail()
    }
    this.getRouteList()
  },
  methods: {
    getRouteList() {
      this.isFinished = false
      fetchList().then(response => {
        if (response.code === 1) {
          response.data.list.forEach(element => {
            this.routeList = response.data.list
            this.isFinished = true
          })
        }
      })
    },
    getDetail() {
      const params = {
        route_id: this.postForm.route_id
      }
      getDetail(params).then(response => {
        const self = this
        if (response.code !== 1) {
          return
        }
        for (var index in response.data) {
          self.postForm[index] = response.data[index]
          if ( index === 'status' || index === 'is_route' ) {
            self.postForm[index] = self.postForm[index] === 1
          }
          if ( index === 'hidden'){
             self.postForm[index] = !(self.postForm[index] === 1)
          }
        }
      })
    },
    onSubmit() {
      const self = this
      this.$refs['postForm'].validate((valid) => {
        if (valid) {
          const form = {
            route_id: this.postForm.route_id,
            route_name: this.postForm.route_name,
            route: this.postForm.route,
            parent_id: Number(this.postForm.parent_id),
            request: this.postForm.request,
            status: this.postForm.status ? 1 : 0,
            is_route: this.postForm.is_route ? 1 : 0,
            component: this.postForm.component,
            path: this.postForm.path,
            name: this.postForm.name,
            redirect: this.postForm.redirect,
            hidden: this.postForm.hidden ? 0 : 1,
            icon: this.postForm.icon,
            extra: this.postForm.extra,
            sort: Number(this.postForm.sort)
          }
          const method = form.route_id > 0 ? 'post' : 'put'
          create(form, method).then(response => {
            if (response.code !== 1) {
              return
            }
            this.$message({
              message: '成功',
              type: 'success',
              onClose: function() {
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

