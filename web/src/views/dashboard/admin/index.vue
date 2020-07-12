<template>
  <div class="dashboard-editor-container">
    <github-corner class="github-corner" />

    <panel-group
      @handleSetLineChartData="handleSetLineChartData"
      :user_count="user_count"
      :audit_task="audit_task"
      :success_task="success_task"
      :fail_task="fail_task"
    />

    <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <line-chart :chart-data="lineChartData" />
    </el-row>
  </div>
</template>

<script>
import { getIndexData } from "@/api/admin";
import GithubCorner from "@/components/GithubCorner";
import PanelGroup from "./components/PanelGroup";
import LineChart from "./components/LineChart";
import RaddarChart from "./components/RaddarChart";
import PieChart from "./components/PieChart";
import BarChart from "./components/BarChart";
import TransactionTable from "./components/TransactionTable";
import TodoList from "./components/TodoList";
import BoxCard from "./components/BoxCard";

const lineChartData = {
  userCount: {
    actualData: [],
    dateData: []
  },
  auditTask: {
    actualData: [],
    dateData: []
  },
  successTask: {
    actualData: [],
    dateData: []
  },
  failTask: {
    actualData: [],
    dateData: []
  }
};

export default {
  name: "DashboardAdmin",
  components: {
    GithubCorner,
    PanelGroup,
    LineChart,
    RaddarChart,
    PieChart,
    BarChart,
    TransactionTable,
    TodoList,
    BoxCard
  },
  data() {
    return {
      lineChartData: lineChartData.userCount,
      user_count: 0,
      audit_task: 0,
      success_task: 0,
      fail_task: 0
    };
  },
  created() {
    this.getIndexData();
  },
  methods: {
    getIndexData() {
      getIndexData().then(response => {
        if (response.code != 1) {
          return;
        }
        //用户总数
        lineChartData.userCount.actualData = response.data.user_count.actual_data;
        lineChartData.userCount.dateData = response.data.user_count.date_data;
        this.user_count = response.data.user_count.total;

        //待审核任务
        lineChartData.auditTask.actualData = response.data.audit_task.actual_data;
        lineChartData.auditTask.dateData = response.data.audit_task.date_data;
        this.audit_task = response.data.audit_task.total;

        //成功任务
        lineChartData.successTask.actualData = response.data.success_task.actual_data;
        lineChartData.successTask.dateData = response.data.success_task.date_data;
        this.success_task = response.data.success_task.total;

        //失败任务
        lineChartData.failTask.actualData = response.data.fail_task.actual_data;
        lineChartData.failTask.dateData = response.data.fail_task.date_data;
        this.fail_task = response.data.fail_task.total;

      });
    },
    handleSetLineChartData(type) {
      this.lineChartData = lineChartData[type];
    }
  }
};
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
