<template>
  <div class="bg-gray-900 flex items-center justify-center min-h-screen">
    <!-- 登录表单 -->
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-sm transform transition-all duration-700 ease-in-out hover:scale-105 opacity-0 translate-x-full animate-fade-in-right">
      <h2 class="text-2xl font-bold mb-6 text-center text-white">登录 🔐</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="username">用户名</label>
          <input
              class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300"
              id="username"
              type="text"
              placeholder="输入用户名"
          >
        </div>
        <div class="mb-6">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="password">密码</label>
          <input
              class="shadow appearance-none border-2 border-gray-700 rounded w-full py-2 px-3 bg-gray-900 text-white mb-3 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500 transition duration-300"
              id="password"
              type="password"
              placeholder="输入密码"
          >
        </div>
        <div class="flex items-center justify-between">
          <button
              class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 transform hover:scale-110 focus:outline-none focus:shadow-outline transition duration-300"
              type="submit"
          >
            登录 🚀
          </button>
          <a href="#" class="inline-block align-baseline font-bold text-sm text-blue-400 hover:text-blue-600 transition duration-300">
            忘记密码? 🔑
          </a>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  name: 'LoginPage',
  methods: {
    ...mapActions(['login']),
    handleLogin() {
      // 假设验证成功
      this.login().then(() => {
        this.$router.push({ name: 'Home' })
      })
    }
  }
}
</script>

<style>
@keyframes fade-in-right {
  0% {
    opacity: 0;
    transform: translateX(100%);
  }
  100% {
    opacity: 1;
    transform: translateX(0);
  }
}

.animate-fade-in-right {
  animation: fade-in-right 1s forwards;
}
</style>
