<template>
  <div class="bg-gray-900 flex items-center justify-center min-h-screen">
    <!-- 登录表单 -->
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-sm">
      <h2 class="text-2xl font-bold mb-6 text-center text-white">登录 🔐</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="username">用户名</label>
          <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="username" type="text" placeholder="输入用户名">
        </div>
        <div class="mb-6">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="password">密码</label>
          <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="输入密码">
        </div>
        <div class="flex items-center justify-between">
          <button class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-700 focus:outline-none focus:shadow-outline" type="submit">
            登录 🚀
          </button>
          <a href="#" class="inline-block align-baseline font-bold text-sm text-blue-400 hover:text-blue-600">
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
