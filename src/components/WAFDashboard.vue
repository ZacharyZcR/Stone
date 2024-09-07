<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen animate-fade-in">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- 卡片 1 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-green-400 text-3xl">✅</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">成功请求</h3>
              <p class="text-2xl animate-number-scroll" :data-target="metrics.success_requests">{{ metrics.success_requests }}</p>
            </div>
          </div>
        </div>
        <!-- 卡片 2 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-red-400 text-3xl">🚫</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">黑名单请求</h3>
              <p class="text-2xl animate-number-scroll" :data-target="metrics.blacklist_requests">{{ metrics.blacklist_requests }}</p>
            </div>
          </div>
        </div>
        <!-- 卡片 3 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-blue-400 text-3xl">📏</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">规则拦截</h3>
              <p class="text-2xl animate-number-scroll" :data-target="metrics.rules_requests">{{ metrics.rules_requests }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">最近7天流量分析</h2>
        <div class="h-80"> <!-- 增加高度 -->
          <BarChart :chartData="chartData" :chartOptions="chartOptions" />
        </div>
      </div>

      <!-- 活动日志 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">最近活动</h2>
        <ul class="space-y-4">
          <li v-for="(log, index) in activityLogs" :key="index" class="flex items-start animate-fade-in-up">
            <span class="text-blue-400 text-2xl mr-4">🕒</span>
            <div>
              <p class="font-bold">{{ log.time }}</p>
              <p>{{ log.message }}</p>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import HeaderPage from './HeaderPage.vue'
import FooterPage from './FooterPage.vue'
import BarChart from './BarChart.vue'
import api from '../api/axiosInstance'

export default {
  name: 'WAFDashboard',
  components: {
    HeaderPage,
    FooterPage,
    BarChart
  },
  setup() {
    const metrics = ref({
      success_requests: 0,
      blacklist_requests: 0,
      rules_requests: 0
    })

    const activityLogs = ref([])

    const fetchMetrics = async () => {
      try {
        const response = await api.get('/firewall/metrics')
        if (Array.isArray(response.data) && response.data.length > 0) {
          const latestMetrics = response.data[0]
          metrics.value = {
            success_requests: latestMetrics.success_requests || 0,
            blacklist_requests: latestMetrics.blacklist_requests || 0,
            rules_requests: latestMetrics.rules_requests || 0
          }
        }
      } catch (error) {
        console.error('获取防火墙指标失败:', error)
      }
    }

    const animateNumbers = () => {
      const elements = document.querySelectorAll('.animate-number-scroll');
      elements.forEach(el => {
        const target = parseInt(el.getAttribute('data-target'), 10);
        let count = 0;
        const increment = target / 100;
        const updateCount = () => {
          count += increment;
          if (count < target) {
            el.textContent = Math.floor(count);
            requestAnimationFrame(updateCount);
          } else {
            el.textContent = target;
          }
        };
        updateCount();
      });
    }

    const fetchActivityLogs = async () => {
      // 这里应该是从后端获取活动日志的逻辑
      // 现在我们使用模拟数据
      activityLogs.value = [
        { time: '10:30 AM', message: '检测到 SQL 注入攻击，已成功阻止。💪' },
        { time: '09:45 AM', message: '检测到 XSS 攻击，已成功阻止。🔒' },
        { time: '08:20 AM', message: '异常流量增加，正在监控中。👀' }
      ]
    }

    const weeklyMetrics = ref([])

    const fetchWeeklyMetrics = async () => {
      try {
        const endDate = new Date()
        const startDate = new Date(endDate)
        startDate.setDate(startDate.getDate() - 6)

        const response = await api.get('/firewall/metrics', {
          params: {
            start_date: startDate.toISOString().split('T')[0],
            end_date: endDate.toISOString().split('T')[0]
          }
        })
        weeklyMetrics.value = response.data
      } catch (error) {
        console.error('获取每周指标失败:', error)
      }
    }

    const chartData = computed(() => {
      // 确保 weeklyMetrics 不为空
      if (weeklyMetrics.value.length === 0) {
        return {
          labels: [],
          datasets: []
        }
      }

      const labels = weeklyMetrics.value.map(m => m.date)
      const datasets = [
        {
          label: '成功请求',
          backgroundColor: 'rgba(16, 185, 129, 0.7)', // 半透明绿色
          borderColor: '#10B981',
          borderWidth: 2,
          data: weeklyMetrics.value.map(m => m.success_requests)
        },
        {
          label: '黑名单请求',
          backgroundColor: 'rgba(239, 68, 68, 0.7)', // 半透明红色
          borderColor: '#EF4444',
          borderWidth: 2,
          data: weeklyMetrics.value.map(m => m.blacklist_requests)
        },
        {
          label: '规则拦截',
          backgroundColor: 'rgba(59, 130, 246, 0.7)', // 半透明蓝色
          borderColor: '#3B82F6',
          borderWidth: 2,
          data: weeklyMetrics.value.map(m => m.rules_requests)
        }
      ]

      return { labels, datasets }
    })

    const chartOptions = {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          position: 'top',
          labels: {
            font: {
              size: 14
            },
            color: '#FFFFFF' // 白色文字
          }
        },
        tooltip: {
          mode: 'index',
          intersect: false,
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          titleColor: '#FFFFFF',
          bodyColor: '#FFFFFF',
          borderColor: '#FFFFFF',
          borderWidth: 1
        }
      },
      scales: {
        x: {
          grid: {
            color: 'rgba(255, 255, 255, 0.1)' // 淡白色网格线
          },
          ticks: {
            color: '#FFFFFF' // 白色文字
          }
        },
        y: {
          beginAtZero: true,
          grid: {
            color: 'rgba(255, 255, 255, 0.1)' // 淡白色网格线
          },
          ticks: {
            color: '#FFFFFF', // 白色文字
            callback: function(value) {
              return value.toLocaleString() // 格式化大数字
            }
          }
        }
      }
    }

    onMounted(() => {
      fetchMetrics().then(() => {
        animateNumbers()
      })
      fetchWeeklyMetrics()
      fetchActivityLogs()
    })

    return {
      metrics,
      activityLogs,
      chartData,
      chartOptions
    }
  }
}
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

.animate-fade-in {
  animation: fade-in 1s ease-out;
}

.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out;
}
</style>
