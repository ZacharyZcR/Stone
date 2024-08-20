<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 日志过滤器 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">日志过滤器 🔍</h2>
        <form class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-gray-300 text-sm font-bold mb-2" for="startDate">开始日期</label>
            <input class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-indigo-500 transition duration-300" id="startDate" type="date">
          </div>
          <div>
            <label class="block text-gray-300 text-sm font-bold mb-2" for="endDate">结束日期</label>
            <input class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-indigo-500 transition duration-300" id="endDate" type="date">
          </div>
          <div class="flex items-end">
            <button class="bg-indigo-500 text-white px-4 py-2 rounded hover:bg-indigo-700 transform hover:scale-105 transition duration-300 w-full" type="button">
              应用过滤器 ✅
            </button>
          </div>
        </form>
      </div>

      <!-- 日志记录 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">日志记录 📜</h2>
        <table class="min-w-full bg-gray-800">
          <thead>
          <tr>
            <th class="py-2 px-4 border-b-2 border-gray-700">时间 ⏰</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">事件类型 🛡️</th>
            <th class="py-2 px-4 border-b-2 border-gray-700">详情 🔍</th>
          </tr>
          </thead>
          <tbody>
          <tr class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
            <td class="py-2 px-4 border-b border-gray-700">2024-08-20 10:30:45</td>
            <td class="py-2 px-4 border-b border-gray-700">SQL 注入</td>
            <td class="py-2 px-4 border-b border-gray-700">检测到 SQL 注入尝试，已阻止。💪</td>
          </tr>
          <tr class="hover:bg-gray-700 transition duration-300 animate-fade-in-up">
            <td class="py-2 px-4 border-b border-gray-700">2024-08-20 09:45:23</td>
            <td class="py-2 px-4 border-b border-gray-700">XSS 攻击</td>
            <td class="py-2 px-4 border-b border-gray-700">检测到 XSS 攻击尝试，已阻止。🔒</td>
          </tr>
          </tbody>
        </table>
      </div>

      <!-- 统计分析 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">统计分析 📊</h2>
        <div class="h-64 bg-gray-700 flex items-center justify-center animate-pulse">
          <span class="text-gray-400">图表占位符 📊</span>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <FooterPage />
  </div>
</template>

<script>
import HeaderPage from './HeaderPage.vue'
import FooterPage from './FooterPage.vue'

export default {
  name: 'LogAnalysis',
  components: {
    HeaderPage,
    FooterPage
  }
}
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
