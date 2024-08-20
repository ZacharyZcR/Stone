<template>
  <div class="bg-gray-900 text-white flex flex-col min-h-screen">
    <HeaderPage />
    <div class="container mx-auto px-4 py-24 flex-1">
      <!-- 欢迎信息 -->
      <div class="text-center mb-12">
        <h1 class="text-5xl font-bold mb-4">欢迎使用 Stone ⛰️ WAF 防火墙管理系统 🚀</h1>
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
        <!-- 卡片 4 -->
        <div class="bg-gradient-to-r from-purple-500 to-indigo-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">📋</span>
            <h3 class="text-2xl font-bold ml-3">规则管理</h3>
          </div>
          <p class="text-gray-100">集中管理所有安全规则，确保最佳实践的实施。</p>
        </div>
        <!-- 卡片 5 -->
        <div class="bg-gradient-to-r from-teal-400 to-cyan-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🕵️</span>
            <h3 class="text-2xl font-bold ml-3">攻击者画像</h3>
          </div>
          <p class="text-gray-100">深入分析攻击者行为，获取详细的攻击者画像。</p>
        </div>
        <!-- 卡片 6 -->
        <div class="bg-gradient-to-r from-gray-500 to-gray-700 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🖥️</span>
            <h3 class="text-2xl font-bold ml-3">系统管理</h3>
          </div>
          <p class="text-gray-100">全面管理系统设置，优化性能和安全性。</p>
        </div>
        <!-- 卡片 7 -->
        <div class="bg-gradient-to-r from-green-500 to-teal-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🔐</span>
            <h3 class="text-2xl font-bold ml-3">IPS 模式</h3>
          </div>
          <p class="text-gray-100">实时入侵防护，自动阻断恶意活动。</p>
        </div>
        <!-- 卡片 8 -->
        <div class="bg-gradient-to-r from-red-500 to-pink-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🔍</span>
            <h3 class="text-2xl font-bold ml-3">IDS 模式</h3>
          </div>
          <p class="text-gray-100">入侵检测，提供详细的安全事件报告。</p>
        </div>
        <!-- 卡片 9 -->
        <div class="bg-gradient-to-r from-blue-500 to-purple-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">⚪⚫</span>
            <h3 class="text-2xl font-bold ml-3">黑白名单</h3>
          </div>
          <p class="text-gray-100">灵活管理访问控制，确保安全与便利的平衡。</p>
        </div>
        <!-- 卡片 10 -->
        <div class="bg-gradient-to-r from-yellow-500 to-orange-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🚫</span>
            <h3 class="text-2xl font-bold ml-3">攻击阻断</h3>
          </div>
          <p class="text-gray-100">自动识别并阻断恶意攻击，保护您的资产。</p>
        </div>
        <!-- 卡片 11 -->
        <div class="bg-gradient-to-r from-indigo-500 to-blue-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">📡</span>
            <h3 class="text-2xl font-bold ml-3">Webhook 告警</h3>
          </div>
          <p class="text-gray-100">实时告警通知，确保您始终掌握最新动态。</p>
        </div>
        <!-- 卡片 12 -->
        <div class="bg-gradient-to-r from-teal-500 to-green-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">⏱️</span>
            <h3 class="text-2xl font-bold ml-3">频率限制</h3>
          </div>
          <p class="text-gray-100">管理请求频率，防止滥用和DDoS攻击。</p>
        </div>
        <!-- 卡片 13 -->
        <div class="bg-gradient-to-r from-purple-500 to-pink-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🔑</span>
            <h3 class="text-2xl font-bold ml-3">2FA 认证</h3>
          </div>
          <p class="text-gray-100">双因素认证，增强账户安全性。</p>
        </div>
        <!-- 卡片 14 -->
        <div class="bg-gradient-to-r from-cyan-500 to-blue-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">🚀</span>
            <h3 class="text-2xl font-bold ml-3">快速部署</h3>
          </div>
          <p class="text-gray-100">快速集成和部署，节省时间和资源。</p>
        </div>
        <!-- 卡片 15 -->
        <div class="bg-gradient-to-r from-orange-500 to-red-500 p-6 rounded-lg shadow-lg transform hover:scale-105 transition duration-500">
          <div class="flex items-center mb-4">
            <span class="text-4xl">⚙️</span>
            <h3 class="text-2xl font-bold ml-3">高度自定义</h3>
          </div>
          <p class="text-gray-100">根据需求进行深度定制，满足各种场景。</p>
        </div>
      </div>
    </div>
    <FooterPage />
  </div>
</template>

<script>
import HeaderPage from './HeaderPage.vue'
import FooterPage from './FooterPage.vue'

export default {
  name: 'HomePage',
  components: {
    HeaderPage,
    FooterPage
  }
}
</script>
