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
          <tr v-for="rule in urlPatterns" :key="rule.Name" class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.Name }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.Regex }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">GET</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transform hover:scale-105 transition duration-300">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700 transform hover:scale-105 transition duration-300">删除 🗑️</button>
            </td>
          </tr>
          <tr v-for="rule in bodyPatterns" :key="rule.Name" class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.Name }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">{{ rule.Regex }}</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">POST</td>
            <td class="py-2 px-4 border-b border-gray-700 text-left">
              <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transform hover:scale-105 transition duration-300">编辑 ✏️</button>
              <button class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-700 transform hover:scale-105 transition duration-300">删除 🗑️</button>
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
            <label class="block text-gray-300 text-sm font-bold mb-2" for="description">描述</label>
            <textarea v-model="newRule.description" class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300" id="description" placeholder="输入规则描述"></textarea>
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
import axios from 'axios';
import HeaderPage from './HeaderPage.vue';
import FooterPage from './FooterPage.vue';

export default {
  name: 'CustomRuleManagement',
  components: {
    HeaderPage,
    FooterPage
  },
  setup() {
    const urlPatterns = ref([]);
    const bodyPatterns = ref([]);
    const newRule = ref({ name: '', description: '' });

    const fetchRules = async () => {
      try {
        const response = await axios.get('http://172.20.2.226:8081/interception-rules');
        urlPatterns.value = response.data.URLPatterns;
        bodyPatterns.value = response.data.BodyPatterns;
      } catch (error) {
        console.error('获取规则失败:', error);
      }
    };

    const addRule = async () => {
      // 这里可以实现添加规则的后端接口
      console.log('添加规则:', newRule.value);
      newRule.value = { name: '', description: '' }; // 清空表单
    };

    onMounted(() => {
      fetchRules();
    });

    return { urlPatterns, bodyPatterns, newRule, addRule };
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
