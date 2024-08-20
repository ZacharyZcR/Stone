<template>
  <div class="bg-gray-100 flex items-center justify-center min-h-screen">
    <!-- 2FA 验证表单 -->
    <div class="bg-white p-8 rounded-lg shadow-lg w-full max-w-sm">
      <h2 class="text-2xl font-bold mb-6 text-center">双因素认证 🔐</h2>
      <p class="mb-4 text-center text-gray-600">请输入您的 2FA 验证码以继续</p>
      <form>
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="2fa-code">验证码</label>
          <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="2fa-code" type="text" placeholder="输入验证码">
        </div>
        <div class="flex items-center justify-between">
          <button class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-700 focus:outline-none focus:shadow-outline" type="button">
            验证 ✅
          </button>
          <a href="#" class="inline-block align-baseline font-bold text-sm text-green-500 hover:text-green-800">
            重新发送验证码 🔄
          </a>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TwoFactorAuth'
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');
</style>
