<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">日志过滤器 🔍</h2>
    <form class="grid grid-cols-1 md:grid-cols-4 gap-4" @submit.prevent="applyFilter">
      <div>
        <label class="block text-gray-300 text-sm font-bold mb-2" for="startDateTime">开始时间</label>
        <input v-model="filters.startDateTime" class="filter-input" id="startDateTime" type="datetime-local">
      </div>
      <div>
        <label class="block text-gray-300 text-sm font-bold mb-2" for="endDateTime">结束时间</label>
        <input v-model="filters.endDateTime" class="filter-input" id="endDateTime" type="datetime-local">
      </div>
      <div>
        <label class="block text-gray-300 text-sm font-bold mb-2" for="ipFilter">IP 地址</label>
        <input v-model="filters.ip" class="filter-input" id="ipFilter" type="text" placeholder="输入 IP 地址">
      </div>
      <div>
        <label class="block text-gray-300 text-sm font-bold mb-2" for="statusFilter">日志状态</label>
        <select v-model="filters.status" class="filter-input" id="statusFilter">
          <option value="">全部</option>
          <option value="passed">正常</option>
          <option value="blocked">拦截</option>
        </select>
      </div>
      <div class="flex items-end col-span-full">
        <button class="filter-button" type="submit">
          应用过滤器 🔍
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue';

export default {
  name: 'LogFilter',
  emits: ['filter-applied'],
  setup(props, { emit }) {
    const filters = ref({
      startDateTime: '',
      endDateTime: '',
      ip: '',
      status: ''
    });

    const applyFilter = () => {
      emit('filter-applied', filters.value);
    };

    return { filters, applyFilter };
  }
};
</script>

<style scoped>
.filter-input {
  @apply shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-indigo-500 transition duration-300;
}

.filter-button {
  @apply bg-blue-500 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transform hover:scale-105 transition duration-300 border-2 border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 w-full;
}
</style>
