<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 日志过滤器 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">日志过滤器 🔍</h2>
        <form class="grid grid-cols-1 md:grid-cols-4 gap-4" @submit.prevent="applyFilter">
          <div>
            <label class="block text-gray-300 text-sm font-bold mb-2" for="startDateTime">开始时间</label>
            <input v-model="filters.startDateTime" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-indigo-500 transition duration-300" id="startDateTime" type="datetime-local">
          </div>
          <div>
            <label class="block text-gray-300 text-sm font-bold mb-2" for="endDateTime">结束时间</label>
            <input v-model="filters.endDateTime" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-indigo-500 transition duration-300" id="endDateTime" type="datetime-local">
          </div>
          <div>
            <label class="block text-gray-300 text-sm font-bold mb-2" for="ipFilter">IP 地址</label>
            <input v-model="filters.ip" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-indigo-500 transition duration-300" id="ipFilter" type="text" placeholder="输入 IP 地址">
          </div>
          <div class="flex items-end">
            <button class="bg-indigo-500 text-white px-4 py-2 rounded hover:bg-indigo-700 transform hover:scale-105 transition duration-300 w-full" type="submit">
              应用过滤器 ✅
            </button>
          </div>
        </form>
      </div>

      <!-- 日志记录 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">日志记录 📜</h2>
        <table class="min-w-full bg-gray-800">
          <thead>
          <tr>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">源 IP</th>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">HTTP 方法</th>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">目标 IP</th>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">URL</th>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">时间戳</th>
            <th class="py-2 px-4 border-b-2 border-gray-700 text-left">操作</th>
          </tr>
          </thead>
          <tbody>
          <tr v-for="log in logs" :key="log._id" class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ log.client_ip }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ log.method || 'N/A' }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ log.target_ip }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ log.url || 'N/A' }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ formatTimestamp(log.timestamp) }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">
              <button @click="viewDetails(log)" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transform hover:scale-105 transition duration-300">查看详情 🔍</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- 统计分析 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">统计分析 📊</h2>
        <div class="h-64 bg-gray-700 flex items-center justify-center animate-pulse">
          <span class="text-gray-400">图表占位符 📊</span>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />

    <!-- 详情模态框 -->
    <div v-if="selectedLog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-gray-800 p-6 rounded-lg shadow-md w-1/2">
        <h2 class="text-2xl font-bold mb-4">流量详情</h2>
        <pre class="text-white">{{ JSON.stringify(selectedLog, null, 2) }}</pre>
        <button @click="selectedLog = null" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700 mt-4">关闭</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import axios from 'axios';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';

export default {
  name: 'LogAnalysis',
  components: {
    HeaderPage,
    FooterPage
  },
  setup() {
    const logs = ref([]);
    const filters = ref({
      startDateTime: '',
      endDateTime: '',
      ip: ''
    });
    const selectedLog = ref(null);

    const formatToRFC3339 = (datetime) => {
      if (!datetime) return '';
      const date = new Date(datetime);
      return date.toISOString();
    };

    const fetchLogs = async () => {
      try {
        const params = {
          startDateTime: formatToRFC3339(filters.value.startDateTime),
          endDateTime: formatToRFC3339(filters.value.endDateTime),
          ip: filters.value.ip
        };
        const response = await axios.get('http://172.20.2.226:8081/logs', { params });
        logs.value = response.data;
      } catch (error) {
        console.error('获取日志失败:', error);
      }
    };

    const applyFilter = () => {
      fetchLogs();
    };

    const formatTimestamp = (timestamp) => {
      return new Date(timestamp).toLocaleString();
    };

    const viewDetails = (log) => {
      selectedLog.value = log;
    };

    return { logs, filters, applyFilter, formatTimestamp, viewDetails, selectedLog };
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out;
}
</style>
