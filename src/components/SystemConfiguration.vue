<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <HeaderPage />

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-8 flex-1 mt-16">
      <!-- 配置表单 -->
      <div class="bg-gray-800 p-6 rounded-lg shadow-md">
        <h2 class="text-2xl font-bold mb-4">系统配置设置 ⚙️</h2>
        <form>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="firewallStatus">防火墙状态</label>
            <select class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="firewallStatus">
              <option>启用 ✅</option>
              <option>禁用 ❌</option>
            </select>
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="logLevel">日志级别</label>
            <select class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="logLevel">
              <option>详细 📋</option>
              <option>普通 📄</option>
              <option>简略 🗒️</option>
            </select>
          </div>
          <div class="mb-4">
            <label class="block text-gray-300 text-sm font-bold mb-2" for="alertEmail">警报邮箱</label>
            <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="alertEmail" type="email" placeholder="输入警报接收邮箱 📧">
          </div>
          <div class="flex items-center justify-between">
            <button class="bg-teal-500 text-white px-4 py-2 rounded hover:bg-teal-700" type="button">
              保存配置 💾
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
import HeaderPage from './HeaderPage.vue'
import FooterPage from './FooterPage.vue'

export default {
  name: 'SystemConfiguration',
  components: {
    HeaderPage,
    FooterPage
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
