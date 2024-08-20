<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen animate-fade-in">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 攻击者信息 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8 transform hover:scale-105 transition duration-500">
        <h2 class="text-2xl font-bold mb-4">攻击者信息 🕵️</h2>
        <p class="text-lg">正在分析的 IP 地址：<span class="font-bold">192.168.1.1</span></p>
      </div>

      <!-- 访问统计 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">访问统计 📊</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div class="bg-blue-900 p-4 rounded-lg shadow transform hover:scale-105 transition duration-500">
            <h3 class="text-xl font-bold">访问次数</h3>
            <p class="text-3xl animate-number-scroll" data-target="1234">0</p>
          </div>
          <div class="bg-red-900 p-4 rounded-lg shadow transform hover:scale-105 transition duration-500">
            <h3 class="text-xl font-bold">攻击次数</h3>
            <p class="text-3xl animate-number-scroll" data-target="56">0</p>
          </div>
          <!-- 更多统计数据 -->
        </div>
      </div>

      <!-- 相关日志 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">相关日志 📜</h2>
        <ul class="space-y-4">
          <li class="flex items-start animate-fade-in-up">
            <span class="text-blue-400 text-2xl mr-4">🕒</span>
            <div>
              <p class="font-bold">2024-08-20 10:30:45</p>
              <p>检测到 SQL 注入攻击，已阻止。</p>
            </div>
          </li>
          <!-- 更多日志记录 -->
        </ul>
      </div>

      <!-- 相关图表 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">相关图表 📈</h2>
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
  name: 'AttackerProfileTracking',
  components: {
    HeaderPage,
    FooterPage
  },
  mounted() {
    this.animateNumbers();
  },
  methods: {
    animateNumbers() {
      const elements = document.querySelectorAll('.animate-number-scroll');
      elements.forEach(el => {
        const target = parseInt(el.getAttribute('data-target'), 10);
        let count = 0;
        const increment = target / 100;
        const updateCount = () => {
          count += increment;
          if (count < target) {
            el.textContent = Math.floor(count);
            requestAnimationFrame(updateCount);
          } else {
            el.textContent = target;
          }
        };
        updateCount();
      });
    }
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

.animate-fade-in {
  animation: fade-in 1s ease-out;
}

.animate-fade-in-up {
  animation: fade-in-up 0.5s ease-out;
}
</style>
