<template>
  <div style="height: 100vh; overflow: hidden;">
    <n-layout has-sider style="height: 100vh;">
      <n-layout-sider
        bordered
        collapse-mode="width"
        :collapsed-width="64"
        :width="240"
        :collapsed="collapsed"
        show-trigger
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
        <div style="padding: 16px; border-bottom: 1px solid var(--n-border-color);">
          <n-space align="center" :size="12">
            <n-icon size="28" color="#18a058">
              <shield-checkmark-outline />
            </n-icon>
            <n-text strong v-if="!collapsed">Stone WAF</n-text>
          </n-space>
        </div>
        
        <n-menu
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :options="menuOptions"
          :value="currentRoute"
          @update:value="handleMenuSelect"
        />
      </n-layout-sider>
      
      <n-layout>
        <n-layout-header 
          bordered 
          style="height: 64px; padding: 0 24px; display: flex; align-items: center;"
        >
          <n-space align="center" justify="space-between" style="width: 100%;">
            <n-breadcrumb>
              <n-breadcrumb-item>{{ currentPageTitle }}</n-breadcrumb-item>
            </n-breadcrumb>
            
            <n-space v-if="isAuthenticated">
              <n-dropdown :options="userMenuOptions" @select="handleUserMenuSelect">
                <n-button text>
                  <n-icon size="18">
                    <person-outline />
                  </n-icon>
                  管理员
                </n-button>
              </n-dropdown>
            </n-space>
          </n-space>
        </n-layout-header>
        
        <n-layout-content 
          content-style="padding: 24px; overflow-y: auto;" 
          style="height: calc(100vh - 64px);"
        >
          <router-view />
        </n-layout-content>
      </n-layout>
    </n-layout>
  </div>
</template>

<script>
import { ref, computed, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { useMessage } from 'naive-ui'
import { 
  HomeOutline,
  ShieldOutline,
  SettingsOutline,
  DocumentTextOutline,
  PersonOutline,
  BarChartOutline,
  EyeOutline,
  ShieldCheckmarkOutline,
  LogOutOutline
} from '@vicons/ionicons5'

export default {
  name: 'MainLayout',
  components: {
    ShieldCheckmarkOutline,
    PersonOutline
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const store = useStore()
    const message = useMessage()
    
    const collapsed = ref(false)
    
    const isAuthenticated = computed(() => store.state.isAuthenticated)
    const currentRoute = computed(() => route.name)
    
    const pageNames = {
      'Home': '主页',
      'WAFDashboard': '仪表盘',
      'AttackerProfileTracking': '攻击者画像',
      'CustomRuleManagement': '规则管理',
      'LogAnalysis': '日志分析',
      'SystemConfiguration': '系统配置',
      'UserManagement': '用户管理'
    }
    
    const currentPageTitle = computed(() => {
      return pageNames[currentRoute.value] || '未知页面'
    })

    const menuOptions = [
      {
        label: '主页',
        key: 'Home',
        icon: () => h(HomeOutline)
      },
      {
        label: '仪表盘',
        key: 'WAFDashboard',
        icon: () => h(BarChartOutline)
      },
      {
        label: '攻击者画像',
        key: 'AttackerProfileTracking',
        icon: () => h(EyeOutline)
      },
      {
        label: '规则管理',
        key: 'CustomRuleManagement',
        icon: () => h(ShieldOutline)
      },
      {
        label: '日志分析',
        key: 'LogAnalysis',
        icon: () => h(DocumentTextOutline)
      },
      {
        label: '系统配置',
        key: 'SystemConfiguration',
        icon: () => h(SettingsOutline)
      },
      {
        label: '用户管理',
        key: 'UserManagement',
        icon: () => h(PersonOutline)
      }
    ]

    const userMenuOptions = [
      {
        label: '登出',
        key: 'logout',
        icon: () => h(LogOutOutline)
      }
    ]

    const handleMenuSelect = (key) => {
      if (!isAuthenticated.value) {
        router.push({ name: 'LoginPage' })
        return
      }
      router.push({ name: key })
    }

    const handleUserMenuSelect = (key) => {
      if (key === 'logout') {
        handleLogout()
      }
    }

    const handleLogout = async () => {
      await store.dispatch('logout')
      message.success('登出成功')
      setTimeout(() => {
        router.push({ name: 'LoginPage' })
      }, 1000)
    }

    return {
      collapsed,
      isAuthenticated,
      currentRoute,
      currentPageTitle,
      menuOptions,
      userMenuOptions,
      handleMenuSelect,
      handleUserMenuSelect,
      handleLogout
    }
  }
}
</script>