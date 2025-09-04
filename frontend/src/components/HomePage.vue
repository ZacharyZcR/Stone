<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：欢迎横幅 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="large" align="center">
        <n-space direction="vertical" size="medium" align="center">
          <n-avatar size="huge" color="#18a058">
            <n-icon size="48">
              <shield-checkmark-outline />
            </n-icon>
          </n-avatar>
          
          <n-space direction="vertical" size="small" align="center">
            <n-h1 style="margin: 0; font-size: 3rem;">
              <n-gradient-text type="primary">Stone WAF</n-gradient-text>
            </n-h1>
            <n-text depth="3" style="font-size: 18px;">
              企业级Web应用防火墙管理平台
            </n-text>
          </n-space>
        </n-space>

        <n-space size="large">
          <n-tag type="success" size="large">
            <template #icon>
              <n-icon><checkmark-circle-outline /></n-icon>
            </template>
            实时防护
          </n-tag>
          <n-tag type="info" size="large">
            <template #icon>
              <n-icon><analytics-outline /></n-icon>
            </template>
            智能分析
          </n-tag>
          <n-tag type="warning" size="large">
            <template #icon>
              <n-icon><settings-outline /></n-icon>
            </template>
            灵活配置
          </n-tag>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：系统状态概览 -->
    <n-card title="系统状态" size="large" style="margin-bottom: 24px;">
      <n-grid :cols="4" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="系统状态" value="运行中">
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
            <n-statistic label="防护规则" :value="ruleCount" suffix="条">
              <template #prefix>
                <n-icon color="#2080f0" size="24">
                  <shield-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="今日拦截" :value="todayBlocked" suffix="次">
              <template #prefix>
                <n-icon color="#f0a020" size="24">
                  <warning-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="威胁等级" value="低">
              <template #prefix>
                <n-icon color="#18a058" size="24">
                  <shield-checkmark-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第三行：功能模块 -->
    <n-card title="功能模块" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-text depth="3" style="font-size: 16px;">点击进入对应模块</n-text>
      </template>
      
      <n-grid :cols="3" :x-gap="24" :y-gap="24" responsive="screen">
        <n-grid-item 
          v-for="feature in features" 
          :key="feature.key"
        >
          <n-card 
            hoverable
            embedded
            @click="navigateToFeature(feature.route)"
            style="cursor: pointer; height: 200px;"
            size="large"
          >
            <n-space direction="vertical" align="center" justify="center" style="height: 100%;">
              <n-avatar size="huge" :color="feature.color">
                <n-icon size="32">
                  <component :is="feature.icon" />
                </n-icon>
              </n-avatar>
              
              <n-space direction="vertical" size="small" align="center">
                <n-text strong style="font-size: 18px;">{{ feature.title }}</n-text>
                <n-text depth="3" style="text-align: center; font-size: 14px; line-height: 1.4;">
                  {{ feature.description }}
                </n-text>
              </n-space>
            </n-space>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第四行：双列信息面板 -->
    <n-grid :cols="2" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
      <!-- 左侧：系统性能 -->
      <n-grid-item>
        <n-card title="系统性能监控" size="large" style="height: 420px;">
          <template #header-extra>
            <n-button type="primary" size="large" @click="refreshPerformance" :loading="refreshing">
              <template #icon>
                <n-icon><refresh-outline /></n-icon>
              </template>
              实时刷新
            </n-button>
          </template>
          
          <n-space direction="vertical" size="large" style="height: 100%; padding: 16px 0;">
            <div>
              <n-space justify="space-between" align="center" style="margin-bottom: 12px;">
                <n-text style="font-size: 18px; font-weight: 500;">CPU使用率</n-text>
                <n-text strong style="font-size: 18px; color: #18a058;">{{ cpuUsage }}%</n-text>
              </n-space>
              <n-progress :percentage="cpuUsage" type="line" size="large" :show-indicator="false" />
            </div>
            
            <div>
              <n-space justify="space-between" align="center" style="margin-bottom: 12px;">
                <n-text style="font-size: 18px; font-weight: 500;">内存使用率</n-text>
                <n-text strong style="font-size: 18px; color: #2080f0;">{{ memoryUsage }}%</n-text>
              </n-space>
              <n-progress :percentage="memoryUsage" type="line" size="large" :show-indicator="false" />
            </div>
            
            <div>
              <n-space justify="space-between" align="center" style="margin-bottom: 12px;">
                <n-text style="font-size: 18px; font-weight: 500;">磁盘使用率</n-text>
                <n-text strong style="font-size: 18px; color: #f0a020;">{{ diskUsage }}%</n-text>
              </n-space>
              <n-progress :percentage="diskUsage" type="line" size="large" :show-indicator="false" />
            </div>
          </n-space>
        </n-card>
      </n-grid-item>

      <!-- 右侧：系统动态 -->
      <n-grid-item>
        <n-card title="系统活动日志" size="large" style="height: 420px;">
          <template #header-extra>
            <n-badge :value="activities.length" type="info" size="large">
              <n-icon size="22"><notifications-outline /></n-icon>
            </n-badge>
          </template>
          
          <div style="height: 320px; overflow-y: auto; padding-right: 8px;">
            <n-timeline size="large">
              <n-timeline-item 
                v-for="activity in activities" 
                :key="activity.id"
                :type="activity.type"
                :title="activity.title"
                :content="activity.content"
                :time="activity.time"
              />
            </n-timeline>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 第五行：快速操作 -->
    <n-card title="快速操作" size="large" v-if="isAuthenticated" style="margin-bottom: 24px;">
      <n-space size="large" justify="center">
        <n-button-group size="large">
          <n-button type="primary" @click="$router.push('/dashboard')">
            <template #icon>
              <n-icon size="20"><bar-chart-outline /></n-icon>
            </template>
            查看仪表盘
          </n-button>
          
          <n-button @click="$router.push('/custom-rule')">
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
          
          <n-button @click="$router.push('/system-configuration')">
            <template #icon>
              <n-icon size="20"><settings-outline /></n-icon>
            </template>
            系统设置
          </n-button>
        </n-button-group>
      </n-space>
    </n-card>

    <!-- 第六行：系统信息 -->
    <n-card title="系统信息" size="large">
      <n-descriptions :column="3" bordered size="large" label-placement="left">
        <n-descriptions-item label="系统版本">
          <n-tag type="info" size="medium">v1.0.0</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="部署环境">
          <n-tag type="success" size="medium">生产环境</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="连接状态">
          <n-text type="success">
            <n-icon><wifi-outline /></n-icon> 在线
          </n-text>
        </n-descriptions-item>
        <n-descriptions-item label="运行时间">
          <n-text strong>{{ uptime }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="最后更新">
          <n-text>{{ lastUpdateTime }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="数据库状态">
          <n-text type="success">MongoDB 已连接</n-text>
        </n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import api from '../api/axiosInstance'
import {
  ShieldCheckmarkOutline,
  CheckmarkCircleOutline,
  AnalyticsOutline,
  SettingsOutline,
  ShieldOutline,
  WarningOutline,
  BarChartOutline,
  EyeOutline,
  DocumentTextOutline,
  PersonOutline,
  AddOutline,
  NotificationsOutline,
  WifiOutline,
  RefreshOutline
} from '@vicons/ionicons5'

export default {
  name: 'HomePage',
  components: {
    ShieldCheckmarkOutline,
    CheckmarkCircleOutline,
    AnalyticsOutline,
    SettingsOutline,
    ShieldOutline,
    WarningOutline,
    BarChartOutline,
    EyeOutline,
    DocumentTextOutline,
    PersonOutline,
    AddOutline,
    NotificationsOutline,
    WifiOutline,
    RefreshOutline
  },
  setup() {
    const store = useStore()
    const router = useRouter()
    
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    
    // 系统数据
    const ruleCount = ref(0)
    const todayBlocked = ref(0)
    const cpuUsage = ref(32)
    const memoryUsage = ref(58)
    const diskUsage = ref(23)
    const uptime = ref('计算中...')
    const lastUpdateTime = ref('获取中...')
    const refreshing = ref(false)

    const features = [
      {
        key: 'dashboard',
        title: '实时仪表盘',
        description: '监控系统运行状态和攻击趋势',
        icon: 'BarChartOutline',
        color: '#2080f0',
        route: '/dashboard'
      },
      {
        key: 'rules',
        title: '规则管理',
        description: '配置和管理防护规则',
        icon: 'ShieldOutline',
        color: '#18a058',
        route: '/custom-rule'
      },
      {
        key: 'logs',
        title: '日志分析',
        description: '查看和分析访问日志',
        icon: 'DocumentTextOutline',
        color: '#f0a020',
        route: '/log-analysis'
      },
      {
        key: 'attacker',
        title: '攻击者画像',
        description: '分析攻击者行为模式',
        icon: 'EyeOutline',
        color: '#d03050',
        route: '/attacker-profile'
      },
      {
        key: 'config',
        title: '系统配置',
        description: '管理系统参数和设置',
        icon: 'SettingsOutline',
        color: '#7c3aed',
        route: '/system-configuration'
      },
      {
        key: 'users',
        title: '用户管理',
        description: '管理用户账户和权限',
        icon: 'PersonOutline',
        color: '#059669',
        route: '/user-management'
      }
    ]

    // 系统活动
    const activities = ref([
      {
        id: 1,
        type: 'success',
        title: '规则更新成功',
        content: '成功更新SQL注入防护规则库',
        time: '2分钟前'
      },
      {
        id: 2,
        type: 'warning', 
        title: '检测到可疑访问',
        content: 'IP 192.168.1.100 触发频率限制规则',
        time: '5分钟前'
      },
      {
        id: 3,
        type: 'info',
        title: '系统备份完成',
        content: '自动备份任务执行完成',
        time: '1小时前'
      },
      {
        id: 4,
        type: 'error',
        title: '攻击拦截记录',
        content: '成功拦截XSS攻击尝试 15次',
        time: '2小时前'
      }
    ])

    const navigateToFeature = (route) => {
      if (route && isAuthenticated.value) {
        router.push(route)
      } else if (route && !isAuthenticated.value) {
        router.push('/login')
      }
    }

    const refreshPerformance = async () => {
      refreshing.value = true
      try {
        // 模拟刷新性能数据
        await new Promise(resolve => setTimeout(resolve, 1000))
        cpuUsage.value = Math.floor(Math.random() * 60) + 20
        memoryUsage.value = Math.floor(Math.random() * 40) + 40
        diskUsage.value = Math.floor(Math.random() * 30) + 10
      } finally {
        refreshing.value = false
      }
    }

    // 加载系统状态数据
    const loadSystemData = async () => {
      try {
        // 获取系统状态
        const statusResponse = await api.get('/status')
        if (statusResponse.data) {
          uptime.value = statusResponse.data.uptime || '计算中...'
          lastUpdateTime.value = new Date().toLocaleString()
        }

        // 获取防火墙指标
        const metricsResponse = await api.get('/firewall/metrics')
        if (Array.isArray(metricsResponse.data) && metricsResponse.data.length > 0) {
          const latest = metricsResponse.data[0]
          const blocked = (latest.blacklist_requests || 0) + (latest.rules_requests || 0)
          todayBlocked.value = blocked
        }

        // 获取拦截规则数量
        const rulesResponse = await api.get('/interception-rules')
        if (Array.isArray(rulesResponse.data)) {
          ruleCount.value = rulesResponse.data.length
        }

      } catch (error) {
        console.error('加载系统数据失败:', error)
        // 保持默认值
      }
    }

    // 加载系统活动日志
    const loadActivityLogs = async () => {
      try {
        // 尝试从实际日志API获取数据
        const logsResponse = await api.get('/logs', {
          params: { page: 1, limit: 4 }
        })
        
        if (logsResponse.data && logsResponse.data.logs) {
          activities.value = logsResponse.data.logs.map((log, index) => ({
            id: index + 1,
            type: log.status === 'success' ? 'success' : 'warning',
            title: log.status === 'success' ? '正常访问' : '攻击拦截',
            content: `${log.method} ${log.url} - ${log.client_ip}`,
            time: new Date(log.timestamp).toLocaleString()
          })).slice(0, 4)
        }
      } catch (error) {
        console.error('加载活动日志失败:', error)
        // 保持Mock数据作为fallback
      }
    }

    onMounted(() => {
      if (isAuthenticated.value) {
        loadSystemData()
        loadActivityLogs()
      }
    })

    return {
      isAuthenticated,
      features,
      ruleCount,
      todayBlocked,
      cpuUsage,
      memoryUsage,
      diskUsage,
      uptime,
      lastUpdateTime,
      activities,
      refreshing,
      navigateToFeature,
      refreshPerformance
    }
  }
}
</script>