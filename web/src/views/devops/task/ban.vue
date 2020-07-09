<template>
  <div class="createPost-container">
    <el-tag
      :key="index"
      v-for="(tag, index) in dynamicTags"
      closable
      type="danger"
      :disable-transitions="false"
      @click="editTag(tag,index)"
      @close="handleClose(tag)"
    >
      <span v-if="index!=num">{{tag}}</span>
      <el-input
        size="small"
        class="custom_input"
        type="text"
        v-model="words"
        v-if="index==num"
        ref="editInput"
        @keyup.enter.native="handleInput(tag,index)"
        @blur="handleInput(tag,index)"
      ></el-input>
    </el-tag>
    <el-input
      class="input-new-tag"
      v-if="inputVisible"
      v-model="inputValue"
      ref="saveTagInput"
      size="small"
      @keyup.enter.native="handleInputConfirm"
      @blur="handleInputConfirm"
    ></el-input>
    <el-button v-else class="button-new-tag" size="small" @click="showInput">+添加命令</el-button>
  </div>
</template>
<style>
.el-tag + .el-tag {
  margin-left: 10px;
}
.button-new-tag {
  margin-left: 10px;
  height: 32px;
  line-height: 30px;
  padding-top: 0;
  padding-bottom: 0;
}
.input-new-tag {
  width: 90px;
  margin-left: 10px;
  vertical-align: bottom;
}
.custom_input {
  width: 80px;
  height: 16px;
  outline: none;
  border: transparent;
  background-color: transparent;
  font-size: 12px;
  color: #f56c6c;
}
</style>

<script>
import { getBanList, editBanList } from "@/api/task";
export default {
  name: "taskBan",
  data() {
    return {
      dynamicTags: [],
      inputVisible: false,
      inputValue: "",
      id: "",
      num: -1,
      words: "",
      listLoading: true
    };
  },
  created() {
    this.getBanList();
  },
  methods: {
    getBanList() {
      this.listLoading = true;
      getBanList().then(response => {
        if (response.code !== 1) {
          return;
        }
        console.log(response)
        const ban_list = response.data ? response.data.ban_list : [];
        const id = response.data ? response.data._id : "";
        this.dynamicTags = Object.values(ban_list);
        this.dynamicTags = this.unique(this.dynamicTags);
        this.id = id;
        this.listLoading = false;
      });
    },
    // 数组去重
    unique(arr) {
      let x = new Set(arr);
      return [...x];
    },
    handleClose(tag) {
      //删除命令
      this.$confirm("您确定要删除该命令吗?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(() => {
          let newBanList = JSON.parse(JSON.stringify(this.dynamicTags));
          newBanList.splice(this.dynamicTags.indexOf(tag), 1);
          let form = {
            ban_list: newBanList,
            id: this.id
          };
          editBanList(form).then(response => {
            if (response.code == 1)  {
              this.dynamicTags.splice(this.dynamicTags.indexOf(tag), 1)
            }
          });
        }).catch(() => {});
    },

    showInput() {
      this.inputVisible = true;
      this.$nextTick(_ => {
        this.$refs.saveTagInput.$refs.input.focus();
      });
    },
    handleInputConfirm() {
      //添加命令
      let inputValue = this.inputValue;

      if (inputValue) {
        if (this.checkUnique(inputValue, -1)) {
          let newBanList = JSON.parse(JSON.stringify(this.dynamicTags));
          newBanList.push(inputValue);
          let form = {
            ban_list: newBanList,
            id: this.id
          };
          editBanList(form).then(response => {
            response.code == 1 ? this.getBanList() : "";
          });
        }
      }
      this.inputVisible = false;
      this.inputValue = "";
    },
    editTag(tag, index) {
      this.num = index;
      this.$nextTick(_ => {
        this.$refs.editInput[0].focus();
      });
      this.words = tag;
    },
    handleInput(tag, index) {
      //修改命令
      let words = this.words;
      if (words && words != this.dynamicTags[index]) {
        if (this.checkUnique(words, index)) {
          let newBanList = JSON.parse(JSON.stringify(this.dynamicTags));
          newBanList[index] = words;
          let form = {
            ban_list: newBanList,
            id: this.id
          };
          editBanList(form).then(response => {
            response.code == 1 ? this.getBanList() : "";
          });
        }
      }

      this.words = "";
      this.num = -1;
    },
    checkUnique(words, index) {
      let i = 0;
      for (i; i < this.dynamicTags.length; i++) {
        if (i != index && words == this.dynamicTags[i]) {
          this.$message({
            message: "存在相同的命令",
            type: "error"
          });
          return false;
        }
      }
      return true;
    }
  }
};
</script>
