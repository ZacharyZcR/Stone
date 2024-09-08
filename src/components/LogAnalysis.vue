<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />

    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <LogFilter @filter-applied="applyFilter" />
      <LogTable :logs="logs" @view-details="viewDetails" />
      <LogPagination
          :current-page="currentPage"
          :total-pages="totalPages"
          :total-count="totalCount"
          :page-size="pageSize"
          @page-changed="changePage"
          @page-size-changed="changePageSize"
      />
      <StatisticsChart class="mt-8" /> <!-- 添加 mt-8 类来增加顶部边距 -->
    </div>

    <FooterPage />

    <LogDetailModal
        v-if="selectedLog"
        :log="selectedLog"
        @close="selectedLog = null"
    />
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import api from '../api/axiosInstance';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';
import LogFilter from './LogFilter.vue';
import LogTable from './LogTable.vue';
import LogPagination from './LogPagination.vue';
import StatisticsChart from './StatisticsChart.vue';
import LogDetailModal from './LogDetailModal.vue';

export default {
  name: 'LogAnalysis',
  components: {
    HeaderPage,
    FooterPage,
    LogFilter,
    LogTable,
    LogPagination,
    StatisticsChart,
    LogDetailModal
  },
  setup() {
    const logs = ref([]);
    const selectedLog = ref(null);
    const currentPage = ref(1);
    const pageSize = ref(20);
    const totalCount = ref(0);
    const filters = ref({});

    const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value));

    const fetchLogs = async () => {
      try {
        const params = {
          ...filters.value,
          page: currentPage.value,
          pageSize: pageSize.value
        };
        const response = await api.get('/logs', { params });
        logs.value = response.data.logs;
        totalCount.value = response.data.totalCount;
        currentPage.value = response.data.page;
      } catch (error) {
        console.error('获取日志失败:', error);
      }
    };

    const applyFilter = (newFilters) => {
      filters.value = newFilters;
      currentPage.value = 1;
      fetchLogs();
    };

    const changePage = (page) => {
      currentPage.value = page;
      fetchLogs();
    };

    const changePageSize = (size) => {
      pageSize.value = size;
      currentPage.value = 1;
      fetchLogs();
    };

    const viewDetails = (log) => {
      selectedLog.value = log;
    };

    // 初始加载
    fetchLogs();

    return {
      logs,
      selectedLog,
      currentPage,
      totalPages,
      totalCount,
      pageSize,
      applyFilter,
      changePage,
      changePageSize,
      viewDetails
    };
  }
};
</script>
