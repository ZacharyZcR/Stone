<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：仪表盘标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#2080f0">
          <n-icon size="48">
            <bar-chart-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="primary">实时仪表盘</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            WAF防护状态监控与流量分析
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：核心指标统计 -->
    <n-card title="今日防护统计" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-button type="primary" @click="fetchMetrics" :loading="refreshing">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新数据
        </n-button>
      </template>
      
      <n-grid :cols="4" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="成功请求" :value="metrics.success_requests" suffix="次">
              <template #prefix>
                <n-icon color="#18a058" size="24">
                  <checkmark-circle-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="黑名单拦截" :value="metrics.blacklist_requests" suffix="次">
              <template #prefix>
                <n-icon color="#d03050" size="24">
                  <ban-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="规则拦截" :value="metrics.rules_requests" suffix="次">
              <template #prefix>
                <n-icon color="#f0a020" size="24">
                  <shield-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="拦截率" :value="blockRate" suffix="%">
              <template #prefix>
                <n-icon color="#2080f0" size="24">
                  <analytics-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第三行：双列布局 - 图表和实时活动 -->
    <n-grid :cols="2" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
      <!-- 左侧：流量趋势图表 -->
      <n-grid-item>
        <n-card title="最近7天流量趋势" size="large" style="height: 480px;">
          <template #header-extra>
            <n-space>
              <n-tag type="success" size="medium">实时更新</n-tag>
              <n-button size="small" quaternary @click="fetchWeeklyMetrics">
                <template #icon>
                  <n-icon><refresh-outline /></n-icon>
                </template>
              </n-button>
            </n-space>
          </template>
          
          <div style="height: 380px;">
            <BarChart v-if="chartData.labels.length" :chartData="chartData" :chartOptions="chartOptions" />
            <n-empty v-else description="暂无图表数据" style="height: 100%; display: flex; align-items: center; justify-content: center;">
              <template #icon>
                <n-icon size="48" color="#d0d0d0">
                  <bar-chart-outline />
                </n-icon>
              </template>
            </n-empty>
          </div>
        </n-card>
      </n-grid-item>

      <!-- 右侧：实时活动日志 -->
      <n-grid-item>
        <n-card title="实时防护日志" size="large" style="height: 480px;">
          <template #header-extra>
            <n-badge :value="activityLogs.length" type="warning">
              <n-icon size="20"><shield-checkmark-outline /></n-icon>
            </n-badge>
          </template>
          
          <div style="height: 380px; overflow-y: auto; padding-right: 8px;">
            <n-timeline size="large">
              <n-timeline-item 
                v-for="(log, index) in activityLogs" 
                :key="index"
                :type="getLogType(log.type)"
                :title="log.title"
                :content="log.message"
                :time="log.time"
              />
            </n-timeline>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 第四行：威胁分析 -->
    <n-card title="威胁情报分析" size="large" style="margin-bottom: 24px;">
      <n-grid :cols="3" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded style="height: 180px;">
            <n-space direction="vertical" align="center" justify="center" style="height: 100%;">
              <n-avatar size="large" color="#d03050">
                <n-icon size="28"><bug-outline /></n-icon>
              </n-avatar>
              <n-space direction="vertical" size="small" align="center">
                <n-text strong style="font-size: 18px;">高危攻击</n-text>
                <n-statistic :value="threatStats.high" suffix="次" style="font-size: 24px;" />
                <n-text depth="3" style="font-size: 12px;">SQL注入、XSS等</n-text>
              </n-space>
            </n-space>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card hoverable embedded style="height: 180px;">
            <n-space direction="vertical" align="center" justify="center" style="height: 100%;">
              <n-avatar size="large" color="#f0a020">
                <n-icon size="28"><warning-outline /></n-icon>
              </n-avatar>
              <n-space direction="vertical" size="small" align="center">
                <n-text strong style="font-size: 18px;">中危攻击</n-text>
                <n-statistic :value="threatStats.medium" suffix="次" style="font-size: 24px;" />
                <n-text depth="3" style="font-size: 12px;">异常请求、扫描</n-text>
              </n-space>
            </n-space>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card hoverable embedded style="height: 180px;">
            <n-space direction="vertical" align="center" justify="center" style="height: 100%;">
              <n-avatar size="large" color="#18a058">
                <n-icon size="28"><shield-checkmark-outline /></n-icon>
              </n-avatar>
              <n-space direction="vertical" size="small" align="center">
                <n-text strong style="font-size: 18px;">正常流量</n-text>
                <n-statistic :value="threatStats.normal" suffix="次" style="font-size: 24px;" />
                <n-text depth="3" style="font-size: 12px;">合法访问请求</n-text>
              </n-space>
            </n-space>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第五行：快速操作面板 -->
    <n-card title="快速操作" size="large">
      <n-space size="large" justify="center">
        <n-button-group size="large">
          <n-button type="primary" @click="$router.push('/custom-rule')">
            <template #icon>
              <n-icon size="20"><add-outline /></n-icon>
            </template>
            添加规则
          </n-button>
          
          <n-button @click="$router.push('/log-analysis')">
            <template #icon>
              <n-icon size="20"><document-text-outline /></n-icon>
            </template>
            查看日志
          </n-button>
          
          <n-button @click="$router.push('/attacker-profile')">
            <template #icon>
              <n-icon size="20"><eye-outline /></n-icon>
            </template>
            攻击者画像
          </n-button>
          
          <n-button @click="$router.push('/system-configuration')">
            <template #icon>
              <n-icon size="20"><settings-outline /></n-icon>
            </template>
            系统配置
          </n-button>
        </n-button-group>
      </n-space>
    </n-card>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import BarChart from './BarChart.vue'
