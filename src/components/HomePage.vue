<template>
  <div class="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 animate-gradient-x text-white flex flex-col min-h-screen">
    <HeaderPage />
    <div class="container mx-auto px-4 py-24 flex-1">
      <!-- 欢迎信息 -->
      <div class="text-center mb-12">
        <h1 class="text-5xl font-bold mb-4 animate-pulse">欢迎使用 Stone ⛰️ WAF 防火墙管理系统 🚀</h1>
        <p class="text-xl text-gray-300">保护您的应用程序免受网络攻击，确保数据安全 🔐。</p>
      </div>

      <!-- 功能卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <div
            v-for="(card, index) in cards"
            :key="index"
            :class="card.bgClass"
            class="p-6 rounded-lg shadow-lg transform hover:scale-110 hover:rotate-6 hover:shadow-2xl transition-all duration-700 ease-in-out"
        >
          <div class="flex items-center mb-4">
            <span class="text-4xl animate-bounce">{{ card.icon }}</span>
            <h3 class="text-2xl font-bold ml-3">{{ card.title }}</h3>
          </div>
          <p class="text-gray-100">{{ card.description }}</p>
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
  },
  data() {
    return {
      cards: [
        {
          icon: '🛡️',
          title: '实时监控',
          description: '通过实时监控，您可以快速识别和响应潜在的安全威胁 ⚠️。',
          bgClass: 'bg-gradient-to-r from-green-400 to-blue-500'
        },
        {
          icon: '⚙️',
          title: '自定义规则',
          description: '灵活的规则配置，帮助您根据需要自定义安全策略 🛠️。',
          bgClass: 'bg-gradient-to-r from-pink-500 to-red-500'
        },
        {
          icon: '📊',
          title: '详细日志',
          description: '访问详细的日志记录，帮助您分析和审计安全事件 📈。',
          bgClass: 'bg-gradient-to-r from-yellow-400 to-orange-500'
        },
        {
          icon: '📋',
          title: '规则管理',
          description: '集中管理所有安全规则，确保最佳实践的实施。',
          bgClass: 'bg-gradient-to-r from-purple-500 to-indigo-500'
        },
        {
          icon: '🕵️',
          title: '攻击者画像',
          description: '深入分析攻击者行为，获取详细的攻击者画像。',
          bgClass: 'bg-gradient-to-r from-teal-400 to-cyan-500'
        },
        {
          icon: '🖥️',
          title: '系统管理',
          description: '全面管理系统设置，优化性能和安全性。',
          bgClass: 'bg-gradient-to-r from-gray-500 to-gray-700'
        },
        {
          icon: '🔐',
          title: 'IPS 模式',
          description: '实时入侵防护，自动阻断恶意活动。',
          bgClass: 'bg-gradient-to-r from-green-500 to-teal-500'
        },
        {
          icon: '🔍',
          title: 'IDS 模式',
          description: '入侵检测，提供详细的安全事件报告。',
          bgClass: 'bg-gradient-to-r from-red-500 to-pink-500'
        },
        {
          icon: '⚪⚫',
          title: '黑白名单',
          description: '灵活管理访问控制，确保安全与便利的平衡。',
          bgClass: 'bg-gradient-to-r from-blue-500 to-purple-500'
        },
        {
          icon: '🚫',
          title: '攻击阻断',
          description: '自动识别并阻断恶意攻击，保护您的资产。',
          bgClass: 'bg-gradient-to-r from-yellow-500 to-orange-500'
        },
        {
          icon: '📡',
          title: 'Webhook 告警',
          description: '实时告警通知，确保您始终掌握最新动态。',
          bgClass: 'bg-gradient-to-r from-indigo-500 to-blue-500'
        },
        {
          icon: '⏱️',
          title: '频率限制',
          description: '管理请求频率，防止滥用和DDoS攻击。',
          bgClass: 'bg-gradient-to-r from-teal-500 to-green-500'
        },
        {
          icon: '🔑',
          title: '2FA 认证',
          description: '双因素认证，增强账户安全性。',
          bgClass: 'bg-gradient-to-r from-purple-500 to-pink-500'
        },
        {
          icon: '🚀',
          title: '快速部署',
          description: '快速集成和部署，节省时间和资源。',
          bgClass: 'bg-gradient-to-r from-cyan-500 to-blue-500'
        },
        {
          icon: '⚙️',
          title: '高度自定义',
          description: '根据需求进行深度定制，满足各种场景。',
          bgClass: 'bg-gradient-to-r from-orange-500 to-red-500'
        }
      ]
    }
  }
}
</script>
