<template>
  <n-layout-header class="header-container" bordered>
    <n-space justify="space-between" align="center" class="header-content">
      <n-h2 class="header-title">
        <n-text type="primary">🔒 Stone ⛰️ WAF 管理面板</n-text>
      </n-h2>
      
      <n-space size="medium">
        <template v-if="!isAuthenticated">
          <n-button 
            text 
            @click="$router.push('/login')"
            class="nav-button"
          >
            🔐 登录
          </n-button>
          <n-button 
            text 
            @click="$router.push('/setup-2fa')"
            class="nav-button"
          >
            📱 注册
          </n-button>
        </template>
        <template v-else>
          <n-button 
            text 
            @click="$router.push('/')"
            class="nav-button"
          >
            🏠 主页
          </n-button>
          <n-button 
            text 
            @click="$router.push('/attacker-profile')"
            class="nav-button"
          >
            🕵️ 攻击者画像
          </n-button>
          <n-button 
            text 
            @click="$router.push('/custom-rule')"
            class="nav-button"
          >
            ⚙️ 规则管理
          </n-button>
          <n-button 
            text 
            @click="$router.push('/log-analysis')"
            class="nav-button"
          >
            📜 日志分析
          </n-button>
          <n-button 
            text 
            @click="$router.push('/system-configuration')"
            class="nav-button"
          >
            ⚙️ 系统配置
          </n-button>
          <n-button 
            text 
            @click="$router.push('/user-management')"
            class="nav-button"
          >
            👥 用户管理
          </n-button>
          <n-button 
            text 
            @click="$router.push('/dashboard')"
            class="nav-button"
          >
            📊 仪表盘
          </n-button>
          <n-button 
            text 
            @click="handleLogout"
            class="nav-button logout-button"
          >
            🚪 登出
          </n-button>
        </template>
      </n-space>
    </n-space>
  </n-layout-header>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { useMessage } from 'naive-ui'

export default {
  name: 'HeaderPage',
  setup() {
    const router = useRouter()
    const store = useStore()
    const message = useMessage()

    const handleLogout = async () => {
      await store.dispatch('logout')
      
      message.success('👋 登出成功！期待您的再次访问！')
      
      setTimeout(() => {
        router.push({ name: 'Home' })
      }, 1500)
    }

    return {
      isAuthenticated: computed(() => store.state.isAuthenticated),
      handleLogout
    }
  }
}
</script>

<style scoped>
.header-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  backdrop-filter: blur(10px);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 64px;
}

.header-title {
  margin: 0;
  color: white;
}

.nav-button {
  color: white;
  transition: all 0.3s ease;
}

.nav-button:hover {
  color: #fbbf24;
  transform: translateY(-1px);
}

.logout-button:hover {
  color: #f87171;
}
</style>
