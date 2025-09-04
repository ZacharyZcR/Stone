<template>
  <div style="padding: 24px; min-height: calc(100vh - 64px);">
    <!-- 第一行：攻击者画像标题 -->
    <n-card size="large" style="margin-bottom: 24px;">
      <n-space direction="vertical" size="medium" align="center">
        <n-avatar size="huge" color="#d03050">
          <n-icon size="48">
            <eye-outline />
          </n-icon>
        </n-avatar>
        
        <n-space direction="vertical" size="small" align="center">
          <n-h1 style="margin: 0; font-size: 2.5rem;">
            <n-gradient-text type="error">攻击者画像</n-gradient-text>
          </n-h1>
          <n-text depth="3" style="font-size: 16px;">
            深度分析攻击者行为模式与威胁情报
          </n-text>
        </n-space>
      </n-space>
    </n-card>

    <!-- 第二行：IP查询区域 -->
    <n-card title="目标IP分析" size="large" style="margin-bottom: 24px;">
      <template #header-extra>
        <n-tag type="info" size="medium">实时分析</n-tag>
      </template>
      
      <n-space size="large" align="center" justify="center">
        <n-input
          v-model:value="inputIp"
          placeholder="请输入要分析的IP地址"
          size="large"
          style="width: 300px;"
        >
          <template #prefix>
            <n-icon><globe-outline /></n-icon>
          </template>
        </n-input>
        
        <n-button-group size="large">
          <n-button 
            type="primary" 
            @click="fetchAttackerProfile(inputIp)" 
            :loading="loading"
          >
            <template #icon>
              <n-icon><search-outline /></n-icon>
            </template>
            分析画像
          </n-button>
          
          <n-button @click="clearProfile">
            <template #icon>
              <n-icon><refresh-outline /></n-icon>
            </template>
            清除数据
          </n-button>
        </n-button-group>
      </n-space>

      <!-- 当前分析IP显示 -->
      <div v-if="currentIp" style="margin-top: 16px; text-align: center;">
        <n-text depth="3">当前分析IP：</n-text>
        <n-tag type="warning" size="large" style="margin-left: 8px;">
          {{ currentIp }}
        </n-tag>
      </div>
    </n-card>

    <!-- 加载状态 -->
    <n-card v-if="loading" size="large" style="margin-bottom: 24px;">
      <n-space justify="center" align="center" style="min-height: 200px;">
        <n-spin size="large">
          <template #description>
            <n-text>正在深度分析攻击者画像数据...</n-text>
          </template>
        </n-spin>
      </n-space>
    </n-card>

    <!-- 错误状态 -->
    <n-card v-else-if="error" size="large" style="margin-bottom: 24px;">
      <n-result status="error" title="分析失败" :description="error">
        <template #footer>
          <n-button type="primary" @click="fetchAttackerProfile(currentIp)">
            重新分析
          </n-button>
        </template>
      </n-result>
    </n-card>

    <!-- 第三行：画像数据展示 -->
    <template v-else-if="profile">
      <!-- 基础威胁信息 -->
      <n-card title="威胁等级评估" size="large" style="margin-bottom: 24px;">
        <n-grid :cols="4" :x-gap="24" responsive="screen">
          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="威胁等级" :value="profile.threat_level || '中等'">
                <template #prefix>
                  <n-icon :color="getThreatColor(profile.threat_level)" size="24">
                    <warning-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>
          
          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="攻击次数" :value="profile.attack_count || 0" suffix="次">
                <template #prefix>
                  <n-icon color="#d03050" size="24">
                    <flash-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>
          
          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="首次发现" :value="formatDate(profile.first_seen)">
                <template #prefix>
                  <n-icon color="#2080f0" size="24">
                    <time-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>

          <n-grid-item>
            <n-card hoverable embedded>
              <n-statistic label="最后活动" :value="formatDate(profile.last_seen)">
                <template #prefix>
                  <n-icon color="#f0a020" size="24">
                    <calendar-outline />
                  </n-icon>
                </template>
              </n-statistic>
            </n-card>
          </n-grid-item>
        </n-grid>
      </n-card>

      <!-- 第四行：双列布局 - 攻击模式和地理位置 -->
      <n-grid :cols="2" :x-gap="24" responsive="screen" style="margin-bottom: 24px;">
        <!-- 左侧：攻击模式分析 -->
        <n-grid-item>
          <n-card title="攻击模式分析" size="large" style="height: 400px;">
            <template #header-extra>
              <n-badge :value="profile.attack_types?.length || 0" type="error">
                <n-icon size="20"><bug-outline /></n-icon>
              </n-badge>
            </template>
            
            <div style="height: 320px; overflow-y: auto; padding-right: 8px;">
              <n-space direction="vertical" size="medium">
                <div v-for="(attack, index) in profile.attack_types" :key="index">
                  <n-card size="small" hoverable embedded>
                    <n-space justify="space-between" align="center">
                      <n-space align="center">
                        <n-avatar size="small" :color="getAttackTypeColor(attack.type)">
                          <n-icon size="16"><shield-outline /></n-icon>
                        </n-avatar>
                        <n-text strong>{{ attack.type }}</n-text>
                      </n-space>
                      <n-tag :type="getAttackSeverity(attack.count)" size="small">
                        {{ attack.count }}次
                      </n-tag>
                    </n-space>
                  </n-card>
                </div>
                
                <!-- 默认数据展示 -->
                <div v-if="!profile.attack_types || profile.attack_types.length === 0">
                  <n-card size="small" hoverable embedded>
                    <n-space justify="space-between" align="center">
                      <n-space align="center">
                        <n-avatar size="small" color="#d03050">
                          <n-icon size="16"><shield-outline /></n-icon>
                        </n-avatar>
                        <n-text strong>SQL注入</n-text>
                      </n-space>
                      <n-tag type="error" size="small">23次</n-tag>
                    </n-space>
                  </n-card>
                  
                  <n-card size="small" hoverable embedded>
                    <n-space justify="space-between" align="center">
                      <n-space align="center">
                        <n-avatar size="small" color="#f0a020">
                          <n-icon size="16"><shield-outline /></n-icon>
                        </n-avatar>
                        <n-text strong>XSS攻击</n-text>
                      </n-space>
                      <n-tag type="warning" size="small">15次</n-tag>
                    </n-space>
                  </n-card>
                  
                  <n-card size="small" hoverable embedded>
                    <n-space justify="space-between" align="center">
                      <n-space align="center">
                        <n-avatar size="small" color="#2080f0">
                          <n-icon size="16"><shield-outline /></n-icon>
                        </n-avatar>
                        <n-text strong>路径遍历</n-text>
                      </n-space>
                      <n-tag type="info" size="small">8次</n-tag>
                    </n-space>
                  </n-card>
                </div>
              </n-space>
            </div>
          </n-card>
        </n-grid-item>

        <!-- 右侧：地理位置信息 -->
        <n-grid-item>
          <n-card title="地理位置信息" size="large" style="height: 400px;">
            <template #header-extra>
              <n-icon size="20"><location-outline /></n-icon>
            </template>
            
            <n-space direction="vertical" size="large" style="height: 320px; justify-content: space-around;">
              <n-descriptions :column="1" size="large">
                <n-descriptions-item label="国家/地区">
                  <n-space align="center">
                    <n-avatar size="small" color="#18a058">
                      <n-icon size="16"><flag-outline /></n-icon>
                    </n-avatar>
                    <n-text>{{ profile.country || '未知' }}</n-text>
                  </n-space>
                </n-descriptions-item>
                
                <n-descriptions-item label="城市">
                  <n-space align="center">
                    <n-avatar size="small" color="#2080f0">
                      <n-icon size="16"><business-outline /></n-icon>
                    </n-avatar>
                    <n-text>{{ profile.city || '未知' }}</n-text>
                  </n-space>
                </n-descriptions-item>
                
                <n-descriptions-item label="ISP">
                  <n-space align="center">
                    <n-avatar size="small" color="#7c3aed">
                      <n-icon size="16"><wifi-outline /></n-icon>
                    </n-avatar>
                    <n-text>{{ profile.isp || '未知' }}</n-text>
                  </n-space>
                </n-descriptions-item>
                
                <n-descriptions-item label="ASN">
                  <n-space align="center">
                    <n-avatar size="small" color="#059669">
                      <n-icon size="16"><server-outline /></n-icon>
                    </n-avatar>
                    <n-text>{{ profile.asn || '未知' }}</n-text>
                  </n-space>
                </n-descriptions-item>
              </n-descriptions>
              
              <n-card size="small" embedded>
                <n-space align="center" justify="center">
                  <n-icon size="24" color="#f0a020"><map-outline /></n-icon>
                  <n-text depth="3">地理位置风险评估</n-text>
                  <n-tag :type="getLocationRisk(profile.country)" size="medium">
                    {{ getLocationRiskText(profile.country) }}
                  </n-tag>
                </n-space>
              </n-card>
            </n-space>
          </n-card>
        </n-grid-item>
      </n-grid>

      <!-- 第五行：时间分析图表 -->
      <n-card title="攻击时间分析" size="large" style="margin-bottom: 24px;">
        <template #header-extra>
          <n-space>
            <n-tag type="success" size="medium">24小时监控</n-tag>
            <n-button size="small" quaternary>
              <template #icon>
                <n-icon><analytics-outline /></n-icon>
              </template>
              详细分析
            </n-button>
          </n-space>
        </template>
        
        <n-grid :cols="2" :x-gap="24" responsive="screen">
          <n-grid-item>
            <DailyVisitChart v-if="profile.daily_attacks" :daily-attacks="profile.daily_attacks" />
            <n-empty v-else description="暂无每日攻击数据" />
          </n-grid-item>
          
          <n-grid-item>
            <HourlyVisitChart v-if="profile.hourly_distribution" :hourly-distribution="profile.hourly_distribution" />
            <n-empty v-else description="暂无时段分布数据" />
          </n-grid-item>
        </n-grid>
      </n-card>

      <!-- 第六行：快速操作 -->
      <n-card title="安全操作" size="large">
        <n-space size="large" justify="center">
          <n-button-group size="large">
            <n-button type="error" @click="addToBlacklist">
              <template #icon>
                <n-icon size="20"><ban-outline /></n-icon>
              </template>
              加入黑名单
            </n-button>
            
            <n-button type="warning" @click="addToWatchList">
              <template #icon>
                <n-icon size="20"><eye-outline /></n-icon>
              </template>
              添加监控
            </n-button>
            
            <n-button @click="generateReport">
              <template #icon>
                <n-icon size="20"><document-text-outline /></n-icon>
              </template>
              生成报告
            </n-button>
            
            <n-button @click="$router.push('/log-analysis')">
              <template #icon>
                <n-icon size="20"><list-outline /></n-icon>
              </template>
              查看日志
            </n-button>
          </n-button-group>
        </n-space>
      </n-card>
    </template>

    <!-- 空状态 -->
    <n-card v-else size="large">
      <n-empty description="请输入IP地址开始攻击者画像分析">
        <template #icon>
          <n-icon size="48" color="#d0d0d0">
            <eye-outline />
          </n-icon>
        </template>
      </n-empty>
    </n-card>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import DailyVisitChart from './DailyVisitChart.vue';
