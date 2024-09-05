<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 规则列表 -->
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
              <button @click="deleteRule(rule)" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700 transform hover:scale-105 transition duration-300 ml-2">删除 🗑️</button>
            </td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- 添加规则 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md transform transition-all duration-500 hover:shadow-2xl">
        <h2 class="text-2xl font-bold mb-4">添加新规则 ➕</h2>
        <form @submit.prevent="addRule">
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleName">规则名称</label>
            <input v-model="newRule.name" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300" id="ruleName" type="text" placeholder="输入规则名称">
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleRegex">规则正则</label>
            <input v-model="newRule.regex" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300" id="ruleRegex" type="text" placeholder="输入规则正则">
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="ruleMethod">HTTP 方法</label>
            <select v-model="newRule.method" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300" id="ruleMethod">
              <option value="GET">GET</option>
              <option value="POST">POST</option>
            </select>
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700 transform hover:scale-105 transition duration-300" type="submit">
              添加规则 ➕
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import api from '../api/axiosInstance'; // 导入配置好的 Axios 实例
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';

export default {
  name: 'CustomRuleManagement',
  components: {
    HeaderPage,
    FooterPage
  },
  setup() {
    const rules = ref([]);
    const newRule = ref({ name: '', regex: '', method: 'GET' });

    const fetchRules = async () => {
      try {
        const response = await api.get('/interception-rules'); // 使用 Axios 实例
        rules.value = response.data.rules || [];
      } catch (error) {
        console.error('获取规则失败:', error);
      }
    };

    const addRule = async () => {
      try {
        await api.post('/interception-rules', newRule.value); // 使用 Axios 实例
        await fetchRules(); // 重新获取规则列表
        newRule.value = { name: '', regex: '', method: 'GET' }; // 清空表单
      } catch (error) {
        console.error('添加规则失败:', error);
      }
    };

    const deleteRule = async (rule) => {
      try {
        await api.delete(`/interception-rules/${rule.name}`); // 使用 Axios 实例
        await fetchRules(); // 重新获取规则列表
      } catch (error) {
        console.error('删除规则失败:', error);
      }
    };

    onMounted(() => {
      fetchRules();
    });

    return {
      rules,
      newRule,
      addRule,
      deleteRule
    };
  }
};
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out;
}
</style>
