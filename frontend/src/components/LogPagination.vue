<template>
  <div class="mt-4 flex flex-wrap justify-between items-center">
    <div class="flex items-center space-x-2 mb-2 sm:mb-0">
      <span class="text-gray-300">每页显示:</span>
      <select
          :value="pageSize"
          @change="$emit('page-size-changed', parseInt($event.target.value))"
          class="page-size-select"
      >
        <option :value="10">10</option>
        <option :value="20">20</option>
        <option :value="50">50</option>
        <option :value="100">100</option>
      </select>
    </div>
    <div class="flex items-center space-x-2">
      <button
          @click="$emit('page-changed', 1)"
          :disabled="currentPage === 1"
          class="pagination-button"
      >
        首页
      </button>
      <button
          @click="$emit('page-changed', currentPage - 1)"
          :disabled="currentPage === 1"
          class="pagination-button"
      >
        上一页
      </button>
      <span class="text-gray-300">
        第 {{ currentPage }} 页，共 {{ totalPages }} 页 (总记录数: {{ totalCount }})
      </span>
      <button
          @click="$emit('page-changed', currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="pagination-button"
      >
        下一页
      </button>
      <button
          @click="$emit('page-changed', totalPages)"
          :disabled="currentPage === totalPages"
          class="pagination-button"
      >
        末页
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'LogPagination',
  props: {
    currentPage: {
      type: Number,
      required: true
    },
    totalPages: {
      type: Number,
      required: true
    },
    totalCount: {
      type: Number,
      required: true
    },
    pageSize: {
      type: Number,
      required: true
    }
  },
  emits: ['page-changed', 'page-size-changed']
};
</script>

