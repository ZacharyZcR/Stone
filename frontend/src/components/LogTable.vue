<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">日志记录 📜</h2>
    <table class="min-w-full bg-gray-800">
      <thead>
      <tr>
        <th class="table-header">源 IP</th>
        <th class="table-header">HTTP 方法</th>
        <th class="table-header">目标 IP</th>
        <th class="table-header">URL</th>
        <th class="table-header">状态</th>
        <th class="table-header">时间戳</th>
        <th class="table-header">操作</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="log in logs" :key="log._id" class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
        <td class="table-cell">{{ log.client_ip }}</td>
        <td class="table-cell">{{ log.method || 'N/A' }}</td>
        <td class="table-cell">{{ log.target_ip }}</td>
        <td class="table-cell">{{ log.url || 'N/A' }}</td>
        <td class="table-cell">
            <span :class="log.status === 'success' ? 'text-green-500' : 'text-red-500'">
              {{ log.status === 'success' ? '正常 ✅' : '拦截 🚫' }}
            </span>
        </td>
        <td class="table-cell">{{ formatTimestamp(log.timestamp) }}</td>
        <td class="table-cell">
          <button @click="$emit('view-details', log)" class="details-button">
            查看详情 🔍
          </button>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'LogTable',
  props: {
    logs: {
      type: Array,
      required: true
    }
  },
  emits: ['view-details'],
  setup() {
    const formatTimestamp = (timestamp) => {
      return new Date(timestamp).toLocaleString();
    };

    return { formatTimestamp };
  }
};
</script>

<style scoped>
.table-header {
  @apply py-2 px-4 border-b-2 border-gray-700 text-left;
}

.table-cell {
  @apply py-2 px-4 border-b border-gray-700 text-left;
}

.details-button {
  @apply bg-blue-500 text-white px-3 py-1 rounded-lg hover:bg-blue-700 transform hover:scale-105 transition duration-300 border-2 border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 text-sm;
}
</style>
