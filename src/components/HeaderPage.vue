<template>
  <nav class="bg-gradient-to-r from-gray-700 via-gray-800 to-gray-900 p-4 shadow-md fixed w-full z-10 transition-all duration-500">
    <div class="container mx-auto flex justify-between items-center">
      <div class="text-2xl font-bold text-white">🔒 Stone ⛰️ WAF 管理面板</div>
      <div class="space-x-4">
        <template v-if="!isAuthenticated">
          <router-link to="/login" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              登录 🔐
            </button>
          </router-link>
          <router-link to="/setup-2fa" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              注册 📱
            </button>
          </router-link>
        </template>
        <template v-else>
          <router-link to="/" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              主页 🏠
            </button>
          </router-link>
          <router-link to="/attacker-profile" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              攻击者画像追踪 🕵️
            </button>
          </router-link>
          <router-link to="/custom-rule" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              规则管理 ⚙️
            </button>
          </router-link>
          <router-link to="/log-analysis" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              日志分析 📜
            </button>
          </router-link>
          <router-link to="/system-configuration" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              系统配置 ⚙️
            </button>
          </router-link>
          <router-link to="/user-management" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              用户管理 👥
            </button>
          </router-link>
          <router-link to="/dashboard" v-slot="{ navigate }">
            <button
                @click="navigate"
                class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
            >
              仪表盘 📊
            </button>
          </router-link>
          <button
              @click="handleLogout"
              class="text-white hover:text-yellow-400 transition duration-300 transform hover:scale-105"
          >
            登出 🚪
          </button>
        </template>
      </div>
    </div>

    <!-- 添加 PopupNotification 组件 -->
    <PopupNotification
        v-if="showNotification"
        :message="notificationMessage"
        :emoji="notificationEmoji"
        :type="notificationType"
        @close="showNotification = false"
    />
  </nav>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import PopupNotification from './PopupNotification.vue'

export default {
  name: 'HeaderPage',
  components: {
    PopupNotification
  },
  setup() {
    const router = useRouter()
    const store = useStore()

    const showNotification = ref(false)
    const notificationMessage = ref('')
    const notificationEmoji = ref('')
    const notificationType = ref('success')

    const handleLogout = async () => {
      await store.dispatch('logout')
      router.push({ name: 'Home' })

      notificationMessage.value = '登出成功！期待您的再次访问！'
      notificationEmoji.value = '👋'
      notificationType.value = 'success'
      showNotification.value = true

    }

    return {
      isAuthenticated: computed(() => store.state.isAuthenticated),
      handleLogout,
      showNotification,
      notificationMessage,
      notificationEmoji,
      notificationType
    }
  }
}
</script>
