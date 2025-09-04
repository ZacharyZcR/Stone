<template>
  <div class="bg-gray-800 p-8 rounded-lg shadow-md">
    <h2 class="text-3xl font-bold mb-6 text-white">统计分析 📊</h2>

    <!-- 过滤器部分 -->
    <div class="mb-8 grid grid-cols-1 md:grid-cols-4 gap-6">
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-2" for="startDate">开始日期</label>
        <input v-model="startDate" type="date" id="startDate" class="w-full bg-gray-700 text-white rounded-md p-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-300">
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-2" for="endDate">结束日期</label>
        <input v-model="endDate" type="date" id="endDate" class="w-full bg-gray-700 text-white rounded-md p-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-300">
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-300 mb-2" for="accessType">访问类型</label>
        <select v-model="accessType" id="accessType" class="w-full bg-gray-700 text-white rounded-md p-3 focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-300">
          <option value="">全部</option>
          <option value="passed">正常访问</option>
          <option value="blocked">攻击行为</option>
        </select>
      </div>
      <div class="flex items-end">
        <button @click="fetchData" class="w-full bg-blue-500 text-white px-6 py-3 rounded-md hover:bg-blue-600 transform hover:scale-105 transition duration-300 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50">
          更新数据
        </button>
      </div>
    </div>

    <!-- 图表和表格部分 -->
    <div class="flex flex-col lg:flex-row space-y-8 lg:space-y-0 lg:space-x-8">
      <div class="w-full lg:w-1/2 bg-gray-700 p-6 rounded-lg shadow-inner">
        <h3 class="text-xl font-semibold mb-4 text-gray-200">访问分布图</h3>
        <div class="h-80">
          <Pie v-if="computedChartData.labels.length"
               :data="computedChartData"
               :options="chartOptions"
               :key="chartKey" />
          <div v-else class="h-full flex items-center justify-center text-gray-400">
            暂无数据
          </div>
        </div>
      </div>
      <div class="w-full lg:w-1/2 bg-gray-700 p-6 rounded-lg shadow-inner">
        <h3 class="text-xl font-semibold mb-4 text-gray-200">TOP 10 访问 IP</h3>
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
            <tr>
              <th class="pb-3 text-gray-300 border-b border-gray-600">IP地址</th>
              <th class="pb-3 text-gray-300 border-b border-gray-600">访问次数</th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="(item, index) in top10Data" :key="index" class="hover:bg-gray-600 transition duration-200">
              <td class="py-3 text-gray-200">{{ item._id }}</td>
              <td class="py-3 text-gray-200">{{ item.count }}</td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { Pie } from 'vue-chartjs';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import api from '../api/axiosInstance';

ChartJS.register(ArcElement, Tooltip, Legend);

export default {
  name: 'StatisticsChart',
  components: { Pie },
  setup() {
    const startDate = ref(getBeijingDate(new Date(Date.now() - 7 * 24 * 60 * 60 * 1000)));
    const endDate = ref(getBeijingDate(new Date()));
    const accessType = ref('');
    const chartData = ref({
      labels: [],
      datasets: [{
        backgroundColor: [],
        data: []
      }]
    });
    const chartKey = ref(0);
    const top10Data = ref([]);

    // 将日期转换为北京时间的函数
    function getBeijingDate(date) {
      const beijingTime = new Date(date.getTime() + (8 * 60 * 60 * 1000));
      return beijingTime.toISOString().substr(0, 10);
    }

    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: false
        },
        tooltip: {
          callbacks: {
            label: function(context) {
              let label = context.label || '';
              if (label) {
                label += ': ';
              }
              if (context.parsed !== null) {
                label += context.parsed + ' 次访问';
              }
              return label;
            }
          }
        }
      }
    };

    const computedChartData = computed(() => {
      return {
        labels: chartData.value.labels,
        datasets: [{
          backgroundColor: generateColors(chartData.value.labels.length),
          data: chartData.value.datasets[0].data
        }]
      };
    });

    const fetchData = async () => {
      try {
        // 将日期转换为 YYYY-MM-DD 格式
        const formatDate = (date) => {
          const d = new Date(date);
          return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
        };

        const startDateTime = formatDate(startDate.value);
        const endDateTime = formatDate(endDate.value);

        const response = await api.get('/ip-stats', {
          params: {
            startDateTime: startDateTime,
            endDateTime: endDateTime,
            status: accessType.value
          }
        });
        console.log('API Response:', response.data);
        
        // 安全处理API响应数据
        if (response.data && response.data.ipStats && Array.isArray(response.data.ipStats)) {
          top10Data.value = response.data.ipStats.slice(0, 10);
        } else {
          console.warn('IP统计数据为空或格式不正确，使用默认数据');
          top10Data.value = [];
        }

        chartData.value = {
          labels: top10Data.value.map(item => item._id || '未知IP'),
          datasets: [{
            data: top10Data.value.map(item => item.count || 0)
          }]
        };
        chartKey.value += 1; // 强制重新渲染图表
      } catch (error) {
        console.error('获取IP统计数据失败:', error);
        // API调用失败时设置空数据
        top10Data.value = [];
        chartData.value = {
          labels: [],
          datasets: [{
            data: []
          }]
        };
      }
    };

    const generateColors = (count) => {
      const colors = [];
      for (let i = 0; i < count; i++) {
        colors.push(`hsl(${(i * 360) / count}, 70%, 50%)`);
      }
      return colors;
    };

    onMounted(fetchData);

    return {
      startDate,
      endDate,
      accessType,
      computedChartData,
      chartOptions,
      fetchData,
      chartKey,
      top10Data
    };
  }
};
</script>

<style scoped>
/* 可以添加一些自定义样式 */
.animate-number-scroll {
  transition: all 0.5s ease-out;
}
</style>
