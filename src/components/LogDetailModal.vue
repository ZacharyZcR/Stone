<template>
  <div v-if="log" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 overflow-auto py-10">
    <div class="bg-gray-800 rounded-lg shadow-xl w-11/12 max-w-4xl max-h-full overflow-hidden">
      <div class="flex justify-between items-center bg-gray-700 px-6 py-4">
        <h2 class="text-2xl font-bold text-white">流量详情</h2>
        <button @click="$emit('close')" class="text-gray-300 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      <div class="p-6 overflow-y-auto" style="max-height: calc(100vh - 200px);">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-gray-700 p-4 rounded-lg">
            <h3 class="text-lg font-semibold mb-2 text-blue-300">基本信息</h3>
            <p><span class="text-gray-400">客户端 IP:</span> <span class="text-white">{{ log.client_ip }}</span></p>
            <p><span class="text-gray-400">目标 IP:</span> <span class="text-white">{{ log.target_ip }}</span></p>
            <p><span class="text-gray-400">时间戳:</span> <span class="text-white">{{ formatTimestamp(log.timestamp) }}</span></p>
            <p><span class="text-gray-400">状态:</span> <span :class="getStatusClass(log.status)">{{ log.status }}</span></p>
          </div>
          <div class="bg-gray-700 p-4 rounded-lg">
            <h3 class="text-lg font-semibold mb-2 text-green-300">请求详情</h3>
            <p><span class="text-gray-400">HTTP 方法:</span> <span class="text-white">{{ log.method || 'N/A' }}</span></p>
            <p><span class="text-gray-400">URL:</span> <span class="text-white">{{ log.url || 'N/A' }}</span></p>
            <p v-if="log.body"><span class="text-gray-400">请求体:</span> <span class="text-white">{{ truncate(log.body, 100) }}</span></p>
          </div>
        </div>
        <div class="mt-6 bg-gray-700 p-4 rounded-lg">
          <h3 class="text-lg font-semibold mb-2 text-yellow-300">请求头</h3>
          <div v-if="log.headers" class="grid grid-cols-1 md:grid-cols-2 gap-2">
            <p v-for="(value, key) in log.headers" :key="key" class="break-all">
              <span class="text-gray-400">{{ key }}:</span> <span class="text-white">{{ Array.isArray(value) ? value.join(', ') : value }}</span>
            </p>
          </div>
          <p v-else class="text-gray-400">无请求头信息</p>
        </div>
        <div v-if="log.error" class="mt-6 bg-red-900 bg-opacity-50 p-4 rounded-lg">
          <h3 class="text-lg font-semibold mb-2 text-red-300">错误信息</h3>
          <p class="text-red-100">{{ log.error }}</p>
        </div>
      </div>
      <div class="bg-gray-700 px-6 py-4 flex justify-end">
        <button @click="$emit('close')" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 transition duration-300">
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LogDetailModal',
  props: {
    log: {
      type: Object,
      required: true
    }
  },
  methods: {
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleString();
    },
    getStatusClass(status) {
      return status === 'success' ? 'text-green-400' : 'text-red-400';
    },
    truncate(str, n) {
      return (str && str.length > n) ? str.substr(0, n-1) + '...' : str;
    }
  }
}
</script>

<style scoped>
.bg-gray-800 {
  transition: all 0.3s ease-out;
}
</style>