import api from '../api/axiosInstance'
import {
  BarChartOutline,
  CheckmarkCircleOutline,
  BanOutline,
  ShieldOutline,
  AnalyticsOutline,
  RefreshOutline,
  ShieldCheckmarkOutline,
  BugOutline,
  WarningOutline,
  AddOutline,
  DocumentTextOutline,
  EyeOutline,
  SettingsOutline
} from '@vicons/ionicons5'

export default {
  name: 'WAFDashboard',
  components: {
    BarChart,
    BarChartOutline,
    CheckmarkCircleOutline,
    BanOutline,
    ShieldOutline,
    AnalyticsOutline,
    RefreshOutline,
    ShieldCheckmarkOutline,
    BugOutline,
    WarningOutline,
    AddOutline,
    DocumentTextOutline,
    EyeOutline,
    SettingsOutline
  },
  setup() {
    const metrics = ref({
      success_requests: 0,
      blacklist_requests: 0,
      rules_requests: 0
    })

    const refreshing = ref(false)
    const activityLogs = ref([])
    const threatStats = ref({
      high: 0,
      medium: 0,
      normal: 0
    })

    const blockRate = computed(() => {
      const total = metrics.value.success_requests + metrics.value.blacklist_requests + metrics.value.rules_requests
      if (total === 0) return 0
      const blocked = metrics.value.blacklist_requests + metrics.value.rules_requests
      return Math.round((blocked / total) * 100)
    })

    const fetchMetrics = async () => {
      refreshing.value = true
      try {
        const response = await api.get('/firewall/metrics')
        if (Array.isArray(response.data) && response.data.length > 0) {
          const latestMetrics = response.data[0]
          // 注意：后端返回的字段名是snake_case
          const successRequests = latestMetrics.success_requests || 0
          const blacklistRequests = latestMetrics.blacklist_requests || 0
          const rulesRequests = latestMetrics.rules_requests || 0
          
          metrics.value = {
            success_requests: successRequests,
            blacklist_requests: blacklistRequests,
            rules_requests: rulesRequests
          }

          // 更新威胁统计
          threatStats.value = {
            high: rulesRequests, // 规则拦截视为高危
            medium: blacklistRequests, // 黑名单拦截视为中危
            normal: successRequests // 成功请求视为正常
          }
        }
      } catch (error) {
        console.error('获取防火墙指标失败:', error)
        // 保持默认数据
      } finally {
        refreshing.value = false
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
      try {
        // 从API获取最新日志数据
        const response = await api.get('/logs', {
          params: { page: 1, limit: 10 }
        })
        
        if (response.data && response.data.logs) {
          activityLogs.value = response.data.logs.map(log => {
            const isBlocked = log.status !== 'success'
            const timeAgo = formatTimeAgo(new Date(log.timestamp))
            
            return {
              type: isBlocked ? 'error' : 'success',
              title: isBlocked ? '攻击拦截' : '正常访问', 
              message: `${log.client_ip} - ${log.method} ${log.url}${log.error ? ' (' + log.error + ')' : ''}`,
              time: timeAgo
            }
          }).slice(0, 6) // 只显示最新6条
        } else {
          // 如果API无数据，使用少量默认数据
          activityLogs.value = [
            { 
              type: 'info', 
              title: '系统启动', 
              message: 'Stone WAF 防护系统已启动', 
              time: '刚才' 
            }
          ]
        }
      } catch (error) {
        console.error('获取活动日志失败:', error)
        // 使用简化的默认数据
        activityLogs.value = [
          { 
            type: 'info', 
            title: '系统运行中', 
            message: '暂无日志数据，系统正常运行', 
            time: '刚才' 
          }
        ]
      }
    }
    
    // 格式化时间为相对时间
    const formatTimeAgo = (date) => {
      const now = new Date()
      const diffMs = now - date
      const diffMins = Math.floor(diffMs / 60000)
      const diffHours = Math.floor(diffMins / 60)
      const diffDays = Math.floor(diffHours / 24)
      
      if (diffMins < 1) return '刚才'
      if (diffMins < 60) return `${diffMins}分钟前`
      if (diffHours < 24) return `${diffHours}小时前`
      return `${diffDays}天前`
    }

    const getLogType = (type) => {
      const typeMap = {
        error: 'error',
        warning: 'warning', 
        info: 'info',
        success: 'success'
      }
      return typeMap[type] || 'default'
    }

    const weeklyMetrics = ref([])

    const fetchWeeklyMetrics = async () => {
      try {
        const endDate = new Date()
        endDate.setHours(endDate.getHours() + 8) // 调整为北京时间
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

    onMounted(() => {
      fetchMetrics().then(() => {
        animateNumbers()
      })
      fetchWeeklyMetrics()
      fetchActivityLogs()
    })

    return {
      metrics,
      refreshing,
      activityLogs,
      threatStats,
      blockRate,
      chartData,
      chartOptions: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: true,
            position: 'top'
          }
        },
        scales: {
          y: {
            beginAtZero: true
          }
        }
      },
      fetchMetrics,
      fetchWeeklyMetrics,
      getLogType
    }
  }
}
</script>

<style scoped>
</style>
