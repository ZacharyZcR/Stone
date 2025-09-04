<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：页面标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#7c3aed">
          <n-icon size="48">
            <settings-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="primary">系统配置</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            系统运行状态监控，IP访问控制与安全配置管理
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：系统运行状态 -->
    <n-card title="系统运行状态" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-button type="primary" @click="fetchStatus" :loading="refreshing">
          <template #icon>
            <n-icon><refresh-outline /></n-icon>
          </template>
          刷新状态
        </n-button>
      </template>
      
      <n-spin :show="!statusInfo">
        <n-grid v-if="statusInfo" :cols="4" :x-gap="24" responsive="screen">
          <n-grid-item v-for="(value, key) in displayStatusInfo" :key="key">
            <n-card hoverable embedded>
              <n-statistic 
                :label="formatTitle(key)" 
                :value="formatValue(key, value)"
              >
                <template #prefix>
                  <n-icon :color="getStatusColor(key, value)" size="24">
                    <component :is="getStatusIcon(key)" />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>
        </n-grid>
        
        <n-skeleton v-else text :repeat="8" />
      </n-spin>
    </n-card>

    <!-- 第三行：系统管理操作 -->
    <n-card title="系统管理操作" size="large" style="margin-bottom: 24px;">
      <n-space size="large" justify="center">
        <n-button-group size="large">
          <n-button type="primary" @click="restartService">
            <template #icon>
              <n-icon size="20"><reload-outline /></n-icon>
            </template>
            重启服务
          </n-button>
          
          <n-button @click="clearCache">
            <template #icon>
              <n-icon size="20"><trash-outline /></n-icon>
            </template>
            清理缓存
          </n-button>
          
          <n-button @click="backupData">
            <template #icon>
              <n-icon size="20"><archive-outline /></n-icon>
            </template>
            数据备份
          </n-button>
          
          <n-button @click="exportConfig">
            <template #icon>
              <n-icon size="20"><cloud-download-outline /></n-icon>
            </template>
            导出配置
          </n-button>
        </n-button-group>
      </n-space>
    </n-card>

    <!-- 第四行：IP访问控制 -->
    <n-grid :cols="2" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
      <!-- 左侧：黑名单管理 -->
      <n-grid-item>
        <n-card title="黑名单管理" size="large" style="height: 480px;">
          <template #header-extra>
            <n-space>
              <n-badge :value="blacklist.length" type="error">
                <n-icon size="20"><ban-outline /></n-icon>
              </n-badge>
              <n-button size="small" @click="showAddBlacklistModal = true">
                <template #icon>
                  <n-icon><add-outline /></n-icon>
                </template>
                添加IP
              </n-button>
            </n-space>
          </template>
          
          <div style="height: 380px; overflow-y: auto; padding-right: 8px;">
            <n-list>
              <n-list-item v-for="ip in blacklist" :key="ip">
                <n-space justify="space-between" align="center">
                  <n-space align="center">
                    <n-avatar size="small" color="#d03050">
                      <n-icon size="16"><ban-outline /></n-icon>
                    </n-avatar>
                    <n-text strong>{{ ip }}</n-text>
                  </n-space>
                  <n-button size="small" type="error" @click="removeFromBlacklist(ip)">
                    移除
                  </n-button>
                </n-space>
              </n-list-item>
            </n-list>
            
            <n-empty v-if="blacklist.length === 0" description="黑名单为空" style="height: 100%; display: flex; align-items: center; justify-content: center;">
              <template #icon>
                <n-icon size="48" color="#d0d0d0">
                  <ban-outline />
                </n-icon>
              </template>
            </n-empty>
          </div>
        </n-card>
      </n-grid-item>

      <!-- 右侧：白名单管理 -->
      <n-grid-item>
        <n-card title="白名单管理" size="large" style="height: 480px;">
          <template #header-extra>
            <n-space>
              <n-badge :value="whitelist.length" type="success">
                <n-icon size="20"><checkmark-circle-outline /></n-icon>
              </n-badge>
              <n-button size="small" @click="showAddWhitelistModal = true">
                <template #icon>
                  <n-icon><add-outline /></n-icon>
                </template>
                添加IP
              </n-button>
            </n-space>
          </template>
          
          <div style="height: 380px; overflow-y: auto; padding-right: 8px;">
            <n-list>
              <n-list-item v-for="ip in whitelist" :key="ip">
                <n-space justify="space-between" align="center">
                  <n-space align="center">
                    <n-avatar size="small" color="#18a058">
                      <n-icon size="16"><checkmark-circle-outline /></n-icon>
                    </n-avatar>
                    <n-text strong>{{ ip }}</n-text>
                  </n-space>
                  <n-button size="small" type="success" @click="removeFromWhitelist(ip)">
                    移除
                  </n-button>
                </n-space>
              </n-list-item>
            </n-list>
            
            <n-empty v-if="whitelist.length === 0" description="白名单为空" style="height: 100%; display: flex; align-items: center; justify-content: center;">
              <template #icon>
                <n-icon size="48" color="#d0d0d0">
                  <checkmark-circle-outline />
                </n-icon>
              </template>
            </n-empty>
          </div>
        </n-card>
      </n-grid-item>
    </n-grid>

    <!-- 第五行：配置参数管理 -->
    <n-card title="配置参数管理" size="large">
      <template #header-extra>
        <n-space>
          <n-button @click="loadConfig" :loading="configLoading">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            重新加载
          </n-button>
          <n-button type="primary" @click="saveConfig" :loading="configSaving">
            <template #icon>
              <n-icon><save-outline /></n-icon>
            </template>
            保存配置
          </n-button>
        </n-space>
      </template>
      
      <n-form ref="configFormRef" :model="configData" label-placement="left" label-width="auto">
        <n-grid :cols="2" :x-gap="24">
          <n-grid-item>
            <n-form-item label="服务端口">
              <n-input-number 
                v-model:value="configData.server.port" 
                placeholder="8080" 
                :min="1" 
                :max="65535"
              />
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="防火墙模式">
              <n-select 
                v-model:value="configData.firewall.mode" 
                :options="modeOptions"
              />
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="目标地址">
              <n-input 
                v-model:value="configData.firewall.targetaddress"
                placeholder="localhost:80"
              />
            </n-form-item>
          </n-grid-item>
          
          <n-grid-item>
            <n-form-item label="规则文件">
              <n-input 
                v-model:value="configData.firewall.rulesfile"
                placeholder="rules.yaml"
              />
            </n-form-item>
          </n-grid-item>
        </n-grid>
      </n-form>
    </n-card>
  </div>

    <!-- 添加黑名单IP模态框 -->
    <n-modal v-model:show="showAddBlacklistModal">
      <n-card closable @close="showAddBlacklistModal = false" title="添加黑名单IP">
        <n-form @submit.prevent="addToBlacklist">
          <n-form-item label="IP地址">
            <n-input 
              v-model:value="newBlacklistIP" 
              placeholder="例如: 192.168.1.100"
            />
          </n-form-item>
          <n-space justify="end">
            <n-button @click="showAddBlacklistModal = false">取消</n-button>
            <n-button type="primary" @click="addToBlacklist">添加</n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-modal>

    <!-- 添加白名单IP模态框 -->
    <n-modal v-model:show="showAddWhitelistModal">
      <n-card closable @close="showAddWhitelistModal = false" title="添加白名单IP">
        <n-form @submit.prevent="addToWhitelist">
          <n-form-item label="IP地址">
            <n-input 
              v-model:value="newWhitelistIP" 
              placeholder="例如: 192.168.1.100"
            />
          </n-form-item>
          <n-space justify="end">
            <n-button @click="showAddWhitelistModal = false">取消</n-button>
            <n-button type="primary" @click="addToWhitelist">添加</n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-modal>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useMessage } from 'naive-ui'
