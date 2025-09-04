<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：日志分析标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#52c41a">
          <n-icon size="48">
            <document-text-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="success">日志分析中心</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            WAF访问日志分析与威胁检测
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：日志统计概览 -->
    <n-card title="今日日志统计" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-button type="primary" @click="fetchLogs" :loading="refreshing">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新数据
        </n-button>
      </template>
      
      <n-grid :cols="5" :x-gap="24" responsive="screen">
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="总日志数" :value="todayLogs" suffix="条">
              <template #prefix>
                <n-icon color="#2080f0" size="24">
                  <document-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
        
        <n-grid-item>
          <n-card hoverable embedded>
            <n-statistic label="拦截日志" :value="blockedLogs" suffix="条">
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
            <n-statistic label="正常访问" :value="normalLogs" suffix="条">
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
            <n-statistic label="可疑IP" :value="suspiciousIps" suffix="个">
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
            <n-statistic label="平均响应时间" :value="avgResponseTime" suffix="ms">
              <template #prefix>
                <n-icon color="#722ed1" size="24">
                  <time-outline />
                </n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
      </n-grid>
    </n-card>

    <!-- 第三行：筛选工具栏 -->
    <n-card title="筛选与搜索" size="large" style="margin-bottom: 24px;">
      <n-space size="large">
        <n-space size="medium" align="center">
          <n-text strong>时间范围：</n-text>
          <n-date-picker
            v-model:value="dateRange"
            type="daterange"
            clearable
            @update:value="handleDateRangeChange"
          />
        </n-space>

        <n-space size="medium" align="center">
          <n-text strong>日志级别：</n-text>
          <n-select
            v-model:value="logLevel"
            style="width: 120px"
            :options="[
              { label: '全部', value: '' },
              { label: '错误', value: 'error' },
              { label: '警告', value: 'warning' },
              { label: '信息', value: 'info' }
            ]"
            @update:value="handleFilter"
          />
        </n-space>

        <n-space size="medium" align="center">
          <n-text strong>状态码：</n-text>
          <n-select
            v-model:value="statusFilter"
            style="width: 120px"
            :options="[
              { label: '全部', value: '' },
              { label: '2xx', value: '2' },
              { label: '4xx', value: '4' },
              { label: '5xx', value: '5' }
            ]"
            @update:value="handleFilter"
          />
        </n-space>

        <n-space size="medium" align="center">
          <n-text strong>搜索：</n-text>
          <n-input
            v-model:value="searchQuery"
            placeholder="搜索IP、URL或User-Agent"
            style="width: 250px"
            clearable
            @update:value="handleSearch"
          >
            <template #prefix>
              <n-icon size="16" color="#808080">
                <search-outline />
              </n-icon>
            </template>
          </n-input>
        </n-space>

        <n-space>
          <n-button @click="applyFilters">
            <template #icon>
              <n-icon><funnel-outline /></n-icon>
            </template>
            应用筛选
          </n-button>
          <n-button @click="clearFilters">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            清空筛选
          </n-button>
          <n-button type="primary" @click="exportLogs">
            <template #icon>
              <n-icon><cloud-download-outline /></n-icon>
            </template>
            导出日志
          </n-button>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第四行：双列布局 - 日志列表和趋势分析 -->
    <n-grid :cols="2" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
      <!-- 左侧：日志列表 -->
      <n-grid-item>
        <n-card title="访问日志" size="large" style="height: 700px;">
          <template #header-extra>
            <n-space>
              <n-tag :type="loading ? 'warning' : 'success'">
                {{ loading ? '加载中' : `共 ${totalCount} 条` }}
              </n-tag>
            </n-space>
          </template>
          
          <div style="height: 580px; overflow-y: auto;">
            <n-data-table
              :columns="columns"
              :data="displayLogs"
              :loading="loading"
              :pagination="false"
              striped
              size="small"
              virtual-scroll
              max-height="580"
              :scroll-x="1200"
            />
          </div>
          
          <template #footer>
            <n-space justify="space-between" align="center">
              <n-text depth="3">
                显示第 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, totalCount) }} 条，共 {{ totalCount }} 条
              </n-text>
              <n-pagination
                v-model:page="currentPage"
                v-model:page-size="pageSize"
                :item-count="totalCount"
                :page-sizes="[20, 50, 100]"
                show-size-picker
                show-quick-jumper
                @update:page="changePage"
                @update:page-size="changePageSize"
              />
            </n-space>
          </template>
        </n-card>
      </n-grid-item>

      <!-- 右侧：趋势分析 -->
      <n-grid-item>
        <n-card title="访问趋势分析" size="large" style="height: 700px;">
          <n-space direction="vertical" size="large" style="height: 580px;">
            <!-- 实时统计 -->
            <n-card embedded>
              <n-space direction="vertical" size="medium">
                <n-text strong>最近1小时统计</n-text>
                <n-grid :cols="2" :x-gap="16">
                  <n-grid-item>
                    <n-statistic label="请求总数" :value="recentHourRequests">
                      <template #suffix>条</template>
                    </n-statistic>
                  </n-grid-item>
                  <n-grid-item>
                    <n-statistic label="拦截次数" :value="recentHourBlocked">
                      <template #suffix>次</template>
                    </n-statistic>
                  </n-grid-item>
                </n-grid>
              </n-space>
            </n-card>

            <!-- 热门攻击类型 -->
            <n-card embedded>
              <n-space direction="vertical" size="medium">
                <n-text strong>热门攻击类型</n-text>
                <div v-if="topAttackTypes.length === 0" style="text-align: center; padding: 20px;">
                  <n-text depth="3">暂无攻击记录</n-text>
                </div>
                <n-list v-else hoverable>
                  <n-list-item v-for="attack in topAttackTypes" :key="attack.type">
                    <template #prefix>
                      <n-icon :color="getThreatColor(attack.level)" size="16">
                        <warning-outline />
                      </n-icon>
                    </template>
                    <n-space justify="space-between" style="width: 100%">
                      <n-text>{{ attack.type }}</n-text>
                      <n-tag :type="getThreatType(attack.level)" size="small">
                        {{ attack.count }} 次
                      </n-tag>
                    </n-space>
                  </n-list-item>
                </n-list>
              </n-space>
            </n-card>

            <!-- 高风险IP -->
            <n-card embedded>
              <n-space direction="vertical" size="medium">
                <n-text strong>高风险IP地址</n-text>
                <div v-if="riskIps.length === 0" style="text-align: center; padding: 20px;">
                  <n-text depth="3">暂无高风险IP</n-text>
                </div>
                <n-list v-else hoverable>
                  <n-list-item v-for="ip in riskIps" :key="ip.address">
                    <template #prefix>
                      <n-icon color="#d03050" size="16">
                        <shield-outline />
                      </n-icon>
                    </template>
                    <n-space justify="space-between" style="width: 100%">
                      <n-text code>{{ ip.address }}</n-text>
                      <n-space size="small">
                        <n-tag type="error" size="small">{{ ip.attacks }} 次</n-tag>
                        <n-button 
                          size="tiny" 
                          type="error" 
                          @click="addToBlacklist(ip.address)"
                          :loading="false"
                        >
                          拉黑
                        </n-button>
                      </n-space>
                    </n-space>
                  </n-list-item>
                </n-list>
              </n-space>
            </n-card>
          </n-space>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 第五行：日志详情模态框 -->
    <n-modal v-model:show="showDetailModal">
      <n-card
        style="width: 800px;"
        title="日志详情"
        :bordered="false"
        size="huge"
        role="dialog"
        aria-modal="true"
      >
        <template #header-extra>
          <n-button quaternary circle @click="showDetailModal = false">
            <template #icon>
              <n-icon><close-outline /></n-icon>
            </template>
          </n-button>
        </template>

        <n-descriptions v-if="selectedLog" :column="2" label-placement="left">
          <n-descriptions-item label="时间戳">
            {{ formatTimestamp(selectedLog.timestamp) }}
          </n-descriptions-item>
          <n-descriptions-item label="IP地址">
            <n-text code>{{ selectedLog.ip }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="请求方法">
            <n-tag :type="getMethodType(selectedLog.method)" size="small">
              {{ selectedLog.method }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="状态码">
            <n-tag :type="getStatusType(selectedLog.status)" size="small">
              {{ selectedLog.status }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="请求URL" :span="2">
            <n-text code>{{ selectedLog.url }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="User-Agent" :span="2">
            <n-text code style="word-break: break-all;">{{ selectedLog.user_agent }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="响应大小">
            {{ formatBytes(selectedLog.response_size) }}
          </n-descriptions-item>
          <n-descriptions-item label="响应时间">
            {{ selectedLog.response_time }}ms
          </n-descriptions-item>
          <n-descriptions-item v-if="selectedLog.threat_level" label="威胁等级">
            <n-tag :type="getThreatType(selectedLog.threat_level)" size="medium">
              {{ selectedLog.threat_level }}
            </n-tag>
          </n-descriptions-item>
        </n-descriptions>

        <template #footer>
          <n-space justify="end">
            <n-button @click="copyLog(selectedLog)">
              <template #icon>
                <n-icon><copy-outline /></n-icon>
              </template>
              复制日志
            </n-button>
            <n-button type="error" @click="addToBlacklist(selectedLog.ip)">
              <template #icon>
                <n-icon><shield-outline /></n-icon>
              </template>
              加入黑名单
            </n-button>
          </n-space>
        </template>
      </n-card>
    </n-modal>
  </div>
</template>

<script>
import { ref, computed, h, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import api from '../api/axiosInstance'
import {
  DocumentTextOutline,
  RefreshOutline,
  DocumentOutline,
  BanOutline,
  CheckmarkCircleOutline,
  WarningOutline,
  TimeOutline,
  SearchOutline,
  FunnelOutline,
  CloudDownloadOutline,
  CloseOutline,
  CopyOutline,
  ShieldOutline
} from '@vicons/ionicons5'

export default {
  name: 'LogAnalysis',
  components: {
    DocumentTextOutline,
    RefreshOutline,
    DocumentOutline,
    BanOutline,
    CheckmarkCircleOutline,
    WarningOutline,
    TimeOutline,
    SearchOutline,
    FunnelOutline,
    CloudDownloadOutline,
    CloseOutline,
    CopyOutline,
    ShieldOutline
  },
  setup() {
    const message = useMessage()
    const logs = ref([])
    const selectedLog = ref(null)
    const currentPage = ref(1)
    const pageSize = ref(20)
    const totalCount = ref(0)
    const loading = ref(false)
    const refreshing = ref(false)
    const showDetailModal = ref(false)

    // 筛选条件
    const dateRange = ref(null)
    const logLevel = ref('')
    const statusFilter = ref('')
    const searchQuery = ref('')

    // 统计数据
    const todayLogs = ref(0)
    const blockedLogs = ref(0)
    const normalLogs = ref(0)
    const suspiciousIps = ref(0)
    const avgResponseTime = ref(0)

    // 趋势分析数据
    const topAttackTypes = ref([])
    const riskIps = ref([])
    
    // 实时统计数据
    const recentHourRequests = ref(0)
    const recentHourBlocked = ref(0)

    // 计算统计数据
    const calculateStatistics = (allLogs) => {
      const now = new Date()
      const oneHourAgo = new Date(now.getTime() - 60 * 60 * 1000)
      const today = now.toDateString()
      
      // 最近1小时的数据
      const recentHourLogs = allLogs.filter(log => new Date(log.timestamp) > oneHourAgo)
      recentHourRequests.value = recentHourLogs.length
      recentHourBlocked.value = recentHourLogs.filter(log => log.threat_level === 'high').length
      
      // 今日统计
      todayLogs.value = allLogs.length
      blockedLogs.value = allLogs.filter(log => log.threat_level === 'high').length
      normalLogs.value = allLogs.filter(log => log.status === 200).length
      
      // 计算平均响应时间
      const validResponseTimes = allLogs.filter(log => log.response_time)
      avgResponseTime.value = validResponseTimes.length > 0 
        ? Math.round(validResponseTimes.reduce((sum, log) => sum + log.response_time, 0) / validResponseTimes.length)
        : 0

      // 分析风险IP
      const ipCounts = {}
      allLogs.forEach(log => {
        if (log.threat_level && log.threat_level !== 'low') {
          ipCounts[log.ip] = (ipCounts[log.ip] || 0) + 1
        }
      })
      
      suspiciousIps.value = Object.keys(ipCounts).length
      
      riskIps.value = Object.entries(ipCounts)
        .map(([address, attacks]) => ({ address, attacks }))
        .sort((a, b) => b.attacks - a.attacks)
        .slice(0, 5)

      // 分析攻击类型（基于URL模式）
      const attackPatterns = {
        'SQL注入': /\b(union|select|insert|drop|delete)\b/i,
        'XSS攻击': /<script|javascript:|onerror=|onload=/i,
        '路径遍历': /\.\.\/|\.\.\\|\.\./,
        '管理员访问': /\/admin|\/wp-admin|\/administrator/i,
        '异常请求': /\.(php|asp|jsp)$/i
      }

      const attackCounts = {}
      allLogs.forEach(log => {
        if (log.threat_level) {
          for (const [type, pattern] of Object.entries(attackPatterns)) {
            if (pattern.test(log.url)) {
              attackCounts[type] = (attackCounts[type] || 0) + 1
              break
            }
          }
        }
      })

      topAttackTypes.value = Object.entries(attackCounts)
        .map(([type, count]) => ({ 
          type, 
          count, 
          level: count > 20 ? 'high' : count > 10 ? 'medium' : 'low' 
        }))
        .sort((a, b) => b.count - a.count)
        .slice(0, 5)
    }

    // 显示的日志（API已经分页返回，直接使用）
    const displayLogs = computed(() => {
      return logs.value
    })

    // 表格列定义
    const columns = [
      {
        title: '时间',
        key: 'timestamp',
        width: 160,
        render(row) {
          return formatTimestamp(row.timestamp)
        }
      },
      {
        title: 'IP地址',
        key: 'ip',
        width: 130,
        render(row) {
          return h('n-text', { code: true }, row.ip)
        }
      },
      {
        title: '方法',
        key: 'method',
        width: 70,
        render(row) {
          return h(
            'n-tag',
            { type: getMethodType(row.method), size: 'small' },
            row.method
          )
        }
      },
      {
        title: '状态码',
        key: 'status',
        width: 80,
        render(row) {
          return h(
            'n-tag',
            { type: getStatusType(row.status), size: 'small' },
            row.status
          )
        }
      },
      {
        title: 'URL',
        key: 'url',
        width: 300,
        ellipsis: {
          tooltip: true
        }
      },
      {
        title: '响应时间',
        key: 'response_time',
        width: 90,
        render(row) {
          return `${row.response_time}ms`
        }
      },
      {
        title: '威胁等级',
        key: 'threat_level',
        width: 100,
        render(row) {
          if (!row.threat_level) return '-'
          return h(
            'n-tag',
            { type: getThreatType(row.threat_level), size: 'small' },
            row.threat_level
          )
        }
      },
      {
        title: '操作',
        key: 'actions',
        width: 120,
        render(row) {
          return h(
            'n-space',
            { size: 'small' },
            {
              default: () => [
                h(
                  'n-button',
                  {
                    size: 'small',
                    quaternary: true,
                    onClick: () => viewDetails(row)
                  },
                  '详情'
                ),
                h(
                  'n-button',
                  {
                    size: 'small',
                    type: 'error',
                    quaternary: true,
                    onClick: () => addToBlacklist(row.ip)
                  },
                  '拉黑'
                )
              ]
            }
          )
        }
      }
    ]

    // 数据转换函数 - 将后端数据格式转换为前端期望格式
    const transformLogData = (backendLogs) => {
      return backendLogs.map(log => ({
        id: log._id,
        timestamp: log.timestamp,
        ip: log.client_ip,
        method: log.method,
        url: log.url,
        status: log.status === 'success' ? 200 : (log.error ? 403 : 500),
        user_agent: log.headers?.['User-Agent']?.[0] || 'Unknown',
        response_size: Math.floor(Math.random() * 5000) + 500, // 暂时模拟，后端未提供
        response_time: Math.floor(Math.random() * 100) + 20, // 暂时模拟，后端未提供
        threat_level: log.error && log.error.includes('Blocked by rules') ? 'high' : (log.status === 'failed' ? 'medium' : null)
      }))
    }

    // 获取日志数据
    const fetchLogs = async () => {
      loading.value = true
      try {
        const response = await api.get('/logs', {
          params: {
            page: currentPage.value,
            limit: pageSize.value,
            level: logLevel.value,
            status: statusFilter.value,
            search: searchQuery.value,
            start_date: dateRange.value?.[0],
            end_date: dateRange.value?.[1]
          }
        })
        
        // 转换后端数据格式为前端期望格式
        const transformedLogs = transformLogData(response.data.logs || [])
        logs.value = transformedLogs
        totalCount.value = response.data.totalCount || 0
        
        // 计算统计数据
        calculateStatistics(transformedLogs)
        
      } catch (error) {
        console.error('获取日志失败:', error)
        message.error('获取日志数据失败')
        logs.value = []
        totalCount.value = 0
      } finally {
        loading.value = false
      }
    }

    // 查看详情
    const viewDetails = (log) => {
      selectedLog.value = log
      showDetailModal.value = true
    }

    // 复制日志
    const copyLog = (log) => {
      const logText = JSON.stringify(log, null, 2)
      navigator.clipboard.writeText(logText).then(() => {
        message.success('日志已复制到剪贴板')
      })
    }

    // 应用筛选
    const applyFilters = () => {
      currentPage.value = 1
      fetchLogs()
    }

    // 清空筛选
    const clearFilters = () => {
      dateRange.value = null
      logLevel.value = ''
      statusFilter.value = ''
      searchQuery.value = ''
      currentPage.value = 1
      fetchLogs()
    }

    // 分页处理
    const changePage = (page) => {
      currentPage.value = page
      fetchLogs()
    }

    const changePageSize = (size) => {
      pageSize.value = size
      currentPage.value = 1
      fetchLogs()
    }

    // 事件处理
    const handleDateRangeChange = () => {
      applyFilters()
    }

    const handleFilter = () => {
      applyFilters()
    }

    const handleSearch = () => {
      // 防抖处理
      clearTimeout(handleSearch.timer)
      handleSearch.timer = setTimeout(() => {
        applyFilters()
      }, 500)
    }

    // 导出日志
    const exportLogs = () => {
      const csvContent = [
        ['时间', 'IP地址', '方法', 'URL', '状态码', '响应时间', '威胁等级'].join(','),
        ...logs.value.map(log => [
          log.timestamp,
          log.ip,
          log.method,
          log.url,
          log.status,
          log.response_time,
          log.threat_level || ''
        ].join(','))
      ].join('\n')

      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
      const link = document.createElement('a')
      link.href = URL.createObjectURL(blob)
      link.download = `access_logs_${new Date().toISOString().split('T')[0]}.csv`
      link.click()
      
      message.success('日志导出成功')
    }

    // 添加到黑名单
    const addToBlacklist = async (ip) => {
      try {
        await api.post('/ip-control-rules', { ip, type: 'blacklist' })
        message.success(`IP ${ip} 已添加到黑名单`)
      } catch (error) {
        console.error('添加黑名单失败:', error)
        message.error('添加黑名单失败')
      }
    }

    // 格式化函数
    const formatTimestamp = (timestamp) => {
      return new Date(timestamp).toLocaleString('zh-CN')
    }

    const formatBytes = (bytes) => {
      if (bytes === 0) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }

    const getMethodType = (method) => {
      const types = {
        GET: 'info',
        POST: 'success',
        PUT: 'warning',
        DELETE: 'error'
      }
      return types[method] || 'default'
    }

    const getStatusType = (status) => {
      if (status >= 200 && status < 300) return 'success'
      if (status >= 400 && status < 500) return 'warning'
      if (status >= 500) return 'error'
      return 'info'
    }

    const getThreatType = (level) => {
      const types = {
        low: 'info',
        medium: 'warning',
        high: 'error'
      }
      return types[level] || 'default'
    }

    const getThreatColor = (level) => {
      const colors = {
        low: '#2080f0',
        medium: '#f0a020',
        high: '#d03050'
      }
      return colors[level] || '#666'
    }

    onMounted(() => {
      fetchLogs()
    })

    return {
      // 数据
      logs,
      displayLogs,
      selectedLog,
      currentPage,
      pageSize,
      totalCount,
      loading,
      refreshing,
      showDetailModal,
      dateRange,
      logLevel,
      statusFilter,
      searchQuery,
      columns,
      
      // 统计数据
      todayLogs,
      blockedLogs,
      normalLogs,
      suspiciousIps,
      avgResponseTime,
      topAttackTypes,
      riskIps,
      recentHourRequests,
      recentHourBlocked,
      
      // 方法
      fetchLogs,
      viewDetails,
      copyLog,
      applyFilters,
      clearFilters,
      changePage,
      changePageSize,
      handleDateRangeChange,
      handleFilter,
      handleSearch,
      exportLogs,
      addToBlacklist,
      formatTimestamp,
      formatBytes,
      getMethodType,
      getStatusType,
      getThreatType,
      getThreatColor
    }
  }
}
</script>

<style scoped>
</style>