<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />

    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 系统运行信息 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-6">系统运行信息 📊</h2>
        <div v-if="statusInfo" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatusCard v-for="(value, key) in statusInfo" :key="key" :title="formatTitle(key)" :value="formatValue(key, value)" />
        </div>
        <div v-else class="text-center py-8">
          <p class="text-xl">加载中... ⏳</p>
        </div>
      </div>

      <!-- 黑白名单 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <!-- 黑名单 IP 列表 -->
        <IPList title="黑名单 IP 列表 🚫" :list="blacklist" @add="addToBlacklist" @remove="removeFromBlacklist" />

        <!-- 白名单 IP 列表 -->
        <IPList title="白名单 IP 列表 ✅" :list="whitelist" @add="addToWhitelist" @remove="removeFromWhitelist" />
      </div>

      <!-- 刷新按钮 -->
      <div class="text-center">
        <button @click="fetchStatus" class="bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-6 rounded-full transition duration-300 ease-in-out transform hover:scale-105">
          刷新系统信息 🔄
        </button>
      </div>
    </div>

    <FooterPage />

    <PopupNotification
        v-if="showNotification"
        :message="notificationMessage"
        :emoji="notificationEmoji"
        :type="notificationType"
        @close="showNotification = false"
    />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import api from '../api/axiosInstance';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';
import PopupNotification from './PopupNotification.vue';
import StatusCard from './StatusCard.vue';
import IPList from './IPList.vue';

export default {
  name: 'SystemStatus',
  components: {
    HeaderPage,
    FooterPage,
    PopupNotification,
    StatusCard,
    IPList
  },
  setup() {
    const statusInfo = ref(null);
    const whitelist = ref([]);
    const blacklist = ref([]);
    const showNotification = ref(false);
    const notificationMessage = ref('');
    const notificationEmoji = ref('');
    const notificationType = ref('success');

    const showPopup = (message, emoji, type) => {
      notificationMessage.value = message;
      notificationEmoji.value = emoji;
      notificationType.value = type;
      showNotification.value = true;
    };

    const fetchStatus = async () => {
      try {
        const response = await api.get('/status');
        statusInfo.value = response.data;
        showPopup('系统信息已更新', '✅', 'success');
      } catch (error) {
        console.error('获取系统状态失败:', error);
        showPopup('获取系统状态失败', '❌', 'error');
      }
    };

    const fetchIPControlRules = async () => {
      try {
        const response = await api.get('/ip-control-rules');
        whitelist.value = response.data.Whitelist;
        blacklist.value = response.data.Blacklist;
      } catch (error) {
        console.error('获取IP控制规则失败:', error);
        showPopup('获取IP控制规则失败', '❌', 'error');
      }
    };

    const addToWhitelist = async (ip) => {
      try {
        await api.post('/ip-control-rules', { ip, type: 'whitelist' });
        await fetchIPControlRules(); // 重新获取最新的 IP 列表
        showPopup(`IP ${ip} 已添加到白名单`, '✅', 'success');
      } catch (error) {
        console.error('添加到白名单失败:', error);
        showPopup(`添加 ${ip} 到白名单失败`, '❌', 'error');
      }
    };

    const addToBlacklist = async (ip) => {
      try {
        await api.post('/ip-control-rules', { ip, type: 'blacklist' });
        await fetchIPControlRules(); // 重新获取最新的 IP 列表
        showPopup(`IP ${ip} 已添加到黑名单`, '✅', 'success');
      } catch (error) {
        console.error('添加到黑名单失败:', error);
        showPopup(`添加 ${ip} 到黑名单失败`, '❌', 'error');
      }
    };

    const removeFromWhitelist = async (ip) => {
      try {
        await api.delete(`/ip-control-rules/${ip}`);
        await fetchIPControlRules(); // 重新获取最新的 IP 列表
        showPopup(`IP ${ip} 已从白名单移除`, '✅', 'success');
      } catch (error) {
        console.error('从白名单移除失败:', error);
        showPopup(`从白名单移除 ${ip} 失败`, '❌', 'error');
      }
    };

    const removeFromBlacklist = async (ip) => {
      try {
        await api.delete(`/ip-control-rules/${ip}`);
        await fetchIPControlRules(); // 重新获取最新的 IP 列表
        showPopup(`IP ${ip} 已从黑名单移除`, '✅', 'success');
      } catch (error) {
        console.error('从黑名单移除失败:', error);
        showPopup(`从黑名单移除 ${ip} 失败`, '❌', 'error');
      }
    };

    const formatTitle = (key) => {
      const titles = {
        status: '状态',
        uptime: '运行时间',
        cpu_usage_percent: 'CPU 使用率',
        memory_usage: '内存使用率',
        disk_usage: '磁盘使用率',
        network_in: '网络入流量',
        network_out: '网络出流量',
        load_average: '平均负载',
        open_file_desc: '打开文件描述符数',
        threads: '线程数',
        processes: '进程数'
      };
      return titles[key] || key;
    };

    const formatValue = (key, value) => {
      switch (key) {
        case 'cpu_usage_percent':
        case 'memory_usage':
        case 'disk_usage':
          return value.toFixed(2) + '%';
        case 'network_in':
        case 'network_out':
          return formatBytes(value);
        case 'load_average':
          return value.toFixed(2);
        case 'uptime':
          return formatUptime(value);
        default:
          return value;
      }
    };

    const formatUptime = (uptimeString) => {
      // 解析时间字符串
      const regex = /(?:(\d+)h)?(?:(\d+)m)?(\d+(?:\.\d+)?)s/;
      const match = uptimeString.match(regex);

      if (!match) {
        return uptimeString; // 如果格式不匹配，直接返回原字符串
      }

      const hours = parseInt(match[1] || '0');
      const minutes = parseInt(match[2] || '0');
      const seconds = Math.floor(parseFloat(match[3] || '0'));

      const days = Math.floor(hours / 24);
      const remainingHours = hours % 24;

      let result = '';
      if (days > 0) result += `${days}天 `;
      if (remainingHours > 0) result += `${remainingHours}小时 `;
      if (minutes > 0) result += `${minutes}分钟 `;
      if (seconds > 0 || (days === 0 && remainingHours === 0 && minutes === 0)) result += `${seconds}秒`;

      return result.trim();
    };

    const formatBytes = (bytes) => {
      if (bytes === 0) return '0 Bytes';
      const k = 1024;
      const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    onMounted(() => {
      fetchStatus();
      fetchIPControlRules();
    });

    return {
      statusInfo,
      whitelist,
      blacklist,
      fetchStatus,
      addToWhitelist,
      addToBlacklist,
      removeFromWhitelist,
      removeFromBlacklist,
      formatTitle,
      formatValue,
      showNotification,
      notificationMessage,
      notificationEmoji,
      notificationType
    };
  }
};
</script>
