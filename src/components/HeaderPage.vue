<template>
  <nav class="bg-gray-800 p-4 shadow-md fixed w-full z-10">
    <div class="container mx-auto flex justify-between items-center">
      <div class="text-2xl font-bold">🔒 Stone ⛰️ WAF 管理面板</div>
      <div class="space-x-4">
        <router-link to="/" class="hover:text-yellow-400 transition duration-300">仪表盘 📊</router-link>
        <router-link to="/attacker-profile" class="hover:text-yellow-400 transition duration-300">攻击者画像追踪 🕵️</router-link>
        <router-link to="/custom-rule" class="hover:text-yellow-400 transition duration-300">规则管理 ⚙️</router-link>
        <router-link to="/log-analysis" class="hover:text-yellow-400 transition duration-300">日志分析 📜</router-link>
        <router-link to="/login" class="hover:text-yellow-400 transition duration-300">登录 🔐</router-link>
        <router-link to="/system-configuration" class="hover:text-yellow-400 transition duration-300">系统配置 ⚙️</router-link>
        <router-link to="/user-management" class="hover:text-yellow-400 transition duration-300">用户管理 👥</router-link>
        <router-link to="/dashboard" class="hover:text-yellow-400 transition duration-300">仪表盘 📊</router-link>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'HeaderPage'
}
</script>
