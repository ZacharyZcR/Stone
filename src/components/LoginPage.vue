<template>
  <div class="bg-gray-900 flex items-center justify-center min-h-screen">
    <!-- 登录表单 -->
    <div class="bg-gray-800 p-8 rounded-lg shadow-lg w-full max-w-md transform transition-all duration-700 ease-in-out hover:scale-105 opacity-0 translate-x-full animate-fade-in-right">
      <h2 class="text-3xl font-bold mb-8 text-center text-white">登录 🔐</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-6">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="account">账户</label>
          <input
              v-model="account"
              class="shadow-lg appearance-none border-3 border-gray-700 rounded-lg w-full py-3 px-4 bg-gray-900 text-white leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-300"
              id="account"
              type="text"
              placeholder="输入账户名"
          >
        </div>
        <div class="mb-8">
          <label class="block text-gray-300 text-sm font-bold mb-2" for="code">验证码</label>
          <input
              v-model="code"
              class="shadow-lg appearance-none border-3 border-gray-700 rounded-lg w-full py-3 px-4 bg-gray-900 text-white mb-3 leading-tight focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-300"
              id="code"
              type="text"
              placeholder="输入验证码"
          >
        </div>
        <div class="flex items-center justify-between">
          <button
              class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50 transition duration-300 shadow-lg"
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

    <!-- 弹窗组件 -->
    <PopupNotification
        v-if="showNotification"
        :message="notificationMessage"
        :emoji="notificationEmoji"
        :type="notificationType"
        @close="showNotification = false"
    />
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import PopupNotification from './PopupNotification.vue'

export default {
  name: 'LoginPage',
  components: {
    PopupNotification
  },
  setup() {
    const account = ref('')
    const code = ref('')
    const router = useRouter()
    const store = useStore()

    const showNotification = ref(false)
    const notificationMessage = ref('')
    const notificationEmoji = ref('')
    const notificationType = ref('success')

    const handleLogin = async () => {
      try {
        const success = await store.dispatch('login', {
          account: account.value,
          code: code.value
        })
        if (success) {
          notificationMessage.value = '登录成功！欢迎回来！'
          notificationEmoji.value = '🎉'
          notificationType.value = 'success'
          showNotification.value = true
          setTimeout(() => {
            router.push({ name: 'Home' })
          }, 1500) // 延迟1.5秒后跳转，让用户有时间看到成功消息
        } else {
          throw new Error('登录失败')
        }
      } catch (error) {
        notificationMessage.value = '登录失败，请检查您的验证码和账户信息。'
        notificationEmoji.value = '❌'
        notificationType.value = 'error'
        showNotification.value = true
      }
    }

    return {
      account,
      code,
      handleLogin,
      showNotification,
      notificationMessage,
      notificationEmoji,
      notificationType
    }
  }
}
</script>

<style scoped>
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

.border-3 {
  border-width: 3px;
}
</style>