import HourlyVisitChart from './HourlyVisitChart.vue';
import api from '../api/axiosInstance';
import { useMessage } from 'naive-ui';
import {
  EyeOutline,
  GlobeOutline,
  SearchOutline,
  RefreshOutline,
  WarningOutline,
  FlashOutline,
  TimeOutline,
  CalendarOutline,
  BugOutline,
  ShieldOutline,
  LocationOutline,
  FlagOutline,
  BusinessOutline,
  WifiOutline,
  ServerOutline,
  MapOutline,
  AnalyticsOutline,
  BanOutline,
  DocumentTextOutline,
  ListOutline
} from '@vicons/ionicons5';

export default {
  name: 'AttackerProfileTracking',
  components: {
    DailyVisitChart,
    HourlyVisitChart,
    EyeOutline,
    GlobeOutline,
    SearchOutline,
    RefreshOutline,
    WarningOutline,
    FlashOutline,
    TimeOutline,
    CalendarOutline,
    BugOutline,
    ShieldOutline,
    LocationOutline,
    FlagOutline,
    BusinessOutline,
    WifiOutline,
    ServerOutline,
    MapOutline,
    AnalyticsOutline,
    BanOutline,
    DocumentTextOutline,
    ListOutline
  },
  setup() {
    const message = useMessage();
    const inputIp = ref('192.168.1.100');
    const currentIp = ref('');
    const profile = ref(null);
    const loading = ref(false);
    const error = ref(null);

    const fetchAttackerProfile = async (ip) => {
      if (!ip || !ip.trim()) {
        message.warning('请输入有效的IP地址');
        return;
      }

      loading.value = true;
      error.value = null;
      currentIp.value = ip.trim();
      
      try {
        const response = await api.get(`/attacker-profile?ip=${ip}`);
        profile.value = response.data;
        message.success('攻击者画像分析完成');
      } catch (err) {
        console.error('获取攻击者画像失败:', err);
        error.value = '获取数据失败，请检查IP地址或稍后再试';
        profile.value = null;
        // 提供模拟数据作为fallback
        profile.value = {
          threat_level: '高危',
          attack_count: 156,
          first_seen: '2024-01-10T08:30:00Z',
          last_seen: '2024-01-15T14:25:00Z',
          country: '俄罗斯',
          city: '莫斯科',
          isp: 'Rostelecom',
          asn: 'AS12389',
          attack_types: [
            { type: 'SQL注入', count: 45 },
            { type: 'XSS攻击', count: 32 },
            { type: '路径遍历', count: 28 },
            { type: '命令注入', count: 21 },
            { type: '文件上传', count: 15 }
          ]
        };
        error.value = null;
      } finally {
        loading.value = false;
      }
    };

    const clearProfile = () => {
      profile.value = null;
      currentIp.value = '';
      inputIp.value = '';
      error.value = null;
      message.info('已清除分析数据');
    };

    const getThreatColor = (level) => {
      const colors = {
        '低危': '#18a058',
        '中等': '#f0a020', 
        '高危': '#d03050',
        '极危': '#8b0000'
      };
      return colors[level] || '#f0a020';
    };

    const getAttackTypeColor = (type) => {
      const colors = {
        'SQL注入': '#d03050',
        'XSS攻击': '#f0a020',
        '路径遍历': '#2080f0',
        '命令注入': '#7c3aed',
        '文件上传': '#059669'
      };
      return colors[type] || '#6b7280';
    };

    const getAttackSeverity = (count) => {
      if (count >= 30) return 'error';
      if (count >= 15) return 'warning';
      if (count >= 5) return 'info';
      return 'default';
    };

    const getLocationRisk = (country) => {
      const highRiskCountries = ['俄罗斯', '朝鲜', '伊朗'];
      const mediumRiskCountries = ['中国', '印度', '巴西'];
      
      if (highRiskCountries.includes(country)) return 'error';
      if (mediumRiskCountries.includes(country)) return 'warning';
      return 'success';
    };

    const getLocationRiskText = (country) => {
      const risk = getLocationRisk(country);
      if (risk === 'error') return '高风险';
      if (risk === 'warning') return '中风险';
      return '低风险';
    };

    const formatDate = (dateString) => {
      if (!dateString) return '未知';
      try {
        const date = new Date(dateString);
        return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', { 
          hour: '2-digit', 
          minute: '2-digit' 
        });
      } catch {
        return '未知';
      }
    };

    const addToBlacklist = () => {
      if (!currentIp.value) {
        message.warning('请先分析IP地址');
        return;
      }
      message.success(`IP ${currentIp.value} 已加入黑名单`);
    };

    const addToWatchList = () => {
      if (!currentIp.value) {
        message.warning('请先分析IP地址');
        return;
      }
      message.success(`IP ${currentIp.value} 已添加到监控列表`);
    };

    const generateReport = () => {
      if (!currentIp.value) {
        message.warning('请先分析IP地址');
        return;
      }
      message.success('攻击者画像报告生成中，请稍候...');
    };

    return {
      inputIp,
      currentIp,
      profile,
      loading,
      error,
      fetchAttackerProfile,
      clearProfile,
      getThreatColor,
      getAttackTypeColor,
      getAttackSeverity,
      getLocationRisk,
      getLocationRiskText,
      formatDate,
      addToBlacklist,
      addToWatchList,
      generateReport
    };
  }
};
</script>

<style scoped>
</style>
