<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen animate-fade-in">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 统计卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <!-- 卡片 1 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-blue-400 text-3xl">📈</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">今日攻击</h3>
              <p class="text-2xl animate-number-scroll" data-target="1234">0</p>
            </div>
          </div>
        </div>
        <!-- 卡片 2 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-green-400 text-3xl">✅</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">阻止攻击</h3>
              <p class="text-2xl animate-number-scroll" data-target="1200">0</p>
            </div>
          </div>
        </div>
        <!-- 卡片 3 -->
        <div class="bg-gray-800 p-6 rounded-lg shadow-md transform hover:scale-105 transition duration-500">
          <div class="flex items-center">
            <div class="text-red-400 text-3xl">🚨</div>
            <div class="ml-4">
              <h3 class="text-xl font-bold">异常流量</h3>
              <p class="text-2xl animate-number-scroll" data-target="34">0</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表占位符 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md mb-8">
        <h2 class="text-2xl font-bold mb-4">流量分析图表</h2>
        <div class="h-64 bg-gray-700 flex items-center justify-center animate-pulse">
          <span class="text-gray-400">图表占位符 📊</span>
        </div>
      </div>

      <!-- 活动日志 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">最近活动</h2>
        <ul class="space-y-4">
          <li class="flex items-start animate-fade-in-up">
            <span class="text-blue-400 text-2xl mr-4">🕒</span>
            <div>
              <p class="font-bold">10:30 AM</p>
              <p>检测到 SQL 注入攻击，已成功阻止。💪</p>
            </div>
          </li>
          <li class="flex items-start animate-fade-in-up">
            <span class="text-blue-400 text-2xl mr-4">🕒</span>
            <div>
              <p class="font-bold">09:45 AM</p>
              <p>检测到 XSS 攻击，已成功阻止。🔒</p>
            </div>
          </li>
          <li class="flex items-start animate-fade-in-up">
            <span class="text-blue-400 text-2xl mr-4">🕒</span>
            <div>
              <p class="font-bold">08:20 AM</p>
              <p>异常流量增加，正在监控中。👀</p>
            </div>
          </li>
        </ul>
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
  name: 'WAFDashboard',
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