import api from '../api/axiosInstance'
import {
  SettingsOutline,
  RefreshOutline,
  ReloadOutline,
  TrashOutline,
  ArchiveOutline,
  CloudDownloadOutline,
  BanOutline,
  CheckmarkCircleOutline,
  AddOutline,
  SaveOutline,
  ServerOutline,
  ShieldOutline,
  TimeOutline,
  SpeedometerOutline
} from '@vicons/ionicons5'

export default {
  name: 'SystemConfiguration',
  components: {
    SettingsOutline,
    RefreshOutline,
    ReloadOutline,
    TrashOutline,
    ArchiveOutline,
    CloudDownloadOutline,
    BanOutline,
    CheckmarkCircleOutline,
    AddOutline,
    SaveOutline,
    ServerOutline,
    ShieldOutline,
    TimeOutline,
    SpeedometerOutline
  },
  setup() {
    const message = useMessage()
    
    // 基础状态
    const statusInfo = ref(null)
    const refreshing = ref(false)
    const configLoading = ref(false)
    const configSaving = ref(false)
    
    // IP控制列表
    const blacklist = ref([])
    const whitelist = ref([])
    
    // 模态框状态
    const showAddBlacklistModal = ref(false)
    const showAddWhitelistModal = ref(false)
    const newBlacklistIP = ref('')
    const newWhitelistIP = ref('')
    
    // 配置数据
    const configData = ref({
      server: {
        port: 8080
      },
      firewall: {
        mode: 'main',
        targetaddress: 'localhost:80',
        rulesfile: 'rules.yaml'
      }
    })
    
    // 模式选项
    const modeOptions = [
      { label: '主模式', value: 'main' },
      { label: '镜像模式', value: 'mirror' },
      { label: '监控模式', value: 'monitor' }
    ]
    
    // 计算属性 - 显示状态信息
    const displayStatusInfo = computed(() => {
      if (!statusInfo.value) return {}
      return {
        uptime: statusInfo.value.uptime,
        processes: statusInfo.value.processes,
        memory_usage: statusInfo.value.memory_usage,
        cpu_usage: statusInfo.value.cpu_usage_percent
      }
    })
    
    // 格式化标题
    const formatTitle = (key) => {
      const titleMap = {
        uptime: '运行时间',
        processes: '进程数',
        memory_usage: '内存使用',
        cpu_usage: 'CPU使用率'
      }
      return titleMap[key] || key
    }
    
    // 格式化值
    const formatValue = (key, value) => {
      if (key === 'uptime') {
        // 处理后端返回的时间字符串格式如 "15m50.455147457s"
        if (typeof value === 'string') {
          return value.replace(/(\d+)m(\d+\.\d+)s/, '$1分$2秒')
        }
        return Math.floor(value / 3600) + '小时'
      }
      if (key === 'memory_usage' || key === 'cpu_usage') {
        return Math.round(value) + '%'
      }
      if (key === 'processes') {
        return value + '个'
      }
      return value
    }
    
    // 获取状态颜色
    const getStatusColor = (key, value) => {
      if (key === 'cpu_usage' || key === 'memory_usage') {
        if (value > 80) return '#d03050'
        if (value > 60) return '#f0a020'
        return '#18a058'
      }
      return '#2080f0'
    }
    
    // 获取状态图标
    const getStatusIcon = (key) => {
      const iconMap = {
        uptime: 'TimeOutline',
        connections: 'ServerOutline',
        memory_usage: 'SpeedometerOutline',
        cpu_usage: 'SpeedometerOutline'
      }
      return iconMap[key] || 'ServerOutline'
    }
    
    // 获取系统状态
    const fetchStatus = async () => {
      refreshing.value = true
      try {
        const response = await api.get('/status')
        statusInfo.value = response.data
      } catch (error) {
        console.error('获取系统状态失败:', error)
        message.error('获取系统状态失败')
      } finally {
        refreshing.value = false
      }
    }
    
    // 系统管理操作
    const restartService = () => {
      message.info('重启服务功能待实现')
    }
    
    const clearCache = () => {
      message.info('清理缓存功能待实现')
    }
    
    const backupData = () => {
      message.info('数据备份功能待实现')
    }
    
    const exportConfig = () => {
      message.info('导出配置功能待实现')
    }
    
    // IP管理
    const addToBlacklist = async () => {
      if (newBlacklistIP.value) {
        try {
          await api.post('/ip-control-rules', {
            ip: newBlacklistIP.value,
            type: 'blacklist'
          })
          blacklist.value.push(newBlacklistIP.value)
          newBlacklistIP.value = ''
          showAddBlacklistModal.value = false
          message.success('IP已添加到黑名单')
        } catch (error) {
          console.error('添加黑名单失败:', error)
          message.error('添加黑名单失败')
        }
      }
    }
    
    const removeFromBlacklist = async (ip) => {
      try {
        await api.delete(`/ip-control-rules/${ip}`)
        const index = blacklist.value.indexOf(ip)
        if (index > -1) {
          blacklist.value.splice(index, 1)
        }
        message.success('IP已从黑名单移除')
      } catch (error) {
        console.error('移除黑名单失败:', error)
        message.error('移除黑名单失败')
      }
    }
    
    const addToWhitelist = async () => {
      if (newWhitelistIP.value) {
        try {
          await api.post('/ip-control-rules', {
            ip: newWhitelistIP.value,
            type: 'whitelist'
          })
          whitelist.value.push(newWhitelistIP.value)
          newWhitelistIP.value = ''
          showAddWhitelistModal.value = false
          message.success('IP已添加到白名单')
        } catch (error) {
          console.error('添加白名单失败:', error)
          message.error('添加白名单失败')
        }
      }
    }
    
    const removeFromWhitelist = async (ip) => {
      try {
        await api.delete(`/ip-control-rules/${ip}`)
        const index = whitelist.value.indexOf(ip)
        if (index > -1) {
          whitelist.value.splice(index, 1)
        }
        message.success('IP已从白名单移除')
      } catch (error) {
        console.error('移除白名单失败:', error)
        message.error('移除白名单失败')
      }
    }
    
    // 配置管理
    const loadConfig = async () => {
      configLoading.value = true
      try {
        // 配置暂时显示当前值，后续可以连接真实配置API
        message.success('配置已刷新')
      } catch (error) {
        message.error('配置加载失败')
      } finally {
        configLoading.value = false
      }
    }
    
    const saveConfig = async () => {
      configSaving.value = true
      try {
        // 配置保存功能待后端API支持
        message.success('配置已保存（演示版本）')
      } catch (error) {
        message.error('配置保存失败')
      } finally {
        configSaving.value = false
      }
    }
    
    // 获取IP控制规则
    const fetchIPRules = async () => {
      try {
        const response = await api.get('/ip-control-rules')
        blacklist.value = response.data.Blacklist || []
        whitelist.value = response.data.Whitelist || []
      } catch (error) {
        console.error('获取IP控制规则失败:', error)
        message.error('获取IP控制规则失败')
      }
    }

    // 初始化数据
    onMounted(() => {
      fetchStatus()
      fetchIPRules()
    })
    
    return {
      // 状态
      statusInfo,
      refreshing,
      configLoading,
      configSaving,
      displayStatusInfo,
      
      // IP控制
      blacklist,
      whitelist,
      showAddBlacklistModal,
      showAddWhitelistModal,
      newBlacklistIP,
      newWhitelistIP,
      
      // 配置
      configData,
      modeOptions,
      
      // 方法
      fetchStatus,
      formatTitle,
      formatValue,
      getStatusColor,
      getStatusIcon,
      restartService,
      clearCache,
      backupData,
      exportConfig,
      addToBlacklist,
      removeFromBlacklist,
      addToWhitelist,
      removeFromWhitelist,
      loadConfig,
      saveConfig
    }
  }
}
</script>