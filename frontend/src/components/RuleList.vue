<template>
  <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
    <h2 class="text-2xl font-bold mb-4">现有规则 📜</h2>
    <table class="min-w-full bg-gray-800">
      <thead>
      <tr>
        <th class="py-2 px-4 border-b-2 border-gray-700 text-left">规则名称</th>
        <th class="py-2 px-4 border-b-2 border-gray-700 text-left">规则正则</th>
        <th class="py-2 px-4 border-b-2 border-gray-700 text-left">HTTP 请求</th>
        <th class="py-2 px-4 border-b-2 border-gray-700 text-left">操作</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="rule in rules" :key="rule.name" class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
        <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.name }}</td>
        <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.regex }}</td>
        <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.method }}</td>
        <td class="py-2 px-4 border-b border-gray-700 text-left">
          <button @click="$emit('delete-rule', rule)" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700 transform hover:scale-105 transition duration-300 ml-2">删除 🗑️</button>
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'RuleList',
  props: {
    rules: {
      type: Array,
      required: true
    }
  },
  emits: ['delete-rule']
};
</script>
