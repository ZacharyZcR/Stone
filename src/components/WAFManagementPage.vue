<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <!-- 顶部导航栏 -->
    <nav class="bg-gray-800 p-4 shadow-md fixed w-full z-10">
      <div class="container mx-auto flex justify-between items-center">
        <div class="text-2xl font-bold">🔒 WAF 管理面板</div>
        <div class="space-x-4">
          <a href="#" class="hover:text-yellow-400 transition duration-300">仪表盘 📊</a>
          <a href="#" class="hover:text-yellow-400 transition duration-300">规则管理 ⚙️</a>
          <a href="#" class="hover:text-yellow-400 transition duration-300">日志 📜</a>
          <a href="#" class="hover:text-yellow-400 transition duration-300">设置 🔧</a>
        </div>
      </div>
    </nav>

    <!-- 主体内容 -->
    <div class="container mx-auto px-4 py-24 flex-1">
      <!-- 欢迎信息 -->
      <div class="text-center mb-12">
        <h1 class="text-5xl font-bold mb-4">欢迎使用 WAF 防火墙管理系统 🚀</h1>
        <p class="text-xl text-gray-300">保护您的应用程序免受网络攻击，确保数据安全 🔐。</p>
      </div>

      <!-- 功能卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <!-- 卡片 1 -->
        <div class="bg-gradient-to-r from-green-400 to-blue-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🛡️</span>
            <h3 class="text-2xl font-bold ml-3">实时监控</h3>
          </div>
          <p class="text-gray-100">通过实时监控，您可以快速识别和响应潜在的安全威胁 ⚠️。</p>
        </div>
        <!-- 卡片 2 -->
        <div class="bg-gradient-to-r from-pink-500 to-red-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">⚙️</span>
            <h3 class="text-2xl font-bold ml-3">自定义规则</h3>
          </div>
          <p class="text-gray-100">灵活的规则配置，帮助您根据需要自定义安全策略 🛠️。</p>
        </div>
        <!-- 卡片 3 -->
        <div class="bg-gradient-to-r from-yellow-400 to-orange-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">📊</span>
            <h3 class="text-2xl font-bold ml-3">详细日志</h3>
          </div>
          <p class="text-gray-100">访问详细的日志记录，帮助您分析和审计安全事件 📈。</p>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <footer class="bg-gray-800 text-center py-6">
      <p>&copy; 2024 WAF 防火墙. 保留所有权利. 🌐</p>
    </footer>
  </div>
</template>

<script>
export default {
  name: 'WAFManagementPage'
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
