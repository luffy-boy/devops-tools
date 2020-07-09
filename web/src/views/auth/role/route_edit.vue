<template>
  <el-container>
    <el-header>【{{roleName}}】权限配置</el-header>
    <el-main>
      <el-form v-model="checkedKeys">
        <el-tree
          :data="data"
          show-checkbox
          node-key="route_id"
          :default-expand-all="true"
          :check-strictly="true"
          :default-checked-keys="checkedKeys"
          :props="defaultProps"
          ref="tree"
          @check="handleCheckChange"
        ></el-tree>
        <el-button type="primary" @click="onSubmit">提交</el-button>
        <el-button @click="$router.back(-1)">取消</el-button>
      </el-form>
    </el-main>
  </el-container>
</template>

<style lang="scss" scoped>
@import "@/styles/route-edit.scss";
</style>

<script>
import { fetchList } from "@/api/route";
import { getDetail,routeEdit } from "@/api/role";

export default {
  name: "route_edit",
  data() {
    return {
      roleName: "",
      data: [],
      checkedKeys: [],
      listLoading: true,
      defaultProps: {
        children: "children",
        label: "route_name"
      }
    };
  },
  created() {
    this.getRouteList();
    this.getRoleDetail();
  },
  methods: {
    getRouteList() {
      this.listLoading = true;
      fetchList().then(response => {
        this.data = response.data.list;
        this.listLoading = false;
      });
    },
    getRoleDetail() {
      let params = {
        role_id: this.$route.params.id
      };
      getDetail(params).then(response => {
        let self = this;
        if (response.code !== 1) {
          return;
        }
        self.roleName = response.data.role_name;
        let routeIds = response.data.route_ids.split(",");

        self.checkedKeys = routeIds.map(function(data) {
          return +data;
        });
      });
    },
    handleCheckChange(data) {
      const node = this.$refs.tree.getNode(data.route_id);
      if (node.checked) {
        //设置父节点和子节点都选上
        this.setParentNode(node);
      }
      this.setChildNode(node);
    },
    setParentNode(node) {
      if (node.parent) {
        for (const key in node) {
          if (key === "parent") {
            node[key].checked = true;
            this.setParentNode(node[key]);
          }
        }
      }
    },
    setChildNode(node) {
      if (node.childNodes && node.childNodes.length) {
        node.childNodes.forEach(item => {
          item.checked = node.checked;
          this.setChildNode(item);
        });
      }
    },
    onSubmit() {
      let self = this;

      let aRouteIds = self.$refs.tree.getCheckedKeys();
      aRouteIds = aRouteIds ? aRouteIds : [];
      let sRouteIds = aRouteIds.join(',');
      let form = {
        role_id: self.$route.params.id,
        route_ids: sRouteIds ? sRouteIds : "",
      };
      routeEdit(form, "post").then(response => {
        if (response.code !== 1) {
              return
        }
        this.$message({
          message: '成功',
          type: 'success',
          onClose: function() {
            self.$router.back(-1);
          }
        });
      });
    }
  }
};
</script>