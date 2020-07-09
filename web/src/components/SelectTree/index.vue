<template>
  <div v-if="update">
    <el-popover
      popper-class="selectTree"
      placement="bottom-start"
      transition="fade-in-linear"
      v-model="visible"
      min-width="200"
      trigger="click"
    >
      <el-tree
        :data="data"
        :props="defaultProps"
        empty-text="无数据"
        :default-expand-all="default_expand_all"
        :node-key="defaultProps.id"
        :default-expanded-keys="defaultExpandedKeys"
        :check-on-click-node="true"
        ref="tree1"
        :expand-on-click-node="false"
        :filter-node-method="filterNode"
        :highlight-current="true"
        @node-click="handleNodeClick"
      ></el-tree>
      <el-input
        v-model="filterText"
        @clear="clear"
        :placeholder="placeholder"
        :disabled="disabled"
        slot="reference"
        :clearable="clearable"
        :suffix-icon="icon"
      ></el-input>
    </el-popover>
  </div>
</template>
 <script>
export default {
  name: "selectTree",
  props: {
    value:{
      default: undefined
    },
    data: Array,
    placeholder: {
      type: String,
      default: "请选择"
    },
    disabled: {
      type: Boolean,
      default: false
    },
    clearable: {
      type: Boolean,
      default: true
    },
    filterable: {
      //禁选值
      type: Number,
      default: undefined
    },
    default_expand_all:{
      type: Boolean,
      default: false
    },
    defaultProps: {
      type: Object,
      default() {
        return {
          children: "children",
          label: "label",
          id: "id"
        };
      }
    },
    nodeKey: {
      type: String,
      default: "id"
    }
  },
  data() {
    return {
      defaultExpandedKeys: ["-1"], //默认展开
      filterText: "",
      visible: false,
      icon: "el-icon-arrow-down",
      update: true
    };
  },
  async created() {
    if (this.filterable) {
      this.setFilter(this.data);
    }
  },
  mounted() {
    this.setFilterText();
  },
  watch: {
    value(val) {
      if (!val) {
        //没有值得时候
        this.filterText = "";
      } else {
        if (this.$refs.tree1) {
          this.$refs.tree1.setCurrentKey(val);
          let obj = this.$refs.tree1.getCurrentNode();
          if (obj) {
            this.filterText = obj[this.defaultProps.label];
            return;
          } else {
            let tree = this.$refs.tree1;
            this.$nextTick(() => {
              tree.setCurrentKey(val);
              let obj = tree.getCurrentNode();
              if (obj) {
                this.filterText = obj[this.defaultProps.label];
              }
              return;
            });
          }
        }
      }
    },
    visible(val) {
      if (val === true) {
        this.icon = "el-icon-arrow-up";
      } else {
        this.icon = "el-icon-arrow-down";
      }
    },
    filterable(val) {
      this.update = false;
      this.setFilter(this.data);
      this.$nextTick(() => {
        this.update = true;
      });
    }
  },
  methods: {
    setFilterText() {
      if (!this.value) {
        return;
      } else {
        this.$refs.tree1.setCurrentKey(this.value);
        let obj = this.$refs.tree1.getCurrentNode();
        if (obj) {
          this.filterText = obj[this.defaultProps.label];
        }
      }
    },
    setFilter(arr) {
      arr.map(item => {
        if (item[this.defaultProps.id] == this.filterable) {
          item['disabled'] = true;
          if (item[this.defaultProps.children] && item[this.defaultProps.children].length != 0) {
            this.setDisabled(item[this.defaultProps.children]);
          }
        } else {
          item['disabled'] = false;
          if (item[this.defaultProps.children] && item[this.defaultProps.children].length != 0) {
            this.setFilter(item[this.defaultProps.children]);
          }
        }
      });
    },
    setDisabled(arr) {
      arr.map(item => {
        item['disabled'] = true;
        if (item[this.defaultProps.children] && item[this.defaultProps.children].length != 0) {
          this.setDisabled(item[this.defaultProps.children]);
        }
      });
    },
    filterNode(value, data) {
      if (!value) return true;
      return data.label.indexOf(value) !== -1;
    },
    handleNodeClick(obj, node) {
      if (node.data.disabled === true || this.disabled === true) {
        return;
      }
      this.$refs.tree1.setCurrentKey(obj[this.defaultProps.id]);
      this.$emit("input", obj[this.defaultProps.id]);
      this.filterText = obj[this.defaultProps.label];
      this.visible = false;
    },
    clear() {
      this.$refs.tree1.setCurrentKey("");
      this.$emit("input", "");
    }
  },
  beforeDestroy() {},
  destroyed() {}
};
</script>
 <style lang="scss">
.selectTree {
  max-height: 600px;
  overflow-y: auto;
  overflow-x: hidden;

  div[aria-disabled="true"] {
    .el-tree-node__content {
      cursor: not-allowed;
    }
  }
}
</style>